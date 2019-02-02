package core

import (
	"sync"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type ContestPlayerRankInfo struct {
	problems map[string]*ContestPlayerRankInfoProblem
}
type ContestPlayerRankInfoProblem struct {
	submissions []*ContestPlayerRankInfoSubmission
}
type ContestPlayerRankInfoSubmission struct {
	Done bool
	Score float64
}

type ContestRankComp interface {
	Less(*Contest, *ContestPlayerRankInfo, *ContestPlayerRankInfo) bool
}
type ContestDummyRankComp struct{}
func (ContestDummyRankComp) Less(c *Contest, p1 *ContestPlayerRankInfo, p2 *ContestPlayerRankInfo) bool {
	return false
}
type ContestRankCompMaxScoreSum struct{}
func (ContestRankCompMaxScoreSum) Less(c *Contest, p1 *ContestPlayerRankInfo, p2 *ContestPlayerRankInfo) bool {
	score1 := playerMaxScoreSum(p1)
	score2 := playerMaxScoreSum(p2)
	return score1 < score2
}
func playerMaxScoreSum(p *ContestPlayerRankInfo) float64 {
	var sum float64
	for _, problem := range p.problems {
		var maxScore float64
		for _, s := range problem.submissions {
			if s.Done && s.Score > maxScore {
				maxScore = s.Score
			}
		}
		sum += maxScore
	}
	return sum
}

type ContestRanklist interface {
	Load()
	UpdatePlayer(primitive.ObjectID, *ContestPlayerRankInfo)
	Unload()
}

type ContestDummyRanklist struct{}
func (ContestDummyRanklist) Load() {}
func (ContestDummyRanklist) UpdatePlayer(primitive.ObjectID, *ContestPlayerRankInfo) {}
func (ContestDummyRanklist) Unload() {}

type ContestRealTimeRanklist struct {
	c *Contest
	lock sync.Mutex
	events []contestRealTimeRanklistEvent
}
type contestRealTimeRanklistEvent struct {
	user primitive.ObjectID
	info *ContestPlayerRankInfo
}
func (r *ContestRealTimeRanklist) Load() {
}
func (r *ContestRealTimeRanklist) UpdatePlayer(userId primitive.ObjectID, info *ContestPlayerRankInfo) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.events = append(r.events, contestRealTimeRanklistEvent{
		user: userId,
		info: info,
	})
}
func (r *ContestRealTimeRanklist) Unload() {
}
