package db

import (
	"time"
)

type UserModel struct {
	ID        uint64
	Email     string
	Password  string
	Balance   int64
	CreatedAt time.Time
}

func (UserModel) TableName() string { return "users" }
