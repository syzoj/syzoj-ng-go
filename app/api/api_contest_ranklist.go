package api

import (
	"github.com/valyala/fastjson"

	"github.com/syzoj/syzoj-ng-go/app/core"
)

func Handle_Contest_Ranklist(c *ApiContext) ApiError {
	var err error
	vars := c.Vars()
	contestId := DecodeObjectID(vars["contest_id"])
	if err = c.SessionStart(); err != nil {
		panic(err)
	}

	contest := c.Server().c.GetContestR(contestId)
	if contest == nil {
		return ErrContestNotLoaded
	}
	var ranklistVisible bool
	switch contest.GetRanklistVisibility() {
	case "all":
		ranklistVisible = true
	case "player":
		if c.Session.LoggedIn() {
			player := contest.GetPlayer(c.Session.AuthUserUid)
			if player != nil {
				ranklistVisible = true
			}
		}
	}
	rankcomp := contest.GetRankComp()
	if !ranklistVisible {
		contest.RUnlock()
		return ErrPermissionDenied
	}
	snapshot := contest.GetRanklist()
	contest.RUnlock()
	arena := new(fastjson.Arena)
	result := arena.NewObject()
	var rankcompn string
	switch rankcomp.(type) {
	case core.ContestRankCompMaxScoreSum:
		rankcompn = "maxsum"
	case core.ContestRankCompLastSum:
		rankcompn = "lastsum"
	case core.ContestRankCompACM:
		rankcompn = "acm"
	default:
		rankcompn = "unknown"
	}
	result.Set("rankcomp", arena.NewString(rankcompn))
	ranklist := arena.NewArray()
	for i, entry := range snapshot {
		entryVal := arena.NewObject()
		entryVal.Set("user_id", arena.NewString(EncodeObjectID(entry.UserId)))
		problemsVal := arena.NewArray()
		var j int
		for problemName, problemEntry := range entry.Info.Problems {
			problemVal := arena.NewObject()
			problemVal.Set("name", arena.NewString(problemName))
			problemVal.Set("score", arena.NewNumberFloat64(rankcomp.GetProblemScore(problemEntry)))
			problemsVal.SetArrayItem(j, problemVal)
			j++
		}
		entryVal.Set("problems", problemsVal)
		ranklist.SetArrayItem(i, entryVal)
	}
	result.Set("ranklist", ranklist)
	c.SendValue(result)
	return nil
}
