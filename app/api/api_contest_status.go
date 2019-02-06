package api

import (
	"io"

	"github.com/gorilla/websocket"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/valyala/fastjson"

	"github.com/syzoj/syzoj-ng-go/util"
)

type contestStatusContext struct {
	srv       *ApiServer
	contestId primitive.ObjectID
	userId    primitive.ObjectID

	wsConn *websocket.Conn
}

func Handle_Contest_Status(c *ApiContext) (apiErr ApiError) {
	var err error
	vars := c.Vars()
	contestId := DecodeObjectID(vars["contest_id"])
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	var ct contestStatusContext
	server := c.Server()
	ct.srv = server
	ct.contestId = contestId
	ct.userId = c.Session.AuthUserUid
	ct.wsConn, err = c.UpgradeWebSocket()
	if err != nil {
		panic(err)
	}
	server.wsConnMutex.Lock()
	server.wsConn[ct.wsConn] = struct{}{}
	server.wsConnMutex.Unlock()
	ct.run()
	return
}

func (c *contestStatusContext) run() {
	var err error
	defer c.close()
	arena := new(fastjson.Arena)
	subscriber := util.NewChanSubscriber()
	contest := c.srv.c.GetContestR(c.contestId)
	if contest == nil {
		return
	}
	contest.StatusBroker.Subscribe(subscriber)
	defer contest.StatusBroker.Unsubscribe(subscriber)
	player := contest.GetPlayer(c.userId)
	if player != nil {
		player.Broker.Subscribe(subscriber)
		defer player.Broker.Unsubscribe(subscriber)
	}
	subscriber.Notify()
	contest.RUnlock()
	for {
		select {
		case <-c.srv.ctx.Done():
			return
		case <-subscriber.C:
		}
		arena.Reset()
		msg := arena.NewObject()
		func() {
			contest = c.srv.c.GetContestR(c.contestId)
			if contest == nil {
				return
			}
			defer contest.RUnlock()
			if contest.Running() {
				msg.Set("running", arena.NewTrue())
			} else {
				msg.Set("running", arena.NewFalse())
			}
		}()

		var w io.WriteCloser
		w, err = c.wsConn.NextWriter(websocket.TextMessage)
		if err != nil {
			log.Warning("Failed to write to WebSocket: ", err)
			return
		}
		_, err = w.Write(msg.MarshalTo(nil))
		if err != nil {
			log.Warning("Failed to write to WebSocket: ", err)
			return
		}
		err = w.Close()
	}
}

func (c *contestStatusContext) close() {
	c.wsConn.Close()
	c.srv.wsConnMutex.Lock()
	if c.srv.wsConn != nil {
		delete(c.srv.wsConn, c.wsConn)
	}
	c.srv.wsConnMutex.Unlock()
}
