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
	jwtSecret string,
) http.Handler {
	r := mux.NewRouter()

	r.Use(middleware.Logger)

	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success":true,"data":{"status":"ok"}}`))
	}).Methods(http.MethodGet)

	r.HandleFunc("/videos", videoH.ListVideos).Methods(http.MethodGet)
	r.HandleFunc("/videos/{id}", videoH.GetVideo).Methods(http.MethodGet)
	r.HandleFunc("/categories/", categoryH.ListCategories).Methods(http.MethodGet)
	r.HandleFunc("/categories/{id}", categoryH.GetCategory).Methods(http.MethodGet)
	r.HandleFunc("/categories/{id}/videos", videoH.ListVideosByCategory).Methods(http.MethodGet)

	protected := r.NewRoute().Subrouter()
	protected.Use(middleware.RequireAuth(jwtSecret))

	protected.HandleFunc("/stream/{id}", streamH.StreamVideo).Methods(http.MethodGet)

	protected.HandleFunc("/videos", videoH.CreateVideo).Methods(http.MethodPost)
	protected.HandleFunc("/videos/{id}", videoH.UpdateVideo).Methods(http.MethodPut)
	protected.HandleFunc("/videos/{id}", videoH.DeleteVideo).Methods(http.MethodDelete)

	protected.HandleFunc("/categories", categoryH.CreateCategory).Methods(http.MethodPost)
	protected.HandleFunc("/categories/{id}", categoryH.UpdateCategory).Methods(http.MethodPut)
	protected.HandleFunc("/categories/{id}", categoryH.DeleteCategory).Methods(http.MethodDelete)

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"success":false,"error":{"code":"NOT_FOUND","message":"Route not found"}}`))
	})

	return r
}
