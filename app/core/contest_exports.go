package core

import (
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
