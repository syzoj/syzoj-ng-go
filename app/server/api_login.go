package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

func (s *apiServer) Handle_Login(ctx context.Context, req *model.LoginRequest) (*empty.Empty, error) {
	log.Info(req)
	log.Info(req.Test.Test())
	return &empty.Empty{}, nil
}
