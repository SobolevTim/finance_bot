package expense

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrEmptyRecurrenceRule = errors.New("recurrence rule is required for recurring expenses")
	ErrorExpenseNotFound   = errors.New("expense not found")
)

type Expense struct {
	ID             uuid.UUID       // ID траты
	UserID         uuid.UUID       // ID пользователя
	CategoryID     uuid.UUID       // ID категории
	Ammount        decimal.Decimal // Сумма траты
	Date           time.Time       // Дата траты
	IsRecurring    bool            // Повторяющаяся траты
	RecurrenceRule string          // Правило повторения
	Description    string          // Описание траты
	//CardID 	   uuid.UUID 	   // TODO: добавить поле для карты
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewExpences(userID, categoryID uuid.UUID, amount decimal.Decimal, date time.Time, isRecurring bool, recurrenceRule, description string) (*Expense, error) {
	if isRecurring == true && recurrenceRule == "" {
		return nil, ErrEmptyRecurrenceRule
	}

	return &Expense{
		ID:             uuid.New(), // Генерация нового ID
		UserID:         userID,
		CategoryID:     categoryID,
		Ammount:        amount,
		Date:           date,
		IsRecurring:    isRecurring,
		RecurrenceRule: recurrenceRule,
		Description:    description,
		//CardID:         cardID, // TODO: добавить поле для карты
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
