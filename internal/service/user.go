package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/SobolevTim/finance_bot/internal/domain/user"
	"github.com/google/uuid"
)

// RegisterUser регистрирует нового пользователя
func (s *Service) RegisterUser(
	ctx context.Context,
	telegramID int64,
	userName string,
	firstName string,
	lastName string,
) (*user.User, error) {
	// Преобразование telegramID в строку
	telegramIDStr := strconv.FormatInt(telegramID, 10)

	existingUser, err := s.uR.UserGetByTelegramID(ctx, telegramIDStr)
	if err == nil {
		return existingUser, nil
	}

	if !errors.Is(err, user.ErrUserNotFound) {
		return nil, err
	}

	// Проверка уникальности username
	if _, err := s.uR.UserGetByUserName(ctx, userName); err == nil {
		return nil, user.ErrDuplicateUserName
	}

	newUser, err := user.New(telegramIDStr, userName, firstName, lastName)
	if err != nil {
		return nil, err
	}

	if err := s.uR.UserCreate(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *Service) GetUserByTelegramID(ctx context.Context, telegramID int64) (*user.User, error) {
	telegramIDStr := strconv.FormatInt(telegramID, 10)
	return s.uR.UserGetByTelegramID(ctx, telegramIDStr)
}

// UpdateUserProfile обновляет профиль пользователя
func (s *Service) UpdateUserProfile(
	ctx context.Context,
	userID uuid.UUID,
	userName string,
	firstName string,
	lastName string,
	timezone string,
) error {
	user, err := s.uR.UserGetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := user.UpdateNames(userName, firstName, lastName); err != nil {
		return err
	}

	if err := user.UpdateTimezone(timezone); err != nil {
		return err
	}

	return s.uR.UserUpdate(ctx, user)
}
