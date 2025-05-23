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

func (repo EventPGRepo) GetEvents(namespace, eventType string, limit int) ([]events.Event, error) {
	// Проверка limit
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be positive, got %d", limit)
	}

	// Формирование запроса
	query := `SELECT * FROM ` + repo.table + ` WHERE namespace = $1`
	args := []any{namespace}

	// Добавление фильтра по type если указан
	if eventType != "" {
		query += ` AND type = $2`
		args = append(args, eventType)
	}

	query += ` ORDER BY last_timestamp DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1)
	args = append(args, limit)

	// Выполнение запроса
	events := make([]events.Event, 0)
	err := repo.database.Select(&events, query, args...)
	if err != nil {
		log.Printf("Query: %s, Args: %v", query, args)
		return nil, fmt.Errorf("failed to query events: %w", err)
	}

	return events, nil
}

type WatchedNamespacePGRepo struct {
	database pkg.Database
	table    string
}

func NewWatchedNamespacePGRepository(database pkg.Database) events.WatchedNamespaceRepository {
	return WatchedNamespacePGRepo{
		database: database,
		table:    "watched_namespaces",
	}
}

func (repo WatchedNamespacePGRepo) AddNamespace(namespace string) error {
	query := `
		INSERT INTO ` + repo.table + ` (namespace)
		VALUES ($1)
		ON CONFLICT (namespace) DO NOTHING
	`
	_, err := repo.database.Exec(query, namespace)
	return err
}

func (repo WatchedNamespacePGRepo) RemoveNamespace(namespace string) error {
	query := `
		DELETE FROM ` + repo.table + `
		WHERE namespace = $1
	`
	_, err := repo.database.Exec(query, namespace)
	return err
}

func (repo WatchedNamespacePGRepo) GetAllNamespaces() ([]string, error) {
	query := `
		SELECT namespace FROM ` + repo.table + `
		ORDER BY created_at DESC
	`
	namespaces := make([]string, 0)
	err := repo.database.Select(&namespaces, query)
	return namespaces, err
}
