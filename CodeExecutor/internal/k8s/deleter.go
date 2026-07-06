package k8s

import (
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// DeletePod immediately deletes the named pod (GracePeriodSeconds=0).
// This is called after code execution completes — the pod has done its job
// and a fresh warm replacement will be created by the Deployment controller.
//
// Errors are logged but not returned — a failed delete is non-critical since
// the pod will eventually be garbage-collected or restarted by Kubernetes.
func DeletePod(ctx context.Context, clientset *kubernetes.Clientset, namespace, podName string) {
	gracePeriod := int64(0)
	err := clientset.CoreV1().Pods(namespace).Delete(ctx, podName, metav1.DeleteOptions{
		GracePeriodSeconds: &gracePeriod,
	})
	if err != nil {
		log.Printf("[deleter] failed to delete pod %s: %v (pod may already be terminating)", podName, err)
		return
	}
	log.Printf("[deleter] pod %s deleted", podName)
}
