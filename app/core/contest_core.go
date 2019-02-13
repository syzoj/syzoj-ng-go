package core

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/syzoj/syzoj-ng-go/app/model"
	"github.com/syzoj/syzoj-ng-go/util"
)

type Contest struct {
	// only c and id is populated before load()
	c                *Core
	id               primitive.ObjectID
	lock             sync.RWMutex
	loaded           bool
	schedules        []*contestSchedule
	scheduleTimer    *time.Timer
	updateChan       chan mongo.WriteModel
	playerUpdateChan chan mongo.WriteModel
	context          context.Context
	cancelFunc       func()
	wg               sync.WaitGroup

	running              bool
	startTime            time.Time // for calculating penalty time
	judgeInContest       bool
	submissionPerProblem int
	problems             []*model.ProblemEntry
	nameToProblems       map[string]int

	players map[primitive.ObjectID]*ContestPlayer

	// ranklist is NOT covcered by lock; it has its own synchronization mechanism
	ranklist ContestRanklist
	rankcomp ContestRankComp
	// broker is not covered by lock as well
	StatusBroker *util.Broker
}

type contestSchedule struct {
	typ  string
	done bool
	t    time.Time
}

func (c *Contest) serializeState() interface{} {
	state := bson.D{}
	state = append(state, bson.E{"running", c.running})
	state = append(state, bson.E{"start_time", c.startTime})
	state = append(state, bson.E{"judge_in_contest", c.judgeInContest})
	state = append(state, bson.E{"submission_per_problem", c.submissionPerProblem})
	problems := bson.A{}
	for _, problem := range c.problems {
		problemEntry := bson.D{
			{"name", problem.Name},
			{"problem_id", problem.ProblemId},
		}
		problems = append(problems, problemEntry)
	}
	state = append(state, bson.E{"problems", problems})
	schedules := bson.A{}
	for _, schedule := range c.schedules {
		scheduleD := bson.D{
			{"type", schedule.typ},
			{"done", schedule.done},
			{"start_time", schedule.t},
		}
		schedules = append(schedules, scheduleD)
	}
	state = append(state, bson.E{"schedule", schedules})
	switch c.ranklist.(type) {
	case *ContestRealTimeRanklist:
		state = append(state, bson.E{"ranklist_type", "realtime"})
	case ContestDummyRanklist:
	default:
	}
	switch c.rankcomp.(type) {
	case ContestRankCompMaxScoreSum:
		state = append(state, bson.E{"ranklist_comp", "maxsum"})
	case ContestRankCompLastSum:
		state = append(state, bson.E{"ranklist_comp", "lastsum"})
	case ContestRankCompACM:
		state = append(state, bson.E{"ranklist_comp", "acm"})
	case ContestDummyRankComp:
	default:
	}
	return state
}

func (c *Contest) loadState(state *model.ContestState) {
	c.running = state.Running
	c.startTime = state.StartTime
	c.judgeInContest = state.JudgeInContest
	c.submissionPerProblem = int(state.SubmissionPerProblem)
	c.problems = state.Problems
	c.nameToProblems = make(map[string]int)
	for i, problem := range c.problems {
		c.nameToProblems[problem.Name] = i
	}
	for _, scheduleModel := range state.Schedule {
		schedule := &contestSchedule{
			typ:  scheduleModel.Type,
			done: scheduleModel.Done,
			t:    scheduleModel.StartTime,
		}
		c.schedules = append(c.schedules, schedule)
	}
	switch state.RanklistType {
	case "realtime":
		c.ranklist = &ContestRealTimeRanklist{c: c}
	default:
		c.ranklist = ContestDummyRanklist{}
	}
	switch state.RanklistComp {
	case "maxsum":
		c.rankcomp = ContestRankCompMaxScoreSum{}
	case "lastsum":
		c.rankcomp = ContestRankCompLastSum{}
	case "acm":
		c.rankcomp = ContestRankCompACM{}
	default:
		c.rankcomp = ContestDummyRankComp{}
	}
	c.ranklist.Load()
}

func (c *Contest) saveState() {
	model := mongo.NewUpdateOneModel()
	model.SetFilter(bson.D{{"_id", c.id}})
	model.SetUpdate(bson.D{{"$set", bson.D{{"state", c.serializeState()}}}})
	c.updateChan <- model
}

func (c *Contest) runSchedule(s *contestSchedule) {
	switch s.typ {
	case "start":
		c.running = true
		c.startTime = time.Now()
		log.WithField("contestId", c.id).Debug("Contest started")
	case "stop":
		c.running = false
		log.WithField("contestId", c.id).Debug("Contest stopped")
	}
}

// Call exactly once when the contest gets loaded into memory.
func (c *Contest) load(contestModel *model.Contest) {
	c.lock.Lock()
	c.StatusBroker = util.NewBroker()
	go func() {
		defer c.lock.Unlock()
		log.WithField("contestId", c.id).Info("Loading contest")
		c.context, c.cancelFunc = context.WithCancel(context.Background())
		c.updateChan = make(chan mongo.WriteModel, 100)
		c.playerUpdateChan = make(chan mongo.WriteModel, 1000)
		// Bring up the writer
		c.wg.Add(2)
		go handleWrites(c.context, c.updateChan, c.c.mongodb.Collection("contest"), c.id, &c.wg)
		go handleWrites(c.context, c.playerUpdateChan, c.c.mongodb.Collection("contest_player"), c.id, &c.wg)

		c.loaded = true
		c.wg.Add(1)
		c.loadState(&contestModel.State)
		c.sortSchedules()
		go c.startSchedule()

		// Load all players
		c.players = make(map[primitive.ObjectID]*ContestPlayer)
		var (
			cursor *mongo.Cursor
			err    error
		)
		if cursor, err = c.c.mongodb.Collection("contest_player").Find(c.context, bson.D{{"contest", c.id}}); err != nil {
			log.WithField("contestId", c.id).Warning("Failed to load contest players: ", err)
		}
		for cursor.Next(c.context) {
			var contestPlayerModel = new(model.ContestPlayer)
			if err = cursor.Decode(contestPlayerModel); err != nil {
				log.WithField("contestId", c.id).Warning("Failed to load a contest player: ", err)
			} else {
				c.loadPlayer(contestPlayerModel)
			}
		}
		log.WithField("contestId", c.id).WithField("playerCount", len(c.players)).Debug("Loaded players")
	}()
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
	c.lock.Lock()
	defer c.lock.Unlock()
	c.StatusBroker.Broadcast() // trigger potential deadlocks
	if !c.loaded {
		return
	}
	curTime := time.Now()
	var d time.Duration
	var found bool
	for _, schedule := range c.schedules {
		if !schedule.done {
			if curTime.After(schedule.t) {
				c.runSchedule(schedule)
				schedule.done = true
			} else {
				if !found || d > schedule.t.Sub(curTime) {
					found = true
					d = schedule.t.Sub(curTime)
				}
			}
		}
	}
	if found {
		c.scheduleTimer = time.AfterFunc(d, c.startSchedule)
	} else {
		c.scheduleTimer = nil
	}
	c.saveState()
}

// This must only be called from UnloadContest so that the
// corresponding entry in map is also deleted.
// Call exactly once to unload the contest and save state into disk.
// Blocks until unloading is done.
func (c *Contest) unload() {
	c.lock.Lock()
	defer c.lock.Unlock()
	if !c.loaded {
		log.WithField("contestId", c.id).Error("Double unloading contest")
		return
	}
	c.StatusBroker.Broadcast()
	c.StatusBroker.Close()
	log.WithField("contestId", c.id).Info("Unloading contest")
	c.ranklist.Unload()
	for _, player := range c.players {
		c.unloadPlayer(player)
	}
	c.players = nil
	c.loaded = false
	c.cancelFunc()
	close(c.updateChan)
	close(c.playerUpdateChan)
	if c.scheduleTimer != nil {
		c.scheduleTimer.Stop()
	}
	c.wg.Done()
	c.wg.Wait()
}

func handleWrites(ctx context.Context, ch chan mongo.WriteModel, coll *mongo.Collection, contestId primitive.ObjectID, wg *sync.WaitGroup) {
	defer wg.Done()
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
			if _, err := coll.BulkWrite(ctx, writes); err != nil {
				log.WithField("contestId", contestId).WithField("writeCount", len(writes)).Error("Failed to write to contest model: ", err)
			} else {
				log.WithField("contestId", contestId).WithField("writeCount", len(writes)).Debug("Applied updates")
			}
		}
	}
}
