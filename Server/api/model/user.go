package model

import "context"

type User struct {
	Name     string `json:"name" bson:"name" binding:"required" db:"name"`
	Age      uint8  `json:"age,omitempty" bson:"age" binding:"gte=3,lte=100" db:"age"`
	MobileNo string `json:"mobileNo" bson:"mobileNo" binding:"required" db:"mobileNo"`
	Email    string `json:"email,omitempty" bson:"email" binding:"required" db:"email"`
	Password string `json:"password" bson:"password" binding:"required" db:"password"`
}

type UserService interface {
	CreateUser(ctx context.Context, user User) error
	CreateVideoFile(ctx context.Context, fileName, url, imagePath, clientId, dstPath string) error
}

type UserRepository interface {
	Create(ctx context.Context, user User, password string) error
	CreateVideo(ctx context.Context, fileName, url, imagePath, clientId, dstPath string) error
}
