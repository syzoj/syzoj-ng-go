package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/syzoj/syzoj-ng-go/app/lock"

	"github.com/google/uuid"

	model_user "github.com/syzoj/syzoj-ng-go/app/model/user"
)

func (srv *ApiServer) GetUserByUserName(ctx context.Context, userName string) (uuid.UUID, error) {
	row := srv.db.QueryRowContext(ctx, "SELECT id FROM users WHERE name=$1", userName)
	var userId uuid.UUID
	if err := row.Scan(&userId); err != nil {
		if err == sql.ErrNoRows {
			return userId, UnknownUsernameError
		}
		return userId, err
	}
	return userId, nil
}

func (srv *ApiServer) GetUserAuthInfo(ctx context.Context, userId uuid.UUID) (*model_user.UserAuthInfo, error) {
	row := srv.db.QueryRowContext(ctx, "SELECT auth_info FROM users WHERE id=$1", userId[:])
	var authInfoBytes []byte
	if err := row.Scan(&authInfoBytes); err != nil {
		return nil, err
	}
	authInfo := new(model_user.UserAuthInfo)
	if err := json.Unmarshal(authInfoBytes, authInfo); err != nil {
		return nil, err
	}
	return authInfo, nil
}

func (srv *ApiServer) WithUserShared(ctx context.Context, userId uuid.UUID, handler func(context.Context) error) error {
	key := fmt.Sprintf("user:%s", userId)
	return srv.lockManager.WithLockShared(ctx, key, true, func(ctx context.Context, l lock.SharedLock) error {
		return handler(ctx)
	})
}
