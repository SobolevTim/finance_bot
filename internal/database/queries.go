package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx"
)

// Users структура для хранения данных о пользователях
type Users struct {
	ID            int       // ID пользователя
	TelegramID    int64     // ID пользователя в Telegram
	Username      string    // Имя пользователя
	Email         string    // Почта пользователя
	PasswordHash  string    // Хеш пароля пользователя
	MonthlyBudget int       // Месячный бюджет пользователя
	CreatedAt     time.Time // Дата создания записи
	Notify        bool      // Подписка на уведомления
}

// Expenses структура для хранения данных о расходах
type Expenses struct {
	ID          int       // ID расхода
	UserID      int       // Пользователь, который внес расход
	CategoryID  int       // Категория расхода
	Note        string    // Описание расхода
	Amount      int       // Сумма расхода
	ExpenseDate time.Time // Дата расхода
	CreatedAt   time.Time // Дата создания записи
}

// InsertStartUsers записывает изначальные данные при первом запуске бота
//
// Параметры:
// - user - данные о пользователе
//
// Возвращает ошибку при возникновении проблем с записью данных в БД
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

// SetDayExpense записывает данные о расходах за день и возвращает общую сумму трат за день
// и данные о расходах
//
// Параметры:
// - user - данные о пользователе
// - expense - данные о расходах
//
// Возвращает ошибку при возникновении проблем с обновлением данных в БД
// или возвращает общую сумму трат за день и данные о расходах
func (b *Service) SetDayExpense(user Users, expense Expenses) (int, []Expenses, error) {
	ctx := context.Background()

	// Запись новой траты
	insertQuery := `
		INSERT INTO expenses (user_id, note, amount, expense_date)
		VALUES ((SELECT user_id FROM users WHERE telegram_id = $1), $2, $3, $4)
		RETURNING expense_id, amount, expense_date;
	`
	_, err := b.DB.Exec(ctx, insertQuery, user.TelegramID, expense.Note, expense.Amount, expense.ExpenseDate)
	if err != nil {
		return 0, nil, fmt.Errorf("ошибка при записи новой траты: %w", err)
	}

	// запрос на получения данных о тратах за день
	query := `
		SELECT e.expense_id, e.note, e.amount, e.expense_date
		FROM expenses e
		JOIN users u ON e.user_id = u.user_id
		WHERE u.telegram_id = $1
		AND e.expense_date = $2;
	`

	// Запрос данных о тратах за день
	rows, err := b.DB.Query(ctx, query, user.TelegramID, expense.ExpenseDate)
	if err != nil {
		return 0, nil, fmt.Errorf("ошибка при получении данных о тратах: %w", err)
	}
	defer rows.Close()

	// Подсчет общей суммы трат за день
	var totalAmount int
	var expenses []Expenses
	for rows.Next() {
		var exp Expenses
		if err := rows.Scan(&exp.ID, &exp.Note, &exp.Amount, &exp.ExpenseDate); err != nil {
			return 0, nil, fmt.Errorf("ошибка при считывании данных о тратах: %w", err)
		}
		totalAmount += exp.Amount
		expenses = append(expenses, exp)
	}

	// Завершение итерации по данным о тратах
	if err = rows.Err(); err != nil {
		return 0, nil, fmt.Errorf("ошибка при завершении итерации по данным о тратах: %w", err)
	}

	// Возврат общей суммы трат за день и данных о тратах
	return totalAmount, expenses, nil
}

// UpdateMontlyBudget обновляет данные о месячном бюджете
//
// Параметры:
// - user - данные о пользователе
//
// Возвращает ошибку при возникновении проблем с обновлением данных в БД
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

// GetMontlyBudget получает данные о месячном бюджете
//
// Параметры:
// - user - данные о пользователе
//
// Возвращает ошибку при возникновении проблем с получением данных из БД
// или возвращает данные о месячном бюджете
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

// GetAverageMontlyExpenses получает средние месячные траты
//
// Параметры:
// - user - данные о пользователе
//
// Возвращает ошибку при возникновении проблем с получением данных из БД
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

// GetUserNotify получает данные о настройках уведомлений пользователя
//
// Параметры:
// - user - данные о пользователе
//
// Возвращает ошибку при возникновении проблем с получением данных из БД
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

// UpdateUserNotify обновляет данные о настройках уведомлений пользователя
//
// Параметры:
// - user - данные о пользователе
//
// Возвращает ошибку при возникновении проблем с обновлением данных в БД
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

// GetUsersWitchNotify получает данные о пользователях, которые хотят получать уведомления
//
// Возвращает ошибку при возникновении проблем с получением данных из БД
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

// GetExpenseFromDate получает данные о расходах за определенную дату
//
// Параметры:
// - user - данные о пользователе
// - expese - данные о расходах
//
// Возвращает ошибку при возникновении проблем с получением данных из БД
// или возвращает данные о расходах
func (b *Service) GetExpenseFromDate(user Users, expenseDate time.Time) ([]Expenses, error) {
	ctx := context.Background()

	query := `
		SELECT e.expense_id, e.note, e.amount, e.expense_date
		FROM expenses e
		JOIN users u ON e.user_id = u.user_id
		WHERE u.telegram_id = $1
			AND e.expense_date = $2;
	`

	rows, err := b.DB.Query(ctx, query, user.TelegramID, expenseDate)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении данных о тратах: %w", err)
	}
	defer rows.Close()

	var expenses []Expenses
	for rows.Next() {
		var exp Expenses
		if err := rows.Scan(&exp.ID, &exp.Note, &exp.Amount, &exp.ExpenseDate); err != nil {
			return nil, fmt.Errorf("ошибка при считывании данных о тратах: %w", err)
		}
		expenses = append(expenses, exp)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при завершении итерации по данным о тратах: %w", err)
	}

	return expenses, nil
}
