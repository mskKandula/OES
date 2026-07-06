package model

import (
	"context"

	xlsx "github.com/tealeg/xlsx/v3"
)

type Student struct {
	Id       int    `json:"id,omitempty"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile,omitempty"`
	Password string `json:"password,omitempty"`
	ClientId string `json:"clientId,omitempty"`
}

// CodeSubmitRequest is the payload the student sends to POST /r/executeCode.
type CodeSubmitRequest struct {
	Language  string `json:"language" binding:"required"`
	Code      string `json:"code" binding:"required"`
	Stdin     string `json:"stdin"`
	TimeoutMs int    `json:"timeoutMs"`
}

// CodeSubmitResponse is returned by POST /r/executeCode.
// If Pending=false: result fields are fully populated (200 OK fast path).
// If Pending=true:  only SubmissionId is set (202 Accepted slow path);
//
//	result is delivered via WebSocket Type-6 message.
type CodeSubmitResponse struct {
	SubmissionId string `json:"submissionId"`
	Pending      bool   `json:"pending"`
	Status       string `json:"status,omitempty"`
	Stdout       string `json:"stdout,omitempty"`
	Stderr       string `json:"stderr,omitempty"`
	ExitCode     int    `json:"exitCode,omitempty"`
	DurationMs   int64  `json:"durationMs,omitempty"`
}

// CodeJob is the message envelope published to RabbitMQ code execution queues.
// Consumed by the code-executor microservice.
type CodeJob struct {
	SubmissionId string `json:"submissionId"`
	Language     string `json:"language"`
	Code         string `json:"code"`
	Stdin        string `json:"stdin"`
	TimeoutMs    int    `json:"timeoutMs"`
	UserId       string `json:"userId"`
	ClientId     string `json:"clientId"`
}

type StudentService interface {
	CreateStudents(context.Context, []byte, string) ([]Student, error)
	FetchStudents(context.Context, string) ([]Student, error)
	FetchAndPrepare(context.Context, string, string) (*xlsx.File, error)
	SubmitCode(ctx context.Context, req CodeSubmitRequest, userId, clientId string) (CodeSubmitResponse, error)
}

type StudentRepository interface {
	Create(context.Context, *Student) error
	ReadAll(context.Context, string) ([]Student, error)
}
