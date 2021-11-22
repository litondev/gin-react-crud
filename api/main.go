package main

import (
	"fmt"

	"github.com/litondev/gin-react-crud/api/db"
)

func main() {
	database := db.Database()

	if database == nil {
		fmt.Println("Database Not Connected")
	}

	fmt.Println("Database Connect")
}

// func main() {
// 	r := gin.Default()

// 	r.Use(cors.New(cors.Config{
// 		AllowOrigins:     []string{"http://localhost:3000"},
// 		AllowMethods:     []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodHead, http.MethodDelete, http.MethodOptions},
// 		AllowHeaders:     []string{"Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization"},
// 		ExposeHeaders:    []string{"Content-Length"},
// 		AllowCredentials: true,
// 	}))

// 	r.GET("/", func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"message": "Hellos",
// 		})
// 	})

// 	r.GET("/ping", func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"message": "pongs",
// 		})
// 	})

// 	r.GET("/api/v1/p", func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"message": "Hellos",
// 		})
// 	})

// 	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
// }
