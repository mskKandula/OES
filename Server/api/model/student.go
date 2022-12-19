package model

import (
	"context"

	xlsx "github.com/tealeg/xlsx/v3"
)

type Student struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile,omitempty"`
	Password string `json:"password,omitempty"`
	ClientId string `json:"clientId,omitempty"`
}

type StudentService interface {
	CreateStudents(context.Context, []byte, string) ([]Student, error)
	FetchStudents(context.Context, string) ([]Student, error)
	FetchAndPrepare(context.Context, string, string) (*xlsx.File, error)
}

type StudentRepository interface {
	Create(context.Context, *Student) error
	ReadAll(context.Context, string) ([]Student, error)
}
