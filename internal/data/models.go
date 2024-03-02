package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Tasks TaskModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Tasks: TaskModel{DB: db},
	}
}
