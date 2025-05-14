package database

import (
	"context"
	"time"

	"github.com/SobolevTim/finance_bot/internal/domain/categories"
	"github.com/google/uuid"
)

// CategoriesCreate создает новую категорию
func (r *Repository) CategoriesCreate(ctx context.Context, category *categories.Categories) error {
	r.Logger.Debug("Создание категории", "category", category)
	query := `
		INSERT INTO categories (id, user_id, name, is_default, icon, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	now := time.Now()
	_, err := r.DB.Exec(ctx, query, category.ID, category.UserID, category.Name, category.IsDefault, category.Icon, category.CreatedAt, category.UpdatedAt)
	if err != nil {
		r.Logger.Debug("Ошибка создания категории", "error", err)
		return err
	}
	r.Logger.Debug("Категория создана", "category", category, "timeSinnce", time.Since(now))
	return nil
}

// CategoriesGetByID возвращает категорию по ID
func (r *Repository) CategoriesGetByID(ctx context.Context, id uuid.UUID) (*categories.Categories, error) {
	r.Logger.Debug("Получение категории по ID", "id", id)
	query := `
		SELECT id, user_id, name, is_default, icon, created_at, updated_at
		FROM categories
		WHERE id = $1
	`
	now := time.Now()
	row := r.DB.QueryRow(ctx, query, id)
	c := &categories.Categories{}
	err := row.Scan(&c.ID, &c.UserID, &c.Name, &c.IsDefault, &c.Icon, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		r.Logger.Debug("Ошибка получения категории", "error", err)
		return nil, err
	}
	r.Logger.Debug("Категория получена", "category", c, "timeSinnce", time.Since(now))
	return c, nil
}

// CategoriesGetBuIDs возвращает категории по ID
func (r *Repository) CategoriesGetBuIDs(ctx context.Context, ids []uuid.UUID) ([]*categories.Categories, error) {
	r.Logger.Debug("Получение категорий по ID", "ids", ids)
	query := `
		SELECT id, user_id, name, is_default, icon, created_at, updated_at
		FROM categories
		WHERE id = ANY($1)
	`
	now := time.Now()
	rows, err := r.DB.Query(ctx, query, ids)
	if err != nil {
		r.Logger.Debug("Ошибка получения категорий", "error", err)
		return nil, err
	}
	defer rows.Close()

	cat := make([]*categories.Categories, 0)
	for rows.Next() {
		c := &categories.Categories{}
		err = rows.Scan(&c.ID, &c.UserID, &c.Name, &c.IsDefault, &c.Icon, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			r.Logger.Debug("Ошибка сканирования категории", "error", err)
			return nil, err
		}
		cat = append(cat, c)
	}
	if rows.Err() != nil {
		r.Logger.Debug("Ошибка перебора категорий", "error", rows.Err())
		return nil, rows.Err()
	}
	r.Logger.Debug("Категории получены", "categories", cat, "timeSinnce", time.Since(now))
	return cat, nil
}

// CategoriesGetForUser возвращает список категорий для пользователя
func (r *Repository) CategoriesGetForUser(ctx context.Context, userID uuid.UUID) ([]*categories.Categories, error) {
	r.Logger.Debug("Получение категорий для пользователя", "userID", userID)
	query := `
		SELECT id, user_id, name, is_default, icon, created_at, updated_at
		FROM categories
		WHERE user_id = $1
	`
	now := time.Now()
	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		r.Logger.Debug("Ошибка получения категорий для пользователя", "error", err)
		return nil, err
	}
	defer rows.Close()

	cat := make([]*categories.Categories, 0)
	for rows.Next() {
		c := &categories.Categories{}
		err = rows.Scan(&c.ID, &c.UserID, &c.Name, &c.IsDefault, &c.Icon, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			r.Logger.Debug("Ошибка сканирования категории", "error", err)
			return nil, err
		}
		cat = append(cat, c)
	}
	if rows.Err() != nil {
		r.Logger.Debug("Ошибка перебора категорий", "error", rows.Err())
		return nil, rows.Err()
	}
	r.Logger.Debug("Категории получены", "categories", cat, "timeSinnce", time.Since(now))
	return cat, nil
}

// CategoriesGetDefaults возвращает список базовых категорий
func (r *Repository) CategoriesGetDefaults(ctx context.Context) ([]*categories.Categories, error) {
	r.Logger.Debug("Получение базовых категорий")
	query := `
		SELECT id, user_id, name, is_default, icon, created_at, updated_at
		FROM categories
		WHERE is_default = true
	`
	now := time.Now()
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		r.Logger.Debug("Ошибка получения базовых категорий", "error", err)
		return nil, err
	}
	defer rows.Close()

	cat := make([]*categories.Categories, 0)
	for rows.Next() {
		c := &categories.Categories{}
		err = rows.Scan(&c.ID, &c.UserID, &c.Name, &c.IsDefault, &c.Icon, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			r.Logger.Debug("Ошибка сканирования базовой категории", "error", err)
			return nil, err
		}
		cat = append(cat, c)
	}
	if rows.Err() != nil {
		r.Logger.Debug("Ошибка перебора базовых категорий", "error", rows.Err())
		return nil, rows.Err()
	}
	r.Logger.Debug("Базовые категории получены", "categories", cat, "timeSinnce", time.Since(now))
	return cat, nil
}

// CategoriesGetDefaultsByName возвращает Defaults категорию по именам
func (r *Repository) CategoriesGetDefaultsByName(ctx context.Context, name string) (*categories.Categories, error) {
	r.Logger.Debug("Получение базовой категории по имени", "name", name)
	query := `
		SELECT id, user_id, name, is_default, icon, created_at, updated_at
		FROM categories
		WHERE name = $1 AND is_default = true
	`
	now := time.Now()
	row := r.DB.QueryRow(ctx, query, name)
	c := &categories.Categories{}
	err := row.Scan(&c.ID, &c.UserID, &c.Name, &c.IsDefault, &c.Icon, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		r.Logger.Debug("Ошибка получения базовой категории по имени", "error", err)
		return nil, err
	}
	r.Logger.Debug("Базовая категория получена", "category", c, "timeSinnce", time.Since(now))
	return c, nil
}
