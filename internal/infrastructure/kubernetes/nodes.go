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

		result = append(result, metrics.NodeMetrics{
			NodeName:              nodeMetric.Name,
			CPUUsage:              fmt.Sprintf("%dm", cpuUsage),
			CpuCapacity:           fmt.Sprintf("%dm", cpuCapacity),
			CpuUsagePercentage:    fmt.Sprintf("%f", cpuUsagePercent),
			MemoryUsage:           fmt.Sprintf("%dMi", memoryUsage),
			MemoryCapacity:        fmt.Sprintf("%dMi", memoryCapacity),
			MemoryUsagePercentage: fmt.Sprintf("%f", memoryUsagePercent),
		})
	}
	return result
}
