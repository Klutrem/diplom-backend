package alerts

import "time"

type AlertType string

const (
	AlertTypeAll     AlertType = "all"
	AlertTypeNormal  AlertType = "normal"
	AlertTypeWarning AlertType = "warning"
)

type TelegramAlert struct {
	ID        int64     `json:"id" db:"id"`
	BotToken  string    `json:"bot_token" db:"bot_token"`
	ChatID    string    `json:"chat_id" db:"chat_id"`
	ThreadID  *int      `json:"thread_id,omitempty" db:"thread_id"`
	AlertType AlertType `json:"alert_type" db:"alert_type"`
	Namespace string    `json:"namespace" db:"namespace"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}

// TelegramAlertResponse is used for API responses, excluding sensitive data
type TelegramAlertResponse struct {
	ID        int64     `json:"id"`
	ChatID    string    `json:"chat_id"`
	ThreadID  *int      `json:"thread_id,omitempty"`
	AlertType AlertType `json:"alert_type"`
	Namespace string    `json:"namespace"`
	CreatedAt time.Time `json:"created_at"`
}

type TelegramAlertRepository interface {
	CreateAlert(alert TelegramAlert) error
	UpdateAlert(alert TelegramAlert) error
	DeleteAlert(id int64) error
	GetAlert(id int64) (*TelegramAlert, error)
	GetAlertsByNamespace(namespace string) ([]TelegramAlert, error)
	GetAllAlerts() ([]TelegramAlert, error)
}

type TelegramAlertService interface {
	CreateAlert(alert TelegramAlert) error
	UpdateAlert(alert TelegramAlert) error
	DeleteAlert(id int64) error
	GetAlert(id int64) (*TelegramAlert, error)
	GetAlertsByNamespace(namespace string) ([]TelegramAlert, error)
	GetAllAlerts() ([]TelegramAlert, error)
	SendAlert(alert TelegramAlert, message string) error
}
