package session

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
)

type leveldbSessionService struct {
	mutex sync.Mutex
	db    *leveldb.DB
}

func NewLevelDBSessionService(db *leveldb.DB) (SessionService, error) {
	return &leveldbSessionService{
		db: db,
	}, nil
}

func (s *leveldbSessionService) NewSession() (id uuid.UUID, sess *Session, err error) {
	id, err = uuid.NewRandom()
	if err != nil {
		return
	}

	sess = &Session{}
	sess.Version, err = uuid.NewRandom()
	if err != nil {
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	key := []byte(fmt.Sprintf("sess:%s", id))
	var val []byte
	if val, err = json.Marshal(sess); err != nil {
		return
	}
	if err = s.db.Put(key, val, nil); err != nil {
		return
	}
	return
}

func (s *leveldbSessionService) GetSession(id uuid.UUID) (sess *Session, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	key := []byte(fmt.Sprintf("sess:%s", id))
	var val []byte
	if val, err = s.db.Get(key, nil); err != nil {
		return
	}
	if err = json.Unmarshal(val, &sess); err != nil {
		return
	}
	return
}

func (s *leveldbSessionService) UpdateSession(id uuid.UUID, sess *Session) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	key := []byte(fmt.Sprintf("sess:%s", id))
	var val []byte
	if val, err = s.db.Get(key, nil); err != nil {
		return
	}
	var sess2 Session
	if err = json.Unmarshal(val, &sess2); err != nil {
		return
	}
	if sess.Version != sess2.Version {
		return ConcurrentUpdateError
	}
	if sess.Version, err = uuid.NewRandom(); err != nil {
		return
	}
	if val, err = json.Marshal(sess); err != nil {
		return
	}
	err = s.db.Put(key, val, nil)
	return
}
