package core

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

type ContestPlayer struct {
	modelId primitive.ObjectID
	userId  primitive.ObjectID
	problems map[string]*ContestPlayerProblem
}
type ContestPlayerProblem struct {
	subscriptions []*contestPlayerSubscription
}
type contestPlayerSubscription struct {
	c *Contest
	submissionId primitive.ObjectID
}
func (s *contestPlayerSubscription) HandleNewScore(done bool, score float64) {
	log.WithField("contestId", s.c.id).WithField("done", done).WithField("score", score).Info("Received submission score")
}

func (c *Contest) loadPlayer(contestPlayerModel *model.ContestPlayer) {
	player := new(ContestPlayer)
	player.modelId = contestPlayerModel.Id
	player.userId = contestPlayerModel.User
	player.problems = make(map[string]*ContestPlayerProblem)
	for i, problemEntryModel := range contestPlayerModel.Problems {
		problemEntry := new(ContestPlayerProblem)
		for _, submissionId := range problemEntryModel.Submissions {
			subscription := &contestPlayerSubscription{
				c: c,
				submissionId: submissionId,
			}
			c.c.SubscribeSubmission(submissionId, subscription)
			problemEntry.subscriptions = append(problemEntry.subscriptions, subscription)
		}
		player.problems[i] = problemEntry
	}
	c.players[contestPlayerModel.User] = player
}

func (p *ContestPlayer) unload() {
	for _, problemEntry := range p.problems {
		for _, subscription := range problemEntry.subscriptions {
			subscription.c.c.UnsubscribeSubmission(subscription.submissionId, subscription)
		}
	}
}

func (c *Contest) Register(UserId primitive.ObjectID) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if !c.loaded {
		return ErrContestNotRunning
	}
	_, ok := c.players[UserId]
	if ok {
		return ErrAlreadyRegistered
	}
	id := primitive.NewObjectID()
	player := new(ContestPlayer)
	player.modelId = id
	player.userId = UserId
	c.players[UserId] = player

	model := mongo.NewInsertOneModel()
	model.SetDocument(bson.D{
		{"_id", id},
		{"contest", c.id},
		{"user", UserId},
	})
	c.playerUpdateChan <- model
	return nil
}
