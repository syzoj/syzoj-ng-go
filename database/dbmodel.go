package database

import (
	"context"
	"database/sql"

	"github.com/syzoj/syzoj-ng-go/model"
)

func (t *DatabaseTxn) GetUser(ctx context.Context, ref model.UserRef) (*model.User, error) {
	v := new(model.User)
	err := t.tx.QueryRowContext(ctx, "SELECT user_name, auth FROM user WHERE id=?", ref).Scan(&v.Id, &v.UserName, &v.Auth)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) UpdateUser(ctx context.Context, ref model.UserRef, v *model.User) error {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := t.tx.ExecContext(ctx, "UPDATE user SET user_name=?, auth=? WHERE id=?", v.UserName, v.Auth, v.Id)
	return err
}

func (t *DatabaseTxn) InsertUser(ctx context.Context, v *model.User) error {
	if v.Id == nil {
		ref := model.NewUserRef()
		v.Id = &ref
	}
	_, err := t.tx.ExecContext(ctx, "INSERT INTO user (id, user_name, auth) VALUES (?, ?, ?)", v.Id, v.UserName, v.Auth)
	return err
}

func (t *DatabaseTxn) DeleteUser(ctx context.Context, ref model.UserRef) error {
	_, err := t.tx.ExecContext(ctx, "DELETE FROM user WHERE id=?", ref)
	return err
}

func (t *DatabaseTxn) GetDevice(ctx context.Context, ref model.DeviceRef) (*model.Device, error) {
	v := new(model.Device)
	err := t.tx.QueryRowContext(ctx, "SELECT user, info FROM device WHERE id=?", ref).Scan(&v.Id, &v.User, &v.Info)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) UpdateDevice(ctx context.Context, ref model.DeviceRef, v *model.Device) error {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := t.tx.ExecContext(ctx, "UPDATE device SET user=?, info=? WHERE id=?", v.User, v.Info, v.Id)
	return err
}

func (t *DatabaseTxn) InsertDevice(ctx context.Context, v *model.Device) error {
	if v.Id == nil {
		ref := model.NewDeviceRef()
		v.Id = &ref
	}
	_, err := t.tx.ExecContext(ctx, "INSERT INTO device (id, user, info) VALUES (?, ?, ?)", v.Id, v.User, v.Info)
	return err
}

func (t *DatabaseTxn) DeleteDevice(ctx context.Context, ref model.DeviceRef) error {
	_, err := t.tx.ExecContext(ctx, "DELETE FROM device WHERE id=?", ref)
	return err
}

func (t *DatabaseTxn) GetProblem(ctx context.Context, ref model.ProblemRef) (*model.Problem, error) {
	v := new(model.Problem)
	err := t.tx.QueryRowContext(ctx, "SELECT title FROM problem WHERE id=?", ref).Scan(&v.Id, &v.Title)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) UpdateProblem(ctx context.Context, ref model.ProblemRef, v *model.Problem) error {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := t.tx.ExecContext(ctx, "UPDATE problem SET title=? WHERE id=?", v.Title, v.Id)
	return err
}

func (t *DatabaseTxn) InsertProblem(ctx context.Context, v *model.Problem) error {
	if v.Id == nil {
		ref := model.NewProblemRef()
		v.Id = &ref
	}
	_, err := t.tx.ExecContext(ctx, "INSERT INTO problem (id, title) VALUES (?, ?)", v.Id, v.Title)
	return err
}

func (t *DatabaseTxn) DeleteProblem(ctx context.Context, ref model.ProblemRef) error {
	_, err := t.tx.ExecContext(ctx, "DELETE FROM problem WHERE id=?", ref)
	return err
}

func (t *DatabaseTxn) GetProblemSource(ctx context.Context, ref model.ProblemSourceRef) (*model.ProblemSource, error) {
	v := new(model.ProblemSource)
	err := t.tx.QueryRowContext(ctx, "SELECT data FROM problem_source WHERE id=?", ref).Scan(&v.Id, &v.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) UpdateProblemSource(ctx context.Context, ref model.ProblemSourceRef, v *model.ProblemSource) error {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := t.tx.ExecContext(ctx, "UPDATE problem_source SET data=? WHERE id=?", v.Data, v.Id)
	return err
}

func (t *DatabaseTxn) InsertProblemSource(ctx context.Context, v *model.ProblemSource) error {
	if v.Id == nil {
		ref := model.NewProblemSourceRef()
		v.Id = &ref
	}
	_, err := t.tx.ExecContext(ctx, "INSERT INTO problem_source (id, data) VALUES (?, ?)", v.Id, v.Data)
	return err
}

func (t *DatabaseTxn) DeleteProblemSource(ctx context.Context, ref model.ProblemSourceRef) error {
	_, err := t.tx.ExecContext(ctx, "DELETE FROM problem_source WHERE id=?", ref)
	return err
}

func (t *DatabaseTxn) GetProblemJudger(ctx context.Context, ref model.ProblemJudgerRef) (*model.ProblemJudger, error) {
	v := new(model.ProblemJudger)
	err := t.tx.QueryRowContext(ctx, "SELECT problem, user, type, data FROM problem_judger WHERE id=?", ref).Scan(&v.Id, &v.Problem, &v.User, &v.Type, &v.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) UpdateProblemJudger(ctx context.Context, ref model.ProblemJudgerRef, v *model.ProblemJudger) error {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := t.tx.ExecContext(ctx, "UPDATE problem_judger SET problem=?, user=?, type=?, data=? WHERE id=?", v.Problem, v.User, v.Type, v.Data, v.Id)
	return err
}

func (t *DatabaseTxn) InsertProblemJudger(ctx context.Context, v *model.ProblemJudger) error {
	if v.Id == nil {
		ref := model.NewProblemJudgerRef()
		v.Id = &ref
	}
	_, err := t.tx.ExecContext(ctx, "INSERT INTO problem_judger (id, problem, user, type, data) VALUES (?, ?, ?, ?, ?)", v.Id, v.Problem, v.User, v.Type, v.Data)
	return err
}

func (t *DatabaseTxn) DeleteProblemJudger(ctx context.Context, ref model.ProblemJudgerRef) error {
	_, err := t.tx.ExecContext(ctx, "DELETE FROM problem_judger WHERE id=?", ref)
	return err
}

func (t *DatabaseTxn) GetProblemStatement(ctx context.Context, ref model.ProblemStatementRef) (*model.ProblemStatement, error) {
	v := new(model.ProblemStatement)
	err := t.tx.QueryRowContext(ctx, "SELECT problem, user, data FROM problem_statement WHERE id=?", ref).Scan(&v.Id, &v.Problem, &v.User, &v.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) UpdateProblemStatement(ctx context.Context, ref model.ProblemStatementRef, v *model.ProblemStatement) error {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := t.tx.ExecContext(ctx, "UPDATE problem_statement SET problem=?, user=?, data=? WHERE id=?", v.Problem, v.User, v.Data, v.Id)
	return err
}

func (t *DatabaseTxn) InsertProblemStatement(ctx context.Context, v *model.ProblemStatement) error {
	if v.Id == nil {
		ref := model.NewProblemStatementRef()
		v.Id = &ref
	}
	_, err := t.tx.ExecContext(ctx, "INSERT INTO problem_statement (id, problem, user, data) VALUES (?, ?, ?, ?)", v.Id, v.Problem, v.User, v.Data)
	return err
}

func (t *DatabaseTxn) DeleteProblemStatement(ctx context.Context, ref model.ProblemStatementRef) error {
	_, err := t.tx.ExecContext(ctx, "DELETE FROM problem_statement WHERE id=?", ref)
	return err
}

func (t *DatabaseTxn) GetSubmission(ctx context.Context, ref model.SubmissionRef) (*model.Submission, error) {
	v := new(model.Submission)
	err := t.tx.QueryRowContext(ctx, "SELECT problem_judger, user, data FROM submission WHERE id=?", ref).Scan(&v.Id, &v.ProblemJudger, &v.User, &v.Data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return v, nil
}

func (t *DatabaseTxn) UpdateSubmission(ctx context.Context, ref model.SubmissionRef, v *model.Submission) error {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := t.tx.ExecContext(ctx, "UPDATE submission SET problem_judger=?, user=?, data=? WHERE id=?", v.ProblemJudger, v.User, v.Data, v.Id)
	return err
}

func (t *DatabaseTxn) InsertSubmission(ctx context.Context, v *model.Submission) error {
	if v.Id == nil {
		ref := model.NewSubmissionRef()
		v.Id = &ref
	}
	_, err := t.tx.ExecContext(ctx, "INSERT INTO submission (id, problem_judger, user, data) VALUES (?, ?, ?, ?)", v.Id, v.ProblemJudger, v.User, v.Data)
	return err
}

func (t *DatabaseTxn) DeleteSubmission(ctx context.Context, ref model.SubmissionRef) error {
	_, err := t.tx.ExecContext(ctx, "DELETE FROM submission WHERE id=?", ref)
	return err
}
