package models

import (
	"time"
)

type Product struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"size:25"`

	UserID uint `gorm:"not null"`
	// relasi dari tabel user

	IsGood bool `gorm:"default=true"`
	Price  int  `gorm:"default=0"`

	Type string `gorm:"size:25;default='on'"`
	// type off and on
	// postgres does have enum type

	Stock       int     `gorm:"default=1"`
	Description *string `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
