package main

import (
	"fmt"

	"github.com/litondev/gin-react-crud/api/config"
	"github.com/litondev/gin-react-crud/api/models"
)

func main() {
	database, err := config.Database()

	if err == nil {
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
