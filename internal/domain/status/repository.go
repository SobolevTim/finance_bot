package status

import "context"

type Repository interface {
	SetStatus(ctx context.Context, tgID string, status string) error
	GetStatus(ctx context.Context, tgID string) (string, error)
}
