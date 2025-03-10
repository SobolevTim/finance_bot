package users_test

import (
	"testing"

	"github.com/SobolevTim/finance_bot/internal/domain/user"
	"github.com/gofrs/uuid"
)

func TestNew_ValidTelegramID(t *testing.T) {
	tg := "123"
	UserName := "test"
	FirstName := "test"
	LastName := "test"
	u, err := user.New(tg, UserName, FirstName, LastName)
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
	u, err := user.New("", "test", "test", "test")
	if err != user.ErrEmptyTelegramID {
		t.Fatalf("expected error %v, got %v", user.ErrEmptyTelegramID, err)
	}
	if u != nil {
		t.Errorf("expected user to be nil, got %v", u)
	}
}

func TestNew_EmptyUserName(t *testing.T) {
	u, err := user.New("123", "", "test", "test")
	if err != user.ErrEmptyUserName {
		t.Fatalf("expected error %v, got %v", user.ErrEmptyUserName, err)
	}
	if u != nil {
		t.Errorf("expected user to be nil, got %v", u)
	}
}

func TestNew_EmptyFirstName(t *testing.T) {
	tg := "123"
	UserName := "test"
	u, err := user.New(tg, UserName, "", "test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u.FirstName != UserName {
		t.Errorf("expected FirstName %v, got %v", UserName, u.FirstName)
	}
}

func TestNew_EmptyLastName(t *testing.T) {
	tg := "123"
	UserName := "testName"
	u, err := user.New(tg, UserName, "test", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u.LastName != UserName {
		t.Errorf("expected LastName %v, got %v", UserName, u.LastName)
	}
}

func TestNew_UserFields(t *testing.T) {
	tg := "123"
	UserName := "testName"
	FirstName := "testFirstName"
	LastName := "testLastName"
	u, err := user.New(tg, UserName, FirstName, LastName)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u.TelegramID != tg {
		t.Errorf("expected TelegramID %v, got %v", tg, u.TelegramID)
	}
	if u.UserName != UserName {
		t.Errorf("expected UserName %v, got %v", UserName, u.UserName)
	}
	if u.FirstName != FirstName {
		t.Errorf("expected FirstName %v, got %v", FirstName, u.FirstName)
	}
	if u.LastName != LastName {
		t.Errorf("expected LastName %v, got %v", LastName, u.LastName)
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
