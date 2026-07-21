package service

import (
	"context"

	"github.com/mskKandula/oes/api/model"
	"github.com/mskKandula/oes/api/pkg/intelligence/pb"
	"golang.org/x/crypto/bcrypt"
)

type commonService struct {
	CommonRepository    model.CommonRepository
	IntelligenceClient  pb.IntelligenceServiceClient
}

type CommonServiceConfig struct {
	CommonRepository   model.CommonRepository
	IntelligenceClient pb.IntelligenceServiceClient
}

func NewCommonService(ssc *CommonServiceConfig) model.CommonService {
	return &commonService{
		CommonRepository:   ssc.CommonRepository,
		IntelligenceClient: ssc.IntelligenceClient,
	}
}

func (cs *commonService) UserLogin(ctx context.Context, userLogin model.UserLogin) (int, string, string, error) {
	id, userType, password, clientId, err := cs.CommonRepository.LoginUser(ctx, userLogin)
	if err != nil {
		return 0, "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(userLogin.Password)); err != nil {
		return 0, "", "", err
	}
	return id, userType, clientId, nil
}

func (cs *commonService) GetRoutes(ctx context.Context, id int, uType string) ([]model.Route, error) {
	routes, err := cs.CommonRepository.ReadRoutes(ctx, id, uType)
	if err != nil {
		return routes, err
	}
	return routes, nil
}

func (cs *commonService) GetVideos(ctx context.Context, clientId string) ([]model.Video, error) {
	videos, err := cs.CommonRepository.ReadVideos(ctx, clientId)
	if err != nil {
		return videos, err
	}
	return videos, nil
}

// Query sends a natural-language query to the Intelligence Agent via gRPC.
// role maps to the JWT userType ("Examiner" or "Student").
// Returns (answer, toolUsed, error).
func (cs *commonService) Query(ctx context.Context, query, role, clientId, userId, contextId string) (string, string, error) {
	resp, err := cs.IntelligenceClient.Ask(ctx, &pb.IntelligenceRequest{
		Query:     query,
		Role:      role,
		ClientId:  clientId,
		UserId:    userId,
		ContextId: contextId,
	})
	if err != nil {
		return "", "", err
	}
	return resp.GetAnswer(), resp.GetToolUsed(), nil
}
