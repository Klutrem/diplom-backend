package kubernetes

import (
	"context"
	"fmt"
	"main/internal/domain/metrics"
	"os"
	"path/filepath"

	"go.uber.org/fx"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

type KubernetesClient struct {
	clientset     *kubernetes.Clientset
	metricsClient *versioned.Clientset
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
	metricsClient, err := versioned.NewForConfig(config)
	if err != nil {
		fmt.Println("Error creating metrics client:", err)
		return nil, err
	}

	return &KubernetesClient{clientset: clientset, metricsClient: metricsClient}, nil
}

func (c *KubernetesClient) GetNodes() ([]metrics.NodeMetrics, error) {
	nodeMetrics, err := c.getNodeMetrics()
	if err != nil {
		return nil, err
	}

	nodes, err := c.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return NodeFromMetrics(nodeMetrics, nodes.Items), nil
}

func (c *KubernetesClient) getNodeMetrics() ([]v1beta1.NodeMetrics, error) {
	nodeMetricsList, err := c.metricsClient.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("Error getting node metrics:", err)
		return nil, err
	}
	return nodeMetricsList.Items, nil
}
