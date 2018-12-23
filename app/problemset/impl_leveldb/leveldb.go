package impl_leveldb

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	
	"github.com/syzoj/syzoj-ng-go/app/problemset"
)

type dbGetter interface {
	Get([]byte, *opt.ReadOptions) ([]byte, error)
}
type dbPutter interface {
	Put([]byte, []byte, *opt.WriteOptions) error
}
type dbDeleter interface {
	Delete([]byte, *opt.WriteOptions) error
}

func (*service) getProblemsetInfo(db dbGetter, id uuid.UUID) (info *problemsetInfo, err error) {
	keyProblemset := []byte(fmt.Sprintf("problemset.regular:%s", id))
	var data []byte
	if data, err = db.Get(keyProblemset, nil); err != nil {
		if err == leveldb.ErrNotFound {
			err = problemset.ErrProblemsetNotFound
		}
		return
	}
	info = new(problemsetInfo)
	if err = json.Unmarshal(data, info); err != nil {
		return
	}
	return
}

func (*service) putProblemsetInfo(db dbPutter, id uuid.UUID, info *problemsetInfo) (err error) {
	keyProblemset := []byte(fmt.Sprintf("problemset.regular:%s", id))
	var data []byte
	if data, err = json.Marshal(info); err != nil {
		return
	}
	if err = db.Put(keyProblemset, data, nil); err != nil {
		return
	}
	return
}

func (*service) getProblemInfo(db dbGetter, id uuid.UUID, name string) (info *problemInfo, err error) {
	keyProblem := []byte(fmt.Sprintf("problemset.regular:%s.problem:%s", id, name))
	var data []byte
	if data, err = db.Get(keyProblem, nil); err != nil {
		return
	}
	info = new(problemInfo)
	if err = json.Unmarshal(data, info); err != nil {
		return
	}
	return
}

func (*service) putProblemInfo(db dbPutter, id uuid.UUID, name string, info *problemInfo) (err error) {
	keyProblem := []byte(fmt.Sprintf("problemset.regular:%s.problem:%s", id, name))
	var data []byte
	if data, err = json.Marshal(info); err != nil {
		return
	}
	if err = db.Put(keyProblem, data, nil); err != nil {
		return
	}
	return
}

func (*service) getRole(db dbGetter, id uuid.UUID, userId uuid.UUID) (role roleInfo, err error) {
	keyRole := []byte(fmt.Sprintf("problemset.regular:%s.role:%s", id, userId))
	var data []byte
	if data, err = db.Get(keyRole, nil); err != nil {
		if err == leveldb.ErrNotFound {
			// Default role is empty
			return "", nil
		}
		return
	}
	role = roleInfo(data)
	return
}

func (*service) putRole(db dbPutter, id uuid.UUID, userId uuid.UUID, role roleInfo) (err error) {
	keyRole := []byte(fmt.Sprintf("problemset.regular:%s.role:%s", id, userId))
	var data []byte = []byte(role)
	if err = db.Put(keyRole, data, nil); err != nil {
		return
	}
	return
}

func (*service) getSubmissionInfo(db dbGetter, id uuid.UUID, submissionId uuid.UUID) (info *submissionInfo, err error) {
	keySubmission := []byte(fmt.Sprintf("problemset.regular:%s.submission:%s", id, submissionId))
	var data []byte
	if data, err = db.Get(keySubmission, nil); err != nil {
		if err == leveldb.ErrNotFound {
			err = problemset.ErrSubmissionNotFound
			return
		}
	}
	info = new(submissionInfo)
	if err = json.Unmarshal(data, info); err != nil {
		return
	}
	return
}

func (*service) putSubmissionInfo(db dbPutter, id uuid.UUID, submissionId uuid.UUID, info *submissionInfo) (err error) {
	keySubmission := []byte(fmt.Sprintf("problemset.regular:%s.submission:%s", id, submissionId))
	var data []byte
	if data, err = json.Marshal(info); err != nil {
		return
	}
	if err = db.Put(keySubmission, data, nil); err != nil {
		return
	}
	return
}
