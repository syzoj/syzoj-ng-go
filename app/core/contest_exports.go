package core

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func (c *Contest) Unlock() {
	c.lock.Unlock()
}

func (c *Contest) RUnlock() {
	c.lock.RUnlock()
}

func (c *Contest) Running() bool {
	return c.running
}

func (c *Contest) GetRankComp() ContestRankComp {
	return c.rankcomp
}

func (c *Contest) GetPlayer(userId primitive.ObjectID) *ContestPlayer {
	return c.players[userId]
}

func (c *Contest) RegisterPlayer(userId primitive.ObjectID) bool {
	_, ok := c.players[userId]
	if ok {
		log.WithField("userId", userId).Debug("RegisterPlayer failed: player already registered")
		return false
	}
	id := primitive.NewObjectID()
	player := new(ContestPlayer)
	player.modelId = id
	player.userId = userId
	player.problems = make(map[string]*ContestPlayerProblem)
	c.players[userId] = player

	model := mongo.NewInsertOneModel()
	model.SetDocument(bson.D{
		{"_id", id},
		{"contest", c.id},
		{"user", userId},
	})
	c.playerUpdateChan <- model
	return true
}

// Possible errors:
// * ErrGeneral
// * ErrTooManySubmissions
func (c *Contest) PlayerSubmission(player *ContestPlayer, name string, submissionId primitive.ObjectID) error {
	if !checkName(name) {
		log.Debug("Contest.PlayerSubmission: Invalid problem name")
		return ErrGeneral
	}
	if _, ok := player.problems[name]; !ok {
		problem := new(ContestPlayerProblem)
		player.problems[name] = problem
		model := mongo.NewUpdateOneModel()
		model.SetFilter(bson.D{{"_id", player.modelId}})
		model.SetUpdate(bson.D{{"$set", bson.D{{"problems." + name, bson.D{{"submissions", bson.A{}}}}}}})
		c.playerUpdateChan <- model
	}
	problemEntry := player.problems[name]
	if len(problemEntry.submissions) >= c.submissionPerProblem {
		return ErrTooManySubmissions
	}
	submission := c.c.GetSubmission(submissionId)
	playerSubmission := &ContestPlayerSubmission{
		c:           c,
		userId:      player.userId,
		submission:  submission,
		penaltyTime: time.Now().Sub(c.startTime),
	}
	submission.Broker.Subscribe(playerSubmission)
	problemEntry.submissions = append(problemEntry.submissions, playerSubmission)
	playerSubmission.Notify()
	if c.judgeInContest {
		go c.c.EnqueueSubmission(submissionId)
	}
	model := mongo.NewUpdateOneModel()
	model.SetFilter(bson.D{{"_id", player.modelId}})
	model.SetUpdate(bson.D{{"$push", bson.D{{"problems." + name + ".submissions", bson.D{
		{"submission_id", submissionId},
		{"penalty_time", playerSubmission.penaltyTime},
	}}}}})
	c.playerUpdateChan <- model
	return nil
}

func (p *ContestPlayer) GetProblems() map[string]*ContestPlayerProblem {
	return p.problems
}

func (p *ContestPlayerProblem) GetSubmissions() []*ContestPlayerSubmission {
	return p.submissions
}
