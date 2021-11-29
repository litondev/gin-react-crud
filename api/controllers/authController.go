package controllers

import (
	"errors"
	"time"	
	"strconv"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"	
	"github.com/litondev/gin-react-crud/api/helpers"
	"github.com/litondev/gin-react-crud/api/models"
	"github.com/litondev/gin-react-crud/api/requests"
	jwt "github.com/appleboy/gin-jwt/v2"
	gomail "gopkg.in/mail.v2"
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

	hash,errHash := helpers.HashPassword(requests.VSignupRequest.Password)

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
	errSendMail := sendEmail(resultUser,c)

	if(errSendMail != nil){
		fmt.Println(errSendMail)
		c.JSON(500,gin.H{
			"message" : "Terjadi Kesalahan",
		})
		return 
	}

	tx.Commit()
	c.JSON(200, gin.H{
		"message": true,
	})
	return 
}

func sendEmail(result map[string]interface{},c *gin.Context) error {
	frontend_url := c.MustGet("FRONTEND_URL").(string)
	token := *result["remember_token"].(*string)

	
	mail := gomail.NewMessage()
	
	mail.SetHeader("From",c.MustGet("MAIL_FROM").(string))
	mail.SetHeader("To", result["email"].(string))
	mail.SetHeader("Subject", "Forgot Password")	
	mail.SetBody("text/html", "<a href='" + frontend_url + "/reset-password?token=" + token + "'>Reset Password</a>")
	
	
	mailPort,errMailPort := strconv.Atoi(c.MustGet("MAIL_PORT").(string))

	if errMailPort != nil {
		return errors.New("Port Email Tidak Valid")
	}

	sendMail := gomail.NewDialer(
		c.MustGet("MAIL_HOST").(string),
		mailPort, 
		c.MustGet("MAIL_USERNAME").(string), 
		c.MustGet("MAIL_PASSWORD").(string),
	)

	sendMail.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	if errSendMail := sendMail.DialAndSend(mail); errSendMail != nil {
		return errors.New("Gagal Mengirim Email")
	}

	return nil
}

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
