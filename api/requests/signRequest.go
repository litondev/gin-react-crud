package requests

type SigninRequest struct {
	Email string `form:"email" json:"email" xml:"email" binding:"required,email"`
 	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

var VSigninRequest SigninRequest	
