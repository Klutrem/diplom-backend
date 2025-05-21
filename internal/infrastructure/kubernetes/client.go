package kubernetes

import (
	"main/internal/domain/events"
	"main/internal/infrastructure/prometheus"
	"main/pkg"
	"os"
	"path/filepath"

	"go.uber.org/fx"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

type KubernetesClient struct {
	clientset        *kubernetes.Clientset
	logger           pkg.Logger
	metricsClient    *versioned.Clientset
	prometheusClient prometheus.PrometheusClient
}

var Module = fx.Module("kubernetes", fx.Provide(NewKubernetesClient), fx.Provide(func(kc *KubernetesClient) events.EventsKubernetesClient { return kc }))

func NewKubernetesClient(logger pkg.Logger, prometheusClient prometheus.PrometheusClient) (*KubernetesClient, error) {
	var config *rest.Config
	var err error

	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		config, err = rest.InClusterConfig()
	} else {
		var kubeconfigPath string
		if os.Getenv("KUBECONFIG_PATH") != "" {
			kubeconfigPath = os.Getenv("KUBECONFIG_PATH")
		} else {
			kubeconfigPath = filepath.Join(homedir.HomeDir(), ".kube", "config")
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	metricsClient, err := versioned.NewForConfig(config)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &KubernetesClient{clientset: clientset, metricsClient: metricsClient, logger: logger, prometheusClient: prometheusClient}, nil
}
