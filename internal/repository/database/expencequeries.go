package database

import (
	"context"
	"time"

	"github.com/SobolevTim/finance_bot/internal/domain/expense"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CreateExpens создает новый расход
// возвращает ошибку, если не удалось создать расход
func (r *Repository) CreateExpens(ctx context.Context, expense *expense.Expense) error {
	r.Logger.Debug("Запись нового расхода в базу данных", "expense", expense)
	query := `INSERT INTO expenses (user_id, category_id, amount, date, is_recurring, recurrence_rule, description) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	now := time.Now()
	_, err := r.DB.Exec(ctx, query, expense.UserID, expense.CategoryID, expense.Ammount, expense.Date, expense.IsRecurring, expense.RecurrenceRule, expense.Description)
	if err != nil {
		r.Logger.Debug("Не удалось создать расход", "error", err)
		return err
	}
	r.Logger.Debug("Расход успешно создан", "expense", expense, "duration", time.Since(now))
	return nil
}

// UpdateExpens обновляет существующий расход
// возвращает ошибку, если не удалось обновить расход
func (r *Repository) UpdateExpens(ctx context.Context, expense *expense.Expense) error {
	r.Logger.Debug("Обновление расхода в базе данных", "expense", expense)
	query := `UPDATE expenses SET category_id = $1, amount = $2, date = $3, is_recurring = $4, recurrence_rule = $5, description = $6 WHERE id = $7`

	now := time.Now()
	_, err := r.DB.Exec(ctx, query, expense.CategoryID, expense.Ammount, expense.Date, expense.IsRecurring, expense.RecurrenceRule, expense.Description, expense.ID)
	if err != nil {
		r.Logger.Debug("Не удалось обновить расход", "error", err)
		return err
	}
	r.Logger.Debug("Расход успешно обновлен", "expense", expense, "duration", time.Since(now))
	return nil
}

// DeleteExpens удаляет существующий расход
// возвращает ошибку, если не удалось удалить расход
func (r *Repository) DeleteExpens(ctx context.Context, id uuid.UUID) error {
	r.Logger.Debug("Удаление расхода из базы данных", "id", id)
	query := `DELETE FROM expenses WHERE id = $1`

	now := time.Now()
	_, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		r.Logger.Debug("Не удалось удалить расход", "error", err)
		return err
	}
	r.Logger.Debug("Расход успешно удален", "id", id, "duration", time.Since(now))
	return nil
}

// GetExpenses возвращает расход по ID
// возвращает ошибку, если не удалось получить расход
func (r *Repository) GetExpenses(ctx context.Context, id uuid.UUID) (*expense.Expense, error) {
	r.Logger.Debug("Получение расхода из базы данных", "id", id)
	query := `SELECT id, user_id, category_id, amount, date, is_recurring, recurrence_rule, description FROM expenses WHERE id = $1`

	now := time.Now()
	row := r.DB.QueryRow(ctx, query, id)
	expense := &expense.Expense{}
	err := row.Scan(&expense.ID, &expense.UserID, &expense.CategoryID, &expense.Ammount, &expense.Date, &expense.IsRecurring, &expense.RecurrenceRule, &expense.Description)
	if err != nil {
		r.Logger.Debug("Не удалось получить расход", "error", err)
		return nil, err
	}
	r.Logger.Debug("Расход успешно получен", "expense", expense, "duration", time.Since(now))
	return expense, nil
}

// GetExpensesByUserID возвращает все расходы по ID пользователя
// возвращает ошибку, если не удалось получить расходы
func (r *Repository) GetExpensesByUserID(ctx context.Context, userID uuid.UUID) ([]*expense.Expense, error) {
	r.Logger.Debug("Получение всех расходов из базы данных по ID пользователя", "userID", userID)
	query := `SELECT id, user_id, category_id, amount, date, is_recurring, recurrence_rule, description FROM expenses WHERE user_id = $1`

	now := time.Now()
	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		r.Logger.Debug("Не удалось получить расходы", "error", err)
		return nil, err
	}
	defer rows.Close()

	expenses := make([]*expense.Expense, 0)
	for rows.Next() {
		expense := &expense.Expense{}
		err := rows.Scan(&expense.ID, &expense.UserID, &expense.CategoryID, &expense.Ammount, &expense.Date, &expense.IsRecurring, &expense.RecurrenceRule, &expense.Description)
		if err != nil {
			r.Logger.Debug("Не удалось получить расход", "error", err)
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	r.Logger.Debug("Расходы успешно получены", "expenses", expenses, "duration", time.Since(now))
	return expenses, nil
}

// GetExpensesByDate возвращает все расходы в диапазоне дат по ID пользователя
// возвращает ошибку, если не удалось получить расходы
func (r *Repository) GetExpensesByDate(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*expense.Expense, error) {
	r.Logger.Debug("Получение всех расходов из базы данных по дате", "userID", userID, "startDate", startDate, "endDate", endDate)
	query := `SELECT id, user_id, category_id, amount, date, is_recurring, recurrence_rule, description FROM expenses WHERE user_id = $1 AND date >= $2 AND date <= $3 ORDER BY date ASC`

	now := time.Now()
	rows, err := r.DB.Query(ctx, query, userID, startDate, endDate)
	if err != nil {
		r.Logger.Debug("Не удалось получить расходы", "error", err)
		return nil, err
	}
	defer rows.Close()

	expenses := make([]*expense.Expense, 0)
	for rows.Next() {
		expense := &expense.Expense{}
		err := rows.Scan(&expense.ID, &expense.UserID, &expense.CategoryID, &expense.Ammount, &expense.Date, &expense.IsRecurring, &expense.RecurrenceRule, &expense.Description)
		if err != nil {
			r.Logger.Debug("Не удалось получить расход", "error", err)
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	r.Logger.Debug("Расходы успешно получены", "expenses", expenses, "duration", time.Since(now))
	return expenses, nil
}

// GetExpensesByDate возвращает сумму расходов в диапазоне дат по ID пользователя
// возвращает ошибку, если не удалось получить расходы
func (r *Repository) GetExpensesByDateForDay(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*expense.Expense, error) {
	r.Logger.Debug("Получение всех расходов из базы данных по дате", "userID", userID, "startDate", startDate, "endDate", endDate)
	query := `SELECT date, SUM(amount) as total_amount FROM expenses WHERE user_id = $1 AND date >= $2 AND date <= $3
		GROUP BY date
		ORDER BY date ASC`

	now := time.Now()
	rows, err := r.DB.Query(ctx, query, userID, startDate, endDate)
	if err != nil {
		r.Logger.Debug("Не удалось получить расходы", "error", err)
		return nil, err
	}
	defer rows.Close()

	expenses := make([]*expense.Expense, 0)
	for rows.Next() {
		var dailyExpense struct {
			Date        time.Time
			TotalAmount decimal.Decimal
		}
		err := rows.Scan(&dailyExpense.Date, &dailyExpense.TotalAmount)
		if err != nil {
			r.Logger.Debug("Не удалось получить расход", "error", err)
			return nil, err
		}
		expenses = append(expenses, &expense.Expense{
			Date:    dailyExpense.Date,
			Ammount: dailyExpense.TotalAmount,
		})
	}
	r.Logger.Debug("Расходы успешно получены", "expenses", &expenses, "duration", time.Since(now))
	return expenses, nil
}

// GetExpensesByTelegramID возвращает все расходы по ID пользователя Telegram
// возвращает ошибку, если не удалось получить расходы
func (r *Repository) GetExpensesByTelegramID(ctx context.Context, telegramID int64) ([]*expense.Expense, error) {
	r.Logger.Debug("Получение всех расходов из базы данных по ID пользователя Telegram", "telegramID", telegramID)
	query := `SELECT e.id, e.user_id, e.category_id, e.amount, e.date, e.is_recurring, e.recurrence_rule, e.description FROM expenses e JOIN users u ON e.user_id = u.id WHERE u.telegram_id = $1`

	now := time.Now()
	rows, err := r.DB.Query(ctx, query, telegramID)
	if err != nil {
		r.Logger.Debug("Не удалось получить расходы", "error", err)
		return nil, err
	}
	defer rows.Close()

	expenses := make([]*expense.Expense, 0)
	for rows.Next() {
		expense := &expense.Expense{}
		err := rows.Scan(&expense.ID, &expense.UserID, &expense.CategoryID, &expense.Ammount, &expense.Date, &expense.IsRecurring, &expense.RecurrenceRule, &expense.Description)
		if err != nil {
			r.Logger.Debug("Не удалось получить расход", "error", err)
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	r.Logger.Debug("Расходы успешно получены", "expenses", expenses, "duration", time.Since(now))
	return expenses, nil
}

// GetExpensesByTelegramIDAndDate возвращает все расходы по ID пользователя Telegram и дате
// возвращает ошибку, если не удалось получить расходы
func (r *Repository) GetExpensesByTelegramIDAndDate(ctx context.Context, telegramID int64, date string) ([]*expense.Expense, error) {
	r.Logger.Debug("Получение всех расходов из базы данных по ID пользователя Telegram и дате", "telegramID", telegramID, "date", date)
	query := `SELECT e.id, e.user_id, e.category_id, e.amount, e.date, e.is_recurring, e.recurrence_rule, e.description FROM expenses e JOIN users u ON e.user_id = u.id WHERE u.telegram_id = $1 AND e.date = $2`

	now := time.Now()
	rows, err := r.DB.Query(ctx, query, telegramID, date)
	if err != nil {
		r.Logger.Debug("Не удалось получить расходы", "error", err)
		return nil, err
	}
	defer rows.Close()

	expenses := make([]*expense.Expense, 0)
	for rows.Next() {
		expense := &expense.Expense{}
		err := rows.Scan(&expense.ID, &expense.UserID, &expense.CategoryID, &expense.Ammount, &expense.Date, &expense.IsRecurring, &expense.RecurrenceRule, &expense.Description)
		if err != nil {
			r.Logger.Debug("Не удалось получить расход", "error", err)
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	r.Logger.Debug("Расходы успешно получены", "expenses", expenses, "duration", time.Since(now))
	return expenses, nil
}
