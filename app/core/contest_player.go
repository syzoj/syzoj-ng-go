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
}

func (c *Contest) loadPlayer(contestPlayerModel *model.ContestPlayer) {
	player := new(ContestPlayer)
	player.modelId = contestPlayerModel.Id
	player.userId = contestPlayerModel.User
	c.players[contestPlayerModel.User] = player
}

func (c *Contest) Register(UserId primitive.ObjectID) error {
	c.lock.Lock()
	defer c.lock.Unlock()
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
