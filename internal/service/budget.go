package service

import (
	"context"
	"time"

	"github.com/SobolevTim/finance_bot/internal/domain/budget"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CreateBudget создает новый бюджет для пользователя
func (s *Service) CreateBudget(
	ctx context.Context,
	userID uuid.UUID,
	amount string,
	currency string,
	startDate, endDate time.Time,
) (*budget.Budget, error) {
	// Преобразование строки в decimal
	amountDec, err := decimal.NewFromString(amount)
	if err != nil {
		return nil, err
	}

	budget, err := budget.New(userID, amountDec, currency, startDate, endDate)
	if err != nil {
		return nil, err
	}

	if err := s.bR.BudgetCreate(ctx, budget); err != nil {
		return nil, err
	}

	return budget, nil
}

// GetCurrentBudget возвращает текущий активный бюджет пользователя
func (s *Service) GetCurrentBudget(ctx context.Context, userID uuid.UUID) (*budget.Budget, error) {
	return s.bR.BudgetGetCurrent(ctx, userID)
}

// AddExpense обновляет бюджет после добавления расхода
func (s *Service) AddExpense(ctx context.Context, budgetID uuid.UUID, amount decimal.Decimal, categoryID uuid.UUID) error {
	budget, err := s.bR.BudgetGetByID(ctx, budgetID)
	if err != nil {
		return err
	}

	if err := budget.UpdateBalance(amount, categoryID); err != nil {
		return err
	}

	return s.bR.BudgetUpdate(ctx, budget)
}
