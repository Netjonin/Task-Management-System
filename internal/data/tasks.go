package data

import (
	"database/sql"
	"errors"
	"time"

	"TMS.netjonin.net/internal/validator"
	//"github.com/lib/pq"
)

type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"-"`
	Status      string    `json:"status"`
	ExpiredAt   time.Time `json:"expired_at"`
	Expired     bool      `json:"expired"`
	Version     int32     `json:"version,string"`
}

func ValidateTask(v *validator.Validator, task *Task) {
	v.Check(task.Title != "", "title", "must be provided")
	v.Check(len(task.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(task.Description != "", "description", "must be provided")
	v.Check(!task.Expired, "expired", "newly created task should be active")
	v.Check(task.Status != "", "status", "newly created task should be in To-Do")
	v.Check(task.ExpiredAt.IsZero(), "expired_at", "must be provided")

	// v.Check(input.Year != 0, "year", "must be provided")
	// v.Check(input.Year >= 1888, "year", "must be greater than 1888")
	// v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	// v.Check(input.Runtime != 0, "runtime", "must be provided")
	// v.Check(input.Runtime > 0, "runtime", "must be a positive integer")
	// v.Check(input.Genres != nil, "genres", "must be provided")
	// v.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 genre")
	// v.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 genres")
	// v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")
}

type TaskModel struct {
	DB *sql.DB
}

func (t TaskModel) Insert(task *Task) error {

	query := `
      INSERT INTO tasks (title, description, status, expired_at, expired)
      VALUES ($1, $2, $3, $4, $5)
      RETURNING id, created_at, version`

	//args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}
	args := []any{task.Title, task.Description, task.Status, task.ExpiredAt, task.Expired}
	return t.DB.QueryRow(query, args...).Scan(&task.ID, &task.CreatedAt, &task.Version)
}

func (t TaskModel) Get(id int64) (*Task, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, title, description, status, created_at, expired_at, version
	FROM tasks
	WHERE id = $1`

	var task Task

	err := t.DB.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		//pq.Array(&movie.Genres),
		&task.ExpiredAt,
		&task.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &task, nil
}

func (t TaskModel) Update(task *Task) error {

	query := `
        UPDATE tasks
        SET title = $1, description = $2, status = $3, expired = $4, expired_at = $5, version = version + 1
        WHERE id = $6
        RETURNING version`

	args := []any{
		task.Title,
		task.Description,
		task.Status,
		//pq.Array(movie.Genres),
		task.Expired,
		task.ExpiredAt,
		task.ID,
	}
	return t.DB.QueryRow(query, args...).Scan(&task.Version)
}

func (t TaskModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
	DELETE FROM tasks
	WHERE id = $1`

	result, err := t.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
