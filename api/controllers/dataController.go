package controllers

import "github.com/gin-gonic/gin"
import "github.com/litondev/gin-react-crud/api/config"
import "github.com/litondev/gin-react-crud/api/models"
import "github.com/litondev/gin-react-crud/api/helpers"
import "github.com/litondev/gin-react-crud/api/requests"
// import jwt "github.com/appleboy/gin-jwt/v2"
import "strconv"
import "fmt"
import "strings"
import "html"
import "github.com/xuri/excelize/v2"

import (
	"bufio"
    "encoding/base64"
    "io/ioutil"
    "os"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"bytes"
    "html/template"
    // "strings"
)

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

func ExportPdfData(c *gin.Context){
	htmlTmp, err := template.ParseFiles("./input.html")
    if err != nil {
        fmt.Println(err)
        return
    }

	buf := new(bytes.Buffer)
    err = htmlTmp.Execute(buf, nil)
    if err != nil {
        fmt.Println(err)
        return
    }

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
    if err != nil {
		fmt.Println(err)
		return;
    }

	// f, err := os.Open("./input.html")
	// if f != nil {
	// 	defer f.Close()
	// }
	// if err != nil {
	// 	fmt.Println(err)
	// 	return;
	// }

	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(buf.String())))

	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.Dpi.Set(300)

	err = pdfg.Create()
	if err != nil {
		fmt.Println(err)
		return;
	}

	err = pdfg.WriteFile("./output.pdf")
	if err != nil {
		fmt.Println(err)
		return;
	}

	fmt.Println("Done")
}

func ExportExcelData(c *gin.Context){
	f := excelize.NewFile()

	f.SetCellValue("Sheet1", "A1","ID")
	f.SetCellValue("Sheet1", "B1","Name")
	f.SetCellValue("Sheet1", "C1","Phone")

	result := []map[string]interface{}{}

	query := config.DB.Model(&models.Data{})	
		query.Find(&result)

    for index, row := range result {
		index := index + 2;
		phone := row["phone"]

		if row["phone"] != nil {
			phone = *row["phone"].(*string)
		}
		f.SetCellValue("Sheet1","A"+strconv.Itoa(index),row["id"])
		f.SetCellValue("Sheet1","B"+strconv.Itoa(index),row["name"])
		f.SetCellValue("Sheet1","C"+strconv.Itoa(index),phone)		
	}
    
    if err := f.SaveAs("./assets/Data.xlsx"); err != nil {
        fmt.Println(err)
    }

	// Open file on disk.
    fs, _ := os.Open("./assets/Data.xlsx")
    
    // Read entire JPG into byte slice.
    reader := bufio.NewReader(fs)
    content, _ := ioutil.ReadAll(reader)
    
    // Encode as base64.
    encoded := base64.StdEncoding.EncodeToString(content)
    
	c.JSON(200,gin.H{
		"message": true,
		"download": encoded,
	})
}

func ImportExcelData(c *gin.Context){
	f, err := excelize.OpenFile("./assets/Data.xlsx")
    if err != nil {
        fmt.Println(err)
        return
    }
    // defer func() {
    //     // Close the spreadsheet.
    //     if err := f.Close(); err != nil {
    //         fmt.Println(err)
    //     }
    // }()
    // Get value from cell by given worksheet name and axis.
    cell, err := f.GetCellValue("Sheet1", "B2")
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(cell)
    // Get all the rows in the Sheet1.
    rows, err := f.GetRows("Sheet1")
    if err != nil {
        fmt.Println(err)
        return
    }

    for _, row := range rows {
        for _, colCell := range row {
            fmt.Print(colCell, "\t")
        }
        fmt.Println()
    }
}
