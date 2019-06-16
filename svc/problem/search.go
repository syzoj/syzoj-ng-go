package main

import (
	"encoding/json"

	"github.com/syzoj/syzoj-ng-go/svc/problem/model"
	"github.com/valyala/fasthttp"
)

type searchReq struct {
	Query string `json:"query"`
}

type searchResp struct {
	Hits []*searchRespHit `json:"hit"`
}

type searchRespHit struct {
	Id      string            `json:"id"`
	Problem *model.ProblemDoc `json:"problem"`
}

type esSearchReq struct {
	Query struct {
		QueryString string `json:"query_string"`
	} `json:"query"`
}

type esSearchResp struct {
	Hits struct {
		Hits []struct {
			Id     string            `json:"_id"`
			Source *model.ProblemDoc `json:"_source"`
		}
	}
}

func (m *Main) postSearch(ctx *fasthttp.RequestCtx) {
	args := ctx.URI().QueryArgs()
	size := args.GetUintOrZero("size")
	if size == 0 {
		size = 10
	}
	if size > 50 {
		size = 50
	}
	from := args.GetUintOrZero("from")
	var req searchReq
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		m.handleBadRequest(ctx, err)
		return
	}
	resp, err := m.es.Search(m.es.Search.WithContext(ctx), m.es.Search.WithIndex("problem"), m.es.Search.WithQuery(req.Query), m.es.Search.WithFrom(from), m.es.Search.WithSize(size), m.es.Search.WithSourceExcludes("statement"))
	if err != nil {
		m.handleError(ctx, err)
		return
	}
	defer resp.Body.Close()
	log.Info(resp)
	var esResp esSearchResp
	if err := json.NewDecoder(resp.Body).Decode(&esResp); err != nil {
		m.handleError(ctx, err)
		return
	}
	res := &searchResp{}
	for _, hit := range esResp.Hits.Hits {
		res.Hits = append(res.Hits, &searchRespHit{
			Id:      hit.Id,
			Problem: hit.Source,
		})
	}
	m.sendBody(ctx, res)
}
