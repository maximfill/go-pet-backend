package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"github.com/go-chi/chi/v5"

	todoservice "github.com/maximfill/go-pet-backend/internal/service/todo"
)

type createTodoRequest struct {
	Title string `json:"title"`
}

type TodoHandler struct {
	service *todoservice.Service
}

func NewTodoHandler(service *todoservice.Service) *TodoHandler {
	return &TodoHandler{service: service}
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest

	// 1. Читаем JSON из body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// 2. Вызываем service
	// временно: userID = 1
	id, err := h.service.CreateTodo(r.Context(), 1, req.Title)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// 3. Отдаём HTTP-ответ
	w.WriteHeader(http.StatusCreated) // 201
	fmt.Fprintf(w, "todo created with id %d", id)
}

func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.GetTodosByUser(r.Context(), 1)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return

	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

type updateTodoRequest struct {
	Completed bool `json:"completed"`
}

func (h *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req updateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.service.SetCompleted(r.Context(), id, req.Completed); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteTodo(r.Context(), id); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}