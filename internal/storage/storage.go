package storage

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(dbPath string) (*Storage, error) {
	_, err := os.Stat(dbPath)
	install := os.IsNotExist(err)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if install {
		log.Println("Creating new database and table...")

		createTableQuery := `
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat TEXT
		);`
		_, err = db.Exec(createTableQuery)
		if err != nil {
			return nil, err
		}

		createIndexQuery := `
		CREATE INDEX idx_date ON scheduler(date);`
		_, err = db.Exec(createIndexQuery)
		if err != nil {
			return nil, err
		}

		log.Println("Database and table created successfully.")
	}

	return &Storage{DB: db}, nil
}

func (s *Storage) Close() error {
	return s.DB.Close()
}
