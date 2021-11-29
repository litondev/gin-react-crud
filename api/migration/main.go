package main

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
	"github.com/litondev/gin-react-crud/api/config"
	"github.com/litondev/gin-react-crud/api/models"
)

func main() {
	errEnv := godotenv.Load()

	if errEnv != nil {
		fmt.Println(errEnv)
		os.Exit(1)
	}

	var dsn map[string]string = map[string]string{
		"DB_HOST" : os.Getenv("DB_HOST"),
		"DB_PORT" : os.Getenv("DB_PORT"),
		"DB_NAME" : os.Getenv("DB_NAME"),
		"DB_USER" : os.Getenv("DB_USER"),
		"DB_PASSWORD" : os.Getenv("DB_PASSWORD"),
	}

	database, errDb := config.Database(dsn)

	if errDb == nil {
		tableUser := database.Migrator().HasTable(&models.User{})
		if tableUser == true {
			database.Migrator().DropTable(&models.User{})
		}

		tableProduct := database.Migrator().HasTable(&models.Product{})
		if tableProduct == true {
			database.Migrator().DropTable(&models.Product{})
		}

		tableData := database.Migrator().HasTable(&models.Data{})
		if tableData == true {
			database.Migrator().DropTable(&models.Data{})
		}

		database.AutoMigrate(&models.User{}, &models.Product{}, &models.Data{})

		fmt.Println("Success Migrate")
	}
}
