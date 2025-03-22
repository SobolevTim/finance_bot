package service

import (
	"context"
	"strconv"
	"time"

	"github.com/SobolevTim/finance_bot/internal/domain/budget"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (s *Service) UpdateBudgetByTgID(ctx context.Context, tgID int64, amount string) (*budget.Budget, error) {
	// Преобразование строки в decimal
	amountDec, err := decimal.NewFromString(amount)
	if err != nil {
		return nil, err
	}

	budget, err := s.GetBudgetByTgID(ctx, tgID)
	if err != nil {
		return nil, err
	}

	if budget != nil {
		budget.Amount = amountDec
		err := s.bR.BudgetUpdate(ctx, budget)
		if err != nil {
			return nil, err
		}
		return budget, nil
	} else {
		budget, err := s.CreateBudget(ctx, tgID, amount, "RUB")
		if err != nil {
			return nil, err
		}
		return budget, nil
	}
}

// CreateBudget создает новый бюджет для пользователя
func (s *Service) CreateBudget(
	ctx context.Context,
	tgID int64,
	amount string,
	currency string,
) (*budget.Budget, error) {
	// Преобразование строки в decimal
	amountDec, err := decimal.NewFromString(amount)
	// Преобразование telegramID в строку
	telegramID := strconv.FormatInt(tgID, 10)
	user, err := s.uR.UserGetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, budget.ErrUserNotFound
	}
	userID := user.ID
	startDate := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	endDate := startDate.AddDate(0, 1, -1)

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

// UpdateBudget обновляет бюджет пользователя
func (s *Service) UpdateBudget(
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

	budget, err := s.bR.BudgetGetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	budget.Amount = amountDec
	budget.Currency = currency
	budget.StartDate = startDate
	budget.EndDate = endDate

	if err := s.bR.BudgetUpdate(ctx, budget); err != nil {
		return nil, err
	}

	return budget, nil
}

func (s *Service) GetBudgetByTgID(ctx context.Context, tgID int64) (*budget.Budget, error) {
	telegramID := strconv.FormatInt(tgID, 10)
	return s.bR.BudgetGetByTgID(ctx, telegramID)
}
