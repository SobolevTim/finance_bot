package service

import (
	"time"

	"github.com/SobolevTim/finance_bot/internal/domain/budget"
	"github.com/SobolevTim/finance_bot/internal/domain/categories"
	"github.com/SobolevTim/finance_bot/internal/domain/expense"
	"github.com/SobolevTim/finance_bot/internal/domain/status"
	"github.com/SobolevTim/finance_bot/internal/domain/user"
)

// Service реализует бизнес-логику budget
type Service struct {
	uR user.Repository
	bR budget.Repository
	sR status.Repository
	eR expense.Repository
	cR categories.Repository
}

type ExpenseDTO struct {
	ID          string    // ID траты
	UserID      string    // ID пользователя
	CategoryID  string    // ID категории
	Amount      float64   // Сумма
	Date        time.Time // Дата
	IsRecurring bool      // Повторяющаяся
	Recurrence  string    // Периодичность
	Description string    // Описание
}

func NewService(userRepo user.Repository,
	budgetRepo budget.Repository,
	statusRepo status.Repository,
	expenseRepo expense.Repository,
	categoriesRepo categories.Repository,
) *Service {
	return &Service{
		uR: userRepo,
		bR: budgetRepo,
		sR: statusRepo,
		eR: expenseRepo,
		cR: categoriesRepo,
	}
}
