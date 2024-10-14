package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"todo-app/internal/handlers"
	"todo-app/internal/storage"
)

func main() {
	dbFile := getDBFilePath()

	dbStorage, err := storage.NewStorage(dbFile)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer dbStorage.Close()

	handler := &handlers.Handler{Storage: dbStorage}
	webDir := "./web"

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	fileServer := http.FileServer(http.Dir(webDir))

	http.Handle("/", fileServer)
	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	http.HandleFunc("/api/tasks", handler.TasksHandler)
	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.AddTaskHandler(dbStorage).ServeHTTP(w, r)
		case http.MethodGet:
			handlers.GetTaskHandler(dbStorage).ServeHTTP(w, r)
		case http.MethodPut:
			handlers.UpdateTaskHandler(dbStorage).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Printf("Starting server on port %s...", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getDBFilePath() string {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = filepath.Join("storage", "scheduler.db")
	}
	return dbFile
}
