package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	// for User struct

	ID     uint   `gorm:"primaryKey;autoIncrement"`
	Name   string `gorm:"size:25"`
	UserID uint   `gorm:"not null"`
	IsGood bool   `gorm:"default=true"`
	Price  int    `gorm:"default=0"`

	Type string `gorm:"size:25"`
	// type off and on
	// postgres does have enum type

	Stock       int    `gorm:"default=1"`
	Description string `gorm:"type:text;default=null"`

	User User
	// foreign key

	CreatedAt time.Time
	UpdatedAt time.Time
}
