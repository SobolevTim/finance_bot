package service

import (
	"context"

	"github.com/SobolevTim/finance_bot/internal/domain/user"
)

type UserService struct {
	repo user.Repository
}

func NewUserService(repo user.Repository) *UserService {
	return &UserService{repo: repo}
}

// RegisterUser обрабатывает логику регистрации
func (s *UserService) RegisterUser(ctx context.Context, telegramID, UserName, FirstName, LastName string) (*user.User, error) {
	existingUser, err := s.repo.GetByTelegramID(ctx, telegramID)
	if err == nil {
		return existingUser, nil // Пользователь уже существует
	}

	newUser, err := user.New(telegramID, UserName, FirstName, LastName)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
