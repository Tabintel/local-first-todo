package database

import (
	"database/sql"
	_ "modernc.org/sqlite" // Import the sqlite driver
)

var db *sql.DB

// Init initializes the SQLite database connection.
func Init(dbPath string) error {
	var err error
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	// Ensure the database connection is successful
	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

// CreateTaskTable creates the tasks table in the database if it does not exist.
func CreateTaskTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			done BOOLEAN DEFAULT FALSE
		)
	`
	_, err := db.Exec(query)
	return err
}

// GetTasks retrieves all tasks from the database.
func GetTasks() ([]Task, error) {
	rows, err := db.Query("SELECT id, content, done FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Content, &task.Done)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// AddTask inserts a new task into the database.
func AddTask(content string) error {
	_, err := db.Exec("INSERT INTO tasks (content) VALUES (?)", content)
	return err
}

// Task represents a task in the database.
type Task struct {
	ID     int    `json:"id"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}
