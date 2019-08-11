package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	httplib "github.com/syzoj/syzoj-ng-go/lib/fasthttp"
	"github.com/syzoj/syzoj-ng-go/util"
	"github.com/valyala/fasthttp"
)

func (app *App) getProblemShort(ctx *fasthttp.RequestCtx) {
	name := ctx.UserValue("name").(string)
	all := ctx.UserValue("all").(string)
	var problemUid string
	err := app.dbProblem.QueryRowContext(ctx, "SELECT problem_uid FROM problem_short WHERE name=?", name).Scan(&problemUid)
	if err != nil {
		if err == sql.ErrNoRows {
			httplib.NotFound(ctx, fmt.Errorf("Not found"))
		} else {
			httplib.SendInternalError(ctx, err)
		}
		return
	}
	ctx.Response.Header.Set("Location", fmt.Sprintf("/problem/id/%s%s", problemUid, all))
	ctx.Response.SetStatusCode(fasthttp.StatusMovedPermanently)
}

type ProblemSetShortNameRequest struct {
	Name string `json:"name"`
}

func (app *App) postProblemSetShort(ctx *fasthttp.RequestCtx) {
	problemUid := ctx.UserValue("uid").(string)
	var req ProblemSetShortNameRequest
	if err := httplib.ReadBodyJSON(ctx, &req); err != nil {
		httplib.BadRequest(ctx, err)
		return
	}
	name := req.Name
	if len(name) == 0 {
		httplib.BadRequest(ctx, fmt.Errorf("Missing name field"))
		return
	}
	tx, err := app.dbProblem.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	var orgName string
	err = tx.QueryRowContext(ctx, "SELECT name FROM problem_short WHERE problem_uid=?", problemUid).Scan(&orgName)
	if err != nil && err != sql.ErrNoRows {
		tx.Rollback()
		httplib.SendInternalError(ctx, err)
		return
	}
	log.Infof("orgName=%s name=%s", orgName, name)
	if err == sql.ErrNoRows {
		if _, err := tx.ExecContext(ctx, "INSERT INTO problem_short (problem_uid, name) VALUES (?, ?)", problemUid, name); err != nil {
			tx.Rollback()
			httplib.SendInternalError(ctx, err)
			return
		}
	} else if orgName == name {
		tx.Rollback()
	} else {
		if _, err := tx.ExecContext(ctx, "UPDATE problem_short SET name=? WHERE problem_uid=?", name, problemUid); err != nil {
			tx.Rollback()
			httplib.SendInternalError(ctx, err)
			return
		}
	}
	if err := tx.Commit(); err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	ctx.SetStatusCode(204)
}

type ProblemInfo struct {
	Title     string `json:"title"`
	Statement string `json:"statement"`
}

type ProblemNewRequest struct {
	Info *ProblemInfo `json:"info,omitempty"`
	Tags []string     `json:"tags,omitempty"`
}

func (app *App) postProblemNew(ctx *fasthttp.RequestCtx) {
	var req ProblemNewRequest
	if err := httplib.ReadBodyJSON(ctx, &req); err != nil {
		httplib.BadRequest(ctx, err)
		return
	}
	if req.Info == nil {
		httplib.BadRequest(ctx, fmt.Errorf("Missing info field"))
		return
	}
	sess, err := app.getSession(ctx)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	if sess == nil || sess.CurrentUser == nil {
		httplib.SendError(ctx, "Not logged in")
		return
	}
	userUid := sess.CurrentUser.UserUid
	problemUid := util.RandomString(15)
	problemInfoPayload, err := json.Marshal(req.Info)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	_, err = app.dbProblem.ExecContext(ctx, "INSERT INTO problems (uid, owner_user_uid, info) VALUES (?, ?, ?)", problemUid, userUid, problemInfoPayload)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	httplib.SendJSON(ctx, map[string]interface{}{
		"problem": map[string]interface{}{
			"uid": problemUid,
		},
	})
}

func (app *App) postProblemUploadData(ctx *fasthttp.RequestCtx) {
	problemUid := ctx.UserValue("uid").(string)
	sess, err := app.getSession(ctx)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	if sess == nil || sess.CurrentUser == nil {
		httplib.SendError(ctx, "Not logged in")
		return
	}
	userUid := sess.CurrentUser.UserUid

	err = app.dbProblem.QueryRowContext(ctx, "SELECT uid FROM problems WHERE uid=? AND owner_user_uid=?", problemUid, userUid).Scan(&problemUid)
	if err != nil {
		if err == sql.ErrNoRows {
			httplib.NotFound(ctx, fmt.Errorf("Not found"))
		} else {
			httplib.SendInternalError(ctx, err)
		}
		return
	}

	// Ignore errors
	problemDataUid := util.RandomHex(16)
	app.dbProblem.ExecContext(ctx, "UPDATE problems SET problem_data_uid=? WHERE uid=? AND problem_data_uid IS NULL", problemDataUid, problemUid)
	err = app.dbProblem.QueryRowContext(ctx, "SELECT problem_data_uid FROM problems WHERE uid=?", problemUid).Scan(&problemDataUid)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}

	httplib.SendJSON(ctx, map[string]interface{}{
		"problem": map[string]interface{}{
			"problem_data_uid": problemDataUid,
		},
	})
}

func (app *App) getProblemInfo(ctx *fasthttp.RequestCtx) {
	problemUid := ctx.UserValue("uid").(string)
	val, err := app.waitForCache(ctx, fmt.Sprintf("problem:%s:info", problemUid), time.Second*5, func() {
		app.automationCli.Trigger(map[string]interface{}{
			"tags": []string{"cache/problem/*/info/request"},
			"problem": map[string]interface{}{
				"uid": problemUid,
			},
		})
	})
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	httplib.SendJSON(ctx, map[string]interface{}{
		"data": json.RawMessage(val),
	})
}

type GetProblemResponse struct {
	Problems []*GetProblemResponseProblem `json:"problems"`
}

type GetProblemResponseProblem struct {
	Id    string          `json:"id"`
	Score float64         `json:"score"`
	Info  json.RawMessage `json:"info"`
}

func (app *App) getProblems(ctx *fasthttp.RequestCtx) {
	query := string(ctx.QueryArgs().Peek("query"))
	if query == "" {
		query = "*"
	}
	resp, err := app.esProblem.Search(app.esProblem.Search.WithIndex("problem"), app.esProblem.Search.WithQuery(query))
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	} else if resp.IsError() {
		log.Error(resp)
		httplib.SendError(ctx, "Query failed")
		return
	}
	defer resp.Body.Close()
	var data struct {
		Hits struct {
			Hits []struct {
				Id     string          `json:"_id"`
				Score  float64         `json:"score"`
				Source json.RawMessage `json:"_source"`
			}
		} `json:"hits"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	var wg sync.WaitGroup
	resp2 := &GetProblemResponse{}
	resp2.Problems = make([]*GetProblemResponseProblem, 0)
	for _, hit := range data.Hits.Hits {
		entry := &GetProblemResponseProblem{
			Id:    hit.Id,
			Score: hit.Score,
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			val, err := app.waitForCache(ctx, fmt.Sprintf("problem:%s:summary", hit.Id), time.Second*5, func() {
				app.automationCli.Trigger(map[string]interface{}{
					"tags": []string{"cache/problem/*/summary/request"},
					"problem": map[string]interface{}{
						"uid": hit.Id,
					},
				})
			})
			if err == nil {
				entry.Info = val
			} else {
				log.WithError(err).Debug("Failed to get problem summary")
			}
		}()
		resp2.Problems = append(resp2.Problems, entry)
	}
	wg.Wait()
	httplib.SendJSON(ctx, resp2)
}

type SubmitProblemRequest struct {
	SubmissionInfo *SubmitProblemRequestSubmission `json:"submission_info"`
}
type SubmitProblemRequestSubmission struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

func (app *App) postProblemSubmit(ctx *fasthttp.RequestCtx) {
	problemUid := ctx.UserValue("uid").(string)
	var req SubmitProblemRequest
	err := httplib.ReadBodyJSON(ctx, &req)
	if err != nil {
		httplib.BadRequest(ctx, err)
		return
	}
	sess, err := app.getSession(ctx)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	if sess == nil || sess.CurrentUser == nil {
		httplib.SendError(ctx, "Not logged in")
		return
	}
	err = app.dbProblem.QueryRowContext(ctx, "SELECT uid FROM problems WHERE uid=?", problemUid).Scan(&problemUid)
	if err != nil {
		httplib.NotFound(ctx, fmt.Errorf("Problem not found"))
		return
	}
	if req.SubmissionInfo == nil {
		httplib.BadRequest(ctx, fmt.Errorf("Missing submission field"))
		return
	}
	infoBytes, err := json.Marshal(req.SubmissionInfo)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	submissionUid := util.RandomString(15)
	_, err = app.dbProblem.ExecContext(ctx, "INSERT INTO submissions (uid, problem_uid, info) VALUES (?, ?, ?)", submissionUid, problemUid, infoBytes)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	httplib.SendJSON(ctx, map[string]interface{}{
		"submission": map[string]interface{}{
			"uid": submissionUid,
		},
	})
	// TODO: add it to judge queue
}
