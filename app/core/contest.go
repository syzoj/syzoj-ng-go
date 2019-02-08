package core

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

type ContestOptions struct {
	Rules     ContestRules
	StartTime time.Time
	Duration  time.Duration
}

type ContestRules struct {
	JudgeInContest      bool
	SeeResult           bool
	RejudgeAfterContest bool
	RanklistType        string // realtime, defer, ""
	RanklistVisibility  string
	RanklistComp        string // maxsum, lastsum, acm
}

var ErrInvalidOptions = errors.New("Invalid contest options")

// No locking required
func (c *Core) CreateContest(ctx context.Context, id primitive.ObjectID, options *ContestOptions) (err error) {
	var result *mongo.UpdateResult
	schedule := bson.A{}
	if options.Duration <= 0 {
		log.Debug("CreateContest: Invalid contest options: Duration <= 0")
		return ErrInvalidOptions
	}
	switch options.Rules.RanklistType {
	case "realtime":
	case "":
	default:
		return ErrInvalidOptions
	}
	switch options.Rules.RanklistComp {
	case "maxsum":
	case "lastsum":
	case "acm":
	default:
		return ErrInvalidOptions
	}
	schedule = append(schedule, bson.D{
		{"type", "start"},
		{"done", false},
		{"start_time", options.StartTime},
	})
	schedule = append(schedule, bson.D{
		{"type", "stop"},
		{"done", false},
		{"start_time", options.StartTime.Add(options.Duration)},
	})
	contestD := bson.D{
		{"running", false},
		{"schedule", schedule},
		{"state", ""},
		{"ranklist_type", options.Rules.RanklistType},
		{"ranklist_comp", options.Rules.RanklistComp},
		{"start_time", options.StartTime},
	}
	if result, err = c.mongodb.Collection("problemset").UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{{"$set", bson.D{{"contest", contestD}}}}); err != nil {
		return
	}
	if result.MatchedCount == 0 {
		return errors.New("Problemset not found")
	}
	if err = c.LoadContest(id); err != nil {
		return
	}
	return
}

// Loads a contest into memory. Blocks until the contest is in memory.
// This is currently slow and blocks other operations.
func (c *Core) LoadContest(id primitive.ObjectID) (err error) {
	var contestModel model.Problemset
	if err = c.mongodb.Collection("problemset").FindOne(c.context, bson.D{{"_id", id}}, mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"contest", 1}})).Decode(&contestModel); err != nil {
		return
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, found := c.contests[id]; found {
		log.WithField("contestId", id).Warning("LoadContest: contest already loaded")
		return errors.New("Contest already loaded")
	}
	contest := &Contest{c: c, id: id}
	c.contests[id] = contest
	contest.load(&contestModel)
	return
}

func (c *Core) initContestLocked(ctx context.Context) (err error) {
	c.contests = make(map[primitive.ObjectID]*Contest)
	var cursor *mongo.Cursor
	if cursor, err = c.mongodb.Collection("problemset").Find(ctx, bson.D{{"contest", bson.D{{"$exists", true}}}}, mongo_options.Find().SetProjection(bson.D{{"_id", 1}, {"contest", 1}, {"problems", 1}})); err != nil {
		return
	}
	for cursor.Next(ctx) {
		var contestModel model.Problemset
		if err = cursor.Decode(&contestModel); err != nil {
			return
		}
		contest := &Contest{c: c, id: contestModel.Id}
		c.contests[contestModel.Id] = contest
		contest.load(&contestModel)
	}
	if err = cursor.Err(); err != nil {
		return
	}
	return
}

func (c *Core) unloadAllContestsLocked() error {
	var wg sync.WaitGroup
	for id, contest := range c.contests {
		wg.Add(1)
		go func(contest *Contest) {
			contest.unload()
			wg.Done()
		}(contest)
		delete(c.contests, id)
	}
	wg.Wait()
	return nil
}

// Unloads a contest from memory. Waits until the contest is ready to be reloaded.
// This is currently slow and blocks other operations.
func (c *Core) UnloadContest(id primitive.ObjectID) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	contest := c.contests[id]
	if contest == nil {
		log.WithField("contestId", id).Warning("UnloadContest: contest not loaded")
		return errors.New("Contest is not loaded")
	}
	delete(c.contests, id)
	contest.unload()
	return nil
}

// Call RUnlock() if return value is not nil
func (c *Core) GetContestR(id primitive.ObjectID) *Contest {
	c.lock.RLock()
	contest := c.contests[id]
	c.lock.RUnlock()
	if contest == nil {
		return nil
	}
	contest.lock.RLock()
	if !contest.loaded {
		contest.lock.RUnlock()
		panic("Core: contest unloaded without removing from map")
	}
	return contest
}

// Call Unlock() if return value is not nil
func (c *Core) GetContestW(id primitive.ObjectID) *Contest {
	c.lock.RLock()
	contest := c.contests[id]
	c.lock.RUnlock()
	if contest == nil {
		return nil
	}
	contest.lock.Lock()
	if !contest.loaded {
		contest.lock.Unlock()
		panic("Core: contest unloaded without removing from map")
	}
	return contest
}
