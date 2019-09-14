package app

import (
	"context"
	"strings"
)

// Handles aggregated sums from stats package.
func (a *App) saveCounter(ctx context.Context, key string, val int64) {
	ind := strings.Index(key, ":")
	if ind == -1 {
		return
	}
	switch key[:ind] {
	case "user.problems":
		key = key[ind+1:]
		const SQLAddUserProblems = "UPDATE `users` SET `problem_count`=`problem_count`+? WHERE `uid`=?"
		_, err := a.Db.ExecContext(ctx, SQLAddUserProblems, val, key)
		if err != nil {
			log.WithError(err).WithField("uid", key).WithField("val", val).Warning("failed to update user problem count")
			return
		}
	default:
		log.Warning("unrecognized counter key", key)
	}
}
