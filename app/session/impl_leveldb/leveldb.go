package impl_leveldb

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/syzoj/syzoj-ng-go/app/session"
)

var log = logrus.StandardLogger()

type leveldbService struct {
	mutex sync.Mutex
	db    *leveldb.DB
}

func NewLevelDBSessionService(db *leveldb.DB) (session.Service, error) {
	srv := &leveldbService{db: db}
	return srv, nil
}

func (s *leveldbService) NewSession() (id uuid.UUID, sess *session.Session, err error) {
	id, err = uuid.NewRandom()
	if err != nil {
		return
	}

	sess = &session.Session{}
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
		sess = nil
		return
	}
	if err = s.db.Put(key, val, nil); err != nil {
		return
	}
	return
}

func (s *leveldbService) GetSession(id uuid.UUID) (sess *session.Session, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	key := []byte(fmt.Sprintf("sess:%s", id))
	var val []byte
	if val, err = s.db.Get(key, nil); err != nil {
		if err == leveldb.ErrNotFound {
			err = session.ErrSessionNotFound
		}
		return
	}
	if err = json.Unmarshal(val, &sess); err != nil {
		sess = nil
		return
	}
	if sess.Expiry.Before(time.Now()) {
		s.db.Delete(key, nil)
		err = session.ErrSessionNotFound
		return
	}
	return
}

func (s *leveldbService) UpdateSession(id uuid.UUID, sess *session.Session) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	key := []byte(fmt.Sprintf("sess:%s", id))
	var val []byte
	if val, err = s.db.Get(key, nil); err != nil {
		return
	}
	var sess2 session.Session
	if err = json.Unmarshal(val, &sess2); err != nil {
		return
	}
	if sess.Version != sess2.Version {
		return session.ErrConcurrentUpdate
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

func (s *leveldbService) Close() error {
	return nil
}

func (s *leveldbService) GarbageCollect() error {
	iter := s.db.NewIterator(util.BytesPrefix([]byte("sess:")), nil)
	defer iter.Release()
	for iter.Next() {
		key, val := iter.Key(), iter.Value()
		var sess session.Session
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
	return nil
}
