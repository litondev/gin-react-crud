package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	// digunakan ketika membuat relasi

	ID            uint    `gorm:"primaryKey;autoIncrement"`
	Name          string  `gorm:"size:25"`
	Email         string  `gorm:"unique;size:25"`
	Password      string  `gorm:"type:text"`
	RememberToken *string `gorm:"type:text"`

	Photo *string `gorm:"size:25"`
	// *string membuat photo boleh null dan jika diakses gorm maka nilainya akan nil
	// dan tidak akan bisa di masukan zero value "",0 atau false

	Role string `gorm:"size:25;default:'user'"`
	// admin,user

	// Product   Product
	// membuat relasi has One ke table product

	Product []Product
	// membuat relasi has Many ke table product

	CreatedAt time.Time
	UpdatedAt time.Time
}
