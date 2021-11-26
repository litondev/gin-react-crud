package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gin-gonic/gin"
	"fmt"
	// "reflect"
)

type Login struct {
	User     string `form:"user" json:"user"  binding:"required" validate:"required,email" message:"tes"`
	Password string `form:"password" json:"password" binding:"required"  validate:"required" message:"test"`
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	}
	return fe.Error() // default error
}

func Signin(c *gin.Context){
	var json Login
	if err := c.ShouldBind(&json); err != nil {
		// errs := err.(validator.ValidationErrors)
		// for _, e := range errs {
			// can translate each error one at a time.
			//keys := reflect.TypeOf(e).Kind()
			// fmt.Println(e.Error())
			// fmt.Println(keys)
			//for _,p := range(e) {
			//	fmt.Println(p)
			//}
		// }

		// fmt.Println(reflect.TypeOf(errs).Kind())
		// keys := reflect.ValueOf(errs).MapKeys()
		// fmt.Println(keys) // [a b c]
		// var t *validator.FieldError;
		// t = &errs[0];
		// fmt.Println(*t["Key"])
		//fmt.Println(reflect.TypeOf(errs[0]).Kind())
			
			
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {		
			fmt.Println(e.Namespace())
			fmt.Println(e.Field())
			fmt.Println(e.StructNamespace())
			fmt.Println(e.StructField())
			fmt.Println(e.Tag())
			fmt.Println(e.ActualTag())
			fmt.Println(e.Kind())
			fmt.Println(e.Type())
			fmt.Println(e.Value())
			fmt.Println(e.Param())
			fmt.Println(msgForTag(e))
		}

		c.JSON(422, gin.H{"error": errs})
		return
	}
	return;

	// helpers.Validate(c,&requests.VSigninRequest)
	//c.JSON(200,gin.H{
		// "message" : "Signin",
	// });
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

