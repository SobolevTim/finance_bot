package user

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

var (
	telegramIDRegex = regexp.MustCompile(`^[0-9]+$`)
	timezoneRegex   = regexp.MustCompile(`^UTC[+-]\d{1,2}$`)
)

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrEmptyTelegramID         = errors.New("telegram ID cannot be empty")
	ErrEmptyUserName           = errors.New("username cannot be empty")
	ErrEmptyFirstName          = errors.New("first name cannot be empty")
	ErrInvalidTelegramIDFormat = errors.New("invalid telegram ID format")
	ErrDuplicateTelegramID     = errors.New("duplicate telegram ID")
	ErrDuplicateUserName       = errors.New("duplicate username")
	ErrInvalidTimezoneFormat   = errors.New("timezone must be in UTC±XX format")
)

// User представляет сущность пользователя системы
type User struct {
	ID         uuid.UUID
	TelegramID string
	UserName   string
	FirstName  string
	LastName   string
	Timezone   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// New создает нового пользователя с валидацией
func New(telegramID, userName, firstName, lastName string) (*User, error) {
	// Валидация обязательных полей
	if telegramID == "" {
		return nil, ErrEmptyTelegramID
	}
	if userName == "" {
		return nil, ErrEmptyUserName
	}
	if firstName == "" {
		return nil, ErrEmptyFirstName
	}

	if !telegramIDRegex.MatchString(telegramID) {
		return nil, ErrInvalidTelegramIDFormat
	}

	return &User{
		ID:         uuid.New(),
		TelegramID: telegramID,
		UserName:   userName,
		FirstName:  firstName,
		LastName:   lastName,
		Timezone:   "UTC+3", // Значение по умолчанию
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}, nil
}

// UpdateNames обновляет имя пользователя
func (u *User) UpdateNames(userName, firstName, lastName string) error {
	if userName == "" {
		return ErrEmptyUserName
	}
	if firstName == "" {
		return ErrEmptyFirstName
	}

	u.UserName = userName
	u.FirstName = firstName
	u.LastName = lastName
	u.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateTimezone обновляет временную зону пользователя
func (u *User) UpdateTimezone(tz string) error {
	if !timezoneRegex.MatchString(tz) {
		return ErrInvalidTimezoneFormat
	}

	u.Timezone = tz
	u.UpdatedAt = time.Now().UTC()
	return nil
}
