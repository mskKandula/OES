package datasources

import (
	"fmt"

	redis "github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/mskKandula/oes/code-executor/internal/config"
	k8sclient "github.com/mskKandula/oes/code-executor/internal/k8s"
	"github.com/mskKandula/oes/code-executor/internal/util/dalhelper"
)

// DataSources holds all initialised external connections for the code-executor.
// RabbitMQ channels are intentionally excluded — they are per-goroutine and
// opened by the caller directly from RabbitConn.
type DataSources struct {
	RabbitConn *amqp.Connection
	Redis      *redis.Client
	K8sClient  *kubernetes.Clientset
	K8sCfg     *rest.Config
}

// Init initialises all external connections required by the code-executor.
// Returns an error if any connection cannot be established.
func Init(cfg *config.Config) (*DataSources, error) {
	rabbitConn, err := dalhelper.GetRabbitMQConnection(cfg.RabbitMQDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	rdb, err := dalhelper.GetRedisConnection(cfg.RedisDSN)
	if err != nil {
		rabbitConn.Close()
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	clientset, k8sCfg, err := k8sclient.BuildClient()
	if err != nil {
		rabbitConn.Close()
		rdb.Close()
		return nil, fmt.Errorf("failed to build Kubernetes client: %w", err)
	}

	return &DataSources{
		RabbitConn: rabbitConn,
		Redis:      rdb,
		K8sClient:  clientset,
		K8sCfg:     k8sCfg,
	}, nil
}

// Close releases all external connections.
func (ds *DataSources) Close() {
	ds.RabbitConn.Close()
	ds.Redis.Close()
}
