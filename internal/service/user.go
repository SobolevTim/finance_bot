package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/SobolevTim/finance_bot/internal/domain/budget"
	"github.com/SobolevTim/finance_bot/internal/domain/user"
)

type UserService struct {
	userRepo   user.Repository
	budgetRepo budget.Repository
}

func NewUserService(userRepo user.Repository, budgetRepo budget.Repository) *UserService {
	return &UserService{userRepo: userRepo, budgetRepo: budgetRepo}
}

// RegisterUser обрабатывает логику регистрации
func (s *UserService) RegisterUser(ctx context.Context, telegramID, UserName, FirstName, LastName string) (*user.User, *budget.Budget, error) {
	existingUser, err := s.userRepo.GetByTelegramID(ctx, telegramID)
	if err == nil {
		budget, err := s.budgetRepo.GetBudgetByTelegramID(ctx, telegramID)
		if err != nil {
			return nil, nil, fmt.Errorf("ошибка при получении бюджета: %w", err)
		}
		return existingUser, budget, nil
	}

	newUser, err := user.New(telegramID, UserName, FirstName, LastName)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка при создании пользователя: %w", err)
	}

	defaultBudget, err := budget.NewDefaultBudget(telegramID)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка при создании default budget: %w", err)
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, nil, fmt.Errorf("ошибка при запись в бд пользователя: %w", err)
	}

	if err := s.budgetRepo.CreateBudget(ctx, defaultBudget); err != nil {
		return nil, nil, fmt.Errorf("ошибка при запись в бд бюджета: %w", err)
	}

	return newUser, defaultBudget, nil
}

func (s *UserService) UpdateBudget(ctx context.Context, id int64, amount int64) error {
	tgID := strconv.FormatInt(id, 10)
	_, err := s.budgetRepo.GetBudgetByTelegramID(ctx, tgID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении бюджета: %w", err)
	}
	newBudget, err := budget.NewBudget(tgID, amount, "RUB")
	if err != nil {
		return fmt.Errorf("ошибка при обновлении бюджета: %w", err)
	}
	if err := s.budgetRepo.UpdateBudget(ctx, newBudget); err != nil {
		return fmt.Errorf("ошибка при обновлении бюджета: %w", err)
	}
	return nil
}
