package handler

import (
	"github.com/mskKandula/oes/api/model"
)

type Handler struct {
	UserService    model.UserService
	StudentService model.StudentService
	CommonService  model.CommonService
}

func NewHandler(us model.UserService, ss model.StudentService, cs model.CommonService) *Handler {
	return &Handler{UserService: us, StudentService: ss, CommonService: cs}
}
