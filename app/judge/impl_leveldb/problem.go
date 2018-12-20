package impl_leveldb

import (
	"io/ioutil"
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

type problemInfo struct {
	Title     string    `json:"title"`
	Statement string    `json:"statement"`
	Token     string    `json:"token"`
	Owner     uuid.UUID `json:"owner"`
}

func (*judgeService) getProblem(db dbGetter, problemId uuid.UUID) (problem *problemInfo, err error) {
	var data []byte
	keyProblem := []byte(fmt.Sprintf("judge.problem:%s", problemId))
	if data, err = db.Get(keyProblem, nil); err != nil {
		if err == leveldb.ErrNotFound {
			err = judge.ErrProblemNotExist
		}
		return
	}
	problem = new(problemInfo)
	if err = json.Unmarshal(data, problem); problem != nil {
		return
	}
	return
}

func (*judgeService) putProblem(db dbPutter, problemId uuid.UUID, problem *problemInfo) (err error) {
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
	var _info problemInfo
	_info.Title = info.Title
	_info.Owner = info.Owner
	if id, err = uuid.NewRandom(); err != nil {
		return
	}
	var tokenBytes [16]byte
	if _, err = rand.Read(tokenBytes[:]); err != nil {
		return
	}
	_info.Token = hex.EncodeToString(tokenBytes[:])
	if err = os.MkdirAll(filepath.Join(s.dataPath, id.String()), 0755); err != nil {
		return
	}
	if err = s.putProblem(s.db, id, &_info); err != nil {
		return
	}
	return
}

func (s *judgeService) GetProblemFullInfo(id uuid.UUID, info *judge.Problem) (err error) {
	var _info *problemInfo
	if _info, err = s.getProblem(s.db, id); err != nil {
		return
	}
	info.Title = _info.Title
	info.Statement = _info.Statement
	info.Token = _info.Token
	info.Owner = _info.Owner
	return
}

func (s *judgeService) GetProblemOwnerInfo(id uuid.UUID, info *judge.Problem) (err error) {
	var _info *problemInfo
	if _info, err = s.getProblem(s.db, id); err != nil {
		return
	}
	info.Owner = _info.Owner
	return
}

func (s *judgeService) GetProblemStatementInfo(id uuid.UUID, info *judge.Problem) (err error) {
	var _info *problemInfo
	if _info, err = s.getProblem(s.db, id); err != nil {
		return
	}
	info.Statement = _info.Statement
	return
}

func (s *judgeService) UpdateProblem(id uuid.UUID, info *judge.Problem) (err error) {
	s.problemLock.Lock()
	defer s.problemLock.Unlock()
	var _info *problemInfo
	if _info, err = s.getProblem(s.db, id); err != nil {
		return
	}
	// TODO: This is a very naive implementation, change it later
	var b []byte
	if b, err = ioutil.ReadFile(filepath.Join(s.dataPath, id.String(), "statement.md")); err != nil {
		return
	}
	_info.Statement = string(b)
	if err = s.putProblem(s.db, id, _info); err != nil {
		return
	}
	return
}

func (s *judgeService) ChangeProblemTitle(id uuid.UUID, info *judge.Problem) (err error) {
	s.problemLock.Lock()
	defer s.problemLock.Unlock()
	var _info *problemInfo
	if _info, err = s.getProblem(s.db, id); err != nil {
		return
	}
	_info.Title = info.Title
	if err = s.putProblem(s.db, id, _info); err != nil {
		return
	}
	return
}

func (s *judgeService) ResetProblemToken(id uuid.UUID, info *judge.Problem) (err error) {
	s.problemLock.Lock()
	defer s.problemLock.Unlock()
	var _info *problemInfo
	if _info, err = s.getProblem(s.db, id); err != nil {
		return
	}
	var tokenBytes [16]byte
	if _, err = rand.Read(tokenBytes[:]); err != nil {
		return
	}
	_info.Token = hex.EncodeToString(tokenBytes[:])
	if err = s.putProblem(s.db, id, _info); err != nil {
		return
	}
	info.Token = _info.Token
	return
}

func (s *judgeService) DeleteProblem(id uuid.UUID) (err error) {
	if err = s.deleteProblem(s.db, id); err != nil {
		return
	}
	return
}
