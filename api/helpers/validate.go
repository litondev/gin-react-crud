package helpers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strconv"
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
	debug := c.MustGet("DEBUG").(string)
	appDebug, _ := strconv.ParseBool(debug)
	
	if errValidate := c.ShouldBind(validators); errValidate != nil {
		if err, ok := errValidate.(validator.ValidationErrors); ok {		
			return errors.New(msgForTag(err[0]))
		} 
		
		if(appDebug == true){
			return errors.New(errValidate.Error())			
		}

		return errors.New("Terjadi Kesalahan")
	}

	return nil
}
