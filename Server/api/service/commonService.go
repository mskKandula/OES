package service

import (
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

func (cs *commonService) UserLogin(userLogin model.UserLogin) (int, string, error) {
	id, userType, password, err := cs.CommonRepository.LoginUser(userLogin)
	if err != nil {
		return 0, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(userLogin.Password)); err != nil {
		return 0, "", err
	}
	return id, userType, nil

}

func (cs *commonService) GetRoutes(id int, uType string) ([]model.Route, error) {

	routes, err := cs.CommonRepository.ReadRoutes(id, uType)
	if err != nil {
		return routes, err
	}

	return routes, nil

}

func (cs *commonService) GetVideos() ([]model.Video, error) {

	videos, err := cs.CommonRepository.ReadVideos()
	if err != nil {
		return videos, err
	}
	return videos, nil
}
