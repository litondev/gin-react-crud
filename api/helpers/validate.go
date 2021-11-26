package helpers

import "github.com/gin-gonic/gin"
import "fmt"

func Validate(c *gin.Context,validator interface{}){
	if err := c.ShouldBindJSON(validator); err != nil {
		fmt.Println(err)
		c.JSON(422,gin.H{
			"message": err.Error(),
		})
		return;
	}
}