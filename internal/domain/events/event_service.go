package events

import (
	"context"
	"main/pkg"

	"go.uber.org/fx"
)

type EventRepository interface {
	SaveEvent(event Event) error
	GetEvents(namespace string, limit int) ([]Event, error)
}

type EventService struct {
	logger     pkg.Logger
	k8sClient  EventsKubernetesClient
	repository EventRepository
}

var Module = fx.Module("events",
	fx.Provide(NewEventService),
)

func NewEventService(logger pkg.Logger, k8sClient EventsKubernetesClient, repo EventRepository) EventService {
	svc := EventService{
		logger:     logger,
		k8sClient:  k8sClient,
		repository: repo,
	}
	go svc.StartWatching(context.Background(), "default")
	return svc
}

func (s *EventService) StartWatching(ctx context.Context, namespace string) error {
	s.logger.Info("Starting event watching in namespace", namespace)
	eventChan, err := s.k8sClient.WatchEvents(ctx, namespace)
	if err != nil {
		s.logger.Errorf("Failed to start watching events: %v", err)
		return err
	}

	go func() {
		for event := range eventChan {
			s.logger.Info("Received event in namespace", namespace, "event name:", event.Name)
			if err := s.repository.SaveEvent(event); err != nil {
				s.logger.Errorf("Failed to save event %s: %v", event.Name, err)
			}
		}
		s.logger.Info("Event channel closed on namespace", namespace)
	}()

	return nil
}
func (s *EventService) GetEvents(namespace string, limit int) ([]Event, error) {
	return s.repository.GetEvents(namespace, limit)
}
