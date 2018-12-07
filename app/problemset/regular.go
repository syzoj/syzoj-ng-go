package problemset

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
)

type regularProblemsetProvider struct {
	s *problemsetService
}
type RegularCreateRequest struct {
	OwnerId uuid.UUID
}
type RegularAddTraditionalProblemRequest struct {
	UserId uuid.UUID
	Name   string
	Info   traditionalProblemInfo
}
type RegularAddTraditionalProblemResponse struct {
	ProblemId uuid.UUID
}

func newRegularProblemsetProvider(s *problemsetService) ProblemsetServiceProvider {
	return &regularProblemsetProvider{s: s}
}

func (p *regularProblemsetProvider) NewProblemset(data interface{}) (id uuid.UUID, err error) {
	req := data.(*RegularCreateRequest)
	if id, err = uuid.NewRandom(); err != nil {
		return
	}
	var trans *leveldb.Transaction
	if trans, err = p.s.db.OpenTransaction(); err != nil {
		return
	}
	defer trans.Discard()
	keyType := []byte(fmt.Sprintf("problemset.type:%s", id))
	if err = trans.Put(keyType, []byte("regular"), nil); err != nil {
		return
	}
	keyRole := []byte(fmt.Sprintf("problemset.regular.role:%s", req.OwnerId))
	if err = trans.Put(keyRole, []byte("admin"), nil); err != nil {
		return
	}
	err = trans.Commit()
	return
}

func (p *regularProblemsetProvider) InvokeProblemset(id uuid.UUID, req interface{}, resp interface{}) error {
	switch v := req.(type) {
	case *RegularAddTraditionalProblemRequest:
		keyRole := []byte(fmt.Sprintf("problemset.regular.role:%s", v.UserId))
		if val, err := p.s.db.Get(keyRole, nil); err != nil {
			return err
		} else if string(val) != "admin" {
			return ErrPermissionDenied
		}
		return p.doAddProblem(id, v, resp.(*RegularAddTraditionalProblemResponse))
	default:
		return ErrOperationNotSupported
	}
}

func (p *regularProblemsetProvider) doAddProblem(id uuid.UUID, req *RegularAddTraditionalProblemRequest, resp *RegularAddTraditionalProblemResponse) (err error) {
	problemId := uuid.New()
	var trans *leveldb.Transaction
	if trans, err = p.s.db.OpenTransaction(); err != nil {
		return
	}
	defer trans.Discard()
	keyProblemName := []byte(fmt.Sprintf("problemset.regular.problemname:%s", req.Name))
	var has bool
	if has, err = trans.Has(keyProblemName, nil); has {
		return ErrDuplicateProblemName
	} else if err != nil {
		return
	}
	if err = trans.Put(keyProblemName, problemId[:], nil); err != nil {
		return
	}
	keyProblem := []byte(fmt.Sprintf("problemset.regular.problem:%s", problemId))
	if has, err = trans.Has(keyProblem, nil); has {
		return ErrDuplicateUUID
	} else if err != nil {
		return
	}
	dataBuffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(dataBuffer)
	if err = encoder.Encode(req.Info); err != nil {
		return
	}
	if err = trans.Put(keyProblem, dataBuffer.Bytes(), nil); err != nil {
		return
	}
	err = trans.Commit()
	return
}
