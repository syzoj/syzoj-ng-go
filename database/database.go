package database

import (
    "database/sql"
    "context"
    "sync/atomic"
    "race"
    "runtime/debug"
    
    "github.com/sirupsen/logrus"

    "github.com/syzoj/syzoj-ng-go/model"
)

type Database struct {
    db *sql.DB
}

type DatabaseTxn struct {
    tx *sql.Tx
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
        ReadOnly: false,
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
        ReadOnly: true,
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

func (t *DatabaseTxn) SetUser(ctx context.Context, ref model.UserRef, v *model.User) error {
    if v.GetId() != ref {
        panic("ref and v does not match")
    }
    _, err := t.tx.ExecContext(ctx, "UPDATE user SET username=?, auth=? WHERE id=?", v.UserName, v.Auth, v.Id)
    return err
}
