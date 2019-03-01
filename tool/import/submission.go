package tool_import

import (
	"context"
	"database/sql"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

type submission struct {
	id         int64
	userId     sql.NullInt64
	code       sql.NullString
	language   sql.NullString
	status     sql.NullString
	score      sql.NullFloat64
	totalTime  sql.NullInt64
	submitTime sql.NullInt64
	problemId  sql.NullInt64
	typ        sql.NullInt64
}

func (i *importer) readSubmissions(submissions chan<- *submission) {
	var err error
	var rows *sql.Rows
	if rows, err = i.db.Query("SELECT id, user_id, code, language, status, score, total_time, submit_time, problem_id, type FROM judge_state"); err != nil {
		log.WithError(err).Error("Error importing submissions from MySQL")
		close(submissions)
		return
	}
	for rows.Next() {
		s := new(submission)
		err = rows.Scan(&s.id, &s.userId, &s.code, &s.language, &s.status, &s.score, &s.totalTime, &s.submitTime, &s.problemId, &s.typ)
		if err != nil {
			log.WithError(err).Error("Error reading problem")
			err = nil
		}
		submissions <- s
	}
	close(submissions)
}

func (i *importer) writeSubmissions(submissions <-chan *submission) {
	for s := range submissions {
		if !s.typ.Valid || s.typ.Int64 != 0 {
			continue // skip non-normal submission
		}
		submissionModel := new(model.Submission)
		submissionModel.Id = model.NewObjectIDProto()
		submissionModel.Public = proto.Bool(true)
		if userId, found := i.userId[s.userId.Int64]; found {
			submissionModel.User = model.ObjectIDProto(userId)
		}
		if problemId, found := i.problemId[s.problemId.Int64]; found {
			submissionModel.Problem = model.ObjectIDProto(problemId)
		}
		submissionModel.Content = new(model.SubmissionContent)
		if s.language.Valid {
			submissionModel.Content.Language = proto.String(s.language.String)
		}
		if s.code.Valid {
			submissionModel.Content.Code = proto.String(s.code.String)
		}
		if s.submitTime.Valid {
			submissionModel.SubmitTime, _ = ptypes.TimestampProto(time.Unix(s.submitTime.Int64, 0))
		}
		doneResults := map[string]struct{}{
			"Accepted":              {},
			"Compile Error":         {},
			"File Error":            {},
			"Judgement Failed":      {},
			"Memory Limit Exceeded": {},
			"No Testdata":           {},
			"Partially Correct":     {},
			"Runtime Error":         {},
			"System Error":          {},
			"Time Limit Exceeded":   {},
			"Wrong Answer":          {},
			// Unknown and waiting omitted
		}
		if _, found := doneResults[s.status.String]; found {
			submissionModel.Result = new(model.SubmissionResult)
			if s.status.Valid {
				submissionModel.Result.Status = proto.String(s.status.String)
			}
			if s.score.Valid {
				submissionModel.Result.Score = proto.Float64(s.score.Float64)
			}
		}
		if _, err := i.mongodb.Collection("submission").InsertOne(context.Background(), submissionModel); err != nil {
			log.WithField("id", s.id).WithError(err).Warning("Error inserting submission")
		}
	}
}
