package controllers

import (	
	"os"
	"fmt"	
	"strconv"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/litondev/gin-react-crud/api/helpers"
	"github.com/litondev/gin-react-crud/api/models"
	"github.com/litondev/gin-react-crud/api/requests"
	"path/filepath"
	"github.com/disintegration/imaging"
	"gorm.io/gorm"
)

func UpdateProfilData(c *gin.Context){
	errValidate := helpers.Validate(c, &requests.VUpdateProfilData)

	if errValidate != nil {
		c.JSON(422, gin.H{
			"message": errValidate.Error(),
		})
		return
	}

	database := c.MustGet("DB").(*gorm.DB);	

	tx := database.Begin()

	claims := jwt.ExtractClaims(c)

	var ID uint = uint(claims["sub"].(float64))

	resultSearchEmail := map[string]interface{}{}
	
	querySearchEmail := database.Model(&models.User{})
		querySearchEmail.Where("email = ?",requests.VUpdateProfilData.Email)
		querySearchEmail.Not("id = ?",ID)
		querySearchEmail.First(&resultSearchEmail)

	if len(resultSearchEmail) > 0 {
		tx.Rollback()		
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
		tx.Rollback()		
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
		hash, errHash := helpers.HashPassword(requests.VUpdateProfilData.Password)

		if(errHash != nil){
			tx.Rollback()
			fmt.Println(errHash.Error())		
			c.JSON(500,gin.H{
				"message" : "Terjadi Kesalahan",
			})
			return
		}

		updateUser.Password = hash
	}
		
	queryUpdateUser := database.Model(&models.User{})
		queryUpdateUser.Where("id = ?",ID)
		queryUpdateUser.Updates(&updateUser);

	if queryUpdateUser.Error != nil {
		tx.Rollback()
		fmt.Println(queryUpdateUser.Error)
		c.JSON(500, gin.H{
			"message": "Terjadi Kesalahan",
		})
		return
	}

	tx.Commit()
	c.JSON(200, gin.H{
		"message" : true,
	})
	return 
}

func UpdateProfilPhoto(c *gin.Context){
	file, errGetFile := c.FormFile("photo")

	if errGetFile != nil {
		fmt.Println(errGetFile.Error())		
		c.JSON(500, gin.H{
			"message": "Terjadi Kesalahan",
		})
		return
	}

	if(false == (file.Header["Content-Type"][0] != "image/jpeg" || file.Header["Content-Type"][0] != "image/png")){
		fmt.Println(file.Header["Content-Type"][0])
		c.JSON(500, gin.H{
			"message": "Gambar tidak valid",
		})
	 	return
	}	

	claims := jwt.ExtractClaims(c)

	var ID uint = uint(claims["sub"].(float64))

    var stringID string = strconv.FormatUint(uint64(ID),10)

	filename := stringID + "-" + file.Filename;
	/* FOLDER USER HARUS ADA */
	pathname := filepath.Base("") + "/assets/users/" + filename

	if _,errFileExists := os.Stat(pathname); errFileExists == nil {
		errRemoveFile := os.Remove(pathname)
		if errRemoveFile != nil {
			fmt.Println(errRemoveFile.Error())
			c.JSON(500,gin.H{
				"message" : "Terjadi Kesalahan",
			})
			return 
		}
	}

	if errUploadFile := c.SaveUploadedFile(file, pathname); errUploadFile != nil {				
		fmt.Println(errUploadFile.Error())
		c.JSON(500, gin.H{
			"message": "Terjadi Kesalahan",
		})
		return
	}

	openFile, errOpenFile := imaging.Open(pathname)

	if(errOpenFile != nil){
		fmt.Println(errOpenFile.Error())
		c.JSON(500, gin.H{
			"message": "Terjadi Kesalahan",
		})
		return
	}

	resizeFile := imaging.Resize(openFile, 128, 128, imaging.Lanczos)
	
	errRessizeFile := imaging.Save(resizeFile, pathname)

	if(errRessizeFile != nil){
		fmt.Println(errRessizeFile.Error())
		c.JSON(200,gin.H{
			"message": "Terjadi Kesalahan",
		})	
		return 
	}

	database := c.MustGet("DB").(*gorm.DB);	

	tx := database.Begin()

	queryUpdateUser := database.Model(&models.User{})
		queryUpdateUser.Select("photo")
		queryUpdateUser.Where("id = ?",ID)
		queryUpdateUser.Update("photo",filename)

	if queryUpdateUser.Error != nil {
		tx.Rollback()
		fmt.Println(queryUpdateUser.Error)
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
