package problemset

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"sync"

	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"

	judge_traditional "github.com/syzoj/syzoj-ng-go/app/judge/traditional"
)

type regularProblemsetProvider struct {
	s    *problemsetService
	lock sync.Mutex
}
type RegularCreateRequest struct {
	OwnerId uuid.UUID
}
type RegularAddTraditionalProblemRequest struct {
	UserId uuid.UUID
	Name   string
	Info   struct{}
}
type RegularAddTraditionalProblemResponse struct {
	ProblemId uuid.UUID
}
type RegularSubmitTraditionalProblemRequest struct {
	ProblemId uuid.UUID
	UserId    uuid.UUID
	Language  string
	Code      string
}
type RegularSubmitTraditionalProblemResponse struct {
	SubmissionId uuid.UUID
}
type regularTraditionalProblemSubmission struct {
	UserId   uuid.UUID
	Language string
	Code     string
}
type regularProblemInfo struct {
	Type string
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
	keyType := []byte(fmt.Sprintf("problemset:%s.type", id))
	if err = trans.Put(keyType, []byte("regular"), nil); err != nil {
		return
	}
	keyRole := []byte(fmt.Sprintf("problemset:%s.regular.role:%s", id, req.OwnerId))
	if err = trans.Put(keyRole, []byte("admin"), nil); err != nil {
		return
	}
	err = trans.Commit()
	return
}

func (p *regularProblemsetProvider) InvokeProblemset(id uuid.UUID, req interface{}) (interface{}, error) {
	// TODO Lock based on id
	p.lock.Lock()
	defer p.lock.Unlock()
	switch v := req.(type) {
	case *RegularAddTraditionalProblemRequest:
		keyRole := []byte(fmt.Sprintf("problemset:%s.regular.role:%s", id, v.UserId))
		if val, err := p.s.db.Get(keyRole, nil); err != nil {
			return nil, err
		} else if string(val) != "admin" {
			return nil, ErrPermissionDenied
		}
		return p.doAddTraditionalProblem(id, v)
	case *RegularSubmitTraditionalProblemRequest:
		return p.doSubmitProblem(id, v)
	case *judge_traditional.TraditionalSubmissionResultMessage:
		return nil, p.handleTraditionalSubmissionResult(id, v)
	default:
		return nil, ErrOperationNotSupported
	}
}

func (p *regularProblemsetProvider) Close() error {
	return nil
}

func (p *regularProblemsetProvider) doAddTraditionalProblem(id uuid.UUID, req *RegularAddTraditionalProblemRequest) (resp *RegularAddTraditionalProblemResponse, err error) {
	if !checkProblemName(req.Name) {
		return nil, ErrInvalidProblemName
	}
	var problemId uuid.UUID
	if problemId, err = uuid.NewRandom(); err != nil {
		return
	}
	var trans *leveldb.Transaction
	if trans, err = p.s.db.OpenTransaction(); err != nil {
		return
	}
	defer trans.Discard()
	keyProblemName := []byte(fmt.Sprintf("problemset:%s.problemname:%s", id, req.Name))
	var has bool
	if has, err = trans.Has(keyProblemName, nil); has {
		return nil, ErrDuplicateProblemName
	} else if err != nil {
		return
	}
	if err = trans.Put(keyProblemName, problemId[:], nil); err != nil {
		return
	}
	keyProblem := []byte(fmt.Sprintf("problemset:%s.problem:%s", id, problemId))
	if has, err = trans.Has(keyProblem, nil); has {
		return nil, ErrDuplicateUUID
	} else if err != nil {
		return
	}
	var info regularProblemInfo
	info.Type = "traditional"
	var data []byte
	if data, err = json.Marshal(info); err != nil {
		return
	}
	if err = trans.Put(keyProblem, data, nil); err != nil {
		return
	}
	err = trans.Commit()
	resp = &RegularAddTraditionalProblemResponse{ProblemId: problemId}
	return
}

func (p *regularProblemsetProvider) doSubmitProblem(id uuid.UUID, req *RegularSubmitTraditionalProblemRequest) (resp *RegularSubmitTraditionalProblemResponse, err error) {
	var zeroId uuid.UUID
	if id == zeroId {
		return nil, ErrAnonymousSubmission
	}

	keyProblem := []byte(fmt.Sprintf("problemset:%s.problem:%s", id, req.ProblemId))
	var data []byte
	if data, err = p.s.db.Get(keyProblem, nil); err != nil {
		if err == leveldb.ErrNotFound {
			err = ErrProblemNotFound
		}
		return
	}
	var info regularProblemInfo
	if err = json.Unmarshal(data, &info); err != nil {
		return
	}
	if info.Type != "traditional" {
		return nil, ErrProblemNotFound
	}

	var submissionId uuid.UUID
	if submissionId, err = uuid.NewRandom(); err != nil {
		return
	}

	var submission regularTraditionalProblemSubmission
	submission.UserId = req.UserId
	submission.Language = req.Language
	submission.Code = req.Code
	var submissionData []byte
	if submissionData, err = json.Marshal(submission); err != nil {
		return
	}
	keySubmission := []byte(fmt.Sprintf("problemset:%s.submission:%s", id, submissionId))
	if err = p.s.db.Put(keySubmission, submissionData, &opt.WriteOptions{Sync: true}); err != nil {
		return
	}
	resp = &RegularSubmitTraditionalProblemResponse{SubmissionId: submissionId}
	p.queueTraditionalSubmission(id, submissionId)
	return
}

func (p *regularProblemsetProvider) queueTraditionalSubmission(id uuid.UUID, submissionId uuid.UUID) {
	var err error
	defer func() {
		if err != nil {
			log.Printf("Regular problemset %s: Failed to queue submission %s: %s\n", id, submissionId, err)
		}
	}()
	keySubmission := []byte(fmt.Sprintf("problemset:%s.submission:%s", id, submissionId))
	var submissionData []byte
	if submissionData, err = p.s.db.Get(keySubmission, nil); err != nil {
		return
	}
	var submission regularTraditionalProblemSubmission
	if err = json.Unmarshal(submissionData, &submission); err != nil {
		return
	}
	if err = p.s.ts.QueueSubmission(id, submissionId, &judge_traditional.TraditionalSubmission{
		Language: submission.Language,
		Code:     submission.Code,
	}); err != nil {
		return
	}
}

func (p *regularProblemsetProvider) handleTraditionalSubmissionResult(id uuid.UUID, result *judge_traditional.TraditionalSubmissionResultMessage) error {
	// TODO: save result
	log.Printf("Regular problemset %s: Received result for submission %s: %s\n", id, result.SubmissionId, result.Result.Status)
	return nil
}

var problemNameRegexp = regexp.MustCompile("^[0-9A-Z]{1,16}")

func checkProblemName(problemName string) bool {
	return problemNameRegexp.MatchString(problemName)
}
