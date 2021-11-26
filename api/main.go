package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/litondev/gin-react-crud/api/controllers"
)

func main() {
	// load env file
	err := godotenv.Load()

	// check is not an error
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// get debug mode
	debug := os.Getenv("APP_DEBUG")
	appDebug, err := strconv.ParseBool(debug)

	// check is not an error
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// check type data
	// fmt.Println(reflect.TypeOf(test))
	// fmt.Println(reflect.TypeOf(appDebug))

	// set debug mode
	if appDebug == true {
		gin.SetMode("debug")
	} else {
		gin.SetMode("release")
	}

	// log console to file
	// f, _ := os.Create("gin.log")
	// gin.DefaultWriter = io.MultiWriter(f)

	// intial gin
	r := gin.Default()

	// set static assets
	r.Static("/assets", "./assets")
	// try access
	// http://localhost:8000/assets/images/logo.png

	// Cors Middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodHead, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// When Unexpected Error Happend
	r.Use(globalRecover)

	/* Routes */
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})	

	v1 := r.Group("/api/v1")
	{
		v1.GET("/status",func(c *gin.Context){
			c.JSON(200,gin.H{
				"version" : "v1",
				"message" : "Oke",
			})
		})

		v1Auth := v1.Group("/auth")
		{
			v1Auth.GET("/signin",controllers.Signin)
			v1Auth.GET("/signup",controllers.Signup)
			v1Auth.GET("/forgot-password",controllers.ForgotPassword)
		}
	}
	/* Routes */

	// running server
	r.Run(os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT"))
}

// When Unexpected Error Happend
func globalRecover(c *gin.Context) {
	defer func(c *gin.Context) {
		if rec := recover(); rec != nil {
			f, _ := os.OpenFile(os.Getenv("APP_LOGGER_LOCATION"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

			logger := log.New(f, "Error : ", log.LstdFlags)

			logger.Println(time.Now().String())
			logger.Println(rec)

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Terjadi Kesalahan",
			})
		}
	}(c)
	c.Next()
}
