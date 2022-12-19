package model

import "context"

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Route struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Description string `json:"description"`
}

type Video struct {
	Name          string `json:"name"`
	VideoUrl      string `json:"videoUrl"`
	ThumbnailPath string `json:"thumbnailPath"`
	Description   string `json:"description"`
}

type CommonService interface {
	UserLogin(context.Context, UserLogin) (int, string, string, error)
	GetRoutes(context.Context, int, string) ([]Route, error)
	GetVideos(context.Context, string) ([]Video, error)
}

type CommonRepository interface {
	LoginUser(context.Context, UserLogin) (int, string, string, string, error)
	ReadRoutes(context.Context, int, string) ([]Route, error)
	ReadVideos(context.Context, string) ([]Video, error)
}
