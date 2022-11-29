package model

type User struct {
	Name     string `json:"name" bson:"name" binding:"required" db:"name"`
	Age      uint8  `json:"age,omitempty" bson:"age" binding:"gte=3,lte=100" db:"age"`
	MobileNo string `json:"mobileNo" bson:"mobileNo" binding:"required" db:"mobileNo"`
	Email    string `json:"email,omitempty" bson:"email" binding:"required" db:"email"`
	Password string `json:"password" bson:"password" binding:"required" db:"password"`
}

type UserService interface {
	CreateUser(user User) error
	CreateVideoFile(fileName, url, imagePath, clientId string) error
	EncodeVideoFile(filepath string) error
}

type UserRepository interface {
	Create(user User, password string) error
	CreateVideo(fileName, url, imagePath, clientId string) error
	EncodeVideo(fileName string) error
}
