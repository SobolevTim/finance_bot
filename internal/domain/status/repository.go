package status

import "context"

type Repository interface {
	SetStatus(ctx context.Context, tgID string, status string) error
	GetStatus(ctx context.Context, tgID string) (string, error)
	SetExpenseStatus(ctx context.Context, ChatID string, ExpenseEntry *ExpenseEntry) error
	GetExpenseStatus(ctx context.Context, ChatID string) (*ExpenseEntry, error)
	Delete(ctx context.Context, tgID string) error
}
