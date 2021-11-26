package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/litondev/gin-react-crud/api/helpers"
	"github.com/litondev/gin-react-crud/api/requests"
	"github.com/litondev/gin-react-crud/api/models"
	"github.com/litondev/gin-react-crud/api/config"
	"fmt"
	// "encoding/json"
)

func Signin(c *gin.Context){	
	err := helpers.Validate(c,&requests.VSigninRequest)

	fmt.Println(err);

	if(err != nil){		
		c.JSON(422,gin.H{
			"message" : err.Error(),
		})
		return;
	}
	
	database, _ := config.Database()
	result := map[string]interface{}{}
	database.Model(&models.User{}).First(&result)
	fmt.Println(result)
	c.JSON(200,gin.H{
		"message" : "Signin",
		"db" : result,
	});
	return;
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

