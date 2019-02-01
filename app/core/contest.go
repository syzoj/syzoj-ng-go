package core

import (
	"context"
	"errors"
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
	RanklistType        string
	RanklistVisibility  string
}

var ErrInvalidOptions = errors.New("Invalid contest options")

func (c *Core) CreateContest(ctx context.Context, id primitive.ObjectID, options *ContestOptions) (err error) {
	var result *mongo.UpdateResult
	if result, err = c.mongodb.Collection("problemset").UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{{"$set", bson.D{{"contest", bson.D{
		{"running", false},
		{"schedule", bson.A{}},
		{"state", ""},
	}}}}}); err != nil {
		return
	}
	if result.MatchedCount == 0 {
		return errors.New("Problemset not found")
	}
	go func() {
		var contestModel model.Problemset
		if err = c.mongodb.Collection("problemset").FindOne(context.Background(), bson.D{{"_id", id}}, mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"contest", 1}})).Decode(&contestModel); err != nil {
			return
		}
		c.lock.Lock()
		defer c.lock.Unlock()
		if _, found := c.contests[id]; found {
			log.WithField("contestId", id).Warning("Called CreateContest() on existing contest")
			c.unloadContest(id)
		}
		c.loadContest(contestModel.Id, &contestModel.Contest)
	}()
	return
}

func (c *Core) initContest(ctx context.Context) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.contests = make(map[primitive.ObjectID]*Contest)
	var cursor mongo.Cursor
	if cursor, err = c.mongodb.Collection("problemset").Find(ctx, bson.D{{"contest", bson.D{{"$exists", true}}}}, mongo_options.Find().SetProjection(bson.D{{"_id", 1}, {"contest", 1}})); err != nil {
		return
	}
	for cursor.Next(ctx) {
		var contestModel model.Problemset
		if err = cursor.Decode(&contestModel); err != nil {
			return
		}
		c.loadContest(contestModel.Id, &contestModel.Contest)
	}
	if err = cursor.Err(); err != nil {
		return
	}
	return
}

func (c *Core) ReloadContest(ctx context.Context, id primitive.ObjectID) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	var contestModel model.Problemset
	if err = c.mongodb.Collection("problemset").FindOne(ctx, bson.D{{"contest", bson.D{{"$exists", true}}}, {"_id", id}}).Decode(&contestModel); err != nil {
		return
	}
	c.reloadContest(id, &contestModel.Contest)
	return
}

func (c *Core) unloadAllContests() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	for id := range c.contests {
		c.unloadContest(id)
	}
	return nil
}

func (c *Core) unloadContest(id primitive.ObjectID) {
	contest := c.contests[id]
	if contest == nil {
		log.WithField("contestId", id).Warning("unloadContest: contest doesn't exist")
		return
	}
	contest.unload()
	delete(c.contests, id)
}

func (c *Core) loadContest(id primitive.ObjectID, contestModel *model.Contest) {
	contest := &Contest{c: c, id: id}
	c.contests[id] = contest
	contest.load(contestModel)
}

func (c *Core) reloadContest(id primitive.ObjectID, contestModel *model.Contest) {
	if _, ok := c.contests[id]; ok {
		c.unloadContest(id)
	}
	c.loadContest(id, contestModel)
}
