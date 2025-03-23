package service

import (
	"context"
	"strconv"

	"github.com/SobolevTim/finance_bot/internal/domain/status"
)

func (s *Service) GetStatus(ctx context.Context, id int64) (string, error) {
	if id == 0 {
		return "", status.ErrEmptyTelegramID
	}
	tgID := strconv.FormatInt(id, 10)
	status, err := s.sR.GetStatus(ctx, tgID)
	if err != nil {
		return "", err
	}
	return status, nil
}

func (s *Service) SetStatus(ctx context.Context, id int64, statusValue string) error {
	if id == 0 {
		return status.ErrEmptyTelegramID
	}
	tgID := strconv.FormatInt(id, 10)

	err := s.sR.SetStatus(ctx, tgID, statusValue)
	if err != nil {
		return err
	}
	return nil
}
