package service

import (
	"github.com/SobolevTim/finance_bot/internal/domain/budget"
	"github.com/SobolevTim/finance_bot/internal/domain/status"
	"github.com/SobolevTim/finance_bot/internal/domain/user"
)

// Service реализует бизнес-логику budget
type Service struct {
	uR user.Repository
	bR budget.Repository
	sR status.Repository
}

func NewService(userRepo user.Repository,
	budgetRepo budget.Repository,
	statusRepo status.Repository,
) *Service {
	return &Service{
		uR: userRepo,
		bR: budgetRepo,
		sR: statusRepo,
	}
}
