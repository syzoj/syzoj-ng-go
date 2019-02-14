package core

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/syzoj/syzoj-ng-go/app/model"
	"github.com/syzoj/syzoj-ng-go/util"
)

type ContestPlayer struct {
	modelId  primitive.ObjectID
	userId   primitive.ObjectID
	problems map[string]*ContestPlayerProblem
	Broker   *util.Broker
}
type ContestPlayerProblem struct {
	id          int
	submissions []*ContestPlayerSubmission
}
type ContestPlayerSubmission struct {
	c           *Contest
	submission  *Submission
	userId      primitive.ObjectID
	penaltyTime time.Duration
	submitTime  time.Time
	done        bool
	score       float64
	loaded      bool
}

func (s *ContestPlayerSubmission) Notify() {
	go func() {
		s.submission.Lock.RLock()
		done, score := s.submission.Done, s.submission.Score
		s.submission.Lock.RUnlock()

		s.c.lock.Lock()
		defer s.c.lock.Unlock()
		if !s.c.loaded {
			return
		}
		s.c.log.WithField("done", done).WithField("score", score).Debug("Received submission score")
		s.done = done
		s.score = score
		p := s.c.players[s.userId]
		p.Broker.Broadcast() // notify player submission status change
		s.c.updatePlayerRankInfo(p)
	}()
}
func (s *ContestPlayerSubmission) GetRankInfo() *ContestPlayerRankInfoSubmission {
	return &ContestPlayerRankInfoSubmission{
		Done:        s.done,
		Score:       s.score,
		PenaltyTime: s.penaltyTime,
	}
}
func (c *Contest) updatePlayerRankInfo(player *ContestPlayer) {
	rankInfo := new(ContestPlayerRankInfo)
	rankInfo.Problems = make(map[string]*ContestPlayerRankInfoProblem)
	for key, problem := range player.problems {
		problemInfo := new(ContestPlayerRankInfoProblem)
		for _, submission := range problem.submissions {
			submissionInfo := submission.GetRankInfo()
			problemInfo.Submissions = append(problemInfo.Submissions, submissionInfo)
		}
		rankInfo.Problems[key] = problemInfo
	}
	c.UpdatePlayer(player.userId, rankInfo)
}

// The logic is duplicated in RegisterPlayer
func (c *Contest) loadPlayer(contestPlayerModel *model.ContestPlayer) {
	player := new(ContestPlayer)
	player.Broker = util.NewBroker()
	player.modelId = contestPlayerModel.Id
	player.userId = contestPlayerModel.User
	player.problems = make(map[string]*ContestPlayerProblem)
	for i, problemEntryModel := range contestPlayerModel.Problems {
		problemEntry := new(ContestPlayerProblem)
		for _, submissionModel := range problemEntryModel.Submissions {
			csubmission := c.c.GetSubmission(submissionModel.SubmissionId)
			submission := &ContestPlayerSubmission{
				c:           c,
				userId:      player.userId,
				submission:  csubmission,
				penaltyTime: submissionModel.PenaltyTime,
				submitTime:  submissionModel.SubmitTime,
			}
			csubmission.Broker.Subscribe(submission)
			submission.Notify()
			problemEntry.submissions = append(problemEntry.submissions, submission)
		}
		player.problems[i] = problemEntry
	}
	c.players[player.userId] = player
	c.updatePlayerRankInfo(player)
}

func (*Contest) unloadPlayer(p *ContestPlayer) {
	p.Broker.Close()
	for _, problemEntry := range p.problems {
		for _, submission := range problemEntry.submissions {
			submission.submission.Broker.Unsubscribe(submission)
		}
	}
}
