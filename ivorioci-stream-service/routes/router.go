package routes

import (
	"net/http"

	"ivorioci-stream-service/handlers"
	"ivorioci-stream-service/middleware"

	"github.com/gorilla/mux"
)

func New(
	videoH *handlers.VideoHandler,
	categoryH *handlers.CategoryHandler,
	streamH *handlers.StreamHandler,
	thumbnailPath string,
	jwtSecret string,
) http.Handler {
	r := mux.NewRouter()
	r.Use(middleware.Logger)

	// Raccourci pour envelopper un handler avec le middleware JWT
	auth := middleware.RequireAuth(jwtSecret)
	protected := func(fn http.HandlerFunc) http.Handler {
		return auth(fn)
	}

	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success":true,"data":{"status":"ok"}}`))
	}).Methods(http.MethodGet)

	thumbFS := http.FileServer(http.Dir(thumbnailPath))
	r.PathPrefix("/thumbnails/").Handler(
		http.StripPrefix("/thumbnails/", thumbFS),
	).Methods(http.MethodGet)

	r.HandleFunc("/videos", videoH.ListVideos).Methods(http.MethodGet)
	r.HandleFunc("/videos/{id}", videoH.GetVideo).Methods(http.MethodGet)
	r.HandleFunc("/categories", categoryH.ListCategories).Methods(http.MethodGet)
	r.HandleFunc("/categories/{id}", categoryH.GetCategory).Methods(http.MethodGet)
	r.HandleFunc("/categories/{id}/videos", videoH.ListVideosByCategory).Methods(http.MethodGet)

	r.Handle("/stream/{id}", protected(streamH.StreamVideo)).Methods(http.MethodGet)

	r.Handle("/videos", protected(videoH.CreateVideo)).Methods(http.MethodPost)
	r.Handle("/videos/{id}", protected(videoH.UpdateVideo)).Methods(http.MethodPut)
	r.Handle("/videos/{id}", protected(videoH.DeleteVideo)).Methods(http.MethodDelete)

	r.Handle("/categories", protected(categoryH.CreateCategory)).Methods(http.MethodPost)
	r.Handle("/categories/{id}", protected(categoryH.UpdateCategory)).Methods(http.MethodPut)
	r.Handle("/categories/{id}", protected(categoryH.DeleteCategory)).Methods(http.MethodDelete)

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"success":false,"error":{"code":"NOT_FOUND","message":"Route not found"}}`))
	})

	return r
}
