package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"todo-app/internal/storage"
	"todo-app/internal/handlers"
)

func main() {
	dbFile := getDBFilePath()

	dbStorage, err := storage.NewStorage(dbFile)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer dbStorage.Close()

	webDir := "./web"

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	fileServer := http.FileServer(http.Dir(webDir))
	
	http.Handle("/", fileServer)
	http.HandleFunc("/api/nextdate", handlers.NextDateHandler)

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
