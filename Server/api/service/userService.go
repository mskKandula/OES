package service

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/mskKandula/oes/api/model"
	"github.com/mskKandula/oes/api/pkg/questgen/pb"
)

type userService struct {
	UserRepository model.UserRepository
	QuestgenClient pb.QuestGenServiceClient
}

// UserServiceCOnfig will hold repositories that will eventually be injected into this
// this service layer
type UserServiceConfig struct {
	UserRepository model.UserRepository
	QuestgenClient pb.QuestGenServiceClient
}

func NewUserService(usc *UserServiceConfig) model.UserService {
	return &userService{
		UserRepository: usc.UserRepository,
		QuestgenClient: usc.QuestgenClient,
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

func (us *userService) CreateVideoFile(ctx context.Context, fileName, url, imagePath, clientId, dstPath string) error {
	if err := us.UserRepository.CreateVideo(ctx, fileName, url, imagePath, clientId, dstPath); err != nil {
		return err
	}
	return nil

}

func (us *userService) GenQuestion(ctx context.Context, requestData string) (string, error) {
	r, err := us.QuestgenClient.QuestGen(ctx, &pb.QuestGenRequest{Request: requestData})
	if err != nil {
		return "", err
	}

	return r.GetResponse(), nil

}

// func (us *userService) EncodeVideoFile(fileName string) error {
// 	if err := us.UserRepository.EncodeVideo(fileName); err != nil {
// 		return err
// 	}
// 	return nil

// }
