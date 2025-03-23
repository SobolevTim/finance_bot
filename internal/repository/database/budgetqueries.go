package database

import (
	"context"
	"time"

	"github.com/SobolevTim/finance_bot/internal/domain/budget"
	"github.com/google/uuid"
)

func (r *Repository) BudgetCreate(ctx context.Context, budget *budget.Budget) error {
	r.Logger.Debug("Создание бюджета", "budget", budget)
	query := `
		INSERT INTO budgets (id, user_id, amount, currency, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	now := time.Now()
	_, err := r.DB.Exec(ctx, query, budget.ID, budget.UserID, budget.Amount, budget.Currency, budget.StartDate, budget.EndDate, budget.CreatedAt, budget.UpdatedAt)
	if err != nil {
		r.Logger.Debug("Ошибка создания бюджета", "error", err)
		return err
	}
	r.Logger.Debug("Бюджет создан", "budget", budget, "timeSinnce", time.Since(now))
	return nil
}

func (r *Repository) BudgetGetByID(ctx context.Context, id uuid.UUID) (*budget.Budget, error) {
	r.Logger.Debug("Получение бюджета по ID", "id", id)
	query := `
		SELECT id, user_id, amount, currency, start_date, end_date, created_at, updated_at
		FROM budgets
		WHERE id = $1
	`
	now := time.Now()
	row := r.DB.QueryRow(ctx, query, id)
	b := &budget.Budget{}
	err := row.Scan(&b.ID, &b.UserID, &b.Amount, &b.Currency, &b.StartDate, &b.EndDate, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		r.Logger.Debug("Ошибка получения бюджета", "error", err)
		return nil, err
	}
	r.Logger.Debug("Бюджет получен", "budget", b, "timeSinnce", time.Since(now))
	return b, nil
}

func (r *Repository) BudgetGetCurrent(ctx context.Context, userID uuid.UUID) (*budget.Budget, error) {
	r.Logger.Debug("Получение текущего бюджета", "userID", userID)
	query := `
		SELECT id, user_id, amount, currency, start_date, end_date, created_at, updated_at
		FROM budgets
		WHERE user_id = $1 AND start_date <= NOW() AND end_date >= NOW()
	`
	now := time.Now()
	row := r.DB.QueryRow(ctx, query, userID)
	b := &budget.Budget{}
	err := row.Scan(&b.ID, &b.UserID, &b.Amount, &b.Currency, &b.StartDate, &b.EndDate, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		r.Logger.Debug("Ошибка получения текущего бюджета", "error", err)
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	r.Logger.Debug("Текущий бюджет получен", "budget", b, "timeSinnce", time.Since(now))
	return b, nil
}

func (r *Repository) BudgetUpdate(ctx context.Context, budget *budget.Budget) error {
	r.Logger.Debug("Обновление бюджета", "budget", budget)
	query := `
		UPDATE budgets
		SET amount = $2, currency = $3, start_date = $4, end_date = $5, updated_at = $6
		WHERE id = $1
	`
	now := time.Now()
	_, err := r.DB.Exec(ctx, query, budget.ID, budget.Amount, budget.Currency, budget.StartDate, budget.EndDate, budget.UpdatedAt)
	if err != nil {
		r.Logger.Debug("Ошибка обновления бюджета", "error", err)
		return err
	}
	r.Logger.Debug("Бюджет обновлен", "budget", budget, "timeSinnce", time.Since(now))
	return nil
}

func (r *Repository) BudgetDelete(ctx context.Context, id uuid.UUID) error {
	r.Logger.Debug("Удаление бюджета", "id", id)
	query := `
		DELETE FROM budgets
		WHERE id = $1
	`
	now := time.Now()
	_, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		r.Logger.Debug("Ошибка удаления бюджета", "error", err)
		return err
	}
	r.Logger.Debug("Бюджет удален", "id", id, "timeSinnce", time.Since(now))
	return err
}

func (r *Repository) BudgetGetByTgID(ctx context.Context, tgID string) (*budget.Budget, error) {
	r.Logger.Debug("Получение бюджета по telegramID", "telegramID", tgID)
	query := `
		SELECT b.id, b.user_id, b.amount, b.currency, b.start_date, b.end_date, b.created_at, b.updated_at
		FROM budgets b
		JOIN users u ON b.user_id = u.id
		WHERE u.telegram_id = $1 AND b.start_date <= NOW() AND b.end_date >= NOW()
	`
	now := time.Now()
	row := r.DB.QueryRow(ctx, query, tgID)
	b := &budget.Budget{}
	err := row.Scan(&b.ID, &b.UserID, &b.Amount, &b.Currency, &b.StartDate, &b.EndDate, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		r.Logger.Debug("Ошибка получения бюджета по tgID", "error", err)
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	r.Logger.Debug("Бюджет получен", "budget", b, "timeSinnce", time.Since(now))
	return b, nil
}
