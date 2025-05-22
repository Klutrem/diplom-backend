package database

import (
	"fmt"
	"main/internal/domain/alerts"
	"main/pkg"
)

type TelegramAlertPGRepo struct {
	database pkg.Database
	table    string
}

func NewTelegramAlertPGRepository(database pkg.Database) alerts.TelegramAlertRepository {
	return TelegramAlertPGRepo{
		database: database,
		table:    "telegram_alerts",
	}
}

func (repo TelegramAlertPGRepo) CreateAlert(alert alerts.TelegramAlert) error {
	query := `
		INSERT INTO ` + repo.table + ` (
			bot_token, chat_id, thread_id, alert_type, namespace
		)
		VALUES (
			:bot_token, :chat_id, :thread_id, :alert_type, :namespace
		)
		RETURNING id
	`
	rows, err := repo.database.NamedQuery(query, alert)
	if err != nil {
		return fmt.Errorf("failed to create alert: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&alert.ID); err != nil {
			return fmt.Errorf("failed to get created alert ID: %w", err)
		}
	}
	return nil
}

func (repo TelegramAlertPGRepo) UpdateAlert(alert alerts.TelegramAlert) error {
	query := `
		UPDATE ` + repo.table + `
		SET bot_token = :bot_token,
			chat_id = :chat_id,
			thread_id = :thread_id,
			alert_type = :alert_type,
			namespace = :namespace
		WHERE id = :id
	`
	_, err := repo.database.NamedExec(query, alert)
	if err != nil {
		return fmt.Errorf("failed to update alert: %w", err)
	}
	return nil
}

func (repo TelegramAlertPGRepo) DeleteAlert(id int64) error {
	query := `
		DELETE FROM ` + repo.table + `
		WHERE id = $1
	`
	_, err := repo.database.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete alert: %w", err)
	}
	return nil
}

func (repo TelegramAlertPGRepo) GetAlert(id int64) (*alerts.TelegramAlert, error) {
	query := `
		SELECT * FROM ` + repo.table + `
		WHERE id = $1
	`
	var alert alerts.TelegramAlert
	err := repo.database.Get(&alert, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}
	return &alert, nil
}

func (repo TelegramAlertPGRepo) GetAlertsByNamespace(namespace string) ([]alerts.TelegramAlert, error) {
	query := `
		SELECT * FROM ` + repo.table + `
		WHERE namespace = $1
		ORDER BY created_at DESC
	`
	var alerts []alerts.TelegramAlert
	err := repo.database.Select(&alerts, query, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts by namespace: %w", err)
	}
	return alerts, nil
}

func (repo TelegramAlertPGRepo) GetAllAlerts() ([]alerts.TelegramAlert, error) {
	query := `
		SELECT * FROM ` + repo.table + `
		ORDER BY created_at DESC
	`
	var alerts []alerts.TelegramAlert
	err := repo.database.Select(&alerts, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all alerts: %w", err)
	}
	return alerts, nil
}
