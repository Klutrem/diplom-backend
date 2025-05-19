package database

import (
	"fmt"
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
	query := `
		SELECT * FROM ` + repo.table + `
	`
	args := []interface{}{limit}
	var conditions []string
	if namespace != "" {
		conditions = append(conditions, `namespace = $`+fmt.Sprint(len(args)+1))
		args = append(args, namespace)
	}

	if len(conditions) > 0 {
		query += ` WHERE ` + conditions[0]
	}

	query += `
		ORDER BY last_timestamp DESC
		LIMIT $1
	`

	events := make([]events.Event, 0)
	err := repo.database.Select(&events, query, args...)
	if err != nil {
		return nil, err
	}
	return events, nil
}
