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

func (c *Contest) GetRanklistVisibility() string {
	return c.ranklistVisibility
}

func (c *Contest) GetRanklist() []ContestRanklistEntry {
	return c.ranklist_w
}

func (c *Contest) GetPlayer(userId primitive.ObjectID) *ContestPlayer {
	return c.players[userId]
}

type ContestProblemEntry struct {
	ProblemId primitive.ObjectID
	Name      string
}

func (c *Contest) GetProblems() []*ContestProblemEntry {
	entries := make([]*ContestProblemEntry, len(c.problems))
	for i, p := range c.problems {
		entries[i] = new(ContestProblemEntry)
		entries[i].ProblemId = p.ProblemId
		entries[i].Name = p.Name
	}
	return entries
}

func (c *Contest) GetProblemByName(name string) *ContestProblemEntry {
	entryId, found := c.nameToProblems[name]
	if !found {
		return nil
	}
	entry := c.problems[entryId]
	return &ContestProblemEntry{
		ProblemId: entry.ProblemId,
		Name:      entry.Name,
	}
}

func (c *Contest) RegisterPlayer(userId primitive.ObjectID) bool {
	_, ok := c.players[userId]
	if ok {
		c.log.WithField("userId", userId).Debug("RegisterPlayer failed: player already registered")
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
	c.log.WithField("playerId", player.userId).WithField("problem", name).WithField("submissionId", submissionId).Info("Player submission")
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

func (s *ContestPlayerSubmission) GetSubmissionId() primitive.ObjectID {
	return s.submission.Id
}
