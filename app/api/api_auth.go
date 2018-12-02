package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/syzoj/syzoj-ng-go/app/lock"

	"github.com/google/uuid"
	model_session "github.com/syzoj/syzoj-ng-go/app/model/session"
	model_user "github.com/syzoj/syzoj-ng-go/app/model/user"
)

type AuthRegisterRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (srv *ApiServer) HandleAuthRegister(w http.ResponseWriter, r *http.Request) {
	var req AuthRegisterRequest
	var err error
	defer func() {
		if err != nil {
			writeError(w, err)
		}
	}()
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}
	if err = srv.lockManager.WithLockExclusive(r.Context(), fmt.Sprintf("username:%s", req.UserName), false, func(ctx context.Context, l lock.ExclusiveLock) error {
		var cnt int
		if err = srv.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE name=$1", req.UserName).Scan(&cnt); err != nil {
			return err
		}
		if cnt >= 1 {
			return DuplicateUserNameError
		}

		userId := uuid.New()
		authInfo := model_user.PasswordAuth(req.Password)
		authInfoBytes, err := json.Marshal(authInfo)
		if err != nil {
			return err
		}
		_, err = srv.db.ExecContext(
			ctx,
			"INSERT INTO users (id, name, auth_info, can_login) VALUES ($1, $2, $3, $4)",
			userId[:],
			req.UserName,
			authInfoBytes,
			true,
		)
		return err
	}); err != nil {
		return
	} else {
		writeResponse(w, "success")
	}
}

type AuthLoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (srv *ApiServer) HandleAuthLogin(w http.ResponseWriter, r *http.Request) {
	var req AuthLoginRequest
	var err error
	defer func() {
		if err != nil {
			writeError(w, err)
		}
	}()
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}
	if err = srv.WithSessionExclusive(r.Context(), w, r, func(ctx context.Context, sessId uuid.UUID, sess *model_session.Session) error {
		userId, err := srv.GetUserByUserName(ctx, req.UserName)
		if err != nil {
			return err
		}
		// Note that user name may change after this operation
		return srv.WithUserShared(ctx, userId, func(ctx context.Context) error {
			authInfo, err := srv.GetUserAuthInfo(ctx, userId)
			if err != nil {
				return err
			}
			if authInfo.UseTwoFactor {
				return TwoFactorNotSupportedError
			}
			if authInfo.PasswordInfo.Verify(req.Password) != nil {
				return PasswordIncorrectError
			}
			sess.AuthUserId = userId
			return srv.SaveSession(ctx, sessId, sess)
		})
	}); err != nil {
		return
	} else {
		writeResponse(w, "success")
	}
}
