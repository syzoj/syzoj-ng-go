package impl_leveldb

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sync"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"

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
	Type        string    `json:"type"`
	UserId      uuid.UUID `json:"user_id"`
	ProblemName string    `json:"problem_name"`
	ProblemId   uuid.UUID `json:"problem_id"`
	// Simulate union type with pointers
	Traditional *judge.TraditionalSubmission `json:"traditional"`
	Complete    bool                         `json:"complete"`
	Result      judge.TaskCompleteInfo       `json:"result"`
}

func NewLevelDBProblemset(db *leveldb.DB, tjudge judge.Service) problemset.Service {
	return &service{db: db, tjudge: tjudge}
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

func (p *service) ListProblem(id uuid.UUID, userId uuid.UUID) (info []problemset.ProblemInfo, err error) {
	keyProblemPrefix := []byte(fmt.Sprintf("problemset.regular:%s.problem:", id))
	var it = p.db.NewIterator(util.BytesPrefix(keyProblemPrefix), nil)
	defer it.Release()
	for it.Next() {
		v := it.Value()
		info = append(info, problemset.ProblemInfo{})
		if err = json.Unmarshal(v, &info[len(info)-1]); err != nil {
			return
		}
	}
	if err = it.Error(); err != nil {
		return
	}
	return
}

func (p *service) SubmitTraditional(id uuid.UUID, userId uuid.UUID, name string, data judge.TraditionalSubmission) (submissionId uuid.UUID, err error) {
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
		Type:        "traditional",
		UserId:      userId,
		ProblemName: name,
		ProblemId:   pinfo.ProblemId,
		Complete:    false,
		Traditional: &data,
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
	info.UserId = sinfo.UserId
	info.ProblemName = sinfo.ProblemName
	info.Complete = sinfo.Complete
	if info.Type == "traditional" {
		info.Traditional = sinfo.Traditional
	}
	info.Result = sinfo.Result
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

func (c *traditionalSubmissionCallback) getFields() logrus.Fields {
	return logrus.Fields{
		"problemsetId": c.id,
		"submissionId": c.submissionId,
		"problemId":    c.sinfo.ProblemId,
	}
}

func (c *traditionalSubmissionCallback) enqueue() {
	if _, err := c.p.tjudge.QueueSubmission(&judge.Submission{
		ProblemId:   c.sinfo.ProblemId,
		Traditional: *c.sinfo.Traditional,
	}, c); err != nil {
		logrus.WithFields(c.getFields()).Warning("problemset: Failed to enqueue traditional submission")
	}
}

func (c *traditionalSubmissionCallback) OnStart(judge.TaskStartInfo) {

}
func (c *traditionalSubmissionCallback) OnProgress(judge.TaskProgressInfo) {

}
func (c *traditionalSubmissionCallback) OnComplete(info judge.TaskCompleteInfo) {
	logrus.WithFields(c.getFields()).Info("Submission completed")
	c.sinfo.Result = info
	c.sinfo.Complete = true
	var err error
	if err = c.p.putSubmissionInfo(c.p.db, c.id, c.submissionId, c.sinfo); err != nil {
		logrus.WithFields(c.getFields()).Warning("problemset: Failed to store submission result")
	}
}
func (c *traditionalSubmissionCallback) OnError(err error) {
	logrus.WithFields(c.getFields()).Warningf("Submission errored: %s\n", err.Error())
}

var problemNameRegexp = regexp.MustCompile("^[0-9A-Z]{1,16}$")

func checkProblemName(problemName string) bool {
	return problemNameRegexp.MatchString(problemName)
}
