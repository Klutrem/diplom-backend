package kubernetes

import (
	"context"
	"fmt"
	"main/internal/domain/metrics"
	"time"

	"github.com/prometheus/common/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *KubernetesClient) GetPods(namespace string) ([]metrics.PodMetrics, error) {
	pods, err := c.clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make([]metrics.PodMetrics, 0, len(pods.Items))
	for _, pod := range pods.Items {
		cpuUsage, err := c.getPodCPUUsage(namespace, pod.Name)
		if err != nil {
			c.logger.Errorf("failed to get CPU usage for pod %s: %v", pod.Name, err)
			continue
		}
		memoryUsage, err := c.getPodMemoryUsage(namespace, pod.Name)
		if err != nil {
			c.logger.Errorf("failed to get memory usage for pod %s: %v", pod.Name, err)
			continue
		}

		// Получаем ресурсы из pod.Spec
		cpuRequest := pod.Spec.Containers[0].Resources.Requests.Cpu().MilliValue()
		cpuLimit := pod.Spec.Containers[0].Resources.Limits.Cpu().MilliValue()
		memoryRequest := pod.Spec.Containers[0].Resources.Requests.Memory().Value() / 1024 / 1024 // MiB
		memoryLimit := pod.Spec.Containers[0].Resources.Limits.Memory().Value() / 1024 / 1024     // MiB

		var cpuUsagePercent *float64
		if cpuLimit > 0 {
			percent := float64(cpuUsage) / float64(cpuLimit) * 100
			cpuUsagePercent = &percent
		}

		var memoryUsagePercent *float64
		if memoryLimit > 0 {
			percent := float64(memoryUsage) / float64(memoryLimit) * 100
			memoryUsagePercent = &percent
		}

		result = append(result, metrics.PodMetrics{
			PodName:            pod.Name,
			Namespace:          pod.Namespace,
			NodeName:           pod.Spec.NodeName,
			Status:             string(pod.Status.Phase),
			StartTime:          pod.Status.StartTime.Format(time.RFC3339),
			CPUUsage:           cpuUsage,
			CPUUsagePercent:    cpuUsagePercent,
			CPUUsageLimit:      cpuLimit,
			CPUUsageRequest:    cpuRequest,
			MemoryUsage:        memoryUsage,
			MemoryUsagePercent: memoryUsagePercent,
			MemoryUsageLimit:   memoryLimit,
			MemoryUsageRequest: memoryRequest,
			RestartCount:       pod.Status.ContainerStatuses[0].RestartCount,
		})
	}
	return result, nil
}

func (c *KubernetesClient) getPodCPUUsage(namespace, podName string) (int64, error) {
	query := fmt.Sprintf(`sum(rate(container_cpu_usage_seconds_total{namespace="%s", pod="%s"}[5m])) * 1000`, namespace, podName)
	value, err := c.prometheusClient.GetMetricValue(query)
	if err != nil {
		return 0, err
	}
	if vector, ok := value.(model.Vector); ok && len(vector) > 0 {
		return int64(vector[0].Value), nil
	}
	return 0, fmt.Errorf("no CPU usage data for pod %s", podName)
}

func (c *KubernetesClient) getPodMemoryUsage(namespace, podName string) (int64, error) {
	query := fmt.Sprintf(`container_memory_working_set_bytes{namespace="%s", pod="%s"}`, namespace, podName)
	value, err := c.prometheusClient.GetMetricValue(query)
	if err != nil {
		return 0, err
	}
	if vector, ok := value.(model.Vector); ok && len(vector) > 0 {
		return int64(vector[0].Value) / 1024 / 1024, nil // Переводим в MiB
	}
	return 0, fmt.Errorf("no memory usage data for pod %s", podName)
}
