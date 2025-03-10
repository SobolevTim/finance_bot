package user

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

var (
	ErrEmptyTelegramID = errors.New("empty telegram id")
	ErrEmptyUserName   = errors.New("empty user name")
)

type User struct {
	ID         uuid.UUID
	TelegramID string
	UserName   string
	FirstName  string
	LastName   string
	CreatedAt  time.Time
	Timezone   string // Для корректного учета времени
}

func New(telegramID, UserName, FirstName, LastName string) (*User, error) {
	if telegramID == "" {
		return nil, ErrEmptyTelegramID
	}
	if UserName == "" {
		return nil, ErrEmptyUserName
	}
	if FirstName == "" {
		FirstName = UserName
	}
	if LastName == "" {
		LastName = UserName
	}
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return &User{
		ID:         uuid,
		TelegramID: telegramID,
		UserName:   UserName,
		FirstName:  FirstName,
		LastName:   LastName,
		CreatedAt:  time.Now().UTC(),
	}, nil
}
