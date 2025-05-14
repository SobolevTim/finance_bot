package expense

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	CreateExpens(ctx context.Context, expense *Expense) error
	UpdateExpens(ctx context.Context, expense *Expense) error
	DeleteExpens(ctx context.Context, id uuid.UUID) error
	GetExpenses(ctx context.Context, id uuid.UUID) (*Expense, error)
	GetExpensesByUserID(ctx context.Context, userID uuid.UUID) ([]*Expense, error)
	GetExpensesByDate(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*Expense, error)
	GetExpensesByDateForDay(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*Expense, error)
	GetExpensesByTelegramID(ctx context.Context, telegramID int64) ([]*Expense, error)
	GetExpensesByTelegramIDAndDate(ctx context.Context, telegramID int64, date string) ([]*Expense, error)
}
