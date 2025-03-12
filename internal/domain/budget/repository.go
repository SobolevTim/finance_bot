package budget

import "context"

type Repository interface {
	CreateBudget(ctx context.Context, budget *Budget) error
	GetBudgetByTelegramID(ctx context.Context, telegramID string) (*Budget, error)
	//UpdateBudgetAmount(ctx context.Context, budget *Budget) error
	//UpdateBudgetCurrency(ctx context.Context, budget *Budget) error
}
