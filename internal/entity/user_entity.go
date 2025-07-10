package entity

import "time"

type User struct {
	ID        uint `gorm:"primaryKey;autoincrement"`
	Email     string
	Password  string
	Balance   uint `gorm:"default:0"`
	CreatedAt time.Time
}
