package impl_leveldb

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/syzoj/syzoj-ng-go/app/judge"
)

type judgeRequest struct {
	Tag       int64     `json:"tag"`
	ProblemId uuid.UUID `json:"problem_id"`
	Language  string    `json:"language"`
	Code      string    `json:"code"`
}

type judgeResponse struct {
	Tag    int64                              `json:"tag"`
	Result judge.TaskCompleteInfo `json:"result"`
}

type judgeClient struct {
	ps       *judgeService
	conn     *websocket.Conn
	clientId uuid.UUID
	messages chan *judgeResponse
	entries  map[int64]struct{}
	isClosed int32
}

func (ps *judgeService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var username string
	var password string
	var ok bool
	if username, password, ok = r.BasicAuth(); !ok {
		http.Error(w, "Authentication required", 403)
		return
	}
	var clientId uuid.UUID
	var err error
	if clientId, err = uuid.Parse(username); err != nil {
		http.Error(w, "Invalid username", 403)
		return
	}
	keyClientToken := []byte(fmt.Sprintf("judge.client:%s", clientId))
	var token []byte
	if token, err = ps.db.Get(keyClientToken, nil); err != nil {
		if err != leveldb.ErrNotFound {
			http.Error(w, "Internal server error", 500)
			panic(err)
		} else {
			http.Error(w, "Invalid client id", 403)
			return
		}
	}
	if password != string(token) {
		http.Error(w, "Client token mismatch", 403)
		return
	}

	conn, err := ps.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Warning: WebSocket upgrade failed: %s\n", err)
		return
	}
	client := &judgeClient{
		ps:       ps,
		conn:     conn,
		messages: make(chan *judgeResponse),
		entries:  make(map[int64]struct{}),
		clientId: clientId,
	}
	client.start()
}

func (c *judgeClient) start() {
	go c.readMessage()
	c.work()
}

func (c *judgeClient) readMessage() {
	defer c.shutdown()
	for {
		var message judgeResponse
		if err := c.conn.ReadJSON(&message); err != nil {
			_, ok := err.(*websocket.CloseError)
			if !ok || websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				log.WithFields(c.getFields()).WithField("error", err).Warning("Unexpected websocket close")
			} else {
				log.WithFields(c.getFields()).WithField("error", err).Info("Client disconnected")
			}
			close(c.messages)
			break
		}
		c.messages <- &message
	}
}

func (c *judgeClient) work() {
	c.ps.closeGroup.Add(1)
	defer c.ps.closeGroup.Done()
	if atomic.LoadInt32(&c.ps.isClosed) == 1 {
		return
	}

	if nc, loaded := c.ps.clients.LoadOrStore(c.clientId, c); loaded {
		log.WithFields(c.getFields()).Warning("Double connection detected, closing both connections")
		nc.(*judgeClient).shutdown()
		c.shutdown()
		return
	}
	defer c.ps.clients.Delete(c.clientId)

	log.WithFields(c.getFields()).Info("Client connected")
	defer func() {
		for id := range c.entries {
			c.ps.judgeQueue <- id
		}
		log.WithFields(c.getFields()).Info("Client closed")
	}()

	for c.isClosed == 0 {
		if len(c.entries) < 2 {
			select {
			case id := <-c.ps.judgeQueue:
				c.addTask(id)
				break
			case <-c.ps.closeChan:
				c.shutdown()
				break
			case msg, ok := <-c.messages:
				if !ok {
					c.shutdown()
				} else {
					c.handleMessage(msg)
				}
				break
			}
		} else {
			select {
			case <-c.ps.closeChan:
				c.shutdown()
				break
			case msg, ok := <-c.messages:
				if !ok {
					c.shutdown()
				} else {
					c.handleMessage(msg)
				}
				break
			}
		}
	}
}

func (c *judgeClient) getFields() logrus.Fields {
	return logrus.Fields{
		"module":   "app/judge/traditional",
		"clientId": c.clientId,
	}
}

func (c *judgeClient) handleMessage(msg *judgeResponse) {
	id := msg.Tag
	if _, ok := c.entries[id]; !ok {
		log.WithFields(c.getFields()).WithField("tag", id).Warning("Invalid tag")
		c.shutdown()
		return
	}

	delete(c.entries, id)
	_data, _ := c.ps.submissions.Load(id)
	data := _data.(*submissionEntry)
	log.WithFields(c.getFields()).WithFields(data.getFields()).Info("Received result from client")
	go data.Callback.OnComplete(msg.Result)
	c.ps.submissions.Delete(id)
}

func (c *judgeClient) addTask(id int64) {
	c.entries[id] = struct{}{}
	_data, _ := c.ps.submissions.Load(id)
	data := _data.(*submissionEntry)
	msg := judgeRequest{
		Tag:       id,
		ProblemId: data.ProblemId,
		Language:  data.Language,
		Code:      data.Code,
	}
	if err := c.conn.WriteJSON(&msg); err != nil {
		log.WithFields(c.getFields()).WithField("tag", id).Warning("Failed to send task, closing connection")
		c.shutdown()
		return
	}
	log.WithFields(c.getFields()).WithFields(data.getFields()).Info("Sent task to client")
}

func (c *judgeClient) shutdown() {
	if atomic.CompareAndSwapInt32(&c.isClosed, 0, 1) {
		c.conn.WriteControl(websocket.CloseMessage, nil, time.Now())
		c.conn.Close()
	}
}
