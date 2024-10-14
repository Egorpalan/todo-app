package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"todo-app/internal/storage"
)

type Handler struct {
	Storage *storage.Storage
}


func (h *Handler) TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			http.Error(w, "Invalid limit value", http.StatusBadRequest)
			return
		}
	}

	tasks, err := h.Storage.GetUpcomingTasks(limit)
	if err != nil {
		log.Printf("Error fetching tasks: %v", err)
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"tasks": tasks,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
