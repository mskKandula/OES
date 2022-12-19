package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/mskKandula/oes/api/model"
)

type userService struct {
	UserRepository model.UserRepository
}

// UserServiceCOnfig will hold repositories that will eventually be injected into this
// this service layer
type UserServiceConfig struct {
	UserRepository model.UserRepository
}

func NewUserService(usc *UserServiceConfig) model.UserService {
	return &userService{
		UserRepository: usc.UserRepository,
	}
}

func (us *userService) CreateUser(ctx context.Context, user model.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	hashedPassword := string(hash)

	if err = us.UserRepository.Create(ctx, user, hashedPassword); err != nil {
		return err
	}

	return nil
}

func (us *userService) CreateVideoFile(ctx context.Context, fileName, url, imagePath, clientId string) error {
	if err := us.UserRepository.CreateVideo(ctx, fileName, url, imagePath, clientId); err != nil {
		return err
	}
	return nil

}
