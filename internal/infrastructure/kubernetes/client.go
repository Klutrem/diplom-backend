package kubernetes

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/fx"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type KubernetesClient struct {
	clientset *kubernetes.Clientset
}

var Module = fx.Module("kubernetes", fx.Provide(NewKubernetesClient))

func NewKubernetesClient() (*KubernetesClient, error) {
	var config *rest.Config
	var err error

	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		config, err = rest.InClusterConfig()
	} else {
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	if err != nil {
		fmt.Println("Error getting config:", err)
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Error creating clientset:", err)
		return nil, err
	}

	return &KubernetesClient{clientset: clientset}, nil
}

func (c *KubernetesClient) GetNodes() ([]string, error) {
	nodes, err := c.clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	nodesList := make([]string, 0, len(nodes.Items))
	for _, node := range nodes.Items {
		nodesList = append(nodesList, node.Name)
	}
	return nodesList, nil
}
