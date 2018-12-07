package session

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type leveldbSessionService struct {
	mutex sync.Mutex
	db    *leveldb.DB
}

func NewLevelDBSessionService(db *leveldb.DB) (SessionService, error) {
	srv := &leveldbSessionService{
		db: db,
	}
	go srv.runGC()
	return srv, nil
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
	sess.Expiry = time.Now().Add(time.Hour * 24 * 30)

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
		if err == leveldb.ErrNotFound {
			err = ErrSessionNotFound
		}
		return
	}
	if err = json.Unmarshal(val, &sess); err != nil {
		return
	}
	if sess.Expiry.Before(time.Now()) {
		s.db.Delete(key, nil)
		err = ErrSessionNotFound
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
	sess.Expiry = time.Now().Add(time.Hour * 24 * 30)
	if sess.Version, err = uuid.NewRandom(); err != nil {
		return
	}
	if val, err = json.Marshal(sess); err != nil {
		return
	}
	err = s.db.Put(key, val, nil)
	return
}

func (s *leveldbSessionService) collectGarbage() {
	iter := s.db.NewIterator(util.BytesPrefix([]byte("sess:")), nil)
	for iter.Next() {
		key, val := iter.Key(), iter.Value()
		var sess Session
		if err := json.Unmarshal(val, &sess); err != nil {
			log.Printf("Failed to unmarshal session %s: %s\n", string(key), err)
			continue
		}
		if time.Now().After(sess.Expiry) {
			log.Printf("Expiring session %s\n", string(key))
			if err := s.db.Delete(key, nil); err != nil {
				// Expect race condition here
				log.Printf("Warning: failed to delete session %s: %s\n", string(key), err)
			}
		}
	}
}

func (s *leveldbSessionService) runGC() {
	for {
		s.collectGarbage()
		time.Sleep(time.Hour)
	}
}
