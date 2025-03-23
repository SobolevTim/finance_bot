package budget_test

import (
	"testing"
	"time"

	"github.com/SobolevTim/finance_bot/internal/domain/budget"
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
