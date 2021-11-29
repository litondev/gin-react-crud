package config

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Database(env map[string]string) (*gorm.DB, error) {
	// Mysql
		// dsn := os.Get("DB_USER") + ":" + 
		// 	os.Get("DB_PASSWORD") + "@tcp(" + 
		// 	os.Get("DB_HOST") + ":" + 
		// 	os.Get("DB_PORT")+")/" + 
		// 	os.Get("DB_NAME") + "|?charset=utf8mb4&parseTime=True&loc=Local"
		// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// Postgres
		// dsn := "host=" + os.Getenv("DB_HOST") +
		// " user=" + os.Getenv("DB_USER") +
		// " password=" + os.Getenv("DB_PASSWORD") +
		// " dbname=" + os.Getenv("DB_NAME") +
		// " port=" + os.Getenv("DB_PORT") +
		// " sslmode=disable TimeZone=Asia/Jakarta"

	dsn := "host=" + env["DB_HOST"] +
		" user=" + env["DB_USER"] + 
		" password=" + env["DB_PASSWORD"] + 
		" dbname=" + env["DB_NAME"] + 
		" port=" + env["DB_PORT"] +
		" sslmode=disable TimeZone=Asia/Jakarta"

	db,errDb := gorm.Open(postgres.New(postgres.Config{
		PreferSimpleProtocol: true, 
		DSN : dsn,
	}))

	if errDb != nil {
		// Print Error
		fmt.Println(errDb)

		// Throw Error
		return nil, errors.New("Can't Connect To Database")
	}

	// Return Gorm Orm
	return db, nil
}