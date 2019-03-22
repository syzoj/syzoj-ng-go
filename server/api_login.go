package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/proto"

	"github.com/syzoj/syzoj-ng-go/model"
)

func (s *apiServer) Handle_Login(ctx context.Context, req *model.LoginRequest) (*empty.Empty, error) {
    b, err := proto.Marshal(req)
    if err != nil {
        return nil, err
    }
    req2 := new(model.LoginRequest)
    err = proto.Unmarshal(b, req2)
    if err != nil {
        return nil, err
    }
    txn, err := s.s.db.OpenTxn(ctx)
    if err != nil {
        return nil, err
    }
    defer txn.Rollback()
    user, err := txn.GetUser(ctx, req2.GetTest())
    log.Info(user, err)
    log.Info(user.Auth)
    user.Auth = &model.UserAuth{PasswordHash: []byte{1, 2, 3, 4}}
    log.Info(txn.SetUser(ctx, req2.GetTest(), user))
    log.Info(txn.Commit())
	return &empty.Empty{}, nil
}
