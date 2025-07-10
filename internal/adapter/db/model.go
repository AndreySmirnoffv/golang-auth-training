package db

import (
	"time"
)

type UserModel struct {
	ID        uint `gorm:"primaryKey"`
	Email     string
	Password  string
	CreatedAt time.Time
	Payments  []PaymentModel
}

func (UserModel) TableName() string { return "users" }

type PaymentModel struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Amount    float64
	CreatedAt time.Time
}

func (PaymentModel) TableName() string { return "payments" }
