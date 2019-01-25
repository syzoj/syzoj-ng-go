package contest

import (
    "errors"
    "context"
    "time"

    "github.com/mongodb/mongo-go-driver/mongo"
    "github.com/mongodb/mongo-go-driver/bson"
    "github.com/mongodb/mongo-go-driver/bson/primitive"

    "github.com/syzoj/syzoj-ng-go/app/model"
)

type ContestOptions struct {
    Rules ContestRules
    StartTime time.Time
    Duration time.Duration
}

type ContestRules struct {
    JudgeInContest bool
    SeeResult bool
    RejudgeAfterContest bool
    RanklistType string
    RanklistVisibility string
}

var ErrInvalidOptions = errors.New("Invalid contest options")

var boolFalse = false
var boolTrue = true

func (c *ContestService) CreateContest(ctx context.Context, id primitive.ObjectID, options *ContestOptions) (err error) {
    var contestModel model.Contest
    contestModel.Running = &boolFalse
    contestModel.Schedule = &[]model.ContestSchedule{}
    state := ""
    contestModel.State = &state
    switch options.Rules.RanklistType {
    case "":
    case "realtime":
    case "defer":
        _ = ""
    default:
        return ErrInvalidOptions
    }

    var result *mongo.UpdateResult
    if result, err = c.mongodb.Collection("problemset").UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{{"$set", bson.D{{"contest", contestModel}}}}); err != nil {
        return
    }
    if result.MatchedCount == 0 {
        panic("Problemset not found")
    }
    return
}
