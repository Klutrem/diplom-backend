package kubernetes

import (
	"fmt"
	"main/internal/domain/metrics"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

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
