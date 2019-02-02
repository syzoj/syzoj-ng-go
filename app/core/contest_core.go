package core

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

type Contest struct {
	// only c and id is populated before load()
	c                *Core
	id               primitive.ObjectID
	lock             sync.Mutex
	running          bool
	loaded           bool
	schedules        []*contestSchedule
	scheduleTimer    *time.Timer
	state            string
	closeChan        chan struct{}
	updateChan       chan mongo.WriteModel
	playerUpdateChan chan mongo.WriteModel

	players map[primitive.ObjectID]*ContestPlayer
}

type contestSchedule struct {
	t time.Time
	f func()
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
	c.playerUpdateChan = make(chan mongo.WriteModel, 1000)
	// Bring up the writer
	go handleWrites(c.updateChan, c.c.mongodb.Collection("problemset"), c.id)
	go handleWrites(c.playerUpdateChan, c.c.mongodb.Collection("contest_player"), c.id)

	// Load all schedules
	for id, scheduleModel := range contestModel.Schedule {
		if scheduleModel.Done {
			continue
		}
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
		c.schedules = append(c.schedules, &contestSchedule{
			t: scheduleModel.StartTime,
			f: f,
		})
	}
	c.sortSchedules()

	// Load all players
	c.players = make(map[primitive.ObjectID]*ContestPlayer)
	var (
		cursor mongo.Cursor
		err    error
	)
	if cursor, err = c.c.mongodb.Collection("contest_player").Find(context.Background(), bson.D{{"contest", c.id}}); err != nil {
		log.WithField("contestId", c.id).Warning("Failed to load contest players: ", err)
	}
	for cursor.Next(context.TODO()) {
		var contestPlayerModel *model.ContestPlayer
		if err = cursor.Decode(&contestPlayerModel); err != nil {
			log.WithField("contestId", c.id).Warning("Failed to load a contest player: ", err)
		} else {
			c.loadPlayer(contestPlayerModel)
		}
	}
	log.WithField("contestId", c.id).WithField("playerCount", len(c.players)).Info("Loaded players")
	go c.startSchedule()
}

type scheduleSorter struct {
	*Contest
}

func (s scheduleSorter) Len() int {
	return len(s.schedules)
}
func (s scheduleSorter) Less(i, j int) bool {
	return s.schedules[i].t.Before(s.schedules[j].t)
}
func (s scheduleSorter) Swap(i, j int) {
	schedule := s.schedules[i]
	s.schedules[i] = s.schedules[j]
	s.schedules[j] = schedule
}
func (c *Contest) sortSchedules() {
	sort.Stable(scheduleSorter{c})
}

func (c *Contest) startSchedule() {
	if len(c.schedules) == 0 {
		return
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	if !c.loaded {
		return
	}
	for len(c.schedules) > 0 && time.Now().After(c.schedules[0].t) {
		c.schedules[0].f()
		c.schedules = c.schedules[1:]
	}
	if len(c.schedules) == 0 {
		return
	}
	d := c.schedules[0].t.Sub(time.Now())
	c.scheduleTimer = time.AfterFunc(d, c.startSchedule)
}

// Call exactly once to unload the contest and save state into disk.
func (c *Contest) unload() {
	c.lock.Lock()
	defer c.lock.Unlock()
	log.WithField("contestId", c.id).Info("Unloading contest")
	c.loaded = false
	close(c.closeChan)
	close(c.updateChan)
	close(c.playerUpdateChan)
	if c.scheduleTimer != nil {
		c.scheduleTimer.Stop()
	}
}

func handleWrites(ch chan mongo.WriteModel, coll *mongo.Collection, contestId primitive.ObjectID) {
	for {
		var writes []mongo.WriteModel
		select {
		case writeModel, ok := <-ch:
			if !ok {
				return
			}
			writes = append(writes, writeModel)
		loop:
			for {
				select {
				case writeModel, ok = <-ch:
					if !ok {
						break loop
					}
					writes = append(writes, writeModel)
				default:
					break loop
				}
			}
			if _, err := coll.BulkWrite(context.Background(), writes); err != nil {
				log.WithField("contestId", contestId).WithField("writeCount", len(writes)).Error("Failed to write to contest model: ", err)
				return
			} else {
				log.WithField("contestId", contestId).WithField("writeCount", len(writes)).Debug("Applied updates")
			}
		}
	}
}
