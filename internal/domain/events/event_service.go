package events

import (
	"context"
	"main/pkg"

	"go.uber.org/fx"
)

type EventRepository interface {
	SaveEvent(event Event) error
	GetEvents(namespace string, eventType string, limit int) ([]Event, error)
}

type WatchedNamespaceRepository interface {
	AddNamespace(namespace string) error
	RemoveNamespace(namespace string) error
	GetAllNamespaces() ([]string, error)
}

type EventService struct {
	logger              pkg.Logger
	k8sClient           EventsKubernetesClient
	repository          EventRepository
	namespaceRepository WatchedNamespaceRepository
	watchedNamespaces   map[string]context.CancelFunc
}

var Module = fx.Module("events",
	fx.Provide(NewEventService),
)

func NewEventService(logger pkg.Logger, k8sClient EventsKubernetesClient, repo EventRepository, namespaceRepo WatchedNamespaceRepository) EventService {
	svc := EventService{
		logger:              logger,
		k8sClient:           k8sClient,
		repository:          repo,
		namespaceRepository: namespaceRepo,
		watchedNamespaces:   make(map[string]context.CancelFunc),
	}

	namespaces, err := namespaceRepo.GetAllNamespaces()
	if err != nil {
		logger.Errorf("Failed to get watched namespaces: %v", err)
	} else {
		for _, namespace := range namespaces {
			if err := svc.StartWatching(context.Background(), namespace); err != nil {
				logger.Errorf("Failed to start watching namespace %s: %v", namespace, err)
			}
		}
	}

	return svc
}

func (s *EventService) StartWatching(ctx context.Context, namespace string) error {
	s.logger.Info("Starting event watching in namespace", namespace)

	watcherCtx, cancel := context.WithCancel(ctx)
	s.watchedNamespaces[namespace] = cancel

	eventChan, err := s.k8sClient.WatchEvents(watcherCtx, namespace)
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

func (s *EventService) StopWatching(namespace string) {
	if cancel, exists := s.watchedNamespaces[namespace]; exists {
		cancel()
		delete(s.watchedNamespaces, namespace)
	}
}

func (s *EventService) AddNamespaceToWatch(ctx context.Context, namespace string) error {
	if err := s.namespaceRepository.AddNamespace(namespace); err != nil {
		return err
	}
	return s.StartWatching(ctx, namespace)
}

func (s *EventService) RemoveNamespaceFromWatch(namespace string) error {
	if err := s.namespaceRepository.RemoveNamespace(namespace); err != nil {
		return err
	}
	s.StopWatching(namespace)
	return nil
}

func (s *EventService) GetWatchedNamespaces() ([]string, error) {
	return s.namespaceRepository.GetAllNamespaces()
}

func (s *EventService) GetEvents(namespace string, eventType string, limit int) ([]Event, error) {
	return s.repository.GetEvents(namespace, eventType, limit)
}
