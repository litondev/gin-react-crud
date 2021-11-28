package requests

type UpdateProfilData struct {
	Name string `form:"name" json:"name" xml:"name" binding:"required"`
	Email string `form:"email" json:"email" xml:"email" binding:"required,email"`
 	Password string `form:"password" json:"password" xml:"password" binding:"max=255"`
	PasswordConfirm string `form:"password_confirm" json:"password_confirm" xml:"password_confirm" binding:"required,max=255"`
}

var VUpdateProfilData UpdateProfilData
