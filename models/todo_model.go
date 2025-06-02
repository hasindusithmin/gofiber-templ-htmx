package models

import (
	"database/sql"
	"errors"
	"time"
)

type Todo struct {
	ID          uint64    `json:"id"`
	CreatedBy   uint64    `json:"created_by"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Status      bool      `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

// DB should be passed or set externally for cleaner architecture.
var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func (t *Todo) GetAllTodos() ([]Todo, error) {
	query := `SELECT id, title, description, status FROM todos WHERE created_by = ? ORDER BY created_at DESC`
	rows, err := db.Query(query, t.CreatedBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (t *Todo) GetNoteById() (*Todo, error) {
	query := `SELECT id, title, description, status, created_at FROM todos WHERE created_by = ? AND id = ?`
	row := db.QueryRow(query, t.CreatedBy, t.ID)

	var todo Todo
	err := row.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &todo, nil
}

func (t *Todo) CreateTodo() (*Todo, error) {
	query := `INSERT INTO todos (created_by, title, description) VALUES (?, ?, ?) RETURNING id, created_by, title, description, status, created_at`
	row := db.QueryRow(query, t.CreatedBy, t.Title, t.Description)

	var newTodo Todo
	if err := row.Scan(&newTodo.ID, &newTodo.CreatedBy, &newTodo.Title, &newTodo.Description, &newTodo.Status, &newTodo.CreatedAt); err != nil {
		return nil, err
	}
	return &newTodo, nil
}

func (t *Todo) UpdateTodo() (*Todo, error) {
	query := `UPDATE todos SET title = ?, description = ?, status = ? WHERE created_by = ? AND id = ? RETURNING id, title, description, status`
	row := db.QueryRow(query, t.Title, t.Description, t.Status, t.CreatedBy, t.ID)

	var updatedTodo Todo
	if err := row.Scan(&updatedTodo.ID, &updatedTodo.Title, &updatedTodo.Description, &updatedTodo.Status); err != nil {
		return nil, err
	}
	return &updatedTodo, nil
}

func (t *Todo) DeleteTodo() error {
	query := `DELETE FROM todos WHERE created_by = ? AND id = ?`
	result, err := db.Exec(query, t.CreatedBy, t.ID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil || affected != 1 {
		return errors.New("expected exactly one row to be affected")
	}
	return nil
}

func ConvertDateTime(tz string, dt time.Time) string {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return dt.Format(time.RFC822Z) // fallback
	}
	return dt.In(loc).Format(time.RFC822Z)
}
