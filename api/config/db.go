package config

// package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Database() (*gorm.DB, error) {
	err := godotenv.Load()

	if err != nil {
		// fmt.Println("Error loading .env file")
		fmt.Println(err)
		return nil, errors.New("Error loading .env file")
		// os.Exit(1)
	}

	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		// fmt.Println("Can't Connect To Database")
		fmt.Println(err)
		return nil, errors.New("Can't Connect To Database")
		os.Exit(1)
	}

	return db, nil
}

// func main() {
// 	_, err := Database()

// 	if err != nil {
// 		fmt.Println("You are can't connect to database")
// 	} else {
// 		fmt.Println("You are connect to database")
// 	}
// }
