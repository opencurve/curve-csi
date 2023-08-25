package k8s

import (
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewK8sClient create kubernetes client.
func NewK8sClient() (client *kubernetes.Clientset, err error) {
	var cfg *rest.Config
	cPath := os.Getenv("KUBERNETES_CONFIG_PATH")
	if cPath != "" {
		cfg, err = clientcmd.BuildConfigFromFlags("", cPath)
		if err != nil {
			return nil, err
		}
	} else {
		cfg, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}
	client, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return client, nil
}
