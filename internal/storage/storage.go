package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	DB *sql.DB
}

type Task struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
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

func (s *Storage) GetUpcomingTasks(limit int) ([]map[string]string, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
	rows, err := s.DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("error querying tasks: %v", err)
	}
	defer rows.Close()

	var tasks []map[string]string

	for rows.Next() {
		var task Task
		// Сканируем данные из базы в структуру Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			log.Printf("Error scanning task: %v", err)
			continue
		}

		// Преобразуем числовое поле ID в строку
		taskMap := map[string]string{
			"id":      strconv.FormatInt(task.ID, 10), // Преобразуем ID в строку
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		}

		tasks = append(tasks, taskMap)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Если задач нет, возвращаем пустой список
	if tasks == nil {
		tasks = []map[string]string{}
	}

	return tasks, nil
}

func (s *Storage) Close() error {
	return s.DB.Close()
}
