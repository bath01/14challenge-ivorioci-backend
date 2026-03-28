package handlers

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"ivorioci-stream-service/models"
	"ivorioci-stream-service/services"
)

type StreamHandler struct {
	videoService *services.VideoService
	storagePath  string
}

func NewStreamHandler(vs *services.VideoService, storagePath string) *StreamHandler {
	return &StreamHandler{videoService: vs, storagePath: storagePath}
}

// GET /stream/{id}
//
// Streams the video file using HTTP range requests (RFC 7233).
// http.ServeContent handles Accept-Ranges, Content-Range, and 206 Partial Content
// transparently, which is the correct approach for video players.
func (h *StreamHandler) StreamVideo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	video, err := h.videoService.GetVideoByID(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}

	if !video.IsPublished {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Video not found")
		return
	}

	absPath := filepath.Join(h.storagePath, filepath.Clean(video.FilePath))
	file, err := os.Open(absPath) //nolint:gosec — path is DB-controlled, not user-controlled
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "Video file not found")
			return
		}
		writeAppError(w, models.ErrInternal)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		writeAppError(w, models.ErrInternal)
		return
	}

	// Increment view count asynchronously — do not block the stream.
	go h.videoService.IncrementViews(context.Background(), id)

	w.Header().Set("Content-Type", video.MimeType)
	w.Header().Set("Cache-Control", "no-cache")
	// http.ServeContent sets Accept-Ranges, Content-Range and returns 206
	// when the client sends a Range header.
	http.ServeContent(w, r, video.Title, stat.ModTime(), file)
}
