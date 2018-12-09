package traditional

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

var log = logrus.StandardLogger()

type traditionalJudgeService struct {
	ps          problemsetCallbackService
	db          *leveldb.DB
	judgeQueue  chan int64
	count       int64
	submissions sync.Map
	upgrader    websocket.Upgrader
	closeChan   chan struct{}
	closeLock   sync.RWMutex
	clients     sync.Map
}

type traditionalSubmissionEntry struct {
	ProblemsetId uuid.UUID
	SubmissionId uuid.UUID
	ProblemId    uuid.UUID
	Language     string
	Code         string
	Result       TraditionalSubmissionResult
}

func (e *traditionalSubmissionEntry) getFields() logrus.Fields {
	return logrus.Fields{
		"problemsetId": e.ProblemsetId,
		"submissionId": e.SubmissionId,
		"problemId":    e.ProblemId,
	}
}

type traditionalJudgeClient struct {
	ps       *traditionalJudgeService
	conn     *websocket.Conn
	clientId uuid.UUID
	messages chan *TraditionalJudgeResponse
	entries  map[int64]struct{}
	isClosed int32
}

func NewTraditionalJudgeService(db *leveldb.DB) (TraditionalJudgeService, error) {
	s := &traditionalJudgeService{
		judgeQueue: make(chan int64, 1000),
		closeChan:  make(chan struct{}),
		upgrader:   websocket.Upgrader{},
		db:         db,
	}
	return s, nil
}

func (ps *traditionalJudgeService) RegisterProblemsetService(s problemsetCallbackService) {
	if ps.ps != nil {
		panic("traditionalJudgeService: RegisterProblemsetService called twice")
	}
	ps.ps = s
}

func (ps *traditionalJudgeService) QueueSubmission(problemsetId uuid.UUID, submissionId uuid.UUID, submissionData *TraditionalSubmission) error {
	var id = atomic.AddInt64(&ps.count, 1)
	entry := &traditionalSubmissionEntry{
		ProblemsetId: problemsetId,
		SubmissionId: submissionId,
		ProblemId:    submissionData.ProblemId,
		Language:     submissionData.Language,
		Code:         submissionData.Code,
	}
	ps.submissions.Store(id, entry)
	select {
	case ps.judgeQueue <- id:
		log.WithFields(entry.getFields()).Debug("Queued submission")
		return nil
	default:
		return ErrQueueFull
	}
}

func (ps *traditionalJudgeService) ack(id int64) {
	_data, _ := ps.submissions.Load(id)
	data := _data.(*traditionalSubmissionEntry)
	go func() {
		ps.closeLock.RLock()
		defer ps.closeLock.RUnlock()
		select {
		case <-ps.closeChan:
			return
		default:
		}
		for {
			if _, err := ps.ps.InvokeProblemset(data.ProblemsetId, &TraditionalSubmissionResultMessage{
				SubmissionId: data.SubmissionId,
				Result:       data.Result,
			}); err != nil {
				log.WithFields(data.getFields()).Warning("Failed to deliver submission result; redelivering in 1 minute\n")
				select {
				case <-ps.closeChan:
					return
				case <-time.After(time.Minute):
				}
			} else {
				log.WithFields(data.getFields()).Debug("Successfully delivered submission result")
				ps.submissions.Delete(id)
				break
			}
		}
	}()
}

func (ps *traditionalJudgeService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	client := &traditionalJudgeClient{ps: ps, conn: conn, messages: make(chan *TraditionalJudgeResponse), entries: make(map[int64]struct{}), clientId: clientId}
	client.start()
}

func (c *traditionalJudgeClient) start() {
	go c.readMessage()
	go c.work()
}

func (c *traditionalJudgeClient) readMessage() {
	defer c.shutdown()
	for {
		var message TraditionalJudgeResponse
		if err := c.conn.ReadJSON(&message); err != nil {
			if atomic.LoadInt32(&c.isClosed) == 1 {
				break
			}
			_, ok := err.(*websocket.CloseError)
			if !ok || websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				log.WithFields(c.getFields()).WithField("error", err).Warning("Unexpected websocket close")
			} else {
				log.WithFields(c.getFields()).WithField("error", err).Info("Client disconnected")
			}
			break
		}
		// TODO: solve race condition here; c.message may have been closed
		c.messages <- &message
	}
}

func (c *traditionalJudgeClient) work() {
	c.ps.closeLock.RLock()
	defer c.ps.closeLock.RUnlock()
	if nc, loaded := c.ps.clients.LoadOrStore(c.clientId, c); loaded {
		log.WithFields(c.getFields()).Warning("Double connection detected, closing both connections")
		nc.(*traditionalJudgeClient).shutdown()
		c.shutdown()
		return
	}
	defer c.ps.clients.Delete(c.clientId)
	log.WithFields(c.getFields()).Info("Client connected")
loop:
	for {
		if len(c.entries) < 2 {
			select {
			case id := <-c.ps.judgeQueue:
				c.addTask(id)
				break
			case <-c.ps.closeChan:
				c.shutdown()
				break loop
			case msg, ok := <-c.messages:
				if !ok {
					break loop
				}
				c.handleMessage(msg)
				break
			}
		} else {
			select {
			case <-c.ps.closeChan:
				c.shutdown()
				break loop
			case msg, ok := <-c.messages:
				if !ok {
					break loop
				}
				c.handleMessage(msg)
				break
			}
		}
	}

	// Return entries to global queue
	for id := range c.entries {
		c.ps.judgeQueue <- id
	}
	log.WithFields(c.getFields()).Info("Client cleanup done")
}

func (c *traditionalJudgeClient) getFields() logrus.Fields {
	return logrus.Fields{
		"module":   "app/judge/traditional",
		"clientId": c.clientId,
	}
}

func (c *traditionalJudgeClient) handleMessage(msg *TraditionalJudgeResponse) {
	id := msg.Tag
	if _, ok := c.entries[id]; !ok {
		log.WithFields(c.getFields()).WithField("tag", id).Warning("Invalid tag")
		c.shutdown()
		return
	}
	delete(c.entries, id)
	_data, _ := c.ps.submissions.Load(id)
	data := _data.(*traditionalSubmissionEntry)
	data.Result = msg.Result
	log.WithFields(c.getFields()).WithFields(data.getFields()).Info("Received result from client")
	c.ps.ack(id)
}

func (c *traditionalJudgeClient) addTask(id int64) {
	c.entries[id] = struct{}{}
	_data, _ := c.ps.submissions.Load(id)
	data := _data.(*traditionalSubmissionEntry)
	msg := TraditionalJudgeMessage{
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

func (c *traditionalJudgeClient) shutdown() {
	if atomic.CompareAndSwapInt32(&c.isClosed, 0, 1) {
		close(c.messages)
		c.conn.WriteControl(websocket.CloseMessage, nil, time.Now())
		c.conn.Close()
	}
}

func (ps *traditionalJudgeService) Close() error {
	close(ps.closeChan)
	ps.closeLock.Lock()
	defer ps.closeLock.Unlock()
	// TODO: Save queue to LevelDB
	return nil
}
