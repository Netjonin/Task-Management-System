package main

import (
	"fmt"
	"net/http"
	"time"

	"TMS.netjonin.net/internal/data"
	"TMS.netjonin.net/internal/validator"
)

func (app *application) createTaskHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description,omitempty"`
		Status      string `json:"status"`
		Expired     bool   `json:"expired"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()

	v.Check(input.Title != "", "title", "must be provided")
	v.Check(len(input.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(input.Description != "", "description", "must be provided")
	v.Check(!input.Expired, "expired", "newly created task should be active")
	v.Check(input.Status != "", "status", "newly created task should be in To-Do")

	// v.Check(input.Year != 0, "year", "must be provided")
	// v.Check(input.Year >= 1888, "year", "must be greater than 1888")
	// v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	// v.Check(input.Runtime != 0, "runtime", "must be provided")
	// v.Check(input.Runtime > 0, "runtime", "must be a positive integer")
	// v.Check(input.Genres != nil, "genres", "must be provided")
	// v.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 genre")
	// v.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 genres")
	
	// v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)

}

func (app *application) showTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	task := data.Task{
		ID:          id,
		Title:       "Laundry",
		Description: "Laundry",
		CreatedAt:   time.Now(),
		Status:      "To-Do",
		ExpiredAt:   time.Now().Local().Add(time.Hour * 12),
		Expired:     false,
		Version:     1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"task": task}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)

	}
}
