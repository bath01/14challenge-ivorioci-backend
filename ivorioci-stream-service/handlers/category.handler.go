package handlers

import (
	"encoding/json"
	"net/http"

	"ivorioci-stream-service/models"
	"ivorioci-stream-service/services"

	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(cs *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: cs}
}

func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryService.GetCategories(r.Context())
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeSuccess(w, http.StatusOK, categories)
}

func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	category, err := h.categoryService.GetCategoryByID(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeSuccess(w, http.StatusOK, category)
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var dto models.CreateCategoryDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}
	if dto.Name == "" || dto.Slug == "" {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "name and slug are required")
		return
	}

	category, err := h.categoryService.CreateCategory(r.Context(), dto)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeSuccess(w, http.StatusCreated, category)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var dto models.UpdateCategoryDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body")
		return
	}

	category, err := h.categoryService.UpdateCategory(r.Context(), id, dto)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeSuccess(w, http.StatusOK, category)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.categoryService.DeleteCategory(r.Context(), id); err != nil {
		writeAppError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
