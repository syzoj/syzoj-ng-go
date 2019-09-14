package model

import (
	"time"
)

type User struct {
	Id           int64     `db:"id"`
	Uid          string    `db:"uid"`
	UserName     string    `db:"username"`
	Email        string    `db:"email"`
	Password     []byte    `db:"password"`
	RegisterTime time.Time `db:"register_time"`
	ProblemCount int64     `db:"problem_count"`
}
