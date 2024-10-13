package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"todo-app/internal/scheduler"
	"todo-app/internal/storage"
)

type TaskRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TaskResponse struct {
	ID    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func AddTaskHandler(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var task TaskRequest
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			http.Error(w, `{"error":"Title is required"}`, http.StatusBadRequest)
			return
		}

		now := time.Now().Format("20060102")
		if task.Date == "" {
			task.Date = now
		}

		_, err = time.Parse("20060102", task.Date)
		if err != nil {
			http.Error(w, `{"error":"Invalid date format, expected YYYYMMDD"}`, http.StatusBadRequest)
			return
		}

		if task.Repeat != "" {
			if task.Repeat[0] == 'm' || task.Repeat[0] == 'w' {
				http.Error(w, `{"error":"Unsupported repeat type: only daily and yearly repeats are allowed"}`, http.StatusBadRequest)
				return
			}
		}

		today := time.Now().Format("20060102")
		if task.Date < today {
			if task.Repeat != "" {
				nextDate, err := scheduler.NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					http.Error(w, fmt.Sprintf(`{"error":"Invalid repeat rule: %v"}`, err), http.StatusBadRequest)
					return
				}
				task.Date = nextDate
			} else {
				task.Date = today
			}
		}

		if task.Date == today && task.Repeat == "d 1" {
			task.Date = today
		}

		query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
		res, err := db.DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
		if err != nil {
			http.Error(w, `{"error":"Failed to insert task"}`, http.StatusInternalServerError)
			return
		}

		id, err := res.LastInsertId()
		if err != nil {
			http.Error(w, `{"error":"Failed to retrieve task ID"}`, http.StatusInternalServerError)
			return
		}

		response := TaskResponse{ID: id}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
