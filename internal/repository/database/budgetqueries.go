package database

import (
	"context"

	"github.com/SobolevTim/finance_bot/internal/domain/budget"
	"github.com/google/uuid"
)

func (r *Repository) BudgetCreate(ctx context.Context, budget *budget.Budget) error {
	query := `
		INSERT INTO budgets (id, user_id, amount, currency, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.DB.Exec(ctx, query, budget.ID, budget.UserID, budget.Amount, budget.Currency, budget.StartDate, budget.EndDate, budget.CreatedAt, budget.UpdatedAt)
	return err
}

func (r *Repository) BudgetGetByID(ctx context.Context, id uuid.UUID) (*budget.Budget, error) {
	query := `
		SELECT id, user_id, amount, currency, start_date, end_date, created_at, updated_at
		FROM budgets
		WHERE id = $1
	`

	row := r.DB.QueryRow(ctx, query, id)
	b := &budget.Budget{}
	err := row.Scan(&b.ID, &b.UserID, &b.Amount, &b.Currency, &b.StartDate, &b.EndDate, &b.CreatedAt, &b.UpdatedAt)
	return b, err
}

func (r *Repository) BudgetGetCurrent(ctx context.Context, userID uuid.UUID) (*budget.Budget, error) {
	query := `
		SELECT id, user_id, amount, currency, start_date, end_date, created_at, updated_at
		FROM budgets
		WHERE user_id = $1 AND end_date IS NULL
	`
	row := r.DB.QueryRow(ctx, query, userID)
	b := &budget.Budget{}
	err := row.Scan(&b.ID, &b.UserID, &b.Amount, &b.Currency, &b.StartDate, &b.EndDate, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return b, nil
}

func (r *Repository) BudgetUpdate(ctx context.Context, budget *budget.Budget) error {
	query := `
		UPDATE budgets
		SET amount = $2, currency = $3, start_date = $4, end_date = $5, updated_at = $6
		WHERE id = $1
	`

	_, err := r.DB.Exec(ctx, query, budget.ID, budget.Amount, budget.Currency, budget.StartDate, budget.EndDate, budget.UpdatedAt)
	return err
}

func (r *Repository) BudgetDelete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM budgets
		WHERE id = $1
	`

	_, err := r.DB.Exec(ctx, query, id)
	return err
}
