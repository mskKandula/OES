package model

import "context"

// Publisher is the contract for publishing async jobs to a message queue.
type Publisher interface {
	PublishMessageWithContext(ctx context.Context, queue string, body []byte) error
}
