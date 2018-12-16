package impl_leveldb

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/syzoj/syzoj-ng-go/app/git"
)

type gitRepoInfo struct {
	HookType string `json:"hook_type"`
	Token    string `json:"token"`
}

type dbGetter interface {
	Get([]byte, *opt.ReadOptions) ([]byte, error)
}
type dbPutter interface {
	Put([]byte, []byte, *opt.WriteOptions) error
}

func (s *service) getRepoInfo(db dbGetter, id uuid.UUID) (info *gitRepoInfo, err error) {
	keyRepo := []byte(fmt.Sprintf("git.repo:%s", id))
	var val []byte
	if val, err = db.Get(keyRepo, nil); err != nil {
		if err == leveldb.ErrNotFound {
			err = git.ErrRepoNotFound
		}
		return
	}
	info = new(gitRepoInfo)
	if err = json.Unmarshal(val, info); err != nil {
		return
	}
	return
}

func (s *service) putRepoInfo(db dbPutter, id uuid.UUID, info *gitRepoInfo) (err error) {
	keyRepo := []byte(fmt.Sprintf("git.repo:%s", id))
	var val []byte
	if val, err = json.Marshal(info); err != nil {
		return
	}
	if err = db.Put(keyRepo, val, nil); err != nil {
		return
	}
	return
}
