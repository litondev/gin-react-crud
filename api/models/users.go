package model

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:25"`
	Email     string `gorm:"unique;size:25"`
	Password  string `gorm:"type:text"`
	Photo     string `gorm:"size:25"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
