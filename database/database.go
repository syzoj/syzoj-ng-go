package database

import (
	"context"
	"database/sql"
	"sync/atomic"

	"github.com/sirupsen/logrus"

	"github.com/syzoj/syzoj-ng-go/model"
)

type Database struct {
	db *sql.DB
}

type DatabaseTxn struct {
	tx   *sql.Tx
	done int32
}

var log = logrus.StandardLogger()

func Open(driverName, dataSourceName string) (*Database, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) OpenTxn(ctx context.Context) (*DatabaseTxn, error) {
	tx, err := d.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	tx2 := &DatabaseTxn{tx: tx}
	go func() {
		<-ctx.Done()
		if atomic.LoadInt32(&tx2.done) == 0 {
			log.Warning("Detected a transaction that wasn't closed before context done")
		}
	}()
	return tx2, nil
}

func (d *Database) OpenReadonlyTxn(ctx context.Context) (*DatabaseTxn, error) {
	tx, err := d.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, err
	}
	tx2 := &DatabaseTxn{tx: tx}
	go func() {
		<-ctx.Done()
		if atomic.LoadInt32(&tx2.done) == 0 {
			log.Warning("Detected a transaction that wasn't closed before context done")
		}
	}()
	return tx2, nil
}

func (t *DatabaseTxn) Commit() error {
	atomic.StoreInt32(&t.done, 1)
	return t.tx.Commit()
}

func (t *DatabaseTxn) Rollback() error {
	atomic.StoreInt32(&t.done, 1)
	return t.tx.Rollback()
}

func (t *DatabaseTxn) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

func (t *DatabaseTxn) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

func (t *DatabaseTxn) GetUser(ctx context.Context, ref model.UserRef) (*model.User, error) {
	user := new(model.User)
	err := t.tx.QueryRowContext(ctx, "SELECT id, username, auth FROM user WHERE id=?", ref).Scan(&user.Id, &user.UserName, &user.Auth)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (t *DatabaseTxn) UpdateUser(ctx context.Context, ref model.UserRef, v *model.User) error {
	if v.Id == nil || v.GetId() != ref {
		panic("ref and v does not match")
	}
	_, err := t.tx.ExecContext(ctx, "UPDATE user SET username=?, auth=? WHERE id=?", v.UserName, v.Auth, v.Id)
	return err
}

func (t *DatabaseTxn) InsertUser(ctx context.Context, v *model.User) error {
	if v.Id == nil {
		ref := model.NewUserRef()
		v.Id = &ref
	}
	_, err := t.tx.ExecContext(ctx, "INSERT INTO user (id, username, auth) VALUES (?, ?, ?)", v.Id, v.UserName, v.Auth)
	return err
}

func (t *DatabaseTxn) DeleteUser(ctx context.Context, ref model.UserRef) error {
	_, err := t.tx.ExecContext(ctx, "DELETE FROM user WHERE id=?", ref)
	return err
}

func (t *DatabaseTxn) GetProblem(ctx context.Context, ref model.ProblemRef) (*model.Problem, error) {
	problem := new(model.Problem)
	err := t.tx.QueryRowContext(ctx, "SELECT id, title FROM problem WHERE id=?", ref).Scan(&problem.Id, &problem.Title)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return problem, nil
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
