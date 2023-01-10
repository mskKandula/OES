package service

import (
	"context"

	"github.com/mskKandula/oes/api/model"
	"golang.org/x/crypto/bcrypt"
)

type commonService struct {
	CommonRepository model.CommonRepository
}

type CommonServiceConfig struct {
	CommonRepository model.CommonRepository
}

func NewCommonService(ssc *CommonServiceConfig) model.CommonService {
	return &commonService{
		CommonRepository: ssc.CommonRepository,
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
