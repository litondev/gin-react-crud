package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/litondev/gin-react-crud/api/config"
	"github.com/litondev/gin-react-crud/api/helpers"
	"github.com/litondev/gin-react-crud/api/models"
)

func main() {
	database, err := config.Database()

	if err == nil {
		for i := 0; i < 3; i++ {
			hash, _ := helpers.HashPassword("password")
			var id string = "user" + strconv.Itoa(i)

			user := models.User{
				Name:     id,
				Email:    id + "@gmail.com",
				Password: hash,
			}

			result := database.Create(&user)

			if result.Error == nil {
				fmt.Println("Success Memasukan " + strconv.Itoa(int(result.RowsAffected)) + " Data")
			} else {
				fmt.Print(result.Error)
				os.Exit(1)
			}
		}

		for i := 0; i < 3; i++ {
			var id string = "data" + strconv.Itoa(i)

			data := models.Data{
				Name:  id,
				Phone: nil,
			}

			result := database.Create(&data)

			if result.Error == nil {
				fmt.Println("Success Memasukan " + strconv.Itoa(int(result.RowsAffected)) + " Data")
			} else {
				fmt.Print(result.Error)
				os.Exit(1)
			}
		}

		for i := 0; i < 3; i++ {
			hash, _ := helpers.HashPassword("password")
			var id string = "user-p-" + strconv.Itoa(i)

			user := models.User{
				Name:     id,
				Email:    id + "@gmail.com",
				Password: hash,

				// Product:  models.Product{Name: "product"},
				// insert has One

				Product: []models.Product{{Name: "product 1"}, {Name: "product 2"}},
				// insert has Many
			}

			result := database.Create(&user)

			if result.Error == nil {
				fmt.Println("Success Memasukan " + strconv.Itoa(int(result.RowsAffected)) + " Data")
			} else {
				fmt.Print(result.Error)
				os.Exit(1)
			}
		}

		fmt.Println("Success Seed")
	}
}