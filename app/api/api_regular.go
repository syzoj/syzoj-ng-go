package api

import (
	"net/http"

	"github.com/syzoj/syzoj-ng-go/app/problemset"

	"github.com/google/uuid"
	"github.com/syzoj/syzoj-ng-go/app/session"
)

type RegularProblemsetCreateRequest struct{}
type RegularProblemsetCreateResponse struct {
	problemsetId uuid.UUID
}

func (srv *ApiServer) HandleRegularProblemsetCreate(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			writeError(w, err)
		}
	}()

	var sess *session.Session
	if _, sess, err = srv.ensureSession(w, r); err != nil {
		return
	}
	if sess.AuthUserId == defaultUserId {
		err = NotLoggedInError
		return
	}

	preq := &problemset.RegularCreateRequest{
		OwnerId: sess.AuthUserId,
	}
	var id uuid.UUID
	if id, err = srv.problemsetService.NewProblemset("regular", preq); err != nil {
		return
	}
	writeResponse(w, RegularProblemsetCreateResponse{
		problemsetId: id,
	})
}
