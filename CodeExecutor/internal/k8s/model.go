package model

// CodeJob is the message envelope consumed from RabbitMQ.
// Its shape must match the CodeJob struct published by oes-server's studentService.
type CodeJob struct {
	SubmissionId string `json:"submissionId"`
	Language     string `json:"language"`   // "python" | "go" | "nodejs"
	Code         string `json:"code"`
	Stdin        string `json:"stdin"`
	TimeoutMs    int    `json:"timeoutMs"`
	UserId       string `json:"userId"`
	ClientId     string `json:"clientId"`
}

// ExecutionResult is published to Redis after pod execution.
// It is consumed by:
//   - oes-server's redis.Subscribe("result:<submissionId>") — fast 200 path
//   - oes-server's WebSocket hub via the "general" pub/sub channel — 202 WebSocket path
type ExecutionResult struct {
	SubmissionId string `json:"submissionId"`
	Status       string `json:"status"`     // "completed" | "failed" | "timeout" | "error"
	Stdout       string `json:"stdout"`
	Stderr       string `json:"stderr"`
	ExitCode     int    `json:"exitCode"`
	DurationMs   int64  `json:"durationMs"`
	UserId       string `json:"userId"`
	ClientId     string `json:"clientId"`
	Pending      bool   `json:"pending"` // always false — result is fully available
}
