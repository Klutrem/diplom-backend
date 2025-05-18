package kubernetes

import (
	"context"
	"fmt"
	"main/internal/domain/metrics"
	"time"

	"github.com/prometheus/common/model"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

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
	nodes, err := c.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nodeMetrics []v1beta1.NodeMetrics

	for _, node := range nodes.Items {
		// Получаем метрики CPU
		cpuQuery := fmt.Sprintf(`sum(rate(node_cpu_seconds_total{node="%s", mode!="idle"}[5m])) * 1000`, node.Name)
		cpuValue, err := c.prometheusClient.GetMetricValue(cpuQuery)
		if err != nil {
			c.logger.Errorf("failed to get CPU usage for node %s: %v", node.Name, err)
			continue
		}
		cpuUsage := int64(0)
		if vector, ok := cpuValue.(model.Vector); ok && len(vector) > 0 {
			cpuUsage = int64(vector[0].Value)
		}

		// Получаем метрики памяти
		memoryQuery := fmt.Sprintf(`node_memory_MemTotal_bytes{node="%s"} - node_memory_MemAvailable_bytes{node="%s"}`, node.Name, node.Name)
		memoryValue, err := c.prometheusClient.GetMetricValue(memoryQuery)
		if err != nil {
			c.logger.Errorf("failed to get memory usage for node %s: %v", node.Name, err)
			continue
		}
		memoryUsage := int64(0)
		if vector, ok := memoryValue.(model.Vector); ok && len(vector) > 0 {
			memoryUsage = int64(vector[0].Value)
		}

		// Создаем NodeMetrics
		nodeMetric := v1beta1.NodeMetrics{
			ObjectMeta: metav1.ObjectMeta{
				Name: node.Name,
			},
			Usage: corev1.ResourceList{
				corev1.ResourceCPU:    *resource.NewMilliQuantity(cpuUsage, resource.DecimalSI),
				corev1.ResourceMemory: *resource.NewQuantity(memoryUsage, resource.BinarySI),
			},
			Timestamp: metav1.Time{Time: time.Now()},
			Window:    metav1.Duration{Duration: 5 * time.Minute},
		}
		nodeMetrics = append(nodeMetrics, nodeMetric)
	}

	return nodeMetrics, nil
}

func NodeFromMetrics(nodeMetrics []v1beta1.NodeMetrics, nodes []corev1.Node) []metrics.NodeMetrics {
	result := make([]metrics.NodeMetrics, 0, len(nodeMetrics))
	for _, nodeMetric := range nodeMetrics {
		var node *corev1.Node
		for _, n := range nodes {
			if n.Name == nodeMetric.Name {
				node = &n
				break
			}
		}
		if node == nil {
			continue
		}

		cpuUsage := nodeMetric.Usage.Cpu().MilliValue() // mCores
		cpuCapacity := node.Status.Capacity.Cpu().MilliValue()
		cpuUsagePercent := float64(cpuUsage) / float64(cpuCapacity) * 100

		memoryUsage := nodeMetric.Usage.Memory().Value() / 1024 / 1024 // MiB
		memoryCapacity := node.Status.Capacity.Memory().Value() / 1024 / 1024
		memoryUsagePercent := float64(memoryUsage) / float64(memoryCapacity) * 100

		roles := getNodeRoles(node)

		result = append(result, metrics.NodeMetrics{
			NodeName:              nodeMetric.Name,
			CPUUsage:              fmt.Sprintf("%dm", cpuUsage),
			CpuCapacity:           fmt.Sprintf("%dm", cpuCapacity),
			CpuUsagePercentage:    fmt.Sprintf("%f", cpuUsagePercent),
			MemoryUsage:           fmt.Sprintf("%dMi", memoryUsage),
			MemoryCapacity:        fmt.Sprintf("%dMi", memoryCapacity),
			MemoryUsagePercentage: fmt.Sprintf("%f", memoryUsagePercent),
			Roles:                 roles,
			Status:                getNodeStatus(node),
		})
	}
	return result
}

func getNodeRoles(node *corev1.Node) []string {
	roles := []string{}
	if _, exists := node.Labels["node-role.kubernetes.io/control-plane"]; exists {
		roles = append(roles, "control-plane")
	}
	if _, exists := node.Labels["node-role.kubernetes.io/worker"]; exists {
		roles = append(roles, "worker")
	}
	return roles
}

func getNodeStatus(node *corev1.Node) string {
	for _, condition := range node.Status.Conditions {
		if condition.Type == "Ready" && condition.Status == "True" {
			return "Ready"
		}
	}
	return "NotReady"
}
