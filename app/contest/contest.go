package contest

import (
    "sync"
    "sync/atomic"
    "time"
    "fmt"
    "math/rand"
    "encoding/hex"

    "github.com/mongodb/mongo-go-driver/mongo"
    "github.com/mongodb/mongo-go-driver/bson"
    "github.com/mongodb/mongo-go-driver/bson/primitive"

    "github.com/syzoj/syzoj-ng-go/app/model"
)

type Contest struct {
    srv *ContestService
    id primitive.ObjectID

    closed int32
    closeChan chan struct{}
    lock sync.Mutex

    state string
    schedules []*model.ContestSchedule
    running bool
    updates chan mongo.WriteModel
    ops chan contestOperation
}

type contestOperation interface {
    execute()
}

func (c *Contest) loadModelWithLock(contestModel *model.Contest) {
    c.running = *contestModel.Running
    for _, scheduleModel := range *contestModel.Schedule {
        schedule := new(model.ContestSchedule)
        *schedule = scheduleModel
        c.schedules = append(c.schedules, schedule)
    }
    c.state = *contestModel.State
}

func (c *Contest) resetState() {
    var stateBytes [16]byte
    if _, err := rand.Read(stateBytes[:]); err != nil {
        panic(err)
    }
    c.state = hex.EncodeToString(stateBytes[:])
}

func (c *Contest) work() {
    go c.run()
}

func (c *Contest) stop() {
    if atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
        close(c.closeChan)
    }
}

func (c *Contest) run() {
    c.lock.Lock()
    defer c.lock.Unlock()
    if atomic.LoadInt32(&c.srv.closed) == 1 {
        return
    }
    c.srv.wg.Add(1)
    defer c.srv.wg.Done()
    loop:
    for {
        timer := time.NewTimer(time.Second)
        select {
        case <-c.closeChan:
            timer.Stop()
            break loop
        case op := <-c.ops:
            op.execute()
        case <-timer.C: {
            for key, schedule := range c.schedules {
                if *schedule.Done {
                    continue
                }
                if schedule.StartTime.Before(time.Now()) {
                    op := schedule.Operation
                    method := op.Lookup("method").String()
                    switch method {
                    case "start": {
                        c.running = true
                        schedule.Done = &[]bool{true}[0]
                        updateModel := mongo.NewUpdateOneModel()
                        updateModel.Filter(bson.D{{"_id", c.id}, {"contest.state", c.state}})
                        c.resetState()
                        updateModel.Update(bson.D{{fmt.Sprintf("contest.schedule.%d.done", key), true}, {"contest.running", false}, {"contest.state", c.state}})
                        c.updates <- updateModel
                    }
                    case "stop": {
                        c.running = false
                        schedule.Done = &[]bool{true}[0]
                        updateModel := mongo.NewUpdateOneModel()
                        updateModel.Filter(bson.D{{"_id", c.id}, {"contest.state", c.state}})
                        c.resetState()
                        updateModel.Update(bson.D{{fmt.Sprintf("contest.schedule.%d.done", key), true}, {"contest.running", false}, {"contest.state", c.state}})
                        c.updates <- updateModel
                    }
                    default: {
                        log.WithField("contest", c.id).Warning("Unknown operation")
                    }}
                }
            }
        }}
    }
    log.WithField("contest", c.id).Info("Unloading contest")
}

func (c *Contest) loadModel(contestModel *model.Contest) {
    if atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
        close(c.closeChan)
    }
    c.lock.Lock()
    defer c.lock.Unlock()
    c.closed = 0
    c.closeChan = make(chan struct{})
    c.updates = make(chan mongo.WriteModel, 1000)
    c.loadModelWithLock(contestModel)
}
