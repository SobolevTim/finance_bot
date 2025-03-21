package budget

import (
	"context"

	"github.com/google/uuid"
)

// Repository определяет методы для работы с бюджетами
type Repository interface {
	BudgetCreate(ctx context.Context, budget *Budget) error
	BudgetGetByID(ctx context.Context, id uuid.UUID) (*Budget, error)
	BudgetGetCurrent(ctx context.Context, userID uuid.UUID) (*Budget, error)
	BudgetUpdate(ctx context.Context, budget *Budget) error
	BudgetDelete(ctx context.Context, id uuid.UUID) error
}
