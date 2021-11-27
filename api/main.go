package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	// "reflect"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/litondev/gin-react-crud/api/controllers"
	// "github.com/litondev/gin-react-crud/api/models"
	// "github.com/litondev/gin-react-crud/api/requests"
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
	/*
	 fmt.Println(reflect.TypeOf(test))
	 fmt.Println(reflect.TypeOf(appDebug))
	*/

	// set debug mode
	if appDebug == true {
		gin.SetMode("debug")
	} else {
		gin.SetMode("release")
	}

	// log console to file
	/*
		f, _ := os.Create("gin.log")
		gin.DefaultWriter = io.MultiWriter(f)
	*/

	// set jwt
	authMiddleware, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:         "Dev",
		Key:           []byte("secret"),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,

		// IdentityKey
		IdentityKey: "sub",

		// Signin
		Authenticator: controllers.Signin,

		// Signin Response
		LoginResponse: controllers.SigninResponse,

		// Logout Response
		LogoutResponse: controllers.Logout,

		// Unauthorized
		Unauthorized: controllers.Unauthorized,

		// Custome Payload
		PayloadFunc: controllers.PayloadFunc,
		// Custome Identity
		IdentityHandler: controllers.IdentityHandler,

		// Refresh Token Response
		RefreshResponse: controllers.RefreshResponse,
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		fmt.Println(errInit.Error())
		os.Exit(1)
	}

	// intial gin
	r := gin.Default()

	// set static assets
	r.Static("/assets", "./assets")
	// try access
	// http://localhost:8000/assets/images/logo.png

	// Cors Middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodHead, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// When Unexpected Error Happend
	r.Use(globalRecover)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"message": "Page not found",
		})
	})

	/* Routes */
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"version": "v1",
				"message": "Oke",
			})
		})

		v1Auth := v1.Group("/auth")
		{
			v1Auth.POST("/signin", authMiddleware.LoginHandler)
			v1Auth.POST("/signup", controllers.Signup)
			v1Auth.POST("/forgot-password", controllers.ForgotPassword)
		}

		v1.Use(authMiddleware.MiddlewareFunc())
		{
			v1.POST("/logout", authMiddleware.LogoutHandler)
			v1.POST("/refresh-token", authMiddleware.RefreshHandler)
			v1.GET("/me", controllers.Me)
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
			fmt.Println(rec)

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Terjadi Kesalahan",
			})
		}
	}(c)
	c.Next()
}
