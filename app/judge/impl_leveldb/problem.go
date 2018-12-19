package impl_leveldb

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/syzoj/syzoj-ng-go/app/judge"
)

func (*judgeService) getProblem(db dbGetter, problemId uuid.UUID) (problem *judge.Problem, err error) {
	var data []byte
	keyProblem := []byte(fmt.Sprintf("judge.problem:%s", problemId))
	if data, err = db.Get(keyProblem, nil); err != nil {
		if err == leveldb.ErrNotFound {
			err = judge.ErrProblemNotExist
		}
		return
	}
	problem = new(judge.Problem)
	if err = json.Unmarshal(data, problem); problem != nil {
		return
	}
	return
}

func (*judgeService) putProblem(db dbPutter, problemId uuid.UUID, problem *judge.Problem) (err error) {
	var data []byte
	keyProblem := []byte(fmt.Sprintf("judge.problem:%s", problemId))
	if data, err = json.Marshal(problem); err != nil {
		return
	}
	if err = db.Put(keyProblem, data, nil); err != nil {
		return
	}
	return
}

func (*judgeService) deleteProblem(db dbDeleter, problemId uuid.UUID) (err error) {
	keyProblem := []byte(fmt.Sprintf("judge.problem:%s", problemId))
	if err = db.Delete(keyProblem, nil); err != nil {
		if err == leveldb.ErrNotFound {
			err = judge.ErrProblemNotExist
		}
		return
	}
	return
}

func (s *judgeService) CreateProblem(info *judge.Problem) (id uuid.UUID, err error) {
	if id, err = uuid.NewRandom(); err != nil {
		return
	}
	var tokenBytes [16]byte
	if _, err = rand.Read(tokenBytes[:]); err != nil {
		return
	}
	info.Token = hex.EncodeToString(tokenBytes[:])
	if err = os.MkdirAll(filepath.Join(s.dataPath, id.String()), 0755); err != nil {
		return
	}
	if err = s.putProblem(s.db, id, info); err != nil {
		return
	}
	return
}

func (s *judgeService) GetProblem(id uuid.UUID) (info *judge.Problem, err error) {
	if info, err = s.getProblem(s.db, id); err != nil {
		return
	}
	return
}

func (s *judgeService) UpdateProblem(id uuid.UUID, info *judge.Problem) (err error) {
	var org_info *judge.Problem
	s.problemLock.Lock()
	defer s.problemLock.Unlock()
	if org_info, err = s.getProblem(s.db, id); err != nil {
		return
	}
	if org_info.Version != info.Version {
		err = judge.ErrConcurrentUpdate
		return
	}
	info.Version++
	if err = s.putProblem(s.db, id, info); err != nil {
		return
	}
	return
}

func (s *judgeService) DeleteProblem(id uuid.UUID) (err error) {
	if err = s.deleteProblem(s.db, id); err != nil {
		return
	}
	return
}
