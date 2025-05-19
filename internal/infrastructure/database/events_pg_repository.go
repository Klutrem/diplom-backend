package database

import (
	"fmt"
	"log"
	"main/internal/domain/events"
	"main/pkg"
)

func NewEventPGRepository(database pkg.Database) events.EventRepository {
	return EventPGRepo{
		database: database,
		table:    "events",
	}
}

type EventPGRepo struct {
	database pkg.Database
	table    string
}

func (repo EventPGRepo) SaveEvent(event events.Event) error {
	query := `
		INSERT INTO ` + repo.table + ` (
			id, namespace, name, reason, message, type, 
			involved_object, first_timestamp, last_timestamp, count
		)
		VALUES (
			:id, :namespace, :name, :reason, :message, :type,
			:involved_object, :first_timestamp, :last_timestamp, :count
		)
	`
	_, err := repo.database.NamedExec(query, event)
	if err != nil {
		return err
	}
	return nil
}

func (repo EventPGRepo) GetEvents(namespace string, limit int) ([]events.Event, error) {
	// Проверка limit
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be positive, got %d", limit)
	}

	query := `SELECT * FROM ` + repo.table + ` WHERE namespace = $1 ORDER BY last_timestamp DESC LIMIT $2`

	// Выполнение запроса
	events := make([]events.Event, 0)
	err := repo.database.Select(&events, query, namespace, limit)
	if err != nil {
		// Логирование запроса и аргументов для отладки
		log.Printf("Failed query: %s, Args: [%s, %d]", query, namespace, limit)
		return nil, fmt.Errorf("failed to query events: %w", err)
	}

	return events, nil
}
