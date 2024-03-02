package data

import (
	"database/sql"
	"time"

	"TMS.netjonin.net/internal/validator"
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
	return nil
}

func (t TaskModel) Get(id int64) (*Task, error) {
	return nil, nil
}

func (t TaskModel) Update(task *Task) error {
	return nil
}

func (t TaskModel) Delete(id int64) error {
	return nil
}


