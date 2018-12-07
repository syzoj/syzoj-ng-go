package problemset

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
)

type ProblemsetService interface {
	NewProblemset(ptype string, data interface{}) (uuid.UUID, error)
	InvokeProblemset(id uuid.UUID, req interface{}, resp interface{}) error
}

type ProblemsetServiceProvider interface {
	NewProblemset(data interface{}) (uuid.UUID, error)
	InvokeProblemset(id uuid.UUID, req interface{}, resp interface{}) error
}

var ErrInvalidProblemsetType = errors.New("Invalid problemset type")
var ErrProblemsetNotFound = errors.New("Problemset not found")
var ErrOperationNotSupported = errors.New("Operation not supported")
var ErrDuplicateProblemName = errors.New("Duplicate problem name")
var ErrDuplicateUUID = errors.New("UUID dupication")
var ErrPermissionDenied = errors.New("Permission denied")

var psetList = map[string]func(*problemsetService) ProblemsetServiceProvider{
	"regular": newRegularProblemsetProvider,
}

type problemsetService struct {
	db       *leveldb.DB
	provider map[string]ProblemsetServiceProvider
}

func NewProblemsetService(db *leveldb.DB) (ProblemsetService, error) {
	s := &problemsetService{db: db}
	s.provider = make(map[string]ProblemsetServiceProvider)
	for key, value := range psetList {
		s.provider[key] = value(s)
	}
	return s, nil
}

func (s *problemsetService) NewProblemset(ptype string, data interface{}) (id uuid.UUID, err error) {
	if provider, ok := s.provider[ptype]; !ok {
		err = ErrInvalidProblemsetType
		return
	} else {
		id, err = provider.NewProblemset(data)
		return
	}
}

func (s *problemsetService) InvokeProblemset(id uuid.UUID, req interface{}, resp interface{}) (err error) {
	key := []byte(fmt.Sprintf("problemset.type:%s", id))
	var val []byte
	if val, err = s.db.Get(key, nil); err != nil {
		if err == leveldb.ErrNotFound {
			err = ErrProblemsetNotFound
		}
		return
	} else {
		if provider, ok := s.provider[string(val)]; !ok {
			// TODO: Handle data inconsistency
			return ErrInvalidProblemsetType
		} else {
			return provider.InvokeProblemset(id, req, resp)
		}
	}
}
