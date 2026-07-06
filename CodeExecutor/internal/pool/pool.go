package pool

import (
	"sync"
)

// Pool is a thread-safe set of available warm pod names for a single language.
// Pod names (not IPs) are used because the Kubernetes pods/exec API identifies
// pods by namespace + name, not by IP.
//
// Lifecycle:
//   - Add(name)     called by k8s watcher when a pod transitions to Ready
//   - Remove(name)  called by k8s watcher when a pod is Deleted or NotReady
//   - Acquire()     called by dispatcher to claim a pod for exclusive execution
//   - Release(name) called by dispatcher if exec fails before pod is deleted
//                   (normally pods are deleted after exec, not released)
//   - Size()        called to set RabbitMQ prefetch = number of available pods
type Pool struct {
	mu   sync.Mutex
	pods []string
}

// New creates an empty Pool.
func New() *Pool {
	return &Pool{pods: make([]string, 0, 8)}
}

// Add adds a pod name to the available pool.
// It is idempotent — adding an already-present name is a no-op.
func (p *Pool) Add(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, existing := range p.pods {
		if existing == name {
			return
		}
	}
	p.pods = append(p.pods, name)
}

// Remove removes a pod name from the pool.
// It is idempotent — removing a non-existent name is a no-op.
func (p *Pool) Remove(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, existing := range p.pods {
		if existing == name {
			// swap-delete: O(1), order unimportant
			p.pods[i] = p.pods[len(p.pods)-1]
			p.pods = p.pods[:len(p.pods)-1]
			return
		}
	}
}

// Acquire pops one pod name from the pool for exclusive use by the dispatcher.
// Returns ("", false) if the pool is empty.
// The caller is responsible for calling Remove (via the k8s watcher on pod delete)
// after the pod is deleted — NOT Release, unless exec fails before deletion.
func (p *Pool) Acquire() (string, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.pods) == 0 {
		return "", false
	}
	// Pop from the end — O(1)
	name := p.pods[len(p.pods)-1]
	p.pods = p.pods[:len(p.pods)-1]
	return name, true
}

// Release returns a pod name to the pool.
// Use only when the pod was acquired but execution failed before pod deletion,
// meaning the pod is still alive and can be reused.
func (p *Pool) Release(name string) {
	p.Add(name)
}

// Size returns the number of currently available pods.
// Used to set RabbitMQ channel prefetch count.
func (p *Pool) Size() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.pods)
}
