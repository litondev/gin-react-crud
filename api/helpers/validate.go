package helpers

import (
		"github.com/gin-gonic/gin"
		"errors"
		// "fmt"
		"github.com/go-playground/validator/v10"
)
		
// 	errs := err.(validator.ValidationErrors)
// 	for _, e := range errs {		
// 		fmt.Println(e.Namespace())
// 		fmt.Println(e.Field())
// 		fmt.Println(e.StructNamespace())
// 		fmt.Println(e.StructField())
// 		fmt.Println(e.Tag())
// 		fmt.Println(e.ActualTag())
// 		fmt.Println(e.Kind())
// 		fmt.Println(e.Type())
// 		fmt.Println(e.Value())
// 		fmt.Println(e.Param())
// 		fmt.Println(msgForTag(e))
// 	}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " Harus Diisi"
	case "email":
		return "Invalid email"
	}
	return fe.Error()
}

func Validate(c *gin.Context,validators interface{}) (error){
	if err := c.ShouldBind(validators); err != nil {
		errs := err.(validator.ValidationErrors)
		return errors.New(msgForTag(errs[0]))
	}

	return nil;
}