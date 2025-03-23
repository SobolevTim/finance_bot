package status

import (
	"errors"
	"time"
)

var (
	ErrEmptyTelegramID = errors.New("empty telegram id")
	ErrEmptyStatus     = errors.New("empty status")
)

type ExpenseEntry struct {
	Date     time.Time
	Amount   float64
	Category string
	Note     string
	Step     string // Текущий шаг: "date", "date_input", "amount", "category", "note", "note_input", "confirm"
}
