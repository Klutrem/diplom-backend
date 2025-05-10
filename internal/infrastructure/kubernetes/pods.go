package kubernetes

import (
	"context"
	"main/internal/domain/metrics"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func (c *KubernetesClient) GetPods(namespace string) ([]metrics.PodMetrics, error) {
	podMetrics, err := c.GetPodMetrics(namespace)
	if err != nil {
		return nil, err
	}

	pods, err := c.clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return podFromMetrics(podMetrics, pods.Items), nil
}

func (c *KubernetesClient) GetPodMetrics(namespace string) ([]v1beta1.PodMetrics, error) {
	podMetricsList, err := c.metricsClient.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return podMetricsList.Items, nil
}

func podFromMetrics(podMetrics []v1beta1.PodMetrics, pods []corev1.Pod) []metrics.PodMetrics {
	result := make([]metrics.PodMetrics, 0, len(podMetrics))
	for _, pod := range pods {
		var podMetric *v1beta1.PodMetrics
		for _, p := range podMetrics {
			if p.Name == pod.Name {
				podMetric = &p
				break
			}
		}
		if podMetric == nil {
			result = append(result, metrics.PodMetrics{
				PodName:      pod.Name,
				Namespace:    pod.Namespace,
				NodeName:     pod.Spec.NodeName,
				Status:       string(pod.Status.Phase),
				StartTime:    pod.Status.StartTime.Format(time.RFC3339),
				RestartCount: pod.Status.ContainerStatuses[0].RestartCount,
			})
			continue
		}

		var cpuUsagePercent float64
		var memoryUsagePercent float64

		cpuUsage := podMetric.Containers[0].Usage.Cpu().MilliValue() // mCores
		cpuCapacity := pod.Status.ContainerStatuses[0].AllocatedResources.Cpu().MilliValue()
		if cpuCapacity > 0 {
			cpuUsagePercent = float64(cpuUsage) / float64(cpuCapacity) * 100
		}

		memoryUsage := podMetric.Containers[0].Usage.Memory().Value() / 1024 / 1024 // MiB
		memoryCapacity := pod.Status.ContainerStatuses[0].AllocatedResources.Memory().Value() / 1024 / 1024
		if memoryCapacity > 0 {
			memoryUsagePercent = float64(memoryUsage) / float64(memoryCapacity) * 100
		}
		result = append(result, metrics.PodMetrics{
			PodName:            pod.Name,
			Namespace:          pod.Namespace,
			NodeName:           pod.Spec.NodeName,
			Status:             string(pod.Status.Phase),
			StartTime:          pod.Status.StartTime.Format(time.RFC3339),
			CPUUsage:           cpuUsage,
			CPUUsagePercent:    &cpuUsagePercent,
			CPUUsageLimit:      cpuCapacity,
			CPUUsageRequest:    pod.Status.ContainerStatuses[0].AllocatedResources.Cpu().MilliValue(),
			MemoryUsage:        memoryUsage,
			MemoryUsagePercent: &memoryUsagePercent,
			MemoryUsageLimit:   memoryCapacity,
			MemoryUsageRequest: pod.Status.ContainerStatuses[0].AllocatedResources.Memory().Value() / 1024 / 1024,
			RestartCount:       pod.Status.ContainerStatuses[0].RestartCount,
		})
	}
	return result
}
