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
	nodes, err := c.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Try to get metrics, but don't fail if we can't
	nodeMetrics, err := c.getNodeMetrics()
	if err != nil {
		c.logger.Errorf("Failed to get node metrics: %v", err)
		// Return nodes without metrics
		return NodeFromMetrics([]v1beta1.NodeMetrics{}, nodes.Items), nil
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
	result := make([]metrics.NodeMetrics, 0, len(nodes))

	// Create a map for quick lookup of metrics by node name
	metricsMap := make(map[string]v1beta1.NodeMetrics)
	for _, metric := range nodeMetrics {
		metricsMap[metric.Name] = metric
	}

	for _, node := range nodes {
		nodeMetric, hasMetrics := metricsMap[node.Name]

		// Initialize with zero values
		cpuUsage := int64(0)
		cpuCapacity := node.Status.Capacity.Cpu().MilliValue()
		cpuUsagePercent := 0.0
		memoryUsage := int64(0)
		memoryCapacity := node.Status.Capacity.Memory().Value() / 1024 / 1024
		memoryUsagePercent := 0.0

		// If we have metrics, use them
		if hasMetrics {
			cpuUsage = nodeMetric.Usage.Cpu().MilliValue()
			if cpuCapacity > 0 {
				cpuUsagePercent = float64(cpuUsage) / float64(cpuCapacity) * 100
			}
			memoryUsage = nodeMetric.Usage.Memory().Value() / 1024 / 1024
			if memoryCapacity > 0 {
				memoryUsagePercent = float64(memoryUsage) / float64(memoryCapacity) * 100
			}
		}

		roles := getNodeRoles(&node)

		result = append(result, metrics.NodeMetrics{
			NodeName:              node.Name,
			CPUUsage:              fmt.Sprintf("%dm", cpuUsage),
			CpuCapacity:           fmt.Sprintf("%dm", cpuCapacity),
			CpuUsagePercentage:    fmt.Sprintf("%f", cpuUsagePercent),
			MemoryUsage:           fmt.Sprintf("%dMi", memoryUsage),
			MemoryCapacity:        fmt.Sprintf("%dMi", memoryCapacity),
			MemoryUsagePercentage: fmt.Sprintf("%f", memoryUsagePercent),
			Roles:                 roles,
			Status:                getNodeStatus(&node),
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

type NodeHistoricalMetrics struct {
	CPUUsage    []MetricPoint `json:"cpu_usage"`
	MemoryUsage []MetricPoint `json:"memory_usage"`
}

func (c *KubernetesClient) GetNodeHistoricalMetrics(nodeName string, start, end time.Time, step time.Duration) (*NodeHistoricalMetrics, error) {
	cpuQuery := fmt.Sprintf(`sum(rate(node_cpu_seconds_total{mode="user",node="%s"}[5m])) * 1000`, nodeName)
	cpuValue, err := c.prometheusClient.GetMetricHistory(cpuQuery, start, end, step)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU usage history: %v", err)
	}

	memoryQuery := fmt.Sprintf(`node_memory_MemTotal_bytes{node="%s"} - node_memory_MemAvailable_bytes{node="%s"}`, nodeName, nodeName)
	memoryValue, err := c.prometheusClient.GetMetricHistory(memoryQuery, start, end, step)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory usage history: %v", err)
	}

	metrics := &NodeHistoricalMetrics{
		CPUUsage:    make([]MetricPoint, 0),
		MemoryUsage: make([]MetricPoint, 0),
	}

	if matrix, ok := cpuValue.(model.Matrix); ok {
		for _, series := range matrix {
			for _, point := range series.Values {
				metrics.CPUUsage = append(metrics.CPUUsage, MetricPoint{
					Timestamp: point.Timestamp.Time(),
					Value:     float64(point.Value),
				})
			}
		}
	}

	if matrix, ok := memoryValue.(model.Matrix); ok {
		for _, series := range matrix {
			for _, point := range series.Values {
				metrics.MemoryUsage = append(metrics.MemoryUsage, MetricPoint{
					Timestamp: point.Timestamp.Time(),
					Value:     float64(point.Value),
				})
			}
		}
	}

	return metrics, nil
}
