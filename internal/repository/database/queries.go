package database

import (
	"context"

	"github.com/SobolevTim/finance_bot/internal/domain/budget"
	"github.com/SobolevTim/finance_bot/internal/domain/user"
)

// Create создает нового пользователя
//
// ctx - контекст
// u - пользователь
//
// Возвращает ошибку
func (s *Repository) Create(ctx context.Context, u *user.User) error {
	// формируем запрос
	query := `INSERT INTO users (id, telegram_id, user_name, first_name, last_name, created_at, timezone) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	s.Logger.Debug("Создание пользователя", "ID", u.ID, "TelegramID", u.TelegramID, "UserName", u.UserName, "FirstName", u.FirstName, "LastName", u.LastName, "CreatedAt", u.CreatedAt, "Timezone", u.Timezone)
	// выполняем запрос
	_, err := s.DB.Exec(ctx, query, u.ID, u.TelegramID, u.UserName, u.FirstName, u.LastName, u.CreatedAt, u.Timezone)
	if err == nil {
		s.Logger.Debug("User created", "ID", u.ID, "TelegramID", u.TelegramID, "UserName", u.UserName, "FirstName", u.FirstName, "LastName", u.LastName, "CreatedAt", u.CreatedAt, "Timezone", u.Timezone)
	}
	return err
}

// GetByTelegramID ищет пользователя по telegramID
//
// ctx - контекст
// telegramID - telegramID пользователя
//
// Возвращает пользователя или ошибку
func (s *Repository) GetByTelegramID(ctx context.Context, telegramID string) (*user.User, error) {
	// формируем запрос
	query := `SELECT id, telegram_id, user_name, first_name, last_name, created_at, timezone
			  FROM users
			  WHERE telegram_id = $1`
	s.Logger.Debug("Поиск пользователя", "TelegramID", telegramID)
	// выполняем запрос
	row := s.DB.QueryRow(ctx, query, telegramID)

	// сканируем результат
	u := &user.User{}
	err := row.Scan(&u.ID, &u.TelegramID, &u.UserName, &u.FirstName, &u.LastName, &u.CreatedAt, &u.Timezone)
	if err != nil {
		return nil, err
	}
	s.Logger.Debug("Пользователй найден", "ID", u.ID, "TelegramID", u.TelegramID, "UserName", u.UserName, "FirstName", u.FirstName, "LastName", u.LastName, "CreatedAt", u.CreatedAt, "Timezone", u.Timezone)
	return u, nil
}

// UpdateTimeZone обновляет часовой пояс пользователя
//
// ctx - контекст
// u - пользователь
//
// Возвращает ошибку
func (s *Repository) UpdateTimeZone(ctx context.Context, u *user.User) error {
	// формируем запрос
	query := `UPDATE users
			  SET timezone = $1
			  WHERE id = $2`
	s.Logger.Debug("Обновление часового пояса", "ID", u.ID, "Timezone", u.Timezone)
	// выполняем запрос
	_, err := s.DB.Exec(ctx, query, u.Timezone, u.ID)
	if err != nil {
		return err
	}
	s.Logger.Debug("Часовой пояс пользователя обновлен", "ID", u.ID, "Timezone", u.Timezone)
	return nil
}

// CreateBudget создает новый бюджет
//
// ctx - контекст
// b - бюджет
//
// Возвращает ошибку
func (s *Repository) CreateBudget(ctx context.Context, b *budget.Budget) error {
	// формируем запрос
	query := `INSERT INTO budgets (id, telegram_id, amount, currency, created_at) 
			  VALUES ($1, $2, $3, $4, $5)`
	s.Logger.Debug("Создание бюджета", "ID", b.ID, "TelegramID", b.TelegramID, "Amount", b.Amount, "Currency", b.Currency, "Date", b.Date)
	// выполняем запрос
	_, err := s.DB.Exec(ctx, query, b.ID, b.TelegramID, b.Amount, b.Currency, b.Date)
	if err == nil {
		s.Logger.Debug("Бюджет записан", "ID", b.ID, "TelegramID", b.TelegramID, "Amount", b.Amount, "Currency", b.Currency, "Date", b.Date)
	}
	return err
}

// GetBudgetByTelegramID ищет бюджет по telegramID
// ctx - контекст
// telegramID - telegramID пользователя
//
// Возвращает бюджет или ошибку
func (s *Repository) GetBudgetByTelegramID(ctx context.Context, telegramID string) (*budget.Budget, error) {
	// TODO: добавить обработку ошибки, если бюджет не найден

	// формируем запрос
	query := `SELECT id, telegram_id, amount, currency, created_at, updated_at
			  FROM budgets
			  WHERE telegram_id = $1`

	s.Logger.Debug("Поиск бюджета", "TelegramID", telegramID)
	// выполняем запрос
	row := s.DB.QueryRow(ctx, query, telegramID)
	b := &budget.Budget{}

	// сканируем результат
	err := row.Scan(&b.ID, &b.TelegramID, &b.Amount, &b.Currency, &b.Date, &b.UpdateDate)
	if err != nil {
		return nil, err
	}
	s.Logger.Debug("Бюджет найден", "ID", b.ID, "TelegramID", b.TelegramID, "Amount", b.Amount, "Currency", b.Currency, "Date", b.Date, "UpdateDate", b.UpdateDate)
	return b, nil
}

// UpdateBudget обновляет бюджет
//
// ctx - контекст
// b - бюджет
//
// Возвращает ошибку
func (s *Repository) UpdateBudget(ctx context.Context, b *budget.Budget) error {
	// формируем запрос
	query := `UPDATE budgets
			  SET amount = $1, currency = $2, updated_at = $3
			  WHERE telegram_id = $4`
	s.Logger.Debug("Обновление бюджета", "ID", b.ID, "TelegramID", b.TelegramID, "Amount", b.Amount, "Currency", b.Currency, "UpdateDate", b.UpdateDate)
	// выполняем запрос
	_, err := s.DB.Exec(ctx, query, b.Amount, b.Currency, b.UpdateDate, b.TelegramID)
	if err != nil {
		return err
	}
	s.Logger.Debug("Бюджет обновлен", "ID", b.ID, "TelegramID", b.TelegramID, "Amount", b.Amount, "Currency", b.Currency, "UpdateDate", b.UpdateDate)
	return nil
}
