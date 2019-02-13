package core

import (
	"sort"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type ContestPlayerRankInfo struct {
	Problems map[string]*ContestPlayerRankInfoProblem
}
type ContestPlayerRankInfoProblem struct {
	Submissions []*ContestPlayerRankInfoSubmission
}
type ContestPlayerRankInfoSubmission struct {
	Done        bool
	Score       float64
	PenaltyTime time.Duration
}
type ContestRanklistEntry struct {
	UserId primitive.ObjectID
	Info   *ContestPlayerRankInfo
}
type ContestRanklistEvent struct {
	userId primitive.ObjectID
	info   *ContestPlayerRankInfo
}

type ContestRankComp interface {
	Less(*Contest, *ContestPlayerRankInfo, *ContestPlayerRankInfo) bool
	GetProblemScore(*ContestPlayerRankInfoProblem) float64
}
type ContestDummyRankComp struct{}

func (ContestDummyRankComp) Less(c *Contest, p1 *ContestPlayerRankInfo, p2 *ContestPlayerRankInfo) bool {
	return false
}
func (ContestDummyRankComp) GetProblemScore(*ContestPlayerRankInfoProblem) float64 {
	return 0
}

type ContestRankCompMaxScoreSum struct{}

func (s ContestRankCompMaxScoreSum) Less(c *Contest, p1 *ContestPlayerRankInfo, p2 *ContestPlayerRankInfo) bool {
	score1 := s.playerMaxScoreSum(p1)
	score2 := s.playerMaxScoreSum(p2)
	return score1 < score2
}
func (s ContestRankCompMaxScoreSum) GetProblemScore(problem *ContestPlayerRankInfoProblem) float64 {
	var maxScore float64
	for _, s := range problem.Submissions {
		if s.Done && s.Score > maxScore {
			maxScore = s.Score
		}
	}
	return maxScore
}
func (s ContestRankCompMaxScoreSum) playerMaxScoreSum(p *ContestPlayerRankInfo) float64 {
	var sum float64
	for _, problem := range p.Problems {
		sum += s.GetProblemScore(problem)
	}
	return sum
}

type ContestRankCompLastSum struct{}

func (s ContestRankCompLastSum) Less(c *Contest, p1 *ContestPlayerRankInfo, p2 *ContestPlayerRankInfo) bool {
	score1 := s.playerLastSum(p1)
	score2 := s.playerLastSum(p2)
	return score1 < score2
}
func (s ContestRankCompLastSum) GetProblemScore(problem *ContestPlayerRankInfoProblem) float64 {
	var score float64
	for _, s := range problem.Submissions {
		if s.Done {
			score = s.Score
		}
	}
	return score
}
func (s ContestRankCompLastSum) playerLastSum(p *ContestPlayerRankInfo) float64 {
	var sum float64
	for _, problem := range p.Problems {
		sum += s.GetProblemScore(problem)
	}
	return sum
}

type ContestRankCompACM struct{}

func (ContestRankCompACM) Less(c *Contest, p1 *ContestPlayerRankInfo, p2 *ContestPlayerRankInfo) bool {
	sum1, penalty1 := playerSumAndPenalty(p1)
	sum2, penalty2 := playerSumAndPenalty(p2)
	return (sum1 < sum2) || (sum1 == sum2 && penalty1 < penalty2)
}
func (ContestRankCompACM) GetProblemScore(*ContestPlayerRankInfoProblem) float64 {
	return 0
}
func playerSumAndPenalty(p *ContestPlayerRankInfo) (float64, time.Duration) {
	var sum float64
	var penalty time.Duration
	for _, problem := range p.Problems {
		var maxScore float64
		var minPenalty time.Duration
		for _, s := range problem.Submissions {
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

func (c *Contest) loadRanklist() {
	c.ranklist_m = make(map[primitive.ObjectID]*ContestPlayerRankInfo)
	if c.ranklist == "realtime" {
		go c.sortRanklist()
	} else {
		go c.clearRanklist()
	}
}

func (c *Contest) UpdatePlayer(userId primitive.ObjectID, info *ContestPlayerRankInfo) {
	c.ranklist_e = append(c.ranklist_e, ContestRanklistEvent{userId: userId, info: info})
}

func (c *Contest) clearRanklist() {
	c.lock.Lock()
	defer c.lock.Unlock()
	if !c.loaded {
		return
	}
	c.ranklist_e = nil
	time.AfterFunc(time.Second, c.clearRanklist)
}

func (c *Contest) sortRanklist() {
	c.lock.Lock()
	if !c.loaded {
		c.lock.Unlock()
		return
	}
	ev := c.ranklist_e
	c.ranklist_e = nil
	c.lock.Unlock()
	if len(ev) != 0 {
		c.log.WithField("count", len(ev)).Debug("Applying ranklist updates")
		for _, e := range ev {
			c.ranklist_m[e.userId] = e.info
		}
		var players []primitive.ObjectID
		for p, _ := range c.ranklist_m {
			players = append(players, p)
		}
		sorter := ranklistSorter{
			c:          c,
			comp:       c.rankcomp,
			players:    players,
			playerInfo: c.ranklist_m,
		}
		sort.Sort(sorter)
		var snapshot []ContestRanklistEntry
		for _, p := range players {
			snapshot = append(snapshot, ContestRanklistEntry{
				UserId: p,
				Info:   c.ranklist_m[p],
			})
		}
		c.log.Debug("Updating ranklist snapshot")
		log.Info(snapshot)
		c.lock.Lock()
		c.ranklist_w = snapshot
		c.lock.Unlock()
	}
	time.AfterFunc(time.Second, c.sortRanklist)
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
