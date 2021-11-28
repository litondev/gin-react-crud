package requests

type ResetPasswordRequest struct {
	Email           string `form:"email" json:"email" xml:"email" binding:"required,email"`
	NewPassword     string `form:"new_password" json:"new_password" xml:" new_pasword" binding:"required"`
	PasswordConfirm string `form:"password_confirm" json:"password_confirm" xml:"password_confirm" binding:"required"`
	Token           string `form:"token" json:"token" xml:"token" binding:"required"`
}

var VResetPasswordRequest ResetPasswordRequest
