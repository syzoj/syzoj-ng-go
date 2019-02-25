package core

import (
	"sync"

	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

type Contest struct {
	c           *Core
	mu          sync.Mutex
	id          primitive.ObjectID
	running     bool
	playersById map[primitive.ObjectID]*ContestPlayer
	players     []*ContestPlayer

	judgeInContest bool
	problems       []*ContestProblem
	problemsByName map[string]*ContestProblem

	submissions  map[primitive.ObjectID]*model.SubmissionResult
	ranklist     *model.ContestRanklist
	ranklistimpl ContestRanklistImpl
	ranklistSema chan struct{}
}

// immutable
type ContestProblem struct {
	name string
	data *model.ContestProblem
}

type ContestPlayer struct {
	ct       *Contest
	modelId  primitive.ObjectID
	userId   primitive.ObjectID
	problems []*ContestPlayerProblem
}

type ContestPlayerProblem struct {
	name        string
	submissions []*ContestPlayerProblemSubmission
}

type ContestPlayerProblemSubmission struct {
	submissionId primitive.ObjectID
}

func (c *Core) initContests() error {
	c.contestsLock.Lock()
	defer c.contestsLock.Unlock()
	c.contests = make(map[primitive.ObjectID]*Contest)
	cursor, err := c.mongodb.Collection("contest").Find(c.context, bson.D{})
	if err != nil {
		return err
	}
	defer cursor.Close(c.context)
	for cursor.Next(c.context) {
		contestModel := new(model.Contest)
		if err := cursor.Decode(contestModel); err != nil {
			return err
		}
		c.contests[model.MustGetObjectID(contestModel.Id)] = c.loadContest(contestModel)
	}
	if err := cursor.Err(); err != nil {
		return err
	}
	return nil
}

func (c *Core) unloadContests() {
	var wg sync.WaitGroup
	for _, contest := range c.contests {
		wg.Add(1)
		go func() {
			defer wg.Done()
			contest.save()
		}()
	}
	wg.Wait()
	return
}

// Returns a snapshot of all contests
func (c *Core) GetContests() []*Contest {
	c.contestsLock.Lock()
	defer c.contestsLock.Unlock()
	contests := make([]*Contest, len(c.contests))
	i := 0
	for _, contest := range c.contests {
		contests[i] = contest
		i++
	}
	return contests
}

func (c *Core) GetContest(id primitive.ObjectID) *Contest {
	c.contestsLock.Lock()
	defer c.contestsLock.Unlock()
	return c.contests[id]
}

func (c *Core) loadContest(contestModel *model.Contest) *Contest {
	ct := new(Contest)
	ct.mu.Lock()
	defer ct.mu.Unlock()
	ct.c = c
	ct.id = model.MustGetObjectID(contestModel.Id)
	ct.running = contestModel.State.GetRunning()
	ct.judgeInContest = contestModel.State.GetJudgeInContest()
	ct.loadPlayers()
	ct.submissions = make(map[primitive.ObjectID]*model.SubmissionResult)
	ct.problems = make([]*ContestProblem, len(contestModel.State.Problems))
	ct.problemsByName = make(map[string]*ContestProblem)
	for i, problemModel := range contestModel.State.Problems {
		ct.problems[i] = &ContestProblem{data: problemModel, name: problemModel.GetName()}
		ct.problemsByName[problemModel.GetName()] = ct.problems[i]
	}
	ct.ranklistSema = make(chan struct{}, 1)
	go c.AddSubmissionHook(contestHook{ct})
	ranklistType := contestModel.State.GetRanklistType()
	ct.ranklistimpl = ranklists[ranklistType]
	if ct.ranklistimpl == nil {
		ct.ranklistimpl = DummyRanklist{}
	}
	ct.notifyUpdateRanklist()
	return ct
}

func (ct *Contest) save() {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	state := new(model.ContestState)
	state.Running = proto.Bool(ct.running)
	state.JudgeInContest = proto.Bool(ct.judgeInContest)
	state.RanklistType = proto.String(ct.ranklistimpl.Type())
	for _, problem := range ct.problems {
		state.Problems = append(state.Problems, problem.data)
	}
	_, err := ct.c.mongodb.Collection("contest").UpdateOne(ct.c.context, bson.D{{"_id", ct.id}}, bson.D{{"$set", bson.D{{"state", state}}}})
	if err != nil {
		log.WithField("contestId", ct.id).WithError(err).Error("Failed to save contest state")
	}

	var wg sync.WaitGroup
	for _, player := range ct.players {
		wg.Add(1)
		go func() {
			defer wg.Done()
			player.save()
		}()
	}
	wg.Wait()
}

func (ct *Contest) loadPlayers() {
	ct.playersById = make(map[primitive.ObjectID]*ContestPlayer)
	cursor, err := ct.c.mongodb.Collection("contest_player").Find(ct.c.context, bson.D{{"contest", ct.id}})
	if err != nil {
		log.WithField("contestId", ct.id).WithError(err).Error("Failed to load contest players")
		return
	}
	defer cursor.Close(ct.c.context)
	for cursor.Next(ct.c.context) {
		playerModel := new(model.ContestPlayer)
		if err = cursor.Decode(playerModel); err != nil {
			log.WithField("contestId", ct.id).WithError(err).Error("Failed to load contest players")
			return
		}
		player := ct.loadPlayerModel(playerModel)
		ct.players = append(ct.players, player)
		ct.playersById[player.userId] = player
	}
	if err = cursor.Err(); err != nil {
		log.WithField("contestId", ct.id).WithError(err).Error("Failed to load contest players")
		return
	}
}

func (ct *Contest) loadPlayerModel(playerModel *model.ContestPlayer) *ContestPlayer {
	player := new(ContestPlayer)
	player.ct = ct
	player.modelId = model.MustGetObjectID(playerModel.Id)
	player.userId = model.MustGetObjectID(playerModel.User)
	return player
}

func (p *ContestPlayer) save() {
	playerModel := new(model.ContestPlayer)
	playerModel.Id = model.ObjectIDProto(p.modelId)
	playerModel.User = model.ObjectIDProto(p.userId)
	playerModel.Contest = model.ObjectIDProto(p.ct.id)
	if _, err := p.ct.c.mongodb.Collection("contest_player").ReplaceOne(p.ct.c.context, bson.D{{"_id", p.modelId}}, playerModel, mongo_options.Replace().SetUpsert(true)); err != nil {
		log.WithField("contestId", p.ct.id).WithField("playerId", p.modelId).WithError(err).Error("Failed to save player")
	}
}

func (ct *Contest) RegisterPlayer(userId primitive.ObjectID) bool {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	if _, found := ct.playersById[userId]; found {
		return false
	}
	player := new(ContestPlayer)
	player.ct = ct
	player.modelId = primitive.NewObjectID()
	player.userId = userId
	ct.players = append(ct.players, player)
	ct.playersById[userId] = player
	return true
}

func (ct *Contest) GetPlayerById(userId primitive.ObjectID) *ContestPlayer {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	return ct.playersById[userId]
}

func (ct *Contest) GetId() primitive.ObjectID {
	return ct.id
}

func (ct *Contest) GetProblems() []*ContestProblem {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	problems := make([]*ContestProblem, len(ct.problems))
	copy(problems, ct.problems)
	return problems
}

func (ct *Contest) GetProblemByName(name string) *ContestProblem {
    return ct.problemsByName[name]
}

func (ct *Contest) IsRunning() bool {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	return ct.running
}

func (p *ContestPlayer) GetUserId() primitive.ObjectID {
	return p.userId
}

func (p *ContestPlayer) GetModelId() primitive.ObjectID {
	return p.modelId
}

func (p *ContestProblem) GetData() *model.ContestProblem {
	return p.data
}

func (p *ContestProblem) GetName() string {
	return p.name
}
