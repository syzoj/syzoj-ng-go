package auth

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type leveldbAuthService struct {
	mutex sync.Mutex
	db    *leveldb.DB
}

func NewLevelDBAuthService(db *leveldb.DB) (AuthService, error) {
	return &leveldbAuthService{
		db: db,
	}, nil
}

func (s *leveldbAuthService) RegisterUser(userName string, password string) (id uuid.UUID, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if id, err = uuid.NewRandom(); err != nil {
		return
	}
	var trans *leveldb.Transaction
	if trans, err = s.db.OpenTransaction(); err != nil {
		return
	}
	defer trans.Discard()

	usernameKey := []byte(fmt.Sprintf("auth.username:%s", userName))
	if _, err = trans.Get(usernameKey, nil); err != leveldb.ErrNotFound {
		if err == nil {
			err = ErrDuplicateUserName
			return
		}
		return
	}
	if err = trans.Put(usernameKey, id[:], nil); err != nil {
		return
	}

	userAuthKey := []byte(fmt.Sprintf("auth.user:%s", id))
	var val []byte
	var authInfo UserAuthInfo = PasswordAuth(password)
	if val, err = json.Marshal(authInfo); err != nil {
		return
	}
	if err = trans.Put(userAuthKey, val, &opt.WriteOptions{Sync: true}); err != nil {
		return
	}
	err = trans.Commit()
	return
}

func (s *leveldbAuthService) LoginUser(userName string, password string) (userId uuid.UUID, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var snapshot *leveldb.Snapshot
	if snapshot, err = s.db.GetSnapshot(); err != nil {
		return
	}
	defer snapshot.Release()

	usernameKey := []byte(fmt.Sprintf("auth.username:%s", userName))
	var val1 []byte
	if val1, err = snapshot.Get(usernameKey, nil); err != nil {
		if err == leveldb.ErrNotFound {
			err = ErrUserNotFound
		}
		return
	}
	var id uuid.UUID
	if id, err = uuid.FromBytes(val1); err != nil {
		return
	}

	userAuthKey := []byte(fmt.Sprintf("auth.user:%s", id))
	var val2 []byte
	if val2, err = snapshot.Get(userAuthKey, nil); err != nil {
		return
	}
	var authInfo UserAuthInfo
	if err = json.Unmarshal(val2, &authInfo); err != nil {
		return
	}
	if authInfo.PasswordInfo.Verify(password) != nil {
		err = ErrPasswordIncorrect
		return
	}

	userId = id
	return
}
