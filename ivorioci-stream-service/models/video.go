package models

import "time"

type Video struct {
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	FilePath     string     `json:"-"` // never exposed in API
	ThumbnailURL string     `json:"thumbnailUrl"`
	CategoryID   *string    `json:"categoryId"`
	Category     *Category  `json:"category,omitempty"`
	Duration     int        `json:"duration"` // seconds
	ViewsCount   int        `json:"viewsCount"`
	FileSize     int64      `json:"fileSize"` // bytes
	MimeType     string     `json:"mimeType"`
	IsPublished  bool       `json:"isPublished"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type CreateVideoDTO struct {
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	FilePath     string  `json:"filePath"`
	ThumbnailURL string  `json:"thumbnailUrl"`
	CategoryID   *string `json:"categoryId"`
	Duration     int     `json:"duration"`
	FileSize     int64   `json:"fileSize"`
	MimeType     string  `json:"mimeType"`
}

type UpdateVideoDTO struct {
	Title        *string `json:"title"`
	Description  *string `json:"description"`
	ThumbnailURL *string `json:"thumbnailUrl"`
	CategoryID   *string `json:"categoryId"`
	IsPublished  *bool   `json:"isPublished"`
}

type VideoListParams struct {
	Page       int
	Limit      int
	Search     string
	CategoryID string
	SortBy     string // "created_at" | "views_count" | "title"
	SortOrder  string // "asc" | "desc"
}

func (p *VideoListParams) Defaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 20
	}
	validSortBy := map[string]bool{"created_at": true, "views_count": true, "title": true}
	if !validSortBy[p.SortBy] {
		p.SortBy = "created_at"
	}
	if p.SortOrder != "asc" && p.SortOrder != "desc" {
		p.SortOrder = "desc"
	}
}

func (p *VideoListParams) Offset() int {
	return (p.Page - 1) * p.Limit
}
