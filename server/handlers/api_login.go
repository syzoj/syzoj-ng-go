package handlers

import (
	"context"
	"database/sql"

	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/crypto/bcrypt"

	"github.com/syzoj/syzoj-ng-go/model"
	"github.com/syzoj/syzoj-ng-go/server"
)

func Get_Login(ctx context.Context) (*model.LoginPage, error) {
	return &model.LoginPage{}, nil
}

func Handle_Login(ctx context.Context, req *model.LoginRequest) (*empty.Empty, error) {
	var err error
	s := server.GetServer(ctx)
	txn, err := s.GetDB().OpenTxn(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to open transaction")
		return nil, server.ErrBusy
	}
	defer txn.Rollback()
	var userRef model.UserRef
	if err = txn.QueryRowContext(ctx, "SELECT id FROM user WHERE username=?", req.GetUserName()).Scan(&userRef); err != nil {
		if err == sql.ErrNoRows {
			return nil, server.ErrUserNotFound
		}
		log.WithError(err).Error("Handle_Login query failed")
		return nil, server.ErrBusy
	}
	var user *model.User
	if user, err = txn.GetUser(ctx, userRef); err != nil || user == nil {
		log.WithError(err).Error("Handle_Login query failed")
		return nil, server.ErrBusy
	}
	if user.Auth == nil {
		log.Warning("Handle_Login: user.Auth is nil")
		return nil, server.ErrBusy
	}
	if bcrypt.CompareHashAndPassword(user.Auth.PasswordHash, []byte(req.GetPassword())) != nil {
		return nil, server.ErrPasswordIncorrect
	}
	return nil, nil
}

func Handle_Register(ctx context.Context, req *model.RegisterRequest) (*empty.Empty, error) {
	var err error
	s := server.GetServer(ctx)
	txn, err := s.GetDB().OpenTxn(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to open transaction")
		return nil, server.ErrBusy
	}
	defer txn.Rollback()
	user := new(model.User)
	if req.UserName == nil || !model.CheckUserName(req.GetUserName()) {
		return nil, server.ErrBadRequest
	}
	user.UserName = req.UserName
	user.Auth = new(model.UserAuth)
	if req.Password == nil {
		return nil, server.ErrBadRequest
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), 0)
	if err != nil {
		log.WithError(err).Error("Failed to generate passowrd")
		return nil, server.ErrBusy
	}
	user.Auth.PasswordHash = passwordHash
	if err = txn.InsertUser(ctx, user); err != nil {
		log.WithError(err).Error("Handle_Register query failed")
		return nil, server.ErrBusy
	}
	if err = txn.Commit(); err != nil {
		log.WithError(err).Error("Handle_Register query failed")
		return nil, server.ErrBusy
	}
	return new(empty.Empty), nil
}
