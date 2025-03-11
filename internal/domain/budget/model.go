package budget

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

var (
	ErrEmptyTelegramID = errors.New("empty telegram id") // Ошибка пустого Telegram ID
	ErrEmptyCurrency   = errors.New("empty currency")    // Ошибка пустой валюты
)

// Budget - структура для хранения информации о бюджете пользователя
type Budget struct {
	ID         uuid.UUID // Идентификатор
	TelegramID string    // Telegram ID пользователя
	Amount     int64     // Сумма * 100
	Currency   string    // Валюта
	Date       time.Time // Дата
	UpdateDate string    // Дата обновления
}

func NewBudget(telegramID string, amount int64, currency string) (*Budget, error) {
	if telegramID == "" {
		return nil, ErrEmptyTelegramID
	}
	if currency == "" {
		return nil, ErrEmptyCurrency
	}
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return &Budget{
		ID:         uuid,
		TelegramID: telegramID,
		Amount:     amount,
		Currency:   currency,
		Date:       time.Now().UTC(),
	}, nil
}

func NewDefaultBudget(telegramID string) (*Budget, error) {
	return NewBudget(telegramID, 0, "RUB")
}
