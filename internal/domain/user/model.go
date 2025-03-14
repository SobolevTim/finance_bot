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

func New(telegramID, userName, firstName, lastName string) (*User, error) {
	if telegramID == "" {
		return nil, ErrEmptyTelegramID
	}
	if userName == "" {
		return nil, ErrEmptyUserName
	}
	if firstName == "" {
		firstName = userName
	}
	if lastName == "" {
		lastName = userName
	}
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return &User{
		ID:         uuid,
		TelegramID: telegramID,
		UserName:   userName,
		FirstName:  firstName,
		LastName:   lastName,
		CreatedAt:  time.Now().UTC(),
	}, nil
}
