package core

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

type Contest struct {
	// only c and id is populated before load()
	c          *Core
	id         primitive.ObjectID
	lock       sync.Mutex
	running    bool
	loaded     bool
	timers     []*time.Timer
	state      string
	closeChan  chan struct{}
	updateChan chan mongo.WriteModel
}

func newState() string {
	var state [16]byte
	rand.Read(state[:])
	return hex.EncodeToString(state[:])
}

// Call exactly once when the contest gets loaded into memory.
func (c *Contest) load(contestModel *model.Contest) {
	c.lock.Lock()
	defer c.lock.Unlock()
	log.WithField("contestId", c.id).Info("Loading contest")
	c.running = contestModel.Running
	c.state = contestModel.State
	c.loaded = true
	c.closeChan = make(chan struct{})
	c.updateChan = make(chan mongo.WriteModel, 100)
	// Bring up the writer
	go c.handleWrites()

	// Load all schedules
	curTime := time.Now()
	for id, scheduleModel := range contestModel.Schedule {
		if scheduleModel.Done {
			continue
		}
		duration := scheduleModel.StartTime.Sub(curTime)
		var f func()
		switch scheduleModel.Type {
		case "start":
			f = func() {
				c.running = true
				model := mongo.NewUpdateOneModel()
				model.SetFilter(bson.D{{"_id", c.id}, {"contest.state", c.state}})
				c.state = newState()
				model.SetUpdate(bson.D{{"$set", bson.D{{fmt.Sprintf("contest.schedule.%d.done", id), true}, {"contest.running", true}, {"contest.state", c.state}}}})
				c.updateChan <- model
				log.WithField("contestId", c.id).Debug("Contest started")
			}
		case "stop":
			f = func() {
				c.running = false
				model := mongo.NewUpdateOneModel()
				model.SetFilter(bson.D{{"_id", c.id}, {"contest.state", c.state}})
				c.state = newState()
				model.SetUpdate(bson.D{{"$set", bson.D{{fmt.Sprintf("contest.schedule.%d.done", id), true}, {"contest.running", false}, {"contest.state", c.state}}}})
				c.updateChan <- model
				log.WithField("contestId", c.id).Debug("Contest stopped")
			}
		}
		f2 := func() {
			c.lock.Lock()
			defer c.lock.Unlock()
			if !c.loaded {
				return
			}
			f()
		}
		if duration <= 0 {
			go f2()
		} else {
			timer := time.AfterFunc(duration, f2)
			c.timers = append(c.timers, timer)
		}
	}
}

// Call exactly once to unload the contest and save state into disk.
func (c *Contest) unload() {
	c.lock.Lock()
	defer c.lock.Unlock()
	log.WithField("contestId", c.id).Info("Unloading contest")
	c.loaded = false
	close(c.closeChan)
	close(c.updateChan)
	for _, timer := range c.timers {
		timer.Stop()
	}
}

func (c *Contest) handleWrites() {
	for {
		var writes []mongo.WriteModel
		select {
		case writeModel, ok := <-c.updateChan:
			if !ok {
				return
			}
			writes = append(writes, writeModel)
		loop:
			for {
				select {
				case writeModel = <-c.updateChan:
					writes = append(writes, writeModel)
				default:
					break loop
				}
			}
			if _, err := c.c.mongodb.Collection("problemset").BulkWrite(context.Background(), writes); err != nil {
				log.WithField("contestId", c.id).WithField("writeCount", len(writes)).Error("Failed to write to contest model: ", err)
				c.unload()
				return
			} else {
				log.WithField("contestId", c.id).WithField("writeCount", len(writes)).Debug("Applied updates")
			}
		}
	}
}
