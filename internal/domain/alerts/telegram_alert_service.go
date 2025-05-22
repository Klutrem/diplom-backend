package alerts

import (
	"fmt"
	"main/pkg"
	"strconv"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/fx"
)

var Module = fx.Module("alerts",
	fx.Provide(NewTelegramAlertService),
)

type telegramAlertService struct {
	logger     pkg.Logger
	repository TelegramAlertRepository
	botCache   map[string]*tgbotapi.BotAPI
	mu         sync.RWMutex
}

func NewTelegramAlertService(logger pkg.Logger, repository TelegramAlertRepository) TelegramAlertService {
	return &telegramAlertService{
		logger:     logger,
		repository: repository,
		botCache:   make(map[string]*tgbotapi.BotAPI),
	}
}

func (s *telegramAlertService) CreateAlert(alert TelegramAlert) error {
	return s.repository.CreateAlert(alert)
}

func (s *telegramAlertService) UpdateAlert(alert TelegramAlert) error {
	return s.repository.UpdateAlert(alert)
}

func (s *telegramAlertService) DeleteAlert(id int64) error {
	return s.repository.DeleteAlert(id)
}

func (s *telegramAlertService) GetAlert(id int64) (*TelegramAlert, error) {
	return s.repository.GetAlert(id)
}

func (s *telegramAlertService) GetAlertsByNamespace(namespace string) ([]TelegramAlert, error) {
	return s.repository.GetAlertsByNamespace(namespace)
}

func (s *telegramAlertService) GetAllAlerts() ([]TelegramAlert, error) {
	return s.repository.GetAllAlerts()
}

func (s *telegramAlertService) getBot(token string) (*tgbotapi.BotAPI, error) {
	s.mu.RLock()
	bot, exists := s.botCache[token]
	s.mu.RUnlock()

	if exists {
		return bot, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if bot, exists = s.botCache[token]; exists {
		return bot, nil
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot instance: %w", err)
	}

	s.botCache[token] = bot
	return bot, nil
}

func (s *telegramAlertService) SendAlert(alert TelegramAlert, message string) error {
	bot, err := s.getBot(alert.BotToken)
	if err != nil {
		return fmt.Errorf("failed to get bot instance: %w", err)
	}

	chatID, err := strconv.ParseInt(alert.ChatID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid chat ID: %w", err)
	}

	msg := tgbotapi.NewMessage(chatID, message)
	if alert.ThreadID != nil {
		msg.ReplyToMessageID = *alert.ThreadID
	}

	_, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}

	return nil
}
