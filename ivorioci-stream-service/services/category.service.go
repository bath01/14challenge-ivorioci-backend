package services

import (
	"context"
	"errors"

	"ivorioci-stream-service/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryService struct {
	db *pgxpool.Pool
}

func NewCategoryService(db *pgxpool.Pool) *CategoryService {
	return &CategoryService{db: db}
}

func (s *CategoryService) GetCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, name, slug, description, created_at, updated_at
		FROM categories
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	if categories == nil {
		categories = []models.Category{}
	}
	return categories, rows.Err()
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id string) (*models.Category, error) {
	var c models.Category
	err := s.db.QueryRow(ctx, `
		SELECT id, name, slug, description, created_at, updated_at
		FROM categories WHERE id = $1
	`, id).Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, models.ErrNotFound
	}
	return &c, err
}

func (s *CategoryService) GetCategoryBySlug(ctx context.Context, slug string) (*models.Category, error) {
	var c models.Category
	err := s.db.QueryRow(ctx, `
		SELECT id, name, slug, description, created_at, updated_at
		FROM categories WHERE slug = $1
	`, slug).Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, models.ErrNotFound
	}
	return &c, err
}

func (s *CategoryService) CreateCategory(ctx context.Context, dto models.CreateCategoryDTO) (*models.Category, error) {
	var exists bool
	_ = s.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM categories WHERE slug = $1)`, dto.Slug).Scan(&exists)
	if exists {
		return nil, models.ErrConflict
	}

	var c models.Category
	err := s.db.QueryRow(ctx, `
		INSERT INTO categories (id, name, slug, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, slug, description, created_at, updated_at
	`, uuid.New().String(), dto.Name, dto.Slug, dto.Description,
	).Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id string, dto models.UpdateCategoryDTO) (*models.Category, error) {
	existing, err := s.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if dto.Name != nil {
		existing.Name = *dto.Name
	}
	if dto.Slug != nil {
		var conflictID string
		_ = s.db.QueryRow(ctx, `SELECT id FROM categories WHERE slug = $1 AND id != $2`, *dto.Slug, id).Scan(&conflictID)
		if conflictID != "" {
			return nil, models.ErrConflict
		}
		existing.Slug = *dto.Slug
	}
	if dto.Description != nil {
		existing.Description = *dto.Description
	}

	var c models.Category
	err = s.db.QueryRow(ctx, `
		UPDATE categories SET name = $1, slug = $2, description = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING id, name, slug, description, created_at, updated_at
	`, existing.Name, existing.Slug, existing.Description, id,
	).Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	return &c, err
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	result, err := s.db.Exec(ctx, `DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}
