package user

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

var (
	ErrEmptyTelegramID = errors.New("empty telegram id")
)

type User struct {
	ID         uuid.UUID
	TelegramID string
	CreatedAt  time.Time
	Timezone   string // Для корректного учета времени
}

func New(telegramID string) (*User, error) {
	if telegramID == "" {
		return nil, ErrEmptyTelegramID
	}
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return &User{
		ID:         uuid,
		TelegramID: telegramID,
		CreatedAt:  time.Now().UTC(),
	}, nil
}
