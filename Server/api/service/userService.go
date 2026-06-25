package service

import (
	"context"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/mskKandula/oes/api/model"
	"github.com/mskKandula/oes/api/pkg/questgen/pb"
)

type userService struct {
	UserRepository model.UserRepository
	QuestgenClient pb.QuestGenServiceClient
	Publisher      model.Publisher
}

// UserServiceConfig holds repositories and dependencies injected into this service layer.
type UserServiceConfig struct {
	UserRepository model.UserRepository
	QuestgenClient pb.QuestGenServiceClient
	Publisher      model.Publisher
}

func NewUserService(usc *UserServiceConfig) model.UserService {
	return &userService{
		UserRepository: usc.UserRepository,
		QuestgenClient: usc.QuestgenClient,
		Publisher:      usc.Publisher,
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

// CreateVideoFile persists the video metadata and then publishes an async encode job.
func (us *userService) CreateVideoFile(ctx context.Context, fileName, url, imagePath, clientId, dstPath string) error {
	if err := us.UserRepository.CreateVideo(ctx, fileName, url, imagePath, clientId); err != nil {
		return err
	}

	// Publish the encode job after a successful DB commit.
	if err := us.Publisher.PublishMessageWithContext(ctx, "encode", []byte(dstPath)); err != nil {
		log.Printf("encode job: failed to publish for %s: %v", fileName, err)
	}

	return nil
}

func (us *userService) GenQuestion(ctx context.Context, requestData, clientId, contextId string) (string, error) {
	r, err := us.QuestgenClient.QuestGen(ctx, &pb.QuestGenRequest{
		Request:   requestData,
		ClientId:  clientId,
		ContextId: contextId,
	})
	if err != nil {
		return "", err
	}

	return r.GetResponse(), nil
}

func (us *userService) AskQuestion(ctx context.Context, question, contextId, clientId string) (string, error) {
	r, err := us.QuestgenClient.AskQuestion(ctx, &pb.AskQuestionRequest{
		Question:  question,
		ContextId: contextId,
		ClientId:  clientId,
	})
	if err != nil {
		return "", err
	}

	return r.GetAnswer(), nil
}

func (us *userService) CreateExam(ctx context.Context, clientId, examName, examType string) (int64, error) {
	examId, err := us.UserRepository.ExamCreation(ctx, clientId, examName, examType)
	if err != nil {
		return 0, err
	}

	return examId, nil
}
