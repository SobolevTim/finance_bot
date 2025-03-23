package service

import (
	"context"

	"github.com/SobolevTim/finance_bot/internal/domain/categories"
)

func (s *Service) GetDefaultCategories(ctx context.Context) ([]*categories.Categories, error) {
	return s.cR.CategoriesGetDefaults(ctx)
}
