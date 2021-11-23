package model

import "time"

type Data struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:25"`
	Phone     string `gorm:"size:25;default=null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
