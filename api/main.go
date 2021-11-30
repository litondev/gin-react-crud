package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/litondev/gin-react-crud/api/config"
	"github.com/litondev/gin-react-crud/api/controllers"
	// "reflect"
	// "path/filepath"
	// "github.com/disintegration/imaging"
	// "github.com/litondev/gin-react-crud/api/models"
	// "github.com/litondev/gin-react-crud/api/requests"
	// "gorm.io/driver/postgres"
	// "gorm.io/gorm"
	// "gorm.io/driver/mysql"
)

func main() {
	// load env file
	errEnv := godotenv.Load()

	// check is not an error
	if errEnv != nil {
		fmt.Println(errEnv)
		os.Exit(1)
	}
	
	// connect to db
	var dsn map[string]string = map[string]string{
		"DB_HOST" : os.Getenv("DB_HOST"),
		"DB_PORT" : os.Getenv("DB_PORT"),
		"DB_NAME" : os.Getenv("DB_NAME"),
		"DB_USER" : os.Getenv("DB_USER"),
		"DB_PASSWORD" : os.Getenv("DB_PASSWORD"),
	}
	
	// get database
	db,errDb :=  config.Database(dsn)

	// check database error or not
	if errDb != nil {
		fmt.Println(errDb)
		os.Exit(1)
	}

	// get debug mode
	debug := os.Getenv("APP_DEBUG")
	appDebug, errDebug := strconv.ParseBool(debug)

	// check is not an error
	if errDebug != nil {
		fmt.Println(errDebug)
		os.Exit(1)
	}

	// set debug mode
	if appDebug == true {
		gin.SetMode("debug")
	} else {
		gin.SetMode("release")
	}

	// set jwt
	jwtAuthMiddleware, errJwtAuthMiddleware := jwt.New(&jwt.GinJWTMiddleware{
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

	if errJwtAuthMiddleware != nil {
		fmt.Println(errJwtAuthMiddleware)
		os.Exit(1)
	}

	errInitJwtAuthMiddleware := jwtAuthMiddleware.MiddlewareInit()

	if errInitJwtAuthMiddleware != nil {
		fmt.Println(errInitJwtAuthMiddleware.Error())
		os.Exit(1)
	}

	// intial gin
	r := gin.Default()
	
	// set static assets
	r.Static("/assets", "./assets")
	// try access
	// http://localhost:8000/assets/images/logo.png
	
	// insert db to context
	r.Use(func(c *gin.Context){
		c.Set("DB",db)
		c.Set("DEBUG",debug)
		c.Set("MAIL_HOST",os.Getenv("MAIL_HOST"))		
		c.Set("MAIL_USERNAME",os.Getenv("MAIL_USERNAME"))
		c.Set("MAIL_PASSWORD",os.Getenv("MAIL_PASSWORD"))
		c.Set("MAIL_PORT",os.Getenv("MAIL_PORT"))
		c.Set("MAIL_FROM",os.Getenv("MAIL_FROM"))
		c.Set("FRONTEND_URL",os.Getenv("FRONTEND_URL"))
		c.Next()
	})

	// When Unexpected Error Happend
	r.Use(func(c *gin.Context) {
		defer func(c *gin.Context) {
			if rec := recover(); rec != nil {
				logFile,logFileError := os.OpenFile(os.Getenv("APP_LOGGER_LOCATION"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

				if(logFileError != nil){
					c.JSON(500,gin.H{
						"message" : "Terjadi Kesalahan",
					})
					return 
				}

				logger := log.New(logFile, "Error : ", log.LstdFlags)
	
				logger.Println(time.Now().String())

				logger.Println(rec)

				fmt.Println(rec)
				
				var message string = "Terjadi Kesalahan";

				if(appDebug == true){
					message = rec.(string)
				}			

				c.JSON(500, gin.H{
					"message": message,
				})
				return
			}
		}(c)
		c.Next()
	})

	// Cors Middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodHead, http.MethodDelete, http.MethodOptions,http.MethodPut},
		AllowHeaders:     []string{"Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"message": "Page not found",
		})
		return 
	})

	/* Routes */
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
		return
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"version": "v1",
				"message": "Oke",
			})
			return
		})
	
		v1Auth := v1.Group("/auth")
		{
			v1Auth.POST("/signin", jwtAuthMiddleware.LoginHandler)
			v1Auth.POST("/signup", controllers.Signup)
			v1Auth.POST("/forgot-password", controllers.ForgotPassword)
			v1Auth.POST("/reset-password",controllers.ResetPassword)
		}
		
		v1.Use(jwtAuthMiddleware.MiddlewareFunc())
		{
			v1.GET("/data",controllers.IndexData)
			v1.POST("/data",controllers.StoreData)
			v1.GET("/data/:id",controllers.ShowData)
			v1.DELETE("/data/:id",controllers.DestoryData)
			v1.PUT("/data/:id",controllers.UpdateData)
			v1.GET("/data/export/pdf",controllers.ExportPdfData)
			v1.GET("/data/export/excel",controllers.ExportExcelData)
			v1.POST("/data/import",controllers.ImportExcelData)

			v1.POST("/refresh-token", jwtAuthMiddleware.RefreshHandler)

			v1.POST("/logout", jwtAuthMiddleware.LogoutHandler)
			v1.GET("/me", controllers.Me)

			v1.PUT("/profil/update",controllers.UpdateProfilData)
			v1.POST("/profil/upload",controllers.UpdateProfilPhoto)

			/* Product */		
			
			/* User */
			
			// CONTOH UPLOAD IMAGE
			// v1.POST("/profil/upload", func(c *gin.Context) {
			// 	file, err := c.FormFile("file")

			// 	if err != nil {
			// 		c.String(200, fmt.Sprintf("get form err: %s", err.Error()))
			// 		return
			// 	}

			// 	// image/jpeg image/jpg png
			// 	// fmt.Println(file.Header["Content-Type"][0])

			// 	filename := filepath.Base("")+"/assets/"+file.Filename

			// 	if err := c.SaveUploadedFile(file, filename); err != nil {				
			// 		c.String(200, fmt.Sprintf("upload file err: %s", err.Error()))
			// 		return
			// 	}
		
			// 	hello, errP := imaging.Open(filename)
			// 	fmt.Println(errP)
			// 	hell := imaging.Resize(hello, 128, 128, imaging.Lanczos)
			// 	errs := imaging.Save(hell, filepath.Base("")+"/assets/out.png")

			// 	fmt.Println(errs)

			// 	e := os.Remove(filename)
			// 	if e != nil {
			// 			log.Fatal(e)
			// 	}					

			// 	c.String(200, fmt.Sprintf("File %s uploaded successfully", file.Filename))
			// })
		}
	}
	/* Routes */

	// running server
	r.Run(os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT"))	
}