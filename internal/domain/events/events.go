package events

import (
	"context"
	"time"
)

type Event struct {
	ID             string    `json:"id" db:"id"`
	Namespace      string    `json:"namespace" db:"namespace"`
	Name           string    `json:"name" db:"name"`
	Reason         string    `json:"reason" db:"reason"`
	Message        string    `json:"message" db:"message"`
	Type           string    `json:"type" db:"type"`
	InvolvedObject string    `json:"involved_object" db:"involved_object"`
	FirstTimestamp time.Time `json:"first_timestamp" db:"first_timestamp"`
	LastTimestamp  time.Time `json:"last_timestamp" db:"last_timestamp"`
	Count          int32     `json:"count" db:"count"`
}

type InvolvedObject struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type EventsKubernetesClient interface {
	WatchEvents(ctx context.Context, namespace string) (chan Event, error)
}
