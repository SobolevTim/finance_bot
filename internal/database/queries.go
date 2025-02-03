package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx"
)

type Users struct {
	ID            int
	TelegramID    int64
	Username      string
	Email         string
	PasswordHash  string
	MonthlyBudget int
	CreatedAt     time.Time
	Notify        bool
}

type Expenses struct {
	ID          int
	UserID      int
	CategoryID  int
	Amount      int
	ExpenseDate time.Time
	CreatedAt   time.Time
}

func (b *Service) InsertStartUsers(user Users) error {
	ctx := context.Background()

	// Запись изначальных данных при первом запуске бота
	query := `
		INSERT INTO users (telegram_id, username) 
		VALUES ($1, $2);
	`
	_, err := b.DB.Exec(ctx, query,
		user.TelegramID,
		user.Username)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса первичной записи данных в users: %w", err)
	}

	return nil
}

func (b *Service) UpdateDayExpense(user Users, expense Expenses) (Expenses, error) {
	var result Expenses
	ctx := context.Background()
	query := `
        SELECT e.expense_id, e.amount, e.expense_date
        FROM expenses e
        JOIN users u ON e.user_id = u.user_id
        LEFT JOIN categories c ON e.category_id = c.category_id
        WHERE u.telegram_id = $1
        AND e.expense_date = $2;
    `
	err := b.DB.QueryRow(ctx, query, user.TelegramID, expense.ExpenseDate).Scan(&result.ID, &result.Amount, &result.ExpenseDate)
	if err != nil {
		if err.Error() == "no rows in result set" {
			insertQuery := `
                INSERT INTO expenses (user_id, amount, expense_date)
                VALUES ((SELECT user_id FROM users WHERE telegram_id = $1), $2, $3)
                RETURNING expense_id, amount, expense_date;
            `
			err = b.DB.QueryRow(ctx, insertQuery, user.TelegramID, 0, expense.ExpenseDate).Scan(&result.ID, &result.Amount, &result.ExpenseDate)
			if err != nil {
				return Expenses{}, fmt.Errorf("ошибка при записи 0 значения: %w", err)
			}
		} else {
			return Expenses{}, fmt.Errorf("ошибка при получении данных: %w", err)
		}
	}
	result.Amount += expense.Amount
	_, err = b.DB.Exec(ctx, "UPDATE expenses SET amount = $1 WHERE expense_id = $2",
		result.Amount,
		result.ID)
	if err != nil {
		return Expenses{}, fmt.Errorf("ошибка при обновлении expenses amount: %w", err)
	}
	return result, nil
}

func (b *Service) UpdateMontlyBudget(user Users) error {
	ctx := context.Background()
	_, err := b.DB.Exec(ctx, "UPDATE users SET monthly_budget = $1 WHERE telegram_id = $2",
		user.MonthlyBudget,
		user.TelegramID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении monthly_budget: %w", err)
	}

	return nil
}

func (b *Service) GetMontlyBudget(user Users) (Users, error) {
	ctx := context.Background()

	err := b.DB.QueryRow(ctx, "SELECT monthly_budget FROM users WHERE telegram_id = $1",
		user.TelegramID).Scan(&user.MonthlyBudget)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Users{}, fmt.Errorf("в БД нет информации о пользователе %d: %w", user.TelegramID, err)
		}
		return Users{}, fmt.Errorf("ошибка при получении данных: %w", err)
	}

	return user, nil
}

func (b *Service) GetAverageMontlyExpenses(user Users) (Expenses, error) {
	var result Expenses
	ctx := context.Background()

	query := `
		SELECT COALESCE(SUM(e.amount), 0) AS total_expenses
		FROM expenses e
		JOIN users u ON e.user_id = u.user_id
		WHERE u.telegram_id = $1
			AND date_trunc('month', expense_date) = date_trunc('month', CURRENT_DATE);
	`
	err := b.DB.QueryRow(ctx, query, user.TelegramID).Scan(&result.Amount)
	if err != nil {
		return Expenses{}, fmt.Errorf("ошибка при получении Average Montly Expenses: %w", err)
	}
	return result, nil
}

func (b *Service) GetUserNotify(user Users) (Users, error) {
	ctx := context.Background()

	query := `
		SELECT notify
		FROM users
		WHERE telegram_id = $1
	`
	err := b.DB.QueryRow(ctx, query, user.TelegramID).Scan(&user.Notify)
	if err != nil {
		return Users{}, fmt.Errorf("ошибка при получении notify: %w", err)
	}
	return user, nil
}

func (b *Service) UpdateUserNotify(user Users) error {
	ctx := context.Background()

	query := `
		UPDATE users
		SET notify = $1
		WHERE telegram_id = $2
	`
	_, err := b.DB.Exec(ctx, query,
		user.Notify,
		user.TelegramID)

	if err != nil {
		return fmt.Errorf("ошибка при обновлении UserNotify: %w", err)
	}
	return nil
}

func (b *Service) GetUsersWitchNotify() ([]Users, error) {
	ctx := context.Background()

	query := `
		SELECT user_id, telegram_id, username
		FROM users
		WHERE notify = TRUE
	`
	rows, err := b.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнения запроса к БД GetUsersWitchNotify: %w", err)
	}
	defer rows.Close()

	var users []Users
	for rows.Next() {
		var user Users
		if err := rows.Scan(&user.ID, &user.TelegramID, &user.Username); err != nil {
			return nil, fmt.Errorf("ошибка при считывании данных GetUsersWitchNotify: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при завершении итерации в GetUsersWitchNotify: %w", err)
	}

	return users, nil

}

func (b *Service) GetExpenseFromDate(user Users, expese Expenses) (Expenses, error) {
	ctx := context.Background()

	query := `
		SELECT e.amount
		FROM expenses e
		JOIN users u ON e.user_id = u.user_id
		WHERE u.telegram_id = $1
			AND e.expense_date = $2;
	`

	err := b.DB.QueryRow(ctx, query, user.TelegramID, expese.ExpenseDate).Scan(&expese.Amount)
	if err != nil {
		if err.Error() == "no rows in result set" {
			expese.Amount = 0
		} else {
			return Expenses{}, fmt.Errorf("ошибка при обновлении GetExpenseFromDate: %w", err)
		}
	}

	return expese, nil
}
