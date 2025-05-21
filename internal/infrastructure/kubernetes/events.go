package kubernetes

import (
	"context"
	"main/internal/domain/events"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c KubernetesClient) WatchEvents(ctx context.Context, namespace string) (chan events.Event, error) {
	c.logger.Info("Starting event watch in namespace", namespace)
	eventChan := make(chan events.Event)
	watcher, err := c.clientset.CoreV1().Events(namespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		c.logger.Errorf("Failed to create watcher: %v", err)
		return nil, err
	}

	go func() {
		defer close(eventChan)
		for {
			select {
			case watchEvent, ok := <-watcher.ResultChan():
				if !ok {
					c.logger.Info("Watcher channel closed, restarting...")
					time.Sleep(1 * time.Second)
					watcher, err = c.clientset.CoreV1().Events(namespace).Watch(ctx, metav1.ListOptions{})
					if err != nil {
						c.logger.Errorf("Failed to restart watcher: %v", err)
						time.Sleep(5 * time.Second)
						continue
					}
					continue
				}
				c.logger.Info("Received event from watcher")
				event, ok := watchEvent.Object.(*corev1.Event)
				if !ok {
					c.logger.Errorf("Unexpected type for event: %T", watchEvent.Object)
					continue
				}
				involvedObject := event.InvolvedObject.Kind + "/" + event.InvolvedObject.Name
				domainEvent := events.Event{
					Namespace:      event.Namespace,
					Name:           event.Name,
					Reason:         event.Reason,
					Message:        event.Message,
					Type:           event.Type,
					FirstTimestamp: event.FirstTimestamp.Time,
					InvolvedObject: involvedObject,
					Count:          event.Count,
					ID:             string(event.UID),
				}
				if event.DeletionTimestamp != nil {
					domainEvent.LastTimestamp = event.DeletionTimestamp.Time
				}
				c.logger.Info("Sending event to channel", domainEvent.Name)
				eventChan <- domainEvent
			case <-ctx.Done():
				c.logger.Info("Context canceled, stopping watcher")
				watcher.Stop()
				return
			}
		}
	}()

	return eventChan, nil
}
