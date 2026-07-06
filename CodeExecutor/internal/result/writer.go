package result

import (
	"context"
	"encoding/json"
	"log"
	"time"

	redis "github.com/go-redis/redis/v8"

	"github.com/mskKandula/oes/code-executor/internal/model"
)

// Write delivers an ExecutionResult back to oes-server via two Redis mechanisms:
//
//  1. PUBLISH result:<submissionId>
//     Fires oes-server's redis.Subscribe("result:<submissionId>") if the HTTP handler
//     is still within its 5-second optimistic wait window → 200 OK fast path.
//
//  2. SET result:<submissionId> EX 300
//     Persists the result for 5 minutes as a fallback for polling or reconnect.
//
//  3. PUBLISH general
//     Picked up by the existing oes-server WebSocket hub (pool.go) which routes
//     it as a Type-6 message to the student's WebSocket connection → 202 slow path.
func Write(ctx context.Context, rdb *redis.Client, result model.ExecutionResult) {
	body, err := json.Marshal(result)
	if err != nil {
		log.Printf("[result/writer] failed to marshal result for %s: %v", result.SubmissionId, err)
		return
	}

	// 1. Publish to per-submission channel — fires the 5s select block in oes-server
	if err := rdb.Publish(ctx, "result:"+result.SubmissionId, string(body)).Err(); err != nil {
		log.Printf("[result/writer] PUBLISH result:%s failed: %v", result.SubmissionId, err)
	}

	// 2. Persist result for fallback access (5-minute TTL)
	if err := rdb.Set(ctx, "result:"+result.SubmissionId, body, 5*time.Minute).Err(); err != nil {
		log.Printf("[result/writer] SET result:%s failed: %v", result.SubmissionId, err)
	}

	// 3. Publish to the existing "general" WebSocket channel so oes-server's WS hub
	//    routes a Type-6 message to the correct student connection.
	//    The WS hub uses userId to route to the right connection (Type-4 pattern).
	wsPayload, _ := json.Marshal(map[string]any{
		"type":         6,
		"userId":       result.UserId,
		"clientId":     result.ClientId,
		"submissionId": result.SubmissionId,
		"status":       result.Status,
		"stdout":       result.Stdout,
		"stderr":       result.Stderr,
		"exitCode":     result.ExitCode,
		"durationMs":   result.DurationMs,
	})
	if err := rdb.Publish(ctx, "general", string(wsPayload)).Err(); err != nil {
		log.Printf("[result/writer] PUBLISH general failed for %s: %v", result.SubmissionId, err)
	}

	log.Printf("[result/writer] result delivered for submissionId=%s status=%s durationMs=%d",
		result.SubmissionId, result.Status, result.DurationMs)
}
