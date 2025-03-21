package budget_test

import (
	"testing"
	"time"

	"github.com/SobolevTim/finance_bot/internal/domain/budget"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestBudget_IsActive(t *testing.T) {
	now := time.Now()
	budget := budget.Budget{
		StartDate: now.AddDate(0, -1, 0),
		EndDate:   now.AddDate(0, 1, 0),
	}

	assert.True(t, budget.IsActive(now), "бюджет должен быть активен")
}

func TestBudget_UpdateBalance(t *testing.T) {
	categoryID := uuid.New()
	budget := budget.Budget{
		Amount: decimal.NewFromInt(1000),
		Categories: map[uuid.UUID]decimal.Decimal{
			categoryID: decimal.NewFromInt(500),
		},
	}

	err := budget.UpdateBalance(decimal.NewFromInt(300), categoryID)
	assert.NoError(t, err, "обновление баланса не должно вызывать ошибку")
	assert.Equal(t, decimal.NewFromInt(700), budget.Amount, "остаток бюджета должен уменьшиться")

	err = budget.UpdateBalance(decimal.NewFromInt(600), categoryID)
	assert.Error(t, err, "превышение лимита категории должно вызывать ошибку")
}
