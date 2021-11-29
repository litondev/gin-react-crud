package controllers


import (
	"github.com/gin-gonic/gin"
	"github.com/litondev/gin-react-crud/api/models"
	"github.com/litondev/gin-react-crud/api/helpers"
	"github.com/litondev/gin-react-crud/api/requests"
	"strconv"
	"fmt"
	"strings"
	"html"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"bufio"
    "encoding/base64"
    "io/ioutil"
    "os"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"bytes"
    "html/template"
	"path/filepath"
)

func IndexData(c *gin.Context){
	database := c.MustGet("DB").(*gorm.DB);	

	page := c.DefaultQuery("page", "1")
	new_page,_ := strconv.Atoi(page)

	per_page := c.DefaultQuery("per_page", "10")
	new_per_page,_ := strconv.Atoi(per_page)

	search := c.DefaultQuery("search","")

	result := []map[string]interface{}{}

	var resultCount int64;

	queryResultCount := database.Model(&models.Data{})
		queryResultCount.Select("id")
		queryResultCount.Count(&resultCount)

	new_page = (int(resultCount) - ( (new_page * new_per_page) - new_per_page) )

	query := database.Model(&models.Data{})
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
	database := c.MustGet("DB").(*gorm.DB);	

	errValidate := helpers.Validate(c, &requests.VDataRequest)

	if errValidate != nil {
		c.JSON(422, gin.H{
			"message": errValidate.Error(),
		})
		return
	}
		
	var phone string = html.EscapeString(strings.Trim(requests.VDataRequest.Phone," "))

	users := &models.Data{
		Name : html.EscapeString(strings.Trim(requests.VDataRequest.Name," ")),
		Phone : &phone,
	}

	queryData := database.Model(&models.Data{})
		queryData.Create(&users)
	
	if queryData.Error != nil {
		fmt.Println(queryData.Error)
		c.JSON(500,gin.H{
			"message" : "Terjadi Kesalahan",
		})
		return 
	}

	c.JSON(200,gin.H{
		"message": true,
	})
}

func ShowData(c *gin.Context){
	database := c.MustGet("DB").(*gorm.DB);	

	id,errGetParam  := strconv.Atoi(c.Param("id"))

	if errGetParam != nil {
		fmt.Println(errGetParam.Error())
		c.JSON(500,gin.H{
			"message" : "Terjadi Kesalahan",
		})
		return 
	}

	resultData := map[string]interface{}{}

	queryData := database.Model(&models.Data{})
		queryData.Select("id","name","phone")
		queryData.Where("id = ?",id)
		queryData.First(&resultData)

	if len(resultData) == 0 {
		c.JSON(404,gin.H{
			"message" : "Not Found",
		})
		return 
	}

	c.JSON(200,gin.H{
		"data" : resultData,
	})
	return 
}

func UpdateData(c *gin.Context){
	errValidate := helpers.Validate(c, &requests.VDataRequest)

	if errValidate != nil {
		c.JSON(422, gin.H{
			"message": errValidate.Error(),
		})
		return
	}

	database := c.MustGet("DB").(*gorm.DB);	

	id,errGetParam  := strconv.Atoi(c.Param("id"))

	if errGetParam != nil {
		fmt.Println(errGetParam.Error())
		c.JSON(500,gin.H{
			"message" : "Terjadi Kesalahan",
		})
		return 
	}

	var phone string = html.EscapeString(strings.Trim(requests.VDataRequest.Phone," "))

	queryData := database.Model(&models.Data{})
		queryData.Select("name","phone")
		queryData.Where("id = ?",id)
	
	resultData := map[string]interface{}{}

		queryData.First(&resultData)

	if len(resultData) == 0 {
		c.JSON(404,gin.H{
			"message" : "Not Found",
		})
		return 
	}

		queryData.Updates(&models.Data{
			Name : html.EscapeString(strings.Trim(requests.VDataRequest.Name," ")),
			Phone : &phone,
		})

	if queryData.Error != nil {
		fmt.Println(queryData.Error)
		c.JSON(500,gin.H{
			"message" : "Terjadi Kesalahan",
		})
		return 
	}

	c.JSON(200,gin.H{
		"message": true,
	})
	return
}

func DestoryData(c *gin.Context){
	database := c.MustGet("DB").(*gorm.DB);	

	id,errGetParam  := strconv.Atoi(c.Param("id"))

	if errGetParam != nil {
		fmt.Println(errGetParam.Error())
		c.JSON(500,gin.H{
			"message" : "Terjadi Kesalahan",
		})
		return 
	}

	queryData := database.Model(&models.Data{})
		queryData.Where("id = ?",id)

	resultData := map[string]interface{}{}

		queryData.First(&resultData)

	if len(resultData) == 0 {
		c.JSON(404,gin.H{
			"message" : "Not Found",
		})
		return 
	}

		queryData.Delete(&models.Data{})

	if queryData.Error != nil {
		fmt.Println(queryData.Error)
		c.JSON(500,gin.H{
			"message" : "Terjadi Kesalahan",
		})
		return 
	}

	c.JSON(200,gin.H{
		"message": true,
	})
	return 
}

func ExportPdfData(c *gin.Context){
	database := c.MustGet("DB").(*gorm.DB);	

	resultData := []map[string]interface{}{}

	DataFields := database.Model(&models.Data{})	
		DataFields.Find(&resultData)

	var dataTemplate = map[string]interface{}{
		"DataFields" : resultData,
	}

	htmlTmp, errHtmlTmp := template.ParseFiles("./template/data.html")

    if errHtmlTmp != nil {
        fmt.Println(errHtmlTmp)
		c.JSON(500,gin.H{
			"message" : "Terjadi Kesalahan",
		})
        return
    }

	buffer := new(bytes.Buffer)

    errBuffer := htmlTmp.Execute(buffer, dataTemplate)

    if errBuffer != nil {
        fmt.Println(errBuffer)
		c.JSON(500,gin.H{
			"message" : "Terjadi Kesalahan",
		})
        return
    }

	pdfg, errWkhtml := wkhtmltopdf.NewPDFGenerator()

    if errWkhtml != nil {
		fmt.Println(errWkhtml)
		c.JSON(500,gin.H{
			"message" : "Terjadi Kesalahan",
		})
		return;
    }

	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(buffer.String())))
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.Dpi.Set(300)

	errPdfgCreate := pdfg.Create()
	if errPdfgCreate != nil {
		fmt.Println(errPdfgCreate)
		c.JSON(500,gin.H{
			"message" : "Terjadi Kesalahan",
		})
		return;
	}

	errWriteFile := pdfg.WriteFile("./assets/output.pdf")
	if errWriteFile != nil {
		fmt.Println(errWriteFile)
		c.JSON(500,gin.H{
			"message" : "Terjadi Kesalahan",
		})
		return;
	}

	fileExcelOpen,errFileExcelOpen := os.Open("./assets/output.pdf")
	
	if errFileExcelOpen != nil {
		fmt.Println(errFileExcelOpen.Error())

		c.JSON(500,gin.H{
			"message" : "Terjadi Kesalahan",
		})
		return 
	}

	reader := bufio.NewReader(fileExcelOpen)

    content, errContent := ioutil.ReadAll(reader)

	if errContent != nil {
		fmt.Println(errContent.Error())
		c.JSON(500,gin.H{
			"message": "Terjadi Kesalahan",
		})
		return 
	}
        
    encoded := base64.StdEncoding.EncodeToString(content)

	c.JSON(200,gin.H{
		"message": true,
		"download": encoded,
	})
	return 
}

func ExportExcelData(c *gin.Context){
	database := c.MustGet("DB").(*gorm.DB);	

	fileExcel := excelize.NewFile()

	fileExcel.SetCellValue("Sheet1", "A1","ID")
	fileExcel.SetCellValue("Sheet1", "B1","Name")
	fileExcel.SetCellValue("Sheet1", "C1","Phone")

	resultData := []map[string]interface{}{}

	queryData := database.Model(&models.Data{})	
		queryData.Find(&resultData)

    for index, row := range resultData {
		index := index + 2;
		phone := row["phone"]

		if row["phone"] != nil {
			phone = *row["phone"].(*string)
		}

		fileExcel.SetCellValue("Sheet1", "A" + strconv.Itoa(index), row["id"])
		fileExcel.SetCellValue("Sheet1", "B" + strconv.Itoa(index), row["name"])
		fileExcel.SetCellValue("Sheet1", "C" + strconv.Itoa(index), phone)		
	}
    
    if errSaveFileExcel := fileExcel.SaveAs("./assets/Data.xlsx"); errSaveFileExcel != nil {
        fmt.Println(errSaveFileExcel.Error())
		c.JSON(500,gin.H{
			"message" : "Terjadi Kesalahan",
		})
		return 
    }

    fileExcelOpen,errFileExcelOpen := os.Open("./assets/Data.xlsx")
	
	if errFileExcelOpen != nil {
		fmt.Println(errFileExcelOpen.Error())

		c.JSON(500,gin.H{
			"message" : "Terjadi Kesalahan",
		})
		return 
	}

    reader := bufio.NewReader(fileExcelOpen)

    content, errContent := ioutil.ReadAll(reader)

	if errContent != nil {
		fmt.Println(errContent.Error())
		c.JSON(500,gin.H{
			"message": "Terjadi Kesalahan",
		})
		return 
	}
        
    encoded := base64.StdEncoding.EncodeToString(content)
    
	c.JSON(200,gin.H{
		"message": true,
		"download": encoded,
	})
	return 
}

func ImportExcelData(c *gin.Context){
	database := c.MustGet("DB").(*gorm.DB);	
	
	file, errGetFile := c.FormFile("excel")

	if errGetFile != nil {
		fmt.Println(errGetFile.Error())		
		c.JSON(500, gin.H{
			"message": "Terjadi Kesalahan",
		})
		return
	}

	// VALIDATION FILE

	pathname := filepath.Base("") + "/assets/ " + file.Filename

	if errUploadFile := c.SaveUploadedFile(file, pathname); errUploadFile != nil {				
		fmt.Println(errUploadFile.Error())
		c.JSON(500, gin.H{
			"message": "Terjadi Kesalahan",
		})
		return
	}

	fileExcel, errFileExcel := excelize.OpenFile(pathname)
	
    if errFileExcel != nil {
        fmt.Println(errFileExcel.Error())
		c.JSON(500,gin.H{
			"message" : "Terjadi Kesalahan",
		})
        return
    }    
	
    // Get all the rows in the Sheet1.
    rows, errGetRow := fileExcel.GetRows("Sheet1")
    if errGetRow != nil {
        fmt.Println(errGetRow.Error())
		c.JSON(500,gin.H{
			"message" : "Terjadi Kesalahan",
		})
        return
    }
	
    for index, row := range rows {
        // for _, colCell := range row {			
        //  	fmt.Print(colCell, "\t")
        // }

		if(index > 0){
			// VALIDATION DATA
			
			var phone *string = &row[1]
			
			resultData := database.Create(&models.Data{
				Name : row[0],
				Phone : phone,
			})	

			fmt.Println(resultData.Error)
		}

        fmt.Println()
    }
	
	c.JSON(200,gin.H{
	 	"data" : len(rows) - 1 ,
	})
}
