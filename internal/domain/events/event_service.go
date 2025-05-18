package events

import (
	"context"
	"main/pkg"

	"go.uber.org/fx"
)

type EventRepository interface {
	SaveEvent(ctx context.Context, event Event) error
	GetEvents(ctx context.Context, namespace string, limit int) ([]Event, error)
}

type EventService struct {
	logger     pkg.Logger
	k8sClient  EventsKubernetesClient
	repository EventRepository
}

var Module = fx.Module("events",
	fx.Provide(NewEventService),
)

func NewEventService(logger pkg.Logger, k8sClient EventsKubernetesClient, repo EventRepository) *EventService {
	return &EventService{
		logger:     logger,
		k8sClient:  k8sClient,
		repository: repo,
	}
}

func (s *EventService) StartWatching(ctx context.Context, namespace string) error {
	eventChan, err := s.k8sClient.WatchEvents(ctx, namespace)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-eventChan:
				if err := s.repository.SaveEvent(ctx, event); err != nil {
					s.logger.Errorf("failed to save event %s: %v", event.Name, err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (s *EventService) GetEvents(ctx context.Context, namespace string, limit int) ([]Event, error) {
	return s.repository.GetEvents(ctx, namespace, limit)
}
