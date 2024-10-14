package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"todo-app/internal/storage"
)

// Handler struct to hold the storage instance
type Handler struct {
	Storage *storage.Storage
}

// TasksHandler handles the /api/tasks endpoint
func (h *Handler) TasksHandler(w http.ResponseWriter, r *http.Request) {
	// Ожидаем GET-запрос
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметры limit (ограничение на количество задач, если передан)
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // Значение по умолчанию
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			http.Error(w, "Invalid limit value", http.StatusBadRequest)
			return
		}
	}

	// Получаем список задач из базы данных через хранилище
	tasks, err := h.Storage.GetUpcomingTasks(limit)
	if err != nil {
		log.Printf("Error fetching tasks: %v", err)
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		return
	}

	// Формируем ответ в формате JSON
	response := map[string]interface{}{
		"tasks": tasks,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
