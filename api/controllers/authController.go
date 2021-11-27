package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/litondev/gin-react-crud/api/helpers"
	"github.com/litondev/gin-react-crud/api/requests"
	"github.com/litondev/gin-react-crud/api/models"
	"github.com/litondev/gin-react-crud/api/config"
	jwt "github.com/appleboy/gin-jwt/v2"	
	"errors"
	"time"
	// "net/http"
	// "fmt"
	// "reflect"
	// "encoding/json"
)

func Signin(c *gin.Context) (interface{},error) {	
	err := helpers.Validate(c,&requests.VSigninRequest)

	if(err != nil){			
		return nil,err;
	}
	
	database, _ := config.Database()

	result := map[string]interface{}{}

	database.Model(&models.User{}).Where("email = ?",requests.VSigninRequest.Email).First(&result)

	if len(result) == 0 {		
		return nil,errors.New("Email tidak ditemukan");
	}

	var isValidPassword bool = helpers.CheckPasswordHash(
		requests.VSigninRequest.Password,
		result["password"].(string),
	)
	
	if(isValidPassword == false) {		
		return nil,errors.New("Password salah");
	}

	return &models.User{
		ID : result["id"].(uint),
		Email : result["email"].(string),
 	},nil
}

func SigninResponse (c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(200, gin.H{		
		"access_token":  token,
		"expire": expire.Format(time.RFC3339),
	})
}

func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"message": message,
	})
	return;
}

func RefreshResponse (c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(200, gin.H{		
		"access_token":  token,
		"expire": expire.Format(time.RFC3339),
	})
}

func Logout (c *gin.Context, code int) {
	c.JSON(200, gin.H{
		"message": "Success",
	})
}

func Me(c *gin.Context){
	claims := jwt.ExtractClaims(c)
	c.JSON(200,gin.H{
		"Id" : uint(claims["sub"].(float64)),	
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
		ID : uint(claims["sub"].(float64)),
	}
}

func Signup(c *gin.Context){
	c.JSON(200,gin.H{
		"message" : "Signup",
	});
}

func ForgotPassword(c *gin.Context){
	c.JSON(200,gin.H{
		"message" : "ForgotPassword",
	});
}
