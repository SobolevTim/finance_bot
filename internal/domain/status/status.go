package status

import "errors"

var (
	ErrEmptyTelegramID = errors.New("empty telegram id")
	ErrEmptyStatus     = errors.New("empty status")
)
