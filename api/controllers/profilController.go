package controllers

import (
	// "crypto/tls"
	"errors"
	"os"
	"fmt"
	// "time"
	// "reflect"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/litondev/gin-react-crud/api/config"
	"github.com/litondev/gin-react-crud/api/helpers"
	"github.com/litondev/gin-react-crud/api/models"
	"github.com/litondev/gin-react-crud/api/requests"
	"path/filepath"
	"github.com/disintegration/imaging"

	// gomail "gopkg.in/mail.v2"
	// "net/http"
	// "fmt"
	// "encoding/json"
)

func UpdateProfilData(c *gin.Context){
	err := helpers.Validate(c, &requests.VUpdateProfilData)

	if err != nil {
		c.JSON(422, gin.H{
			"message": err.Error(),
		})
		return
	}

	claims := jwt.ExtractClaims(c)

	var ID uint = uint(claims["sub"].(float64))

	database := config.DB

	resultSearchEmail := map[string]interface{}{}
	
	querySearchEmail := database.Model(&models.User{})
	querySearchEmail.Where("email = ?",requests.VUpdateProfilData.Email)
	querySearchEmail.Not("id = ?",ID)
	querySearchEmail.First(&resultSearchEmail)

	if len(resultSearchEmail) > 0 {
		c.JSON(500, gin.H{
			"message": "Email telah terpakai",
		})
		return
	}

	resultUser := map[string]interface{}{}
	queryUser := database.Model(&models.User{})
	queryUser.Where("id = ?",ID)
	queryUser.First(&resultUser)

	var isValidPassword bool = helpers.CheckPasswordHash(
		requests.VUpdateProfilData.PasswordConfirm,
		resultUser["password"].(string),
	)

	if isValidPassword == false {
		c.JSON(500, gin.H{
			"message": "Password Konfirmasi Tidak Valid",
		})
		return
	}

	updateUser := &models.User{
		Email : requests.VUpdateProfilData.Email,
		Name : requests.VUpdateProfilData.Name,
	}

	if(requests.VUpdateProfilData.Password != ""){
		hash, _ := helpers.HashPassword(requests.VUpdateProfilData.Password)
		updateUser.Password = hash
	}

		fmt.Println(updateUser)
	queryUpdateUser := database.Model(&models.User{})
		queryUpdateUser.Where("id = ?",ID)
		queryUpdateUser.Updates(&updateUser);

	c.JSON(200, gin.H{
		"message" : true,
	})
}

func UpdateProfilPhoto(c *gin.Context){
	
	file, err := c.FormFile("photo")

	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})

		return
	}

	// if(file.Header["Content-Type"][0] != "image/jpeg" || file.Header["Content-Type"][0] != "image/png"){
	// 	c.JSON(500, gin.H{
	// 		"message": "Image tidak valid",
	// 	})
	// 	return
	// }

	claims := jwt.ExtractClaims(c)

	var ID uint = uint(claims["sub"].(float64))
    var stringID string = strconv.FormatUint(uint64(ID),10)
	/* FOLDER USER HARUS ADA */
	filename := filepath.Base("")+"/assets/users/" + stringID + "-" + file.Filename

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {

	} else {
		e := os.Remove(filename)
		if e != nil {
			fmt.Println(e)
		}	
	}	

	if err := c.SaveUploadedFile(file, filename); err != nil {				
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	openFile, _ := imaging.Open(filename)

	resizeFile := imaging.Resize(openFile, 128, 128, imaging.Lanczos)
	/* FOLDER USER HARUS ADA */
	errResiszeFile := imaging.Save(resizeFile, filepath.Base("")+"/assets/users/" + stringID + "-" + file.Filename)

	fmt.Println(errResiszeFile)

	c.JSON(200, gin.H{
		"message": true,
	})			
}
