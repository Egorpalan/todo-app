package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"todo-app/internal/scheduler"
	"todo-app/internal/storage"
)

type TaskRequest struct {
	ID      string `json:"id"`
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

func GetTaskHandler(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		taskIDStr := r.URL.Query().Get("id")
		if taskIDStr == "" {
			http.Error(w, `{"error": "Не указан идентификатор"}`, http.StatusBadRequest)
			return
		}

		taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
		if err != nil {
			http.Error(w, `{"error": "Неверный формат идентификатора"}`, http.StatusBadRequest)
			return
		}

		task, err := db.GetTaskByID(taskID)
		if err != nil {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"id":      strconv.FormatInt(task.ID, 10), // Преобразуем ID в строку
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, `{"error": "Ошибка при отправке данных"}`, http.StatusInternalServerError)
		}
	}
}

func UpdateTaskHandler(db *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var task TaskRequest
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
			return
		}

		if task.ID == "" {
			http.Error(w, `{"error":"ID is required"}`, http.StatusBadRequest)
			return
		}

		// Проверяем, корректен ли ID
		taskID, err := strconv.ParseInt(task.ID, 10, 64)
		if err != nil {
			http.Error(w, `{"error":"Invalid ID format"}`, http.StatusBadRequest)
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

		// Проверка формата даты
		_, err = time.Parse("20060102", task.Date)
		if err != nil {
			http.Error(w, `{"error":"Invalid date format, expected YYYYMMDD"}`, http.StatusBadRequest)
			return
		}

		if task.Repeat != "" && !isValidRepeat(task.Repeat) {
			http.Error(w, `{"error":"Unsupported repeat type: must start with 'd' or 'y'"}`, http.StatusBadRequest)
			return
		}

		exists, err := db.TaskExists(taskID) 
		if err != nil || !exists {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			return
		}

		err = db.UpdateTask(taskID, task.Date, task.Title, task.Comment, task.Repeat)
		if err != nil {
			http.Error(w, `{"error": "Задача не найдена"}`, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}
}

func isValidRepeat(repeat string) bool {
	if len(repeat) < 2 {
		return false
	}

	if repeat[0] == 'd' || repeat[0] == 'y' {
		if len(repeat) >= 3 && repeat[1] == ' ' {
			if _, err := strconv.Atoi(repeat[2:]); err == nil {
				return true
			}
		}
	}

	return false
}
