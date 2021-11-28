package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	// "reflect"
	// "path/filepath"
	// "github.com/disintegration/imaging"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/litondev/gin-react-crud/api/config"
	"github.com/litondev/gin-react-crud/api/controllers"
	// "github.com/litondev/gin-react-crud/api/models"
	// "github.com/litondev/gin-react-crud/api/requests"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "gorm.io/driver/mysql"
)

var DB *gorm.DB

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

	// Mysql
	// dsn := os.Get("DB_USER") + ":" + 
	// 	os.Get("DB_PASSWORD") + "@tcp(" + 
	// 	os.Get("DB_HOST") + ":" + 
	// 	os.Get("DB_PORT")+")/" + 
	// 	os.Get("DB_NAME") + "|?charset=utf8mb4&parseTime=True&loc=Local"

	// config.DB, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// Postgres
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable TimeZone=Asia/Jakarta"

	config.DB, _ = gorm.Open(postgres.New(postgres.Config{
		PreferSimpleProtocol: true, 
		DSN : dsn,
  	}))

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
	
	// USING DB IN MIDDLEWARE
	/* 
		DB,_ :=  config.Database()

		r.Use(func(c *gin.Context){
			c.Set("DB",DB)
			c.Next()
		})

		// di controller
		// database := c.MustGet("DB").(*gorm.DB);
	*/

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
			v1Auth.POST("/reset-password",controllers.ResetPassword)
		}
		
		v1.Use(authMiddleware.MiddlewareFunc())
		{
			v1.GET("/data",controllers.IndexData)
			v1.POST("/data",controllers.StoreData)
			v1.GET("/data/:id",controllers.ShowData)
			v1.DELETE("/data/:id",controllers.DestoryData)
			v1.PUT("/data/:id",controllers.UpdateData)

			v1.POST("/refresh-token", authMiddleware.RefreshHandler)

			v1.POST("/logout", authMiddleware.LogoutHandler)
			v1.GET("/me", controllers.Me)

			v1.PUT("/profil/update",controllers.UpdateProfilData)
			v1.POST("/profil/upload",controllers.UpdateProfilPhoto)
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
