package controllers

import "github.com/gin-gonic/gin"
import "github.com/litondev/gin-react-crud/api/config"
import "github.com/litondev/gin-react-crud/api/models"
import "github.com/litondev/gin-react-crud/api/helpers"
import "github.com/litondev/gin-react-crud/api/requests"
// import jwt "github.com/appleboy/gin-jwt/v2"
import "strconv"
// import "fmt"
import "strings"
import "html"

func IndexData(c *gin.Context){
	// claims := jwt.ExtractClaims(c)
	// var ID uint = uint(claims["sub"].(float64))
    // var stringID string = strconv.FormatUint(uint64(ID),10)

	page := c.DefaultQuery("page", "1")
	new_page,_ := strconv.Atoi(page)

	per_page := c.DefaultQuery("per_page", "10")
	new_per_page,_ := strconv.Atoi(per_page)

	search := c.DefaultQuery("search","")

	result := []map[string]interface{}{}

	var resultCount int64;
	config.DB.Model(&models.Data{}).Select("id").Count(&resultCount)
	new_page = (int(resultCount) - ((new_page * new_per_page) - new_per_page))

	query := config.DB.Model(&models.Data{})
		query.Select("name","id","phone")
		if search != "" {
			query.Where("name LIKE ?", "%"+search+"%")		
		}
		query.Where("id <= ?",new_page)
		query.Order("id desc")
		query.Limit(new_per_page)
		query.Find(&result)

	c.JSON(200,gin.H{
		"data" : result,
		"per_page" : new_per_page,
	})
}

func StoreData(c *gin.Context){
	err := helpers.Validate(c, &requests.VDataRequest)

	if err != nil {
		c.JSON(422, gin.H{
			"message": err.Error(),
		})
		return
	}
		
	var phone string = html.EscapeString(strings.Trim(requests.VDataRequest.Phone," "))

	users := &models.Data{
		Name : html.EscapeString(strings.Trim(requests.VDataRequest.Name," ")),
		Phone : &phone,
	}

	config.DB.Model(&models.Data{}).Create(&users)

	c.JSON(200,gin.H{
		"message": true,
	})
}

func ShowData(c *gin.Context){
	id,_  := strconv.Atoi(c.Param("id"))

	resultData := map[string]interface{}{}

	query := config.DB.Model(&models.Data{})
		query.Select("id","name","phone")
		query.Where("id = ?",id)
		query.First(&resultData)

	c.JSON(200,gin.H{
		"message": true,
		"data" : resultData,
	})
}

func UpdateData(c *gin.Context){
	err := helpers.Validate(c, &requests.VDataRequest)

	if err != nil {
		c.JSON(422, gin.H{
			"message": err.Error(),
		})
		return
	}

	id,_  := strconv.Atoi(c.Param("id"))

	var phone string = html.EscapeString(strings.Trim(requests.VDataRequest.Phone," "))

	query := config.DB.Model(&models.Data{})
		query.Select("name","phone")
		query.Where("id = ?",id)
		query.Updates(&models.Data{
			Name : html.EscapeString(strings.Trim(requests.VDataRequest.Name," ")),
			Phone : &phone,
		})

	c.JSON(200,gin.H{
		"message": true,
	})
}

func DestoryData(c *gin.Context){
	id,_  := strconv.Atoi(c.Param("id"))

	query := config.DB.Model(&models.Data{})
		query.Where("id = ?",id)
		query.Delete(&models.Data{})

	c.JSON(200,gin.H{
		"message": true,
	})
}
