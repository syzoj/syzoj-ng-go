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

	"github.com/syzoj/syzoj-ng-go/app/judge"
	"github.com/syzoj/syzoj-ng-go/app/problemset"
)

var log = logrus.StandardLogger()

type service struct {
	db     *leveldb.DB
	tjudge judge.Service
	lock   sync.Mutex
}
type problemsetInfo struct{}
type roleInfo string
type problemInfo struct {
	Name      string    `json:"name"`
	Title     string    `json:"title"`
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

func NewLevelDBProblemset(db *leveldb.DB, tjudge judge.Service) problemset.Service {
	return &service{db: db, tjudge: tjudge}
}

func (*service) getProblemsetInfo(db dbGetter, id uuid.UUID) (info *problemsetInfo, err error) {
	keyProblemset := []byte(fmt.Sprintf("problemset.regular:%s", id))
	var data []byte
	if data, err = db.Get(keyProblemset, nil); err != nil {
		if err == leveldb.ErrNotFound {
			err = problemset.ErrProblemsetNotFound
		}
		return
	}
	info = new(problemsetInfo)
	if err = json.Unmarshal(data, info); err != nil {
		return
	}
	return
}

func (*service) putProblemsetInfo(db dbPutter, id uuid.UUID, info *problemsetInfo) (err error) {
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

func (*service) getProblemInfo(db dbGetter, id uuid.UUID, name string) (info *problemInfo, err error) {
	keyProblem := []byte(fmt.Sprintf("problemset.regular:%s.problem:%s", id, name))
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

func (*service) putProblemInfo(db dbPutter, id uuid.UUID, name string, info *problemInfo) (err error) {
	keyProblem := []byte(fmt.Sprintf("problemset.regular:%s.problem:%s", id, name))
	var data []byte
	if data, err = json.Marshal(info); err != nil {
		return
	}
	if err = db.Put(keyProblem, data, nil); err != nil {
		return
	}
	return
}

func (*service) getRole(db dbGetter, id uuid.UUID, userId uuid.UUID) (role roleInfo, err error) {
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

func (*service) putRole(db dbPutter, id uuid.UUID, userId uuid.UUID, role roleInfo) (err error) {
	keyRole := []byte(fmt.Sprintf("problemset.regular:%s.role:%s", id, userId))
	var data []byte = []byte(role)
	if err = db.Put(keyRole, data, nil); err != nil {
		return
	}
	return
}

func (*service) getSubmissionInfo(db dbGetter, id uuid.UUID, submissionId uuid.UUID) (info *submissionInfo, err error) {
	keySubmission := []byte(fmt.Sprintf("problemset.regular:%s.submission:%s", id, submissionId))
	var data []byte
	if data, err = db.Get(keySubmission, nil); err != nil {
		if err == leveldb.ErrNotFound {
			err = problemset.ErrSubmissionNotFound
			return
		}
	}
	info = new(submissionInfo)
	if err = json.Unmarshal(data, info); err != nil {
		return
	}
	return
}

func (*service) putSubmissionInfo(db dbPutter, id uuid.UUID, submissionId uuid.UUID, info *submissionInfo) (err error) {
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

func (p *service) NewProblemset(OwnerId uuid.UUID) (id uuid.UUID, err error) {
	if id, err = uuid.NewRandom(); err != nil {
		return
	}
	var trans *leveldb.Transaction
	if trans, err = p.db.OpenTransaction(); err != nil {
		return
	}
	defer trans.Discard()
	if _, err = p.getProblemsetInfo(trans, id); err != problemset.ErrProblemsetNotFound {
		if err == nil {
			logrus.Warningf("problemset: UUID duplicate: %s, too little entropy\n", id)
			err = problemset.ErrDuplicateUUID
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

func (p *service) AddProblem(id uuid.UUID, userId uuid.UUID, name string, problemId uuid.UUID) (err error) {
	if !checkProblemName(name) {
		err = problemset.ErrInvalidProblemName
		return
	}
	pinfo := &problemInfo{
		Name:      name,
		ProblemId: problemId,
	}
	if err = p.putProblemInfo(p.db, id, name, pinfo); err != nil {
		return
	}
	return
}

func (p *service) ViewProblem(id uuid.UUID, userId uuid.UUID, name string) (info problemset.ProblemInfo, err error) {
	if !checkProblemName(name) {
		err = problemset.ErrInvalidProblemName
		return
	}
	var pinfo *problemInfo
	if pinfo, err = p.getProblemInfo(p.db, id, name); err != nil {
		return
	}
	info.Name = pinfo.Name
	info.Title = pinfo.Title
	info.ProblemId = pinfo.ProblemId
	return
}

func (p *service) SubmitTraditional(id uuid.UUID, userId uuid.UUID, name string, data problemset.TraditionalSubmissionRequest) (submissionId uuid.UUID, err error) {
	if !checkProblemName(name) {
		err = problemset.ErrInvalidProblemName
		return
	}
	var pinfo *problemInfo
	if pinfo, err = p.getProblemInfo(p.db, id, name); err != nil {
		return
	}
	if submissionId, err = uuid.NewRandom(); err != nil {
		return
	}
	var sinfo submissionInfo = submissionInfo{
		Type:      "traditional",
		UserId:    userId,
		ProblemId: pinfo.ProblemId,
		Traditional: &traditionalSubmissionInfo{
			ProblemId: pinfo.ProblemId,
			Language:  data.Language,
			Code:      data.Code,
			Status:    "",
			Complete:  false,
		},
	}
	if err = p.putSubmissionInfo(p.db, id, submissionId, &sinfo); err != nil {
		return
	}
	p.queueSubmissionWithInfo(id, submissionId, &sinfo)
	return
}

func (p *service) ViewSubmission(id uuid.UUID, userId uuid.UUID, submissionId uuid.UUID) (info problemset.SubmissionInfo, err error) {
	var sinfo *submissionInfo
	if sinfo, err = p.getSubmissionInfo(p.db, id, submissionId); err != nil {
		return
	}
	info.Type = sinfo.Type
	return
}

func (p *service) Close() error {
	return nil
}

type traditionalSubmissionCallback struct {
	p            *service
	id           uuid.UUID
	submissionId uuid.UUID
	sinfo        *submissionInfo
}

func (p *service) queueSubmissionWithInfo(id uuid.UUID, submissionId uuid.UUID, sinfo *submissionInfo) {
	callback := &traditionalSubmissionCallback{
		p:            p,
		id:           id,
		submissionId: submissionId,
		sinfo:        sinfo,
	}
	callback.enqueue()
}

func (c *traditionalSubmissionCallback) enqueue() {
	if _, err := c.p.tjudge.QueueSubmission(&judge.Submission{
		Language:  c.sinfo.Traditional.Language,
		Code:      c.sinfo.Traditional.Code,
		ProblemId: c.sinfo.Traditional.ProblemId,
	}, c); err != nil {
		logrus.WithFields(logrus.Fields{
			"id":           c.id,
			"submissionId": c.submissionId,
			"problemId":    c.sinfo.ProblemId,
		}).Warning("problemset: Failed to enqueue traditional submission")
	}
}

func (c *traditionalSubmissionCallback) OnStart(judge.TaskStartInfo) {

}
func (c *traditionalSubmissionCallback) OnProgress(judge.TaskProgressInfo) {

}
func (c *traditionalSubmissionCallback) OnComplete(judge.TaskCompleteInfo) {
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
