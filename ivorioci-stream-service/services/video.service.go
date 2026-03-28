package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"ivorioci-stream-service/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VideoService struct {
	db *pgxpool.Pool
}

func NewVideoService(db *pgxpool.Pool) *VideoService {
	return &VideoService{db: db}
}

func (s *VideoService) GetVideos(ctx context.Context, p models.VideoListParams) ([]models.Video, int, error) {
	p.Defaults()

	args := []interface{}{}
	conditions := []string{"v.is_published = true"}
	i := 1

	if p.Search != "" {
		conditions = append(conditions, fmt.Sprintf(
			"(v.title ILIKE $%d OR v.description ILIKE $%d)", i, i,
		))
		args = append(args, "%"+p.Search+"%")
		i++
	}
	if p.CategoryID != "" {
		conditions = append(conditions, fmt.Sprintf("v.category_id = $%d", i))
		args = append(args, p.CategoryID)
		i++
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	var total int
	err := s.db.QueryRow(ctx,
		fmt.Sprintf(`SELECT COUNT(*) FROM videos v %s`, where),
		args...,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	sortCol := map[string]string{
		"created_at":  "v.created_at",
		"views_count": "v.views_count",
		"title":       "v.title",
	}[p.SortBy]

	args = append(args, p.Limit, p.Offset())
	query := fmt.Sprintf(`
		SELECT v.id, v.title, v.description, v.thumbnail_url, v.category_id,
		       v.duration, v.views_count, v.file_size, v.mime_type, v.is_published,
		       v.created_at, v.updated_at,
		       c.id, c.name, c.slug
		FROM videos v
		LEFT JOIN categories c ON v.category_id = c.id
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, where, sortCol, strings.ToUpper(p.SortOrder), i, i+1)

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	videos := []models.Video{}
	for rows.Next() {
		var v models.Video
		var catID, catName, catSlug *string
		err := rows.Scan(
			&v.ID, &v.Title, &v.Description, &v.ThumbnailURL, &v.CategoryID,
			&v.Duration, &v.ViewsCount, &v.FileSize, &v.MimeType, &v.IsPublished,
			&v.CreatedAt, &v.UpdatedAt,
			&catID, &catName, &catSlug,
		)
		if err != nil {
			return nil, 0, err
		}
		if catID != nil {
			v.Category = &models.Category{ID: *catID, Name: *catName, Slug: *catSlug}
		}
		videos = append(videos, v)
	}
	return videos, total, rows.Err()
}

func (s *VideoService) GetVideoByID(ctx context.Context, id string) (*models.Video, error) {
	var v models.Video
	var catID, catName, catSlug *string

	err := s.db.QueryRow(ctx, `
		SELECT v.id, v.title, v.description, v.file_path, v.thumbnail_url, v.category_id,
		       v.duration, v.views_count, v.file_size, v.mime_type, v.is_published,
		       v.created_at, v.updated_at,
		       c.id, c.name, c.slug
		FROM videos v
		LEFT JOIN categories c ON v.category_id = c.id
		WHERE v.id = $1
	`, id).Scan(
		&v.ID, &v.Title, &v.Description, &v.FilePath, &v.ThumbnailURL, &v.CategoryID,
		&v.Duration, &v.ViewsCount, &v.FileSize, &v.MimeType, &v.IsPublished,
		&v.CreatedAt, &v.UpdatedAt,
		&catID, &catName, &catSlug,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if catID != nil {
		v.Category = &models.Category{ID: *catID, Name: *catName, Slug: *catSlug}
	}
	return &v, nil
}

func (s *VideoService) GetVideosByCategoryID(ctx context.Context, categoryID string, p models.VideoListParams) ([]models.Video, int, error) {
	p.CategoryID = categoryID
	return s.GetVideos(ctx, p)
}

func (s *VideoService) CreateVideo(ctx context.Context, dto models.CreateVideoDTO) (*models.Video, error) {
	mimeType := dto.MimeType
	if mimeType == "" {
		mimeType = "video/mp4"
	}

	var v models.Video
	err := s.db.QueryRow(ctx, `
		INSERT INTO videos (id, title, description, file_path, thumbnail_url, category_id,
		                    duration, file_size, mime_type)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, title, description, file_path, thumbnail_url, category_id,
		          duration, views_count, file_size, mime_type, is_published, created_at, updated_at
	`,
		uuid.New().String(), dto.Title, dto.Description, dto.FilePath, dto.ThumbnailURL,
		dto.CategoryID, dto.Duration, dto.FileSize, mimeType,
	).Scan(
		&v.ID, &v.Title, &v.Description, &v.FilePath, &v.ThumbnailURL, &v.CategoryID,
		&v.Duration, &v.ViewsCount, &v.FileSize, &v.MimeType, &v.IsPublished, &v.CreatedAt, &v.UpdatedAt,
	)
	return &v, err
}

func (s *VideoService) UpdateVideo(ctx context.Context, id string, dto models.UpdateVideoDTO) (*models.Video, error) {
	existing, err := s.GetVideoByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if dto.Title != nil {
		existing.Title = *dto.Title
	}
	if dto.Description != nil {
		existing.Description = *dto.Description
	}
	if dto.ThumbnailURL != nil {
		existing.ThumbnailURL = *dto.ThumbnailURL
	}
	if dto.CategoryID != nil {
		existing.CategoryID = dto.CategoryID
	}
	if dto.IsPublished != nil {
		existing.IsPublished = *dto.IsPublished
	}

	var v models.Video
	err = s.db.QueryRow(ctx, `
		UPDATE videos
		SET title=$1, description=$2, thumbnail_url=$3, category_id=$4, is_published=$5, updated_at=NOW()
		WHERE id=$6
		RETURNING id, title, description, file_path, thumbnail_url, category_id,
		          duration, views_count, file_size, mime_type, is_published, created_at, updated_at
	`,
		existing.Title, existing.Description, existing.ThumbnailURL,
		existing.CategoryID, existing.IsPublished, id,
	).Scan(
		&v.ID, &v.Title, &v.Description, &v.FilePath, &v.ThumbnailURL, &v.CategoryID,
		&v.Duration, &v.ViewsCount, &v.FileSize, &v.MimeType, &v.IsPublished, &v.CreatedAt, &v.UpdatedAt,
	)
	return &v, err
}

func (s *VideoService) DeleteVideo(ctx context.Context, id string) error {
	result, err := s.db.Exec(ctx, `DELETE FROM videos WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return models.ErrNotFound
	}
	return nil
}

func (s *VideoService) IncrementViews(ctx context.Context, id string) {
	_, _ = s.db.Exec(ctx, `UPDATE videos SET views_count = views_count + 1 WHERE id = $1`, id)
}
