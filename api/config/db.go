package config

// Untuk Mengeksekusi Package
/*
 package main
*/

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Database() (*gorm.DB, error) {
	// Load Dot Env
	err := godotenv.Load()

	// Check If Not Err
	if err != nil {
		// Print Error
		fmt.Println(err)
		// Throw Error
		return nil, errors.New("Error loading .env file")
		// Exit Program
		/*
		 os.Exit(1)
		*/
	}

	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		// Print Error
		fmt.Println(err)
		// Throw Error
		return nil, errors.New("Can't Connect To Database")
		// Exit Program
		/* 
			os.Exit(1)
		*/
	}

	// Return Gorm Orm
	return db, nil
}

// Check Db Koneksi
/*
func main() {
	_, err := Database()

	if err != nil {
		fmt.Println("You are can't connect to database")
	} else {
		fmt.Println("You are connect to database")
	}
}
*/