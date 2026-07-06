package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/mskKandula/oes/code-executor/internal/config"
	"github.com/mskKandula/oes/code-executor/internal/datasources"
	"github.com/mskKandula/oes/code-executor/internal/dispatcher"
	k8sclient "github.com/mskKandula/oes/code-executor/internal/k8s"
	"github.com/mskKandula/oes/code-executor/internal/pool"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("code-executor starting")

	cfg := config.Load()
	log.Printf("config loaded: namespace=%s languages=%v", cfg.Namespace, cfg.Languages)

	// ── Data sources ──────────────────────────────────────────────────────────
	ds, err := datasources.Init(cfg)
	if err != nil {
		log.Fatalf("failed to initialise data sources: %v", err)
	}
	defer ds.Close()
	log.Println("all data sources connected")

	// ── Per-language RabbitMQ channels ────────────────────────────────────────
	// amqp.Channel is not goroutine-safe — one channel per goroutine.
	channels := make(map[string]*amqp.Channel, len(cfg.Languages))
	for _, lang := range cfg.Languages {
		ch, err := ds.RabbitConn.Channel()
		if err != nil {
			log.Fatalf("failed to open RabbitMQ channel for %s: %v", lang, err)
		}
		defer ch.Close()
		channels[lang] = ch
	}

	// ── Context with graceful shutdown ────────────────────────────────────────
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigCh
		log.Printf("received signal %v — shutting down", sig)
		cancel()
	}()

	// ── Start one pool watcher + dispatcher per language ──────────────────────
	for _, lang := range cfg.Languages {
		lang := lang // capture for goroutine
		ch := channels[lang]
		p := pool.New()
		queueName := cfg.QueueFor(lang)

		// Watcher: maintains pool and prefetch as pods become Ready/NotReady/Deleted
		go k8sclient.WatchPods(ctx, ds.K8sClient, cfg.Namespace, lang, p, ch)

		// Dispatcher: consumes from queue, execs into warm pods, publishes results
		d := dispatcher.New(lang, queueName, ch, p, ds.K8sCfg, ds.K8sClient, cfg.Namespace, ds.Redis)
		go d.Run(ctx)

		log.Printf("watcher and dispatcher started for language: %s (queue: %s)", lang, queueName)
	}

	// ── Block until shutdown signal ───────────────────────────────────────────
	<-ctx.Done()
	log.Println("code-executor stopped")
}
