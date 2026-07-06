package config

import (
	"os"
	"strings"
)

// Config holds all runtime configuration for the code-executor service.
// All values are sourced from environment variables — no config files.
type Config struct {
	// RabbitMQDSN is the AMQP connection string.
	// e.g. amqp://rabbitmq:rabbitmq@oes-messageq/
	RabbitMQDSN string

	// RedisDSN is the Redis connection string.
	// e.g. redis://oes-cache:6379/0
	RedisDSN string

	// Namespace is the Kubernetes namespace where runner pods live.
	// e.g. oes
	Namespace string

	// Languages is the ordered list of languages this instance handles.
	// e.g. ["python", "go", "nodejs"]
	Languages []string
}

// QueueFor returns the RabbitMQ queue name for the given language.
// Convention: code.execute.<language>
func (c *Config) QueueFor(language string) string {
	return "code.execute." + language
}

// Load reads configuration from environment variables and returns a Config.
// Panics if any required variable is missing.
func Load() *Config {
	cfg := &Config{
		RabbitMQDSN: requireEnv("RABBITMQ_DSN"),
		RedisDSN:    requireEnv("REDIS_DSN"),
		Namespace:   getEnvOrDefault("K8S_NAMESPACE", "oes"),
		Languages:   parseLanguages(getEnvOrDefault("LANGUAGES", "python,go,nodejs")),
	}
	return cfg
}

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("required environment variable not set: " + key)
	}
	return val
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func parseLanguages(raw string) []string {
	parts := strings.Split(raw, ",")
	langs := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			langs = append(langs, p)
		}
	}
	return langs
}
