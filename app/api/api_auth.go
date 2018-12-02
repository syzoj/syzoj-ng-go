package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	model_session "github.com/syzoj/syzoj-ng-go/app/model/session"
)

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
