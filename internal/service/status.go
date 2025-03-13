package service

import (
	"context"
	"strconv"

	"github.com/SobolevTim/finance_bot/internal/domain/status"
)

type StatusMemory struct {
	memoryRepo status.Repository
}

func NewStatusMemory(memoryRepo status.Repository) *StatusMemory {
	return &StatusMemory{memoryRepo: memoryRepo}
}

func (sm *StatusMemory) GetStatus(ctx context.Context, id int64) (string, error) {
	if id == 0 {
		return "", status.ErrEmptyTelegramID
	}
	tgID := strconv.FormatInt(id, 10)
	status, err := sm.memoryRepo.GetStatus(ctx, tgID)
	if err != nil {
		return "", err
	}
	return status, nil
}

func (sm *StatusMemory) SetStatus(ctx context.Context, id int64, statusValue string) error {
	if id == 0 {
		return status.ErrEmptyTelegramID
	}
	tgID := strconv.FormatInt(id, 10)

	err := sm.memoryRepo.SetStatus(ctx, tgID, statusValue)
	if err != nil {
		return err
	}
	return nil
}
