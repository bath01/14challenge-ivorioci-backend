package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"ivorioci-stream-service/models"
	"ivorioci-stream-service/services"
)

type VideoHandler struct {
	videoService    *services.VideoService
	categoryService *services.CategoryService
	storagePath     string // absolute path where video files are stored
	thumbnailPath   string // absolute path where thumbnail images are stored
	publicBaseURL   string // public gateway base URL, e.g. http://localhost:3000/api
}

func NewVideoHandler(
	vs *services.VideoService,
	cs *services.CategoryService,
	storagePath, thumbnailPath, publicBaseURL string,
) *VideoHandler {
	return &VideoHandler{
		videoService:    vs,
		categoryService: cs,
		storagePath:     storagePath,
		thumbnailPath:   thumbnailPath,
		publicBaseURL:   strings.TrimRight(publicBaseURL, "/"),
	}
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

// POST /videos  (multipart/form-data)
//
// Form fields:
//   - title       string  (required)
//   - description string  (optional)
//   - categoryId  string  (optional UUID)
//   - duration    int     (optional, seconds)
//
// Form files:
//   - video       file    (required — video/mp4, video/webm, video/ogg, video/quicktime)
//   - thumbnail   file    (required — image/jpeg, image/png, image/webp)
func (h *VideoHandler) CreateVideo(w http.ResponseWriter, r *http.Request) {
	// ParseMultipartForm allocates up to 32 MB in memory; the rest spills to temp files.
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Unable to parse multipart form")
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "title is required")
		return
	}
	description := r.FormValue("description")
	categoryID := strings.TrimSpace(r.FormValue("categoryId"))
	duration, _ := strconv.Atoi(r.FormValue("duration"))

	// ── Video file ────────────────────────────────────────────────────────────
	videoFile, videoHeader, err := r.FormFile("video")
	if err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "video file is required")
		return
	}
	defer videoFile.Close()

	if videoHeader.Size > maxVideoSize {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "video file exceeds the 2 GB limit")
		return
	}

	videoMIME, err := detectMIME(videoFile)
	if err != nil {
		writeAppError(w, models.ErrInternal)
		return
	}
	videoExt, ok := allowedVideoTypes[videoMIME]
	if !ok {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR",
			fmt.Sprintf("unsupported video format: %s", videoMIME))
		return
	}

	videoFilename := uniqueFilename(videoExt)
	videoDst := filepath.Join(h.storagePath, videoFilename)
	if _, err := saveUploadedFile(videoFile, videoDst); err != nil {
		writeAppError(w, models.ErrInternal)
		return
	}

	// ── Thumbnail file ────────────────────────────────────────────────────────
	thumbFile, thumbHeader, err := r.FormFile("thumbnail")
	if err != nil {
		os.Remove(videoDst)
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "thumbnail file is required")
		return
	}
	defer thumbFile.Close()

	if thumbHeader.Size > maxThumbSize {
		os.Remove(videoDst)
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "thumbnail exceeds the 10 MB limit")
		return
	}

	thumbMIME, err := detectMIME(thumbFile)
	if err != nil {
		os.Remove(videoDst)
		writeAppError(w, models.ErrInternal)
		return
	}
	thumbExt, ok := allowedImageTypes[thumbMIME]
	if !ok {
		os.Remove(videoDst)
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR",
			fmt.Sprintf("unsupported image format: %s", thumbMIME))
		return
	}

	thumbFilename := uniqueFilename(thumbExt)
	thumbDst := filepath.Join(h.thumbnailPath, thumbFilename)
	if _, err := saveUploadedFile(thumbFile, thumbDst); err != nil {
		os.Remove(videoDst)
		writeAppError(w, models.ErrInternal)
		return
	}

	// ── Persist in DB ─────────────────────────────────────────────────────────
	var catIDPtr *string
	if categoryID != "" {
		catIDPtr = &categoryID
	}

	dto := models.CreateVideoDTO{
		Title:        title,
		Description:  description,
		FilePath:     videoFilename,
		ThumbnailURL: h.publicBaseURL + "/thumbnails/" + thumbFilename,
		CategoryID:   catIDPtr,
		Duration:     duration,
		FileSize:     videoHeader.Size,
		MimeType:     videoMIME,
	}

	video, err := h.videoService.CreateVideo(r.Context(), dto)
	if err != nil {
		os.Remove(videoDst)
		os.Remove(thumbDst)
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
