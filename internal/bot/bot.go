package bot

import (
	"fmt"
	"log"

	"github.com/SobolevTim/finance_bot/internal/database"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

// Bot структура для хранения клиента бота
type Bot struct {
	Client *telego.Bot
}

// NewBot создает новый экземпляр бота
func NewBot(token string) (*Bot, error) {
	bot, err := telego.NewBot(token)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запуске бота: %v", err)
	}
	return &Bot{Client: bot}, nil
}

// Запуск бота и обработка сообщений
// Параметры:
//
// - service - сервис для работы с базой данных
// Возвращает ошибку при возникновении проблем с получением обновлений от Telegram
func (b *Bot) Start(service *database.Service) error {
	updates, err := b.Client.UpdatesViaLongPolling(nil)
	if err != nil {
		return fmt.Errorf("ошибка при получении обновлений от telegram: %v", err)
	}

	for update := range updates {
		if update.Message != nil {
			go b.handleMessage(update.Message, service)
		}
	}
	return nil
}

// Обработка сообщений
// Параметры:
//
// - message - сообщение от пользователя
// - service - сервис для работы с базой данных
func (b *Bot) sendMessage(chatID int64, msg string) {
	if _, err := b.Client.SendMessage(tu.Message(tu.ID(chatID), msg)); err != nil {
		log.Printf("Ошибка отправки сообщения: %v", err)
	}
}
