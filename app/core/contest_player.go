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
	id            int
	subscriptions []*contestPlayerSubscription
}
type contestPlayerSubscription struct {
	c           *Contest
	submission  *Submission
	userId      primitive.ObjectID
	penaltyTime time.Duration
	done        bool
	score       float64
	loaded      bool
}

func (s *contestPlayerSubscription) Notify() {
	go func() {
		s.submission.Lock.RLock()
		done, score := s.submission.Done, s.submission.Score
		s.submission.Lock.RUnlock()

		s.c.lock.Lock()
		defer s.c.lock.Unlock()
		if !s.c.loaded {
			return
		}
		log.WithField("contestId", s.c.id).WithField("done", done).WithField("score", score).Info("Received submission score")
		s.done = done
		s.score = score
		p := s.c.players[s.userId]
		p.Broker.Broadcast() // notify player submission status change
		s.c.updatePlayerRankInfo(p)
	}()
}
func (c *Contest) updatePlayerRankInfo(player *ContestPlayer) {
	rankInfo := new(ContestPlayerRankInfo)
	rankInfo.problems = make(map[string]*ContestPlayerRankInfoProblem)
	for key, problem := range player.problems {
		problemInfo := new(ContestPlayerRankInfoProblem)
		for _, subscription := range problem.subscriptions {
			submissionInfo := &ContestPlayerRankInfoSubmission{
				Done:        subscription.done,
				Score:       subscription.score,
				PenaltyTime: subscription.penaltyTime,
			}
			problemInfo.submissions = append(problemInfo.submissions, submissionInfo)
		}
		rankInfo.problems[key] = problemInfo
	}
	c.ranklist.UpdatePlayer(player.userId, rankInfo)
}

func (c *Contest) loadPlayer(contestPlayerModel *model.ContestPlayer) {
	player := new(ContestPlayer)
	player.Broker = util.NewBroker()
	player.modelId = contestPlayerModel.Id
	player.userId = contestPlayerModel.User
	player.problems = make(map[string]*ContestPlayerProblem)
	for i, problemEntryModel := range contestPlayerModel.Problems {
		problemEntry := new(ContestPlayerProblem)
		for _, submission := range problemEntryModel.Submissions {
			csubmission := c.c.GetSubmission(submission.SubmissionId)
			subscription := &contestPlayerSubscription{
				c:           c,
				userId:      player.userId,
				submission:  csubmission,
				penaltyTime: submission.PenaltyTime,
			}
			csubmission.Broker.Subscribe(subscription)
			subscription.Notify()
			problemEntry.subscriptions = append(problemEntry.subscriptions, subscription)
		}
		player.problems[i] = problemEntry
	}
	c.players[player.userId] = player
}

func (*Contest) unloadPlayer(p *ContestPlayer) {
	p.Broker.Close()
	for _, problemEntry := range p.problems {
		for _, subscription := range problemEntry.subscriptions {
			subscription.submission.Broker.Unsubscribe(subscription)
		}
	}
}
