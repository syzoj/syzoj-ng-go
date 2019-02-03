package core

import (
	"sort"
	"sync"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type ContestPlayerRankInfo struct {
	problems map[string]*ContestPlayerRankInfoProblem
}
type ContestPlayerRankInfoProblem struct {
	submissions []*ContestPlayerRankInfoSubmission
}
type ContestPlayerRankInfoSubmission struct {
	Done  bool
	Score float64
	PenaltyTime time.Duration
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

type ContestRankCompACM struct{}

func (ContestRankCompACM) Less(c *Contest, p1 *ContestPlayerRankInfo, p2 *ContestPlayerRankInfo) bool {
	sum1, penalty1 := playerSumAndPenalty(p1)
	sum2, penalty2 := playerSumAndPenalty(p2)
	return (sum1 < sum2) || (sum1 == sum2 && penalty1 < penalty2)
}
func playerSumAndPenalty(p *ContestPlayerRankInfo) (float64, time.Duration) {
	var sum float64
	var penalty time.Duration
	for _, problem := range p.problems {
		var maxScore float64
		var minPenalty time.Duration
		for _, s := range problem.submissions {
			if s.Done && s.Score > maxScore {
				maxScore = s.Score
				minPenalty = s.PenaltyTime
			}
			if s.Score == maxScore && s.PenaltyTime < minPenalty {
				minPenalty = s.PenaltyTime
			}
		}
		if maxScore > 0 {
			sum += maxScore
			penalty += time.Duration(float64(minPenalty) * (maxScore / 100))
		}
	}
	return sum, penalty
}

type ContestRanklist interface {
	Load()
	UpdatePlayer(primitive.ObjectID, *ContestPlayerRankInfo)
	Unload()
}

type ContestDummyRanklist struct{}

func (ContestDummyRanklist) Load()                                                   {}
func (ContestDummyRanklist) UpdatePlayer(primitive.ObjectID, *ContestPlayerRankInfo) {}
func (ContestDummyRanklist) Unload()                                                 {}

type ContestRealTimeRanklist struct {
	c      *Contest
	lock   sync.Mutex
	events []contestRealTimeRanklistEvent

	sorterSemasphore chan struct{}
	rankMutex        sync.RWMutex
	players          map[primitive.ObjectID]*ContestPlayerRankInfo
	lastEvent        int
	ranklist         []primitive.ObjectID
}
type contestRealTimeRanklistEvent struct {
	user primitive.ObjectID
	info *ContestPlayerRankInfo
}

func (r *ContestRealTimeRanklist) Load() {
	r.sorterSemasphore = make(chan struct{}, 1)
	r.players = make(map[primitive.ObjectID]*ContestPlayerRankInfo)
}
func (r *ContestRealTimeRanklist) UpdatePlayer(userId primitive.ObjectID, info *ContestPlayerRankInfo) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.events = append(r.events, contestRealTimeRanklistEvent{
		user: userId,
		info: info,
	})
	go r.sort(len(r.events), r.events)
}
func (r *ContestRealTimeRanklist) sort(l int, events []contestRealTimeRanklistEvent) {
	select {
	case r.sorterSemasphore <- struct{}{}:
	default:
		break
	}
	defer func() {
		<-r.sorterSemasphore
	}()

	func() {
		r.rankMutex.Lock()
		defer r.rankMutex.Unlock()
		for i := r.lastEvent; i < l; i++ {
			event := events[i]
			r.players[event.user] = event.info
		}
		r.lastEvent = l
		r.ranklist = sortPlayers(r.c, r.c.rankcomp, r.players)
	}()
	time.Sleep(time.Millisecond * 100)
}
func (r *ContestRealTimeRanklist) Unload() {
}

type ranklistSorter struct {
	c          *Contest
	comp       ContestRankComp
	players    []primitive.ObjectID
	playerInfo map[primitive.ObjectID]*ContestPlayerRankInfo
}

func (s ranklistSorter) Len() int {
	return len(s.players)
}
func (s ranklistSorter) Swap(i, j int) {
	p := s.players[i]
	s.players[i] = s.players[j]
	s.players[j] = p
}
func (s ranklistSorter) Less(i, j int) bool {
	return s.comp.Less(s.c, s.playerInfo[s.players[i]], s.playerInfo[s.players[j]])
}
func sortPlayers(c *Contest, comp ContestRankComp, playerInfo map[primitive.ObjectID]*ContestPlayerRankInfo) []primitive.ObjectID {
	var players []primitive.ObjectID
	for p := range playerInfo {
		players = append(players, p)
	}
	sorter := ranklistSorter{c: c, comp: comp, players: players, playerInfo: playerInfo}
	sort.Sort(sorter)
	return sorter.players
}
