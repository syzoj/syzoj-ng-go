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
    log.Info(b)
	log.Info(req2)
	log.Info(req2.Test.Test())
    txn, err := s.s.db.OpenTxn(ctx)
    if err != nil {
        return nil, err
    }
    log.Info(txn.GetUser(ctx, req2.GetTest()))
    
	return &empty.Empty{}, nil
}
