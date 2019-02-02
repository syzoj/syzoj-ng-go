package core

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type ContestRegister1 struct {
	UserId    primitive.ObjectID
	ContestId primitive.ObjectID
}
type ContestRegister1Resp struct{}

// Registers to a contest.
// Possible errors:
// * ErrContestNotRunning
// * ErrAlreadyRegistered
// * MongoDB error or context error
func (c *Core) Action_Contest_Register(ctx context.Context, req *ContestRegister1) (*ContestRegister1Resp, error) {
	c.lock.Lock()
	contest, ok := c.contests[req.ContestId]
	c.lock.Unlock()
	if !ok {
		return nil, ErrContestNotRunning
	}
	if err := contest.Register(req.UserId); err != nil {
		return nil, err
	}
	return &ContestRegister1Resp{}, nil
}
