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

	Payments []PaymentModel `gorm:"foreignKey:UserID"`
}

func (UserModel) TableName() string { return "users" }

type PaymentModel struct {
	ID     uint64
	Amount string
	IsPaid bool
	UserID uint64
	User   UserModel
}

func (PaymentModel) TableName() string { return "payments" }
