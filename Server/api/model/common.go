package model

import (
	"context"
	"time"
)

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

type SelfStatus struct {
	StatusMessage string
	ServerTime    time.Time
}

// QueryRequest is the body for POST /r/query — unified intelligence endpoint.
type QueryRequest struct {
	Query     string `json:"query" binding:"required"`
	ContextId string `json:"contextId"` // Optional: topic/subject tag for scoping
}

type CommonService interface {
	UserLogin(context.Context, UserLogin) (int, string, string, error)
	GetRoutes(context.Context, int, string) ([]Route, error)
	GetVideos(context.Context, string) ([]Video, error)
	// Query sends a natural-language query to the Intelligence Agent via gRPC.
	// role maps to the JWT userType ("Examiner" or "Student").
	Query(ctx context.Context, query, role, clientId, userId, contextId string) (string, string, error)
}

type CommonRepository interface {
	LoginUser(context.Context, UserLogin) (int, string, string, string, error)
	ReadRoutes(context.Context, int, string) ([]Route, error)
	ReadVideos(context.Context, string) ([]Video, error)
}
