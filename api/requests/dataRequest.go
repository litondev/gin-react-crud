package requests

type DataRequest struct {	
	Name string `form:"name" json:"name" xml:"name" binding:"required,max=25"`
	Phone string `form:"phone" json:"phone" xml:"phone" binding:"max=25"`	
}

var VDataRequest DataRequest
