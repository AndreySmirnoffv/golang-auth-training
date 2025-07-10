package db

import (
	"time"
)

type UserModel struct {
	ID        uint `gorm:"primaryKey"`
	Email     string
	Password  string
	CreatedAt time.Time
}

func (UserModel) TableName() string { return "users" }
