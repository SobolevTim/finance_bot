package budget

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrNegativeBudget        = errors.New("budget amount cannot be negative")
	ErrInvalidBudgetPeriod   = errors.New("end date must be after start date")
	ErrCategoryNotFound      = errors.New("category not found in budget")
	ErrCategoryLimitExceeded = errors.New("category limit exceeded")
	ErrUserNotFound          = errors.New("user not found")
)

// Budget представляет собой месячный бюджет пользователя
type Budget struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Amount     decimal.Decimal
	Currency   string
	StartDate  time.Time
	EndDate    time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Categories map[uuid.UUID]decimal.Decimal // Лимиты по категориям
}

// New создает новый бюджет с валидацией
func New(userID uuid.UUID, amount decimal.Decimal, currency string, startDate, endDate time.Time) (*Budget, error) {
	if amount.IsNegative() {
		return nil, ErrNegativeBudget
	}

	if endDate.Before(startDate) || endDate.Equal(startDate) {
		return nil, ErrInvalidBudgetPeriod
	}

	return &Budget{
		ID:         uuid.New(),
		UserID:     userID,
		Amount:     amount,
		Currency:   currency,
		StartDate:  startDate,
		EndDate:    endDate,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Categories: make(map[uuid.UUID]decimal.Decimal),
	}, nil
}

// IsActive проверяет, активен ли бюджет на текущий момент
func (b *Budget) IsActive(now time.Time) bool {
	return !now.Before(b.StartDate) && !now.After(b.EndDate)
}

// AddCategory добавляет лимит для категории
func (b *Budget) AddCategory(categoryID uuid.UUID, limit decimal.Decimal) error {
	if limit.IsNegative() {
		return ErrNegativeBudget
	}
	b.Categories[categoryID] = limit
	return nil
}
