package requests

type SignupRequest struct {
	Name string `form:"name" json:"name" xml:"name" binding:"required"`
	Email string `form:"email" json:"email" xml:"email" binding:"required,email"`
 	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

var VSignupRequest SignupRequest	
