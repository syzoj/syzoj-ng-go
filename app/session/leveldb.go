package session

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var log = logrus.StandardLogger()

type leveldbSessionService struct {
	mutex     sync.Mutex
	db        *leveldb.DB
	closeChan chan struct{}
}

func NewLevelDBSessionService(db *leveldb.DB) (SessionService, error) {
	srv := &leveldbSessionService{
		db:        db,
		closeChan: make(chan struct{}),
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

func (s *leveldbSessionService) Close() error {
	close(s.closeChan)
	return nil
}

func (s *leveldbSessionService) collectGarbage() {
	iter := s.db.NewIterator(util.BytesPrefix([]byte("sess:")), nil)
	defer iter.Release()
	for iter.Next() {
		key, val := iter.Key(), iter.Value()
		var sess Session
		if err := json.Unmarshal(val, &sess); err != nil {
			log.WithFields(logrus.Fields{
				"session-key": string(key),
				"error":       err,
			}).Warning("Failed to unmarshal session")
			continue
		}
		if time.Now().After(sess.Expiry) {
			log.WithField("session-key", string(key)).Debug("Expiring session")
			if err := s.db.Delete(key, nil); err != nil {
				// Expect race condition here
				log.WithFields(logrus.Fields{
					"session-key": string(key),
					"error":       err,
				}).Warning("Failed to expire session")
			}
		}
	}
}

func (s *leveldbSessionService) runGC() {
	for {
		s.collectGarbage()
		select {
		case <-s.closeChan:
			return
		case <-time.After(time.Hour):
		}
	}
}
