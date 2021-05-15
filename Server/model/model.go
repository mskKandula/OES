package model

type User struct{
	Name string `json:"name"` 
	Age uint8	`json:"age,omitempty"`
	MobileNo string	`json:"mobileNo"`
	Email string	`json:"email"`
	Password string	`json:"password"`
}