package controllers

import (
	"crypto/tls"
	"errors"
	"fmt"
	"time"
	// "reflect"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/litondev/gin-react-crud/api/config"
	"github.com/litondev/gin-react-crud/api/helpers"
	"github.com/litondev/gin-react-crud/api/models"
	"github.com/litondev/gin-react-crud/api/requests"
	gomail "gopkg.in/mail.v2"
	// "net/http"
	// "fmt"
	// "encoding/json"
)

func Signin(c *gin.Context) (interface{}, error) {
	err := helpers.Validate(c, &requests.VSigninRequest)

	if err != nil {
		return nil, err
	}

	database, _ := config.Database()

	result := map[string]interface{}{}

	database.Model(&models.User{}).Where("email = ?", requests.VSigninRequest.Email).First(&result)

	if len(result) == 0 {
		return nil, errors.New("Email tidak ditemukan")
	}

	var isValidPassword bool = helpers.CheckPasswordHash(
		requests.VSigninRequest.Password,
		result["password"].(string),
	)

	if isValidPassword == false {
		return nil, errors.New("Password salah")
	}

	return &models.User{
		ID:    result["id"].(uint),
		Email: result["email"].(string),
	}, nil
}

func SigninResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(200, gin.H{
		"access_token": token,
		"expire":       expire.Format(time.RFC3339),
	})
}

func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"message": message,
	})
	return
}

func RefreshResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(200, gin.H{
		"access_token": token,
		"expire":       expire.Format(time.RFC3339),
	})
}

func Logout(c *gin.Context, code int) {
	c.JSON(200, gin.H{
		"message": "Success",
	})
}

func Me(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	c.JSON(200, gin.H{
		"Id": uint(claims["sub"].(float64)),
	})
}

func PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*models.User); ok {
		return jwt.MapClaims{
			"sub": v.ID,
		}
	}
	return jwt.MapClaims{}
}

func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)

	return &models.User{
		ID: uint(claims["sub"].(float64)),
	}
}

func Signup(c *gin.Context) {
	err := helpers.Validate(c, &requests.VSignupRequest)

	if err != nil {
		c.JSON(422, gin.H{
			"message": err.Error(),
		})
		return
	}

	database, _ := config.Database()

	hash, _ := helpers.HashPassword(requests.VSigninRequest.Password)

	user := &models.User{
		Name:     requests.VSignupRequest.Name,
		Password: hash,
		Email:    requests.VSignupRequest.Email,
	}

	result := database.Create(&user)

	if result.Error != nil {
		c.JSON(500, gin.H{
			"message": "Terjadi Kesalahan",
		})

		return
	}

	c.JSON(200, gin.H{
		"message": true,
	})
}

func ForgotPassword(c *gin.Context) {
	err := helpers.Validate(c, &requests.VForgotPasswordRequest)

	if err != nil {
		c.JSON(422, gin.H{
			"message": err.Error(),
		})
		return
	}

	database, _ := config.Database()

	result := map[string]interface{}{}

	database.Model(&models.User{}).Where("email = ?", requests.VForgotPasswordRequest.Email).Update("remember_token", "12345").First(&result)

	if len(result) == 0 {
		c.JSON(500, gin.H{
			"message": "Email tidak ditemukan",
		})
		return
	}

	sendEmail(&result)

	c.JSON(200, gin.H{
		"message": "Forgot Password",
	})
}

func sendEmail(result *map[string]interface{}) {

	realResult := *result
	token := *realResult["remember_token"].(*string)

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "from@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", "to@example.com")

	// Set E-Mail subject
	m.SetHeader("Subject", "Gomail test subject")

	// Set E-Mail body. You can set plain text or html with text/html
	// m.SetBody("text/plain", "This is Gomail test body")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>! <a href='http://localhost:3000/reset-password?token="+token+"'>Reset Password</a>")

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.mailtrap.io", 2525, "75999957ca7383", "95a016b68c6448")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return
}

func ResetPassword(c *gin.Context) {
	database, _ := config.Database()

	tx := database.Begin()

	err := helpers.Validate(c, &requests.VResetPasswordRequest)

	if err != nil {
		tx.Rollback()
		c.JSON(422, gin.H{
			"message": err.Error(),
		})
		return
	}	

	// fmt.Println(reflect.TypeOf(requests.VResetPasswordRequest.NewPassword))

	if requests.VResetPasswordRequest.NewPassword != requests.VResetPasswordRequest.PasswordConfirm {
		tx.Rollback()
		c.JSON(422, gin.H{
			"message": "Password tidak sama",
		})
		return
	}

	hash, _ := helpers.HashPassword(requests.VResetPasswordRequest.NewPassword)

	query := database.Model(&models.User{})
	// Select remember_token digunakan ketika update ke nil/null
	query = query.Select("email","password","remember_token")
	query = query.Where("email = ?", requests.VResetPasswordRequest.Email)
	query = query.Where("remember_token = ?",requests.VResetPasswordRequest.Token)

	result := map[string]interface{}{}

	query.First(&result)

	if len(result) == 0 {
		tx.Rollback()
		c.JSON(500, gin.H{
			"message": "Data anda tidak valid",
		})
		return
	}

	query.Updates(&models.User{
		Password : hash,
		RememberToken : nil,
	})

	tx.Commit()
	c.JSON(200, gin.H{
		"message": true,
	})
}
