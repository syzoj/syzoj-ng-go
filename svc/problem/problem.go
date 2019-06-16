package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/syzoj/syzoj-ng-go/svc/problem/model"
	"github.com/syzoj/syzoj-ng-go/util"
	"github.com/valyala/fasthttp"
)

func (m *Main) getProblemId(ctx *fasthttp.RequestCtx) {
	problemId := ctx.UserValue("id").(string)
	prob := &model.Problem{}
	tx, err := m.dbProb.BeginTx(ctx, nil)
	if err != nil {
		m.handleError(ctx, err)
		return
	}
	defer tx.Rollback()
	err = tx.QueryRowContext(ctx, "SELECT title, statement, name FROM problems WHERE id=?", problemId).Scan(&prob.Title, &prob.Statement, &prob.Name)
	if err == sql.ErrNoRows {
		m.handleNotFound(ctx, fmt.Errorf("problem %s not found", problemId))
		return
	}
	if err != nil {
		m.handleError(ctx, err)
		return
	}
	resTags, err := tx.QueryContext(ctx, "SELECT tag FROM problem_tags WHERE problem_id=?", problemId)
	if err != nil {
		m.handleError(ctx, err)
		return
	}
	for resTags.Next() {
		var tag string
		if err := resTags.Scan(&tag); err != nil {
			m.handleError(ctx, err)
			return
		}
		prob.Tags = append(prob.Tags, tag)
	}
	if resTags.Err() != nil {
		m.handleError(ctx, err)
		return
	}
	m.sendBody(ctx, prob)
}

func (m *Main) postProblem(ctx *fasthttp.RequestCtx) {
	var req model.Problem
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		m.handleBadRequest(ctx, err)
		return
	}
	tx, err := m.dbProb.BeginTx(ctx, nil)
	if err != nil {
		m.handleError(ctx, err)
		return
	}
	defer tx.Rollback()
	problemId := util.NewId()
	if _, err := tx.ExecContext(ctx, "INSERT INTO problems (id, title, statement, name) VALUES (?, ?, ?)", problemId, req.Title, req.Statement, req.Name); err != nil {
		m.handleError(ctx, err)
		return
	}
	for _, tag := range req.Tags {
		if _, err := tx.ExecContext(ctx, "INSERT INTO problem_tags (problem_id, tag) VALUES (?, ?)", problemId, tag); err != nil {
			m.handleError(ctx, err)
			return
		}
	}
	ev := &ProblemInsertEvent{Id: problemId, Title: req.Title, Statement: req.Statement, Tags: req.Tags}
	if _, err := tx.ExecContext(ctx, "INSERT INTO problem_events (data) VALUES (?)", encodeEvent(ev)); err != nil {
		m.handleError(ctx, err)
		return
	}
	if err := tx.Commit(); err != nil {
		m.handleError(ctx, err)
		return
	}
	resp := &model.Problem{}
	resp.Id = problemId
	data, err := json.Marshal(resp)
	if err != nil {
		m.handleError(ctx, err)
		return
	}
	ctx.SetBody(data)
}

func (m *Main) putProblemId(ctx *fasthttp.RequestCtx) {
	var req model.Problem
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		m.handleBadRequest(ctx, err)
		return
	}
	tx, err := m.dbProb.BeginTx(ctx, nil)
	if err != nil {
		m.handleError(ctx, err)
		return
	}
	defer tx.Rollback()
	problemId := ctx.UserValue("id").(string)
	res, err := tx.ExecContext(ctx, "UPDATE problems SET title=?, statement=? WHERE id=?", req.Title, req.Statement, problemId)
	if err != nil {
		m.handleError(ctx, err)
		return
	}
	n, err := res.RowsAffected()
	if err != nil {
		m.handleError(ctx, err)
		return
	}
	if n == 0 {
		m.handleNotFound(ctx, fmt.Errorf("problem %s not found", problemId))
		return
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM problem_tags WHERE problem_id = ?", problemId); err != nil {
		m.handleError(ctx, err)
		return
	}
	for _, tag := range req.Tags {
		if _, err := tx.ExecContext(ctx, "INSERT INTO problem_tags (problem_id, tag) VALUES (?, ?)", problemId, tag); err != nil {
			m.handleError(ctx, err)
			return
		}
	}
	ev := &ProblemUpdateEvent{Id: problemId, Title: req.Title, Statement: req.Statement, Tags: req.Tags}
	if _, err := tx.ExecContext(ctx, "INSERT INTO problem_events (data) VALUES (?)", encodeEvent(ev)); err != nil {
		m.handleError(ctx, err)
		return
	}
	if err := tx.Commit(); err != nil {
		m.handleError(ctx, err)
		return
	}
}

func (m *Main) deleteProblemId(ctx *fasthttp.RequestCtx) {
	tx, err := m.dbProb.BeginTx(ctx, nil)
	if err != nil {
		m.handleError(ctx, err)
		return
	}
	defer tx.Rollback()
	problemId := ctx.UserValue("id").(string)
	_, err = tx.ExecContext(ctx, "DELETE FROM problem_tags WHERE id=?", problemId)
	if err != nil {
		m.handleError(ctx, err)
		return
	}
	res, err := tx.ExecContext(ctx, "DELETE FROM problems WHERE id=?", problemId)
	if err != nil {
		m.handleError(ctx, err)
		return
	}
	n, err := res.RowsAffected()
	if err != nil {
		m.handleError(ctx, err)
		return
	}
	if n == 0 {
		m.handleNotFound(ctx, fmt.Errorf("problem %s not found", problemId))
		return
	}
	ev := &ProblemDeleteEvent{Id: problemId}
	if _, err := tx.ExecContext(ctx, "INSERT INTO problem_events (data) VALUES (?)", encodeEvent(ev)); err != nil {
		m.handleError(ctx, err)
		return
	}
	if err := tx.Commit(); err != nil {
		m.handleError(ctx, err)
		return
	}
}
