package kubernetes

import (
	"context"
	"main/internal/domain/events"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c KubernetesClient) WatchEvents(context context.Context, namespace string) (chan events.Event, error) {
	eventChan := make(chan events.Event)
	watcher, err := c.clientset.CoreV1().Events(namespace).Watch(context, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	go func() {
		for watchEvent := range watcher.ResultChan() {
			event, ok := watchEvent.Object.(*corev1.Event)
			if !ok {
				continue
			}
			involvedObject := event.InvolvedObject.Kind + "/" + event.InvolvedObject.Name
			domainEvent := events.Event{
				Namespace:      event.Namespace,
				Name:           event.Name,
				Reason:         event.Reason,
				Message:        event.Message,
				Type:           event.Type,
				Timestamp:      event.CreationTimestamp.Time,
				InvolvedObject: involvedObject,
			}
			eventChan <- domainEvent
		}
	}()

	return eventChan, nil
}
