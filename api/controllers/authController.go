package controllers

import (
	"errors"
	"time"	
	// "crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"	
	"github.com/litondev/gin-react-crud/api/helpers"
	"github.com/litondev/gin-react-crud/api/models"
	"github.com/litondev/gin-react-crud/api/requests"
	jwt "github.com/appleboy/gin-jwt/v2"
	// gomail "gopkg.in/mail.v2"
)

func Signin(c *gin.Context) (interface{}, error) {
	errValidation := helpers.Validate(c, &requests.VSigninRequest)

	if errValidation != nil {
		return nil, errValidation
	}	
	
	database := c.MustGet("DB").(*gorm.DB);	

	resultUser := map[string]interface{}{}

	queryUser := database.Model(&models.User{})
		queryUser.Select("id","password","email")
		queryUser.Where("email = ?", requests.VSigninRequest.Email)
		queryUser.First(&resultUser)

	if len(resultUser) == 0 {
		return nil, errors.New("Email tidak ditemukan")
	}
	
	var isValidPassword bool = helpers.CheckPasswordHash(
		requests.VSigninRequest.Password,
		resultUser["password"].(string),
	)

	if isValidPassword == false {
		return nil, errors.New("Password salah")
	}

	return &models.User{
		ID:    resultUser["id"].(uint),
		Email: resultUser["email"].(string),
	}, nil
}

func SigninResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(200, gin.H{
		"access_token": token,
		"expire":       expire.Format(time.RFC3339),
	})
	return 
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
		"expire": expire.Format(time.RFC3339),
	})
	return 
}

func Logout(c *gin.Context, code int) {
	c.JSON(200, gin.H{
		"message": true,
	})
	return
}

func Me(c *gin.Context) {
	claims := jwt.ExtractClaims(c)

	c.JSON(200, gin.H{
		"id": uint(claims["sub"].(float64)),
	})
	return 
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
	errValidation := helpers.Validate(c, &requests.VSignupRequest)

	if errValidation != nil {
		c.JSON(422, gin.H{
			"message": errValidation.Error(),
		})
		return
	}

	database := c.MustGet("DB").(*gorm.DB);	

	tx := database.Begin()

	hash,errHash := helpers.HashPassword(requests.VSigninRequest.Password)

	if(errHash != nil){
		tx.Rollback()
		fmt.Println(errHash.Error())
		c.JSON(500,gin.H{
			"message": "Terjadi Kesalahan",
		})
		return 
	}

	user := &models.User{
		Name		: requests.VSignupRequest.Name,
		Email		: requests.VSignupRequest.Email,
		Password	: hash,
	}

	resultUser := database.Create(&user)

	if resultUser.Error != nil {
		tx.Rollback()
		fmt.Println(resultUser.Error)
		c.JSON(500, gin.H{
			"message": "Terjadi Kesalahan",
		})
		return
	}

	tx.Commit()
	c.JSON(200, gin.H{
		"message": true,
	})
	return 
}

func ForgotPassword(c *gin.Context) {
	errValidation := helpers.Validate(c, &requests.VForgotPasswordRequest)

	if errValidation != nil {
		c.JSON(422, gin.H{
			"message": errValidation.Error(),
		})
		return
	}

	database := c.MustGet("DB").(*gorm.DB);	

	tx := database.Begin()

	resultUser := map[string]interface{}{}

	queryUser := database.Model(&models.User{})
		queryUser.Where("email = ?", requests.VForgotPasswordRequest.Email)
		queryUser.First(&resultUser)

	if len(resultUser) == 0 {
		tx.Rollback()		
		c.JSON(500, gin.H{
			"message": "Email tidak ditemukan",
		})
		return
	}
	
	queryUser.Update("remember_token", "12345")

	if queryUser.Error != nil {
		tx.Rollback()
		fmt.Println(queryUser.Error)
		c.JSON(500, gin.H{
			"message": "Terjadi Kesalahan",
		})
		return
	}
	
	// Send Email
	// sendEmail(&resultUser)

	tx.Commit()
	c.JSON(200, gin.H{
		"message": true,
	})
	return 
}

// func sendEmail(result *map[string]interface{}) {
// 	realResult := *result
// 	token := *realResult["remember_token"].(*string)

// 	m := gomail.NewMessage()

// 	// Set E-Mail sender
// 	m.SetHeader("From", "from@gmail.com")

// 	// Set E-Mail receivers
// 	m.SetHeader("To", "to@example.com")

// 	// Set E-Mail subject
// 	m.SetHeader("Subject", "Gomail test subject")

// 	// Set E-Mail body. You can set plain text or html with text/html
// 	// m.SetBody("text/plain", "This is Gomail test body")
// 	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>! <a href='http://localhost:3000/reset-password?token="+token+"'>Reset Password</a>")

// 	// Settings for SMTP server
// 	d := gomail.NewDialer("smtp.mailtrap.io", 2525, "75999957ca7383", "95a016b68c6448")

// 	// This is only needed when SSL/TLS certificate is not valid on server.
// 	// In production this should be set to false.
// 	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

// 	// Now send E-Mail
// 	if err := d.DialAndSend(m); err != nil {
// 		fmt.Println(err)
// 		panic(err)
// 	}

// 	return
// }

func ResetPassword(c *gin.Context) {	
	errValidation := helpers.Validate(c, &requests.VResetPasswordRequest)

	if errValidation != nil {
		c.JSON(422, gin.H{
			"message": errValidation.Error(),
		})
		return
	}	

	database := c.MustGet("DB").(*gorm.DB);	

	tx := database.Begin()

	if requests.VResetPasswordRequest.NewPassword != requests.VResetPasswordRequest.PasswordConfirm {
		tx.Rollback()		
		c.JSON(422, gin.H{
			"message": "Password tidak sama",
		})
		return
	}

	hash, errHash := helpers.HashPassword(requests.VResetPasswordRequest.NewPassword)

	if(errHash != nil){
		tx.Rollback()		
		fmt.Println(errHash.Error())
		c.JSON(500,gin.H{
			"message" : "Terjadi Kesalahan",
		})
		return 
	}

	resultUser := map[string]interface{}{}

	query := database.Model(&models.User{})
		// Select remember_token digunakan ketika update ke nil/null
		query.Select("remember_token","password")
		query.Where("email = ?", requests.VResetPasswordRequest.Email)
		query.Where("remember_token = ?",requests.VResetPasswordRequest.Token)
		query.First(&resultUser)

	if len(resultUser) == 0 {
		tx.Rollback()
		c.JSON(500, gin.H{
			"message": "Data anda tidak valid",
		})
		return
	}

	resultQueryUpdateUser := query.Updates(&models.User{
		Password : hash,
		RememberToken : nil,
	})

	if resultQueryUpdateUser.Error != nil {
		tx.Rollback()
		fmt.Println(resultQueryUpdateUser.Error)
		c.JSON(500, gin.H{
			"message": "Terjadi Kesalahan",
		})
		return
	}

	tx.Commit()
	c.JSON(200, gin.H{
		"message": true,
	})
}
