package categories

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CategoriesCreate(ctx context.Context, category *Categories) error
	CategoriesGetByID(ctx context.Context, id uuid.UUID) (*Categories, error)
	CategoriesGetBuIDs(ctx context.Context, ids []uuid.UUID) ([]*Categories, error)
	CategoriesGetForUser(ctx context.Context, userID uuid.UUID) ([]*Categories, error)
	CategoriesGetDefaults(ctx context.Context) ([]*Categories, error)
	CategoriesGetDefaultsByName(ctx context.Context, name string) (*Categories, error)
}
