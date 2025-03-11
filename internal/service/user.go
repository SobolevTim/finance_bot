package service

import (
	"context"
	"fmt"

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
			return nil, nil, fmt.Errorf("get budget by telegram id: %w", err)
		}
		return existingUser, budget, nil
	}

	newUser, err := user.New(telegramID, UserName, FirstName, LastName)
	if err != nil {
		return nil, nil, fmt.Errorf("create new user: %w", err)
	}

	defaultBudget, err := budget.NewDefaultBudget(telegramID)
	if err != nil {
		return nil, nil, fmt.Errorf("create new default budget: %w", err)
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, nil, fmt.Errorf("create user: %w", err)
	}

	if err := s.budgetRepo.CreateBudget(ctx, defaultBudget); err != nil {
		return nil, nil, fmt.Errorf("create budget: %w", err)
	}

	return newUser, defaultBudget, nil
}
