package model

type User struct {
	Name     string `json:"name"`
	Age      uint8  `json:"age,omitempty"`
	MobileNo string `json:"mobileNo"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Student struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile,omitempty"`
	Password string `json:"password"`
}

type Route struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Description string `json:"description"`
}

type BasicDetails struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Video struct {
	Name          string `json:"name"`
	VideoUrl      string `json:"videoUrl"`
	ThumbnailPath string `json:"thumbnailPath"`
	Description   string `json:"description"`
}
