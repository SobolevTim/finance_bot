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

func (s *Service) SetExpenseStatus(ctx context.Context, id int64, expenseEntryDTO *ExpenseEntryDTO) error {
	if id == 0 {
		return status.ErrEmptyTelegramID
	}
	tgID := strconv.FormatInt(id, 10)

	expenseEntry := &status.ExpenseEntry{
		Date:     expenseEntryDTO.Date,
		Amount:   expenseEntryDTO.Amount,
		Category: expenseEntryDTO.Category,
		Note:     expenseEntryDTO.Note,
		Step:     expenseEntryDTO.Step,
	}

	err := s.sR.SetExpenseStatus(ctx, tgID, expenseEntry)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetExpenseStatus(ctx context.Context, id int64) (*ExpenseEntryDTO, error) {
	if id == 0 {
		return nil, status.ErrEmptyTelegramID
	}
	tgID := strconv.FormatInt(id, 10)
	expenseEntry, err := s.sR.GetExpenseStatus(ctx, tgID)
	if err != nil {
		return nil, err
	}
	if expenseEntry == nil {
		return nil, nil
	}

	expenseEntryDTO := &ExpenseEntryDTO{
		Date:     expenseEntry.Date,
		Amount:   expenseEntry.Amount,
		Category: expenseEntry.Category,
		Note:     expenseEntry.Note,
		Step:     expenseEntry.Step,
	}

	return expenseEntryDTO, nil
}

func (s *Service) DeleteStatus(ctx context.Context, id int64) error {
	if id == 0 {
		return status.ErrEmptyTelegramID
	}
	tgID := strconv.FormatInt(id, 10)
	err := s.sR.Delete(ctx, tgID)
	if err != nil {
		return err
	}
	return nil
}
