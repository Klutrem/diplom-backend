package events

import (
	"context"
	"time"
)

type Event struct {
	ID             int64     `json:"id"`
	Namespace      string    `json:"namespace"`
	Name           string    `json:"name"`
	Reason         string    `json:"reason"`
	Message        string    `json:"message"`
	Type           string    `json:"type"`
	Timestamp      time.Time `json:"timestamp"`
	InvolvedObject string    `json:"involved_object"`
}

type EventsKubernetesClient interface {
	WatchEvents(ctx context.Context, namespace string) (chan Event, error)
}
