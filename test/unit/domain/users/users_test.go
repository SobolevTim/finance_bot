package user_test

import (
	"testing"

	"github.com/SobolevTim/finance_bot/internal/domain/user"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	t.Run("valid user", func(t *testing.T) {
		u, err := user.New("123456789", "john_doe", "John", "Doe")
		assert.NoError(t, err)
		assert.Equal(t, "123456789", u.TelegramID)
		assert.Equal(t, "john_doe", u.UserName)
		assert.Equal(t, "UTC+3", u.Timezone)
	})

	t.Run("empty first name", func(t *testing.T) {
		_, err := user.New("123456789", "johndoe", "", "Doe")
		assert.ErrorIs(t, err, user.ErrEmptyFirstName)
	})
}

func TestUpdateNames(t *testing.T) {
	u, _ := user.New("123456789", "old_username", "John", "Doe")

	t.Run("valid update", func(t *testing.T) {
		err := u.UpdateNames("new_username", "Jane", "Smith")
		assert.NoError(t, err)
		assert.Equal(t, "new_username", u.UserName)
		assert.Equal(t, "Jane", u.FirstName)
	})

	t.Run("empty username", func(t *testing.T) {
		err := u.UpdateNames("", "Jane", "Smith")
		assert.ErrorIs(t, err, user.ErrEmptyUserName)
	})
}

func TestUpdateTimezone(t *testing.T) {
	u, _ := user.New("123456789", "johndoe", "John", "Doe")

	t.Run("valid timezone", func(t *testing.T) {
		err := u.UpdateTimezone("UTC+5")
		assert.NoError(t, err)
		assert.Equal(t, "UTC+5", u.Timezone)
	})

	t.Run("invalid format", func(t *testing.T) {
		err := u.UpdateTimezone("Invalid/Timezone")
		assert.ErrorIs(t, err, user.ErrInvalidTimezoneFormat)
	})
}
