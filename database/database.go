package database

import (
	"context"
	"database/sql"
	"sync/atomic"

	"github.com/sirupsen/logrus"
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
