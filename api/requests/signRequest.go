package requests

type SigninRequest struct {
	Email string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}


var VSigninRequest SigninRequest