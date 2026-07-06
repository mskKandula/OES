package dispatcher

import (
	"context"
	"encoding/json"
	"log"

	redis "github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	k8sexec "github.com/mskKandula/oes/code-executor/internal/k8s"
	"github.com/mskKandula/oes/code-executor/internal/model"
	"github.com/mskKandula/oes/code-executor/internal/pool"
	"github.com/mskKandula/oes/code-executor/internal/result"
)

// Dispatcher consumes code execution jobs for a single language from RabbitMQ,
// dispatches each job to an available warm pod via Kubernetes pods/exec,
// deletes the pod after execution, and publishes the result to Redis.
//
// One Dispatcher runs per language (python, go, nodejs).
type Dispatcher struct {
	language  string
	queueName string
	ch        *amqp.Channel
	pool      *pool.Pool
	k8sCfg    *rest.Config
	clientset *kubernetes.Clientset
	namespace string
	rdb       *redis.Client
}

// New creates a Dispatcher for the given language.
func New(
	language string,
	queueName string,
	ch *amqp.Channel,
	p *pool.Pool,
	k8sCfg *rest.Config,
	clientset *kubernetes.Clientset,
	namespace string,
	rdb *redis.Client,
) *Dispatcher {
	return &Dispatcher{
		language:  language,
		queueName: queueName,
		ch:        ch,
		pool:      p,
		k8sCfg:    k8sCfg,
		clientset: clientset,
		namespace: namespace,
		rdb:       rdb,
	}
}

// Run starts the dispatcher consume loop. Blocks until ctx is cancelled.
// Run as a goroutine: go d.Run(ctx)
func (d *Dispatcher) Run(ctx context.Context) {
	// Declare queue idempotently (oes-server also declares it — safe to call twice)
	_, err := d.ch.QueueDeclare(d.queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("[dispatcher/%s] failed to declare queue %s: %v", d.language, d.queueName, err)
	}

	// Initial prefetch = current pool size (updated dynamically by watcher)
	// Setting 0 initially is safe — watcher will call Qos as soon as pods are listed
	if err := d.ch.Qos(d.pool.Size(), 0, false); err != nil {
		log.Printf("[dispatcher/%s] initial Qos failed: %v", d.language, err)
	}

	msgs, err := d.ch.Consume(
		d.queueName, // queue
		"",          // consumer tag — auto-generated
		false,       // auto-ack = false (manual ack after result published)
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		log.Fatalf("[dispatcher/%s] failed to start consumer: %v", d.language, err)
	}

	log.Printf("[dispatcher/%s] consumer started on queue %s", d.language, d.queueName)

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				log.Printf("[dispatcher/%s] message channel closed", d.language)
				return
			}

			var job model.CodeJob
			if err := json.Unmarshal(msg.Body, &job); err != nil {
				log.Printf("[dispatcher/%s] failed to unmarshal job: %v — discarding", d.language, err)
				msg.Nack(false, false) // discard malformed message
				continue
			}

			// Acquire a warm pod from the pool
			podName, ok := d.pool.Acquire()
			if !ok {
				// No pod available right now — NACK and requeue
				// This is a transient state (pod replacement in progress)
				// The message will be redelivered when a pod becomes available
				log.Printf("[dispatcher/%s] no warm pod available for submissionId=%s — requeuing", d.language, job.SubmissionId)
				msg.Nack(false, true) // requeue=true
				continue
			}

			// Execute asynchronously so the dispatcher loop can continue
			// receiving the next message (up to prefetch limit) immediately
			go d.executeJob(ctx, job, podName, msg)

		case <-ctx.Done():
			log.Printf("[dispatcher/%s] context cancelled — shutting down", d.language)
			return
		}
	}
}

// executeJob runs a single code execution job:
//  1. exec into warm pod
//  2. delete pod after execution (success or failure)
//  3. publish result to Redis
//  4. ACK the RabbitMQ message
func (d *Dispatcher) executeJob(ctx context.Context, job model.CodeJob, podName string, msg amqp.Delivery) {
	log.Printf("[dispatcher/%s] executing submissionId=%s in pod=%s", d.language, job.SubmissionId, podName)

	execResult, err := k8sexec.Execute(ctx, d.k8sCfg, d.clientset, d.namespace, podName, job)

	// Always delete the pod after execution — the K8s watcher handles pool.Remove
	// when it receives the DELETED event for this pod
	k8sexec.DeletePod(ctx, d.clientset, d.namespace, podName)

	var res model.ExecutionResult
	res.SubmissionId = job.SubmissionId
	res.UserId = job.UserId
	res.ClientId = job.ClientId
	res.Pending = false

	if err != nil {
		log.Printf("[dispatcher/%s] exec error for submissionId=%s: %v", d.language, job.SubmissionId, err)
		res.Status = "error"
		res.Stderr = err.Error()
		res.ExitCode = 1
	} else {
		res.Stdout = execResult.Stdout
		res.Stderr = execResult.Stderr
		res.ExitCode = execResult.ExitCode
		res.DurationMs = execResult.DurationMs

		switch execResult.ExitCode {
		case 0:
			res.Status = "completed"
		case 124:
			res.Status = "timeout"
		default:
			res.Status = "failed"
		}
	}

	// Publish result to Redis (fires oes-server 5s select + WebSocket delivery)
	result.Write(ctx, d.rdb, res)

	// ACK the RabbitMQ message — job is fully processed
	msg.Ack(false)

	log.Printf("[dispatcher/%s] submissionId=%s completed status=%s", d.language, job.SubmissionId, res.Status)
}
