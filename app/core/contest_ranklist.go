package core

import (
	"sort"

	"github.com/golang/protobuf/proto"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

type ContestRanklistImpl interface {
	Type() string
	getRanklist(ct *Contest) *model.ContestRanklist
}

var ranklists = map[string]ContestRanklistImpl{"max_sum": MaxSumRanklist{}}

func (ct *Contest) notifyUpdateRanklist() {
	select {
	case ct.ranklistSema <- struct{}{}:
		go func() {
			ct.updateRanklist()
			select {
			case <-ct.ranklistSema:
			default:
			}
		}()
	default:
	}
}

func (ct *Contest) updateRanklist() {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	ct.ranklist = ct.ranklistimpl.getRanklist(ct)
}

type DummyRanklist struct{}

func (DummyRanklist) Type() string {
	return ""
}

func (DummyRanklist) getRanklist(*Contest) *model.ContestRanklist {
	return nil
}

type MaxSumRanklist struct{}

type maxSumRanklistSorter struct {
	info map[*ContestPlayer]maxSumRanklistInfo
	list []*ContestPlayer
}

func (s maxSumRanklistSorter) Len() int {
	return len(s.list)
}

func (s maxSumRanklistSorter) Less(i, j int) bool {
	return s.info[s.list[i]].sum < s.info[s.list[j]].sum
}

func (s maxSumRanklistSorter) Swap(i, j int) {
	p := s.list[i]
	s.list[i] = s.list[j]
	s.list[j] = p
}

type maxSumRanklistInfo struct {
	sum float64
}

func (MaxSumRanklist) Type() string { return "max_sum" }

func (r MaxSumRanklist) getRanklist(ct *Contest) *model.ContestRanklist {
	sorter := maxSumRanklistSorter{}
	sorter.info = make(map[*ContestPlayer]maxSumRanklistInfo)
	for _, player := range ct.players {
		var info maxSumRanklistInfo
		for _, problem := range player.problems {
			var maxScore float64
			for _, submission := range problem.submissions {
				submissionId := submission.submissionId
				result := ct.submissions[submissionId]
				if result != nil && result.Score != nil {
					score := result.GetScore()
					if score > maxScore {
						maxScore = score
					}
				}
			}
			info.sum += maxScore
		}
		sorter.info[player] = info
	}
	sorter.list = make([]*ContestPlayer, len(ct.players))
	copy(sorter.list, ct.players)
	sort.Sort(sorter)

	ranklist := new(model.ContestRanklist)
	ranklist.RanklistType = proto.String("max_sum")
	ranklist.Entries = make([]*model.ContestRanklistEntry, len(sorter.list))
	for i, player := range sorter.list {
		entry := new(model.ContestRanklistEntry)
		entry.UserId = model.ObjectIDProto(player.userId)
		entry.ScoreSum = proto.Float64(sorter.info[player].sum)
		ranklist.Entries[i] = entry
	}
	return ranklist
}
