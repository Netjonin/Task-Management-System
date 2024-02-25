package main

import (
	"fmt"
	"net/http"
	"time"

	"TMS.netjonin.net/internal/data"
)

func (app *application) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
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
