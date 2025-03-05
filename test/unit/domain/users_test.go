package users_test

import (
	"testing"

	"github.com/SobolevTim/finance_bot/internal/domain/user"
	"github.com/gofrs/uuid"
)

func TestNew_ValidTelegramID(t *testing.T) {
	tg := "123"
	u, err := user.New(tg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u.TelegramID != tg {
		t.Errorf("expected TelegramID %v, got %v", tg, u.TelegramID)
	}
	if u.ID == uuid.Nil {
		t.Errorf("expected a valid UUID, got %v", u.ID)
	}
	if u.CreatedAt.IsZero() {
		t.Errorf("expected CreatedAt to be set, got zero value")
	}
}

func TestNew_EmptyTelegramID(t *testing.T) {
	u, err := user.New("")
	if err != user.ErrEmptyTelegramID {
		t.Fatalf("expected error %v, got %v", user.ErrEmptyTelegramID, err)
	}
	if u != nil {
		t.Errorf("expected user to be nil, got %v", u)
	}
}

func TestNew_UserFields(t *testing.T) {
	tg := "123"
	u, err := user.New(tg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u.TelegramID != tg {
		t.Errorf("expected TelegramID %v, got %v", tg, u.TelegramID)
	}
	if u.ID == uuid.Nil {
		t.Errorf("expected a valid UUID, got %v", u.ID)
	}
	if u.CreatedAt.IsZero() {
		t.Errorf("expected CreatedAt to be set, got zero value")
	}
	if u.Timezone != "" {
		t.Errorf("expected Timezone to be empty, got %v", u.Timezone)
	}
}
