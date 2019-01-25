package contest

import (
    "sync"
    "context"

    "github.com/mongodb/mongo-go-driver/mongo"
    mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
    "github.com/mongodb/mongo-go-driver/bson"
    "github.com/mongodb/mongo-go-driver/bson/primitive"
    "github.com/sirupsen/logrus"

    "github.com/syzoj/syzoj-ng-go/app/model"
)
var log = logrus.StandardLogger()

type ContestService struct {
    lock sync.Mutex
    mongodb *mongo.Database
    contests map[primitive.ObjectID]*Contest
    closed int32
    wg sync.WaitGroup
}

func NewContestService(mongo *mongo.Client) *ContestService {
    return &ContestService{
        mongodb: mongo.Database("syzoj"),
    }
}

func (c *ContestService) Init(ctx context.Context) (err error) {
    log.Info("Initializing contest service")
    c.contests = make(map[primitive.ObjectID]*Contest)
    c.lock.Lock()
    defer c.lock.Unlock()
    var cursor mongo.Cursor
    if cursor, err = c.mongodb.Collection("problemset").Find(ctx, bson.D{{"contest", bson.D{{"$exists", true}}}}, mongo_options.Find().SetProjection(bson.D{{"contest", 1}})); err != nil {
        return
    }
    for cursor.Next(ctx) {
        var contestModel model.Problemset
        if err = cursor.Decode(&contestModel); err != nil {
            return
        }
        c.loadContest(contestModel.Id, contestModel.Contest)
    }
    if err = cursor.Err(); err != nil {
        return
    }
    return
}

func (c *ContestService) ReloadContest(ctx context.Context, id primitive.ObjectID) (err error) {
    c.lock.Lock()
    defer c.lock.Unlock()
    var contestModel model.Problemset
    if err = c.mongodb.Collection("problemset").FindOne(ctx, bson.D{{"contest", bson.D{{"$exists", true}}}, {"_id", id}}).Decode(&contestModel); err != nil {
        return
    }
    c.loadContest(id, contestModel.Contest)
    return
}

func (c *ContestService) Close() error {
    c.lock.Lock()
    for _, contest := range c.contests {
        contest.stop()
    }
    c.lock.Unlock()
    c.wg.Wait()
    return nil
}

func (c *ContestService) loadContest(id primitive.ObjectID, contestModel *model.Contest) {
    log.WithField("ContestID", id).Info("Loading contest\n")
    contest := &Contest{srv: c, id: id, closeChan: make(chan struct{})}
    contest.loadModel(contestModel)
    c.contests[id] = contest
    contest.work()
}
