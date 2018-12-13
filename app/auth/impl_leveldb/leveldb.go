package impl_leveldb

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"sync"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/syzoj/syzoj-ng-go/app/auth"
)

type leveldbService struct {
	mutex sync.Mutex
	db    *leveldb.DB
}

type authData struct {
	UserName     string `json:"username"`
	PasswordHash string `json:"password"`
}

var userNameRegex = regexp.MustCompile("^[0-9A-Za-z]{3,32}$")

func checkUserName(userName string) bool {
	return userNameRegex.MatchString(userName)
}

func NewLevelDBAuthService(db *leveldb.DB) (auth.Service, error) {
	return &leveldbService{
		db: db,
	}, nil
}

func (s *leveldbService) RegisterUser(userName string, password string) (id uuid.UUID, err error) {
	if !checkUserName(userName) {
		err = auth.ErrInvalidUserName
		return
	}
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

	keyUserName := []byte(fmt.Sprintf("auth.username:%s", userName))
	if _, err = trans.Get(keyUserName, nil); err != leveldb.ErrNotFound {
		if err == nil {
			err = auth.ErrDuplicateUserName
			return
		}
	}
	if err = trans.Put(keyUserName, id[:], nil); err != nil {
		return
	}

	var passwordHash []byte
	if passwordHash, err = bcrypt.GenerateFromPassword([]byte(password), 0); err != nil {
		return
	}
	info := authData{UserName: userName, PasswordHash: base64.StdEncoding.EncodeToString(passwordHash)}
	var valUser []byte
	if valUser, err = json.Marshal(info); err != nil {
		return
	}

	keyUser := []byte(fmt.Sprintf("auth.user:%s", id))
	if err = trans.Put(keyUser, valUser, nil); err != nil {
		return
	}
	err = trans.Commit()
	return
}

func (s *leveldbService) LoginUser(userName string, password string) (userId uuid.UUID, err error) {
	if !checkUserName(userName) {
		err = auth.ErrInvalidUserName
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var snapshot *leveldb.Snapshot
	if snapshot, err = s.db.GetSnapshot(); err != nil {
		return
	}
	defer snapshot.Release()

	keyUserName := []byte(fmt.Sprintf("auth.username:%s", userName))
	var valUserId []byte
	if valUserId, err = snapshot.Get(keyUserName, nil); err != nil {
		if err == leveldb.ErrNotFound {
			err = auth.ErrUserNotFound
		}
		return
	}
	var id uuid.UUID
	if id, err = uuid.FromBytes(valUserId); err != nil {
		return
	}

	keyUser := []byte(fmt.Sprintf("auth.user:%s", id))
	var valUser []byte
	if valUser, err = snapshot.Get(keyUser, nil); err != nil {
		return
	}
	var info authData
	if err = json.Unmarshal(valUser, &info); err != nil {
		return
	}
	var passwordHash []byte
	if passwordHash, err = base64.StdEncoding.DecodeString(info.PasswordHash); err != nil {
		return
	}
	if bcrypt.CompareHashAndPassword(passwordHash, []byte(password)) != nil {
		err = auth.ErrPasswordIncorrect
		return
	}

	userId = id
	return
}

func (s *leveldbService) Close() error {
	return nil
}
