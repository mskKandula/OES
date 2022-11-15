package model

import (
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
	CreateStudents([]byte, string) ([]Student, error)
	FetchStudents() ([]Student, error)
	FetchAndPrepare(string) (*xlsx.File, error)
}

type StudentRepository interface {
	Create(*Student) error
	ReadAll() ([]Student, error)
}
