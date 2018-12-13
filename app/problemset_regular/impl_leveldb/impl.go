package impl_leveldb

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sync"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/syzoj/syzoj-ng-go/app/judge_traditional"
	"github.com/syzoj/syzoj-ng-go/app/problemset_regular"
)

var log = logrus.StandardLogger()

type problemset struct {
	db     *leveldb.DB
	tjudge judge_traditional.Service
	lock   sync.Mutex
}
type problemsetInfo struct{}
type roleInfo string
type problemInfo struct {
	Type      string    `json:"type"`
	ProblemId uuid.UUID `json:"problem_id"`
}
type submissionInfo struct {
	Type      string    `json:"type"`
	UserId    uuid.UUID `json:"user_id"`
	ProblemId uuid.UUID `json:"problem_id"`
	// Simulate union type with pointers
	Traditional *traditionalSubmissionInfo `json:"traditional"`
}
type traditionalSubmissionInfo struct {
	ProblemId uuid.UUID `json:"problem_id"`
	Language  string    `json:"language"`
	Code      string    `json:"code"`
	Complete  bool      `json:"complete"`
	Status    string    `json:"result"`
}

type dbGetter interface {
	Get([]byte, *opt.ReadOptions) ([]byte, error)
}
type dbPutter interface {
	Put([]byte, []byte, *opt.WriteOptions) error
}
type dbDeleter interface {
	Delete([]byte, *opt.WriteOptions) error
}

func NewLevelDBProblemset(db *leveldb.DB, tjudge judge_traditional.Service) problemset_regular.Service {
	return &problemset{db: db, tjudge: tjudge}
}

func (*problemset) getProblemsetInfo(db dbGetter, id uuid.UUID) (info *problemsetInfo, err error) {
	keyProblemset := []byte(fmt.Sprintf("problemset.regular:%s", id))
	var data []byte
	if data, err = db.Get(keyProblemset, nil); err != nil {
		if err == leveldb.ErrNotFound {
			err = problemset_regular.ErrProblemsetNotFound
		}
		return
	}
	info = new(problemsetInfo)
	if err = json.Unmarshal(data, info); err != nil {
		return
	}
	return
}

func (*problemset) putProblemsetInfo(db dbPutter, id uuid.UUID, info *problemsetInfo) (err error) {
	keyProblemset := []byte(fmt.Sprintf("problemset.regular:%s", id))
	var data []byte
	if data, err = json.Marshal(info); err != nil {
		return
	}
	if err = db.Put(keyProblemset, data, nil); err != nil {
		return
	}
	return
}

func (*problemset) getProblemIdByName(db dbGetter, id uuid.UUID, name string) (problemId uuid.UUID, err error) {
	keyProblemName := []byte(fmt.Sprintf("problemset.regular:%s.problem.name:%s", id, name))
	var data []byte
	if data, err = db.Get(keyProblemName, nil); err != nil {
		if err == leveldb.ErrNotFound {
			err = problemset_regular.ErrProblemNotFound
		}
		return
	}
	if problemId, err = uuid.FromBytes(data); err != nil {
		return
	}
	return
}

func (*problemset) putProblemIdToName(db dbPutter, id uuid.UUID, name string, problemId uuid.UUID) (err error) {
	keyProblemName := []byte(fmt.Sprintf("problemset.regular:%s.problem.name:%s", id, name))
	if err = db.Put(keyProblemName, problemId[:], nil); err != nil {
		return
	}
	return
}

func (*problemset) deleteProblemName(db dbDeleter, id uuid.UUID, name string) (err error) {
	keyProblemName := []byte(fmt.Sprintf("problemset.regular:%s.problem.name:%s", id, name))
	if err = db.Delete(keyProblemName, nil); err != nil {
		return
	}
	return
}

func (*problemset) getProblemInfo(db dbGetter, id uuid.UUID, problemId uuid.UUID) (info *problemInfo, err error) {
	keyProblem := []byte(fmt.Sprintf("problemset.regular:%s.problem:%s", id, problemId))
	var data []byte
	if data, err = db.Get(keyProblem, nil); err != nil {
		return
	}
	info = new(problemInfo)
	if err = json.Unmarshal(data, info); err != nil {
		return
	}
	return
}

func (*problemset) putProblemInfo(db dbPutter, id uuid.UUID, problemId uuid.UUID, info *problemInfo) (err error) {
	keyProblem := []byte(fmt.Sprintf("problemset.regular:%s.problem:%s", id, problemId))
	var data []byte
	if data, err = json.Marshal(info); err != nil {
		return
	}
	if err = db.Put(keyProblem, data, nil); err != nil {
		return
	}
	return
}

func (*problemset) getRole(db dbGetter, id uuid.UUID, userId uuid.UUID) (role roleInfo, err error) {
	keyRole := []byte(fmt.Sprintf("problemset.regular:%s.role:%s", id, userId))
	var data []byte
	if data, err = db.Get(keyRole, nil); err != nil {
		if err == leveldb.ErrNotFound {
			// Default role is empty
			return "", nil
		}
		return
	}
	role = roleInfo(data)
	return
}

func (*problemset) putRole(db dbPutter, id uuid.UUID, userId uuid.UUID, role roleInfo) (err error) {
	keyRole := []byte(fmt.Sprintf("problemset.regular:%s.role:%s", id, userId))
	var data []byte = []byte(role)
	if err = db.Put(keyRole, data, nil); err != nil {
		return
	}
	return
}

func (*problemset) getSubmissionInfo(db dbGetter, id uuid.UUID, submissionId uuid.UUID) (info *submissionInfo, err error) {
	keySubmission := []byte(fmt.Sprintf("problemset.regular:%s.submission:%s", id, submissionId))
	var data []byte
	if data, err = db.Get(keySubmission, nil); err != nil {
		if err == leveldb.ErrNotFound {
			err = problemset_regular.ErrSubmissionNotFound
			return
		}
	}
	info = new(submissionInfo)
	if err = json.Unmarshal(data, info); err != nil {
		return
	}
	return
}

func (*problemset) putSubmissionInfo(db dbPutter, id uuid.UUID, submissionId uuid.UUID, info *submissionInfo) (err error) {
	keySubmission := []byte(fmt.Sprintf("problemset.regular:%s.submission:%s", id, submissionId))
	var data []byte
	if data, err = json.Marshal(info); err != nil {
		return
	}
	if err = db.Put(keySubmission, data, nil); err != nil {
		return
	}
	return
}

func (p *problemset) NewProblemset(OwnerId uuid.UUID) (id uuid.UUID, err error) {
	if id, err = uuid.NewRandom(); err != nil {
		return
	}
	var trans *leveldb.Transaction
	if trans, err = p.db.OpenTransaction(); err != nil {
		return
	}
	defer trans.Discard()
	if _, err = p.getProblemsetInfo(trans, id); err != problemset_regular.ErrProblemsetNotFound {
		if err == nil {
			logrus.Warningf("problemset_regular: UUID duplicate: %s, too little entropy\n", id)
			err = problemset_regular.ErrDuplicateUUID
		}
		return
	}
	if err = p.putProblemsetInfo(trans, id, &problemsetInfo{}); err != nil {
		return
	}
	if err = p.putRole(trans, id, OwnerId, "admin"); err != nil {
		return
	}
	if err = trans.Commit(); err != nil {
		return
	}
	return
}

func (p *problemset) AddTraditionalProblem(id uuid.UUID, userId uuid.UUID, name string, problemId uuid.UUID) (err error) {
	if !checkProblemName(name) {
		err = problemset_regular.ErrInvalidProblemName
		return
	}
	var trans *leveldb.Transaction
	if trans, err = p.db.OpenTransaction(); err != nil {
		return
	}
	defer trans.Discard()
	if _, err = p.getProblemIdByName(trans, id, name); err != problemset_regular.ErrProblemNotFound {
		if err == nil {
			err = problemset_regular.ErrDuplicateProblemName
		}
		return
	}
	if err = p.putProblemIdToName(trans, id, name, problemId); err != nil {
		return
	}
	pinfo := &problemInfo{
		Type:      "traditional",
		ProblemId: problemId,
	}
	if err = p.putProblemInfo(trans, id, problemId, pinfo); err != nil {
		return
	}
	if err = trans.Commit(); err != nil {
		return
	}
	return
}

func (p *problemset) ViewProblem(id uuid.UUID, userId uuid.UUID, name string) (info problemset_regular.ProblemInfo, err error) {
	if !checkProblemName(name) {
		err = problemset_regular.ErrInvalidProblemName
		return
	}
	var problemId uuid.UUID
	if problemId, err = p.getProblemIdByName(p.db, id, name); err != nil {
		return
	}
	var pinfo *problemInfo
	if pinfo, err = p.getProblemInfo(p.db, id, problemId); err != nil {
		return
	}
	info.Type = pinfo.Type
	return
}

func (p *problemset) SubmitTraditional(id uuid.UUID, userId uuid.UUID, name string, data problemset_regular.TraditionalSubmissionRequest) (submissionId uuid.UUID, err error) {
	if !checkProblemName(name) {
		err = problemset_regular.ErrInvalidProblemName
		return
	}
	var trans *leveldb.Transaction
	if trans, err = p.db.OpenTransaction(); err != nil {
		return
	}
	defer trans.Discard()
	var problemId uuid.UUID
	if problemId, err = p.getProblemIdByName(trans, id, name); err != nil {
		return
	}
	var pinfo *problemInfo
	if pinfo, err = p.getProblemInfo(trans, id, problemId); err != nil {
		return
	}
	if pinfo.Type != "traditional" {
		err = problemset_regular.ErrProblemNotFound
		return
	}
	if submissionId, err = uuid.NewRandom(); err != nil {
		return
	}
	var sinfo submissionInfo = submissionInfo{
		Type:      "traditional",
		UserId:    userId,
		ProblemId: problemId,
		Traditional: &traditionalSubmissionInfo{
			ProblemId: pinfo.ProblemId,
			Language:  data.Language,
			Code:      data.Code,
			Status:    "",
			Complete:  false,
		},
	}
	if err = p.putSubmissionInfo(trans, id, submissionId, &sinfo); err != nil {
		return
	}
	p.queueSubmissionWithInfo(id, submissionId, &sinfo)
	return
}

func (p *problemset) ViewSubmission(id uuid.UUID, userId uuid.UUID, submissionId uuid.UUID) (info problemset_regular.SubmissionInfo, err error) {
	var sinfo *submissionInfo
	if sinfo, err = p.getSubmissionInfo(p.db, id, submissionId); err != nil {
		return
	}
	info.Type = sinfo.Type
	return
}

func (p *problemset) Close() error {
	return nil
}

type traditionalSubmissionCallback struct {
	p            *problemset
	id           uuid.UUID
	submissionId uuid.UUID
	sinfo        *submissionInfo
}

func (p *problemset) queueSubmissionWithInfo(id uuid.UUID, submissionId uuid.UUID, sinfo *submissionInfo) {
	callback := &traditionalSubmissionCallback{
		p:            p,
		id:           id,
		submissionId: submissionId,
		sinfo:        sinfo,
	}
	callback.enqueue()
}

func (c *traditionalSubmissionCallback) enqueue() {
	if _, err := c.p.tjudge.QueueSubmission(&judge_traditional.Submission{
		Language:  c.sinfo.Traditional.Language,
		Code:      c.sinfo.Traditional.Code,
		ProblemId: c.sinfo.Traditional.ProblemId,
	}, c); err != nil {
		logrus.WithFields(logrus.Fields{
			"id":           c.id,
			"submissionId": c.submissionId,
			"problemId":    c.sinfo.ProblemId,
		}).Warning("problemset_regular: Failed to enqueue traditional submission")
	}
}

func (c *traditionalSubmissionCallback) OnStart(judge_traditional.TaskStartInfo) {

}
func (c *traditionalSubmissionCallback) OnProgress(judge_traditional.TaskProgressInfo) {

}
func (c *traditionalSubmissionCallback) OnComplete(judge_traditional.TaskCompleteInfo) {
	logrus.WithFields(logrus.Fields{
		"id":           c.id,
		"submissionId": c.submissionId,
		"problemId":    c.sinfo.ProblemId,
	}).Info("Submission completed")
}
func (c *traditionalSubmissionCallback) OnError(err error) {
	logrus.WithFields(logrus.Fields{
		"id":           c.id,
		"submissionId": c.submissionId,
		"problemId":    c.sinfo.ProblemId,
		"err":          err,
	}).Info("Submission errored")
}

var problemNameRegexp = regexp.MustCompile("^[0-9A-Z]{1,16}")

func checkProblemName(problemName string) bool {
	return problemNameRegexp.MatchString(problemName)
}
