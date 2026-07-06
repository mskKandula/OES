package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// BuildClient creates a Kubernetes clientset using in-cluster config.
// This works automatically when the code-executor pod has a mounted ServiceAccount
// token (set via serviceAccountName in the Deployment spec).
//
// For local development outside a cluster, set KUBECONFIG env var or use
// rest.NewConfig(rest.DefaultKubernetesUserAgent()) with a kubeconfig file.
func BuildClient() (*kubernetes.Clientset, *rest.Config, error) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, nil, err
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, nil, err
	}

	return clientset, cfg, nil
}
