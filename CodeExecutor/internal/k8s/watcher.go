package k8s

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"

	"github.com/mskKandula/oes/code-executor/internal/pool"
)

// WatchPods watches Kubernetes pods with the label selector
//
//	app.kubernetes.io/component=code-runner,app.kubernetes.io/language=<language>
//
// in the given namespace. It maintains the pool and updates the RabbitMQ channel
// prefetch count in response to pod readiness changes.
//
// This function blocks until ctx is cancelled. Run it as a goroutine.
func WatchPods(
	ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string,
	language string,
	p *pool.Pool,
	ch *amqp.Channel,
) {
	labelSelector := "app.kubernetes.io/component=code-runner,app.kubernetes.io/language=" + language

	// ── Initial sync ──────────────────────────────────────────────────────────
	// Pre-populate the pool with pods that are already Ready before starting the
	// watch. This ensures the dispatcher has work to do immediately on startup.
	podList, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Printf("[watcher/%s] failed to list pods: %v", language, err)
	} else {
		for i := range podList.Items {
			if isPodReady(&podList.Items[i]) {
				p.Add(podList.Items[i].Name)
				log.Printf("[watcher/%s] pre-populated pool with pod %s", language, podList.Items[i].Name)
			}
		}
		updatePrefetch(language, ch, p)
	}

	// ── Continuous watch ──────────────────────────────────────────────────────
	for {
		watcher, err := clientset.CoreV1().Pods(namespace).Watch(ctx, metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			log.Printf("[watcher/%s] watch error (will retry): %v", language, err)
			// ctx cancellation exits the outer loop; other errors retry
			select {
			case <-ctx.Done():
				return
			default:
				continue
			}
		}

		log.Printf("[watcher/%s] watch started", language)

		for event := range watcher.ResultChan() {
			pod, ok := event.Object.(*corev1.Pod)
			if !ok {
				continue
			}

			switch event.Type {
			case watch.Added, watch.Modified:
				if isPodReady(pod) {
					p.Add(pod.Name)
					log.Printf("[watcher/%s] pod %s Ready → added to pool (size=%d)", language, pod.Name, p.Size())
				} else {
					p.Remove(pod.Name)
					log.Printf("[watcher/%s] pod %s NotReady → removed from pool (size=%d)", language, pod.Name, p.Size())
				}
			case watch.Deleted:
				p.Remove(pod.Name)
				log.Printf("[watcher/%s] pod %s Deleted → removed from pool (size=%d)", language, pod.Name, p.Size())
			}

			updatePrefetch(language, ch, p)
		}

		// watcher channel closed (API server reset the watch) — restart
		select {
		case <-ctx.Done():
			return
		default:
			log.Printf("[watcher/%s] watch channel closed, restarting", language)
		}
	}
}

// isPodReady returns true if the pod has a Ready condition with status True.
func isPodReady(pod *corev1.Pod) bool {
	if pod.Status.Phase != corev1.PodRunning {
		return false
	}
	for _, cond := range pod.Status.Conditions {
		if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

// updatePrefetch sets the RabbitMQ channel's QoS prefetch count to the current
// pool size. This ensures RabbitMQ never delivers more jobs than there are
// available pods to execute them.
func updatePrefetch(language string, ch *amqp.Channel, p *pool.Pool) {
	size := p.Size()
	if err := ch.Qos(size, 0, false); err != nil {
		log.Printf("[watcher/%s] failed to update prefetch to %d: %v", language, size, err)
	} else {
		log.Printf("[watcher/%s] prefetch updated to %d", language, size)
	}
}
