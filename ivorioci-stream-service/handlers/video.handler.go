package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"ivorioci-stream-service/models"
	"ivorioci-stream-service/services"
)

type VideoHandler struct {
	videoService    *services.VideoService
	categoryService *services.CategoryService
}

func NewVideoHandler(vs *services.VideoService, cs *services.CategoryService) *VideoHandler {
	return &VideoHandler{videoService: vs, categoryService: cs}
}

// GET /videos
func (h *VideoHandler) ListVideos(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	params := models.VideoListParams{
		Page:       queryInt(r, "page", 1),
		Limit:      queryInt(r, "limit", 20),
		Search:     q.Get("search"),
		CategoryID: q.Get("categoryId"),
		SortBy:     q.Get("sortBy"),
		SortOrder:  q.Get("sortOrder"),
	}

	videos, total, err := h.videoService.GetVideos(r.Context(), params)
	if err != nil {
		writeAppError(w, err)
		return
	}

	params.Defaults()
	totalPages := total / params.Limit
	if total%params.Limit != 0 {
		totalPages++
	}

	writeSuccess(w, http.StatusOK, models.PaginatedData{
		Items:      videos,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	})
}

// GET /videos/{id}
func (h *VideoHandler) GetVideo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	video, err := h.videoService.GetVideoByID(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeSuccess(w, http.StatusOK, video)
}

// GET /categories/{id}/videos
func (h *VideoHandler) ListVideosByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := mux.Vars(r)["id"]

	// Verify the category exists
	if _, err := h.categoryService.GetCategoryByID(r.Context(), categoryID); err != nil {
		writeAppError(w, err)
		return
	}

	params := models.VideoListParams{
		Page:      queryInt(r, "page", 1),
		Limit:     queryInt(r, "limit", 20),
		SortBy:    r.URL.Query().Get("sortBy"),
		SortOrder: r.URL.Query().Get("sortOrder"),
	}

	videos, total, err := h.videoService.GetVideosByCategoryID(r.Context(), categoryID, params)
	if err != nil {
		writeAppError(w, err)
		return
	}

	params.Defaults()
	totalPages := total / params.Limit
	if total%params.Limit != 0 {
		totalPages++
	}

	writeSuccess(w, http.StatusOK, models.PaginatedData{
		Items:      videos,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	})
}

// POST /videos
func (h *VideoHandler) CreateVideo(w http.ResponseWriter, r *http.Request) {
	var dto models.CreateVideoDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}
	if dto.Title == "" || dto.FilePath == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "title and filePath are required")
		return
	}

	video, err := h.videoService.CreateVideo(r.Context(), dto)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeSuccess(w, http.StatusCreated, video)
}

// PUT /videos/{id}
func (h *VideoHandler) UpdateVideo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var dto models.UpdateVideoDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	video, err := h.videoService.UpdateVideo(r.Context(), id, dto)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeSuccess(w, http.StatusOK, video)
}

// DELETE /videos/{id}
func (h *VideoHandler) DeleteVideo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.videoService.DeleteVideo(r.Context(), id); err != nil {
		writeAppError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
