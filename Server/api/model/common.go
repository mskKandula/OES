package model

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
	UserLogin(UserLogin) (int, string, error)
	GetRoutes(int, string) ([]Route, error)
	GetVideos() ([]Video, error)
}

type CommonRepository interface {
	LoginUser(UserLogin) (int, string, string, error)
	ReadRoutes(int, string) ([]Route, error)
	ReadVideos() ([]Video, error)
}
