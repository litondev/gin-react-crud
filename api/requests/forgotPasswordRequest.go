package requests

type ForgotPasswordRequest struct {
	Email string `form:"email" json:"email" xml:"email" binding:"required,email"`
}

var VForgotPasswordRequest ForgotPasswordRequest
