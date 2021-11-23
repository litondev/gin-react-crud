package models

import "time"

type Data struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"size:25"`

	Phone *string `gorm:"size:25"`
	// *string membuat phone boleh null dan jika diakses gorm maka nilainya akan nil
	// dan tidak akan bisa di masukan zero value "",0 atau false

	CreatedAt time.Time
	UpdatedAt time.Time
}
