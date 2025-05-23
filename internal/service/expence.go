package service

import (
	"context"
	"strconv"
	"time"

	"github.com/SobolevTim/finance_bot/internal/domain/expense"
	"github.com/SobolevTim/finance_bot/internal/domain/user"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (s *Service) CreateExpensByTelegramID(ctx context.Context, telegramID int64, amount float64, date time.Time, description string) error {
	// Преобразование int64 в строку
	telegramIDStr := strconv.FormatInt(telegramID, 10)
	// Получение пользователя по telegramID
	u, err := s.uR.UserGetByTelegramID(ctx, telegramIDStr)
	if err != nil {
		return user.ErrUserNotFound
	}

	if u == nil {
		return user.ErrUserNotFound
	}

	// Преобразование строки в decimal
	amountDec := decimal.NewFromFloat(amount)

	// Получение категории по умолчанию
	c, err := s.cR.CategoriesGetDefaultsByName(ctx, "Прочее")
	if err != nil {
		return err
	}

	// Создание новой траты
	newExpens, err := expense.NewExpences(u.ID, c.ID, amountDec, date, false, "", description)
	if err != nil {
		return err
	}
	// Сохранение траты в базе данных
	err = s.eR.CreateExpens(ctx, newExpens)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetExpenses(ctx context.Context, telegramID int64, startDate, endDate time.Time) ([]*ExpenseDTO, error) {
	// Преобразование int64 в строку
	telegramIDStr := strconv.FormatInt(telegramID, 10)
	// Получение пользователя по telegramID
	u, err := s.uR.UserGetByTelegramID(ctx, telegramIDStr)
	if err != nil {
		return nil, user.ErrUserNotFound
	}

	if u == nil {
		return nil, user.ErrUserNotFound
	}

	// Получение трат за период
	expenses, err := s.eR.GetExpensesByDateForDay(ctx, u.ID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Преобразование трат в DTO
	expensesDTO := make([]*ExpenseDTO, 0, len(expenses))
	for _, e := range expenses {
		expensesDTO = append(expensesDTO, &ExpenseDTO{
			ID:          e.ID.String(),
			UserID:      e.UserID.String(),
			CategoryID:  e.CategoryID.String(),
			Amount:      e.Ammount.InexactFloat64(),
			Date:        e.Date,
			IsRecurring: e.IsRecurring,
			Recurrence:  e.RecurrenceRule,
			Description: e.Description,
		})
	}

	return expensesDTO, nil
}

func (s *Service) AddExpense(ctx context.Context, telegramID int64, amount float64, date time.Time, category, description string) error {
	// Преобразование int64 в строку
	telegramIDStr := strconv.FormatInt(telegramID, 10)
	// Получение пользователя по telegramID
	u, err := s.uR.UserGetByTelegramID(ctx, telegramIDStr)
	if err != nil {
		return user.ErrUserNotFound
	}

	if u == nil {
		return user.ErrUserNotFound
	}

	// Преобразование строки в decimal
	amountDec := decimal.NewFromFloat(amount)

	// Получение категории по умолчанию
	c, err := s.cR.CategoriesGetDefaultsByName(ctx, category)
	if err != nil {
		return err
	}

	// Создание новой траты
	newExpens, err := expense.NewExpences(u.ID, c.ID, amountDec, date, false, "", description)
	if err != nil {
		return err
	}
	// Сохранение траты в базе данных
	err = s.eR.CreateExpens(ctx, newExpens)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetExpensesByMonth(ctx context.Context, telegramID int64) ([]*ExpenseDTO, float64, error) {
	// Преобразование int64 в строку
	telegramIDStr := strconv.FormatInt(telegramID, 10)
	// Получение пользователя по telegramID
	u, err := s.uR.UserGetByTelegramID(ctx, telegramIDStr)
	if err != nil {
		return nil, 0, user.ErrUserNotFound
	}

	if u == nil {
		return nil, 0, user.ErrUserNotFound
	}

	// startDate - первый день текущего месяца
	startDate := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	// endDate - последний день текущего месяца
	endDate := startDate.AddDate(0, 1, -1)

	// Получение трат за текущий месяц
	expenses, err := s.eR.GetExpensesByDate(ctx, u.ID, startDate, endDate)
	if err != nil {
		return nil, 0, err
	}

	if len(expenses) == 0 {
		return nil, 0, expense.ErrorExpenseNotFound
	}

	// Получение Категорий с иконками
	var ids []uuid.UUID

	for _, e := range expenses {
		ids = append(ids, e.CategoryID)
	}

	categories, err := s.cR.CategoriesGetBuIDs(ctx, ids)
	sum := 0.0

	// Преобразование трат в DTO
	expensesDTO := make([]*ExpenseDTO, 0, len(expenses))
	for _, e := range expenses {
		var icon string
		var category string
		for _, c := range categories {
			if c.ID == e.CategoryID {
				icon = c.Icon
				category = c.Name
				break
			}
		}
		sum += e.Ammount.InexactFloat64()
		expensesDTO = append(expensesDTO, &ExpenseDTO{
			ID:           e.ID.String(),
			UserID:       e.UserID.String(),
			CategoryID:   e.CategoryID.String(),
			CategoryIcon: icon,
			Category:     category,
			Amount:       e.Ammount.InexactFloat64(),
			Date:         e.Date,
			IsRecurring:  e.IsRecurring,
			Recurrence:   e.RecurrenceRule,
			Description:  e.Description,
		})
	}

	return expensesDTO, sum, nil
}
