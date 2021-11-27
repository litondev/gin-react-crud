package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"

	// "fmt"
	"github.com/go-playground/validator/v10"
)

/*
	SANGAT PENTING JIKA MENGUNAKAN VALIDASI HARUS ADA BODYNYA JIKA PADA REQUEST CLIENT
	JIKA TIDAK AKAN EFO ERROR

	DAN JIKA ADA APPLICATION/JSONNYA MAKA JSON HARUS ADA ISINYA
*/

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

func Validate(c *gin.Context, validators interface{}) error {
	if err := c.ShouldBind(validators); err != nil {
		if _, ok := err.(validator.FieldError); ok {
			errs := err.(validator.ValidationErrors)
			return errors.New(msgForTag(errs[0]))
		} else {
			// return errors.New(err.Error())
			return errors.New("Terjadi Kesalahan")
		}
	}

	return nil
}
