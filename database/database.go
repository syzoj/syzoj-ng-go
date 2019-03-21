package database

import (
    "database/sql"
    "context"

    "github.com/syzoj/syzoj-ng-go/model"
)

type Database struct {
    db *sql.DB
}

type DatabaseTxn struct {
    tx *sql.Tx
}

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
    return &DatabaseTxn{tx: tx}, nil
}

func (d *Database) OpenReadonlyTxn(ctx context.Context) (*DatabaseTxn, error) {
    tx, err := d.db.BeginTx(ctx, &sql.TxOptions{
        Isolation: sql.LevelSerializable,
        ReadOnly: true,
    })
    if err != nil {
        return nil, err
    }
    return &DatabaseTxn{tx: tx}, nil
}

func (t *DatabaseTxn) Commit() error {
    return t.tx.Commit()
}

func (t *DatabaseTxn) Rollback() error {
    return t.tx.Rollback()
}

func (t *DatabaseTxn) GetUser(ctx context.Context, ref model.UserRef) (*model.User, error) {
    user := new(model.User)
    err := t.tx.QueryRowContext(ctx, "SELECT id, username FROM user WHERE id=?", ref).Scan(&user.Id, &user.UserName)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return user, nil
}
