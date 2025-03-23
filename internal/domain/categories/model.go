package categories

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	MaxNameLength            = 100
	DefaultIcon              = "📁"
	ErrNameTooLong           = errors.New("category name exceeds maximum length")
	ErrCategoryNotFound      = errors.New("category not found")
	ErrEmptyName             = errors.New("category name cannot be empty")
	ErrDuplicateName         = errors.New("category with this name already exists")
	ErrDeleteDefaultCategory = errors.New("cannot delete default category")
	ErrUpdateDefaultCategory = errors.New("cannot update default category")
	ErrCategoryInUse         = errors.New("category is used in existing expenses")
)

type Categories struct {
	ID        uuid.UUID // ID категории
	UserID    uuid.UUID // ID пользователя
	Name      string    // Название категории
	IsDefault bool      // Признак категории по умолчанию
	Icon      string    // Иконка категории
	CreatedAt time.Time // Дата создания
	UpdatedAt time.Time // Дата обновления
}

// New создает новую категорию с валидацией
func New(userID uuid.UUID, name string, isDefault bool) (*Categories, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}

	return &Categories{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		IsDefault: isDefault,
		Icon:      DefaultIcon,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func validateName(name string) error {
	if name == "" {
		return ErrEmptyName
	}

	if len(name) > MaxNameLength {
		return ErrNameTooLong
	}

	return nil
}
