package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"TMS.netjonin.net/internal/data"
	"TMS.netjonin.net/internal/validator"
)

//var task *data.Task

func (app *application) createTaskHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title       string    `json:"title"`
		Description string    `json:"description,omitempty"`
		Status      string    `json:"status"`
		ExpiredAt   time.Time `json:"expired_at"`
		Expired     bool      `json:"expired"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	//num := len(store) + 1
	task := &data.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		ExpiredAt:   input.ExpiredAt,
		Expired:     input.Expired,
	}

	v := validator.New()

	if data.ValidateTask(v, task); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// store[num] = *task
	// fmt.Fprintf(w, "%+v\n", len(store))

	err = app.models.Tasks.Insert(task)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/tasks/%d", task.ID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"task": task}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// for _, v := range store {
	// 	if v.ID == id {
	// 		task = &data.Task{
	// 			ID:          v.ID,
	// 			Title:       v.Title,
	// 			Description: v.Description,
	// 			CreatedAt:   v.CreatedAt,
	// 			Status:      v.Status,
	// 			ExpiredAt:   v.ExpiredAt,
	// 			Expired:     v.Expired,
	// 			Version:     v.Version,
	// 		}

	// 	}
	// }

	task, err := app.models.Tasks.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"task": task}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	task, err := app.models.Tasks.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Status      string    `json:"status"`
		Expired     bool      `json:"expired"`
		ExpiredAt   time.Time `json:"expired_at"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Copy the values from the request body to the appropriate fields of the movie
	// record.
	task.Title = input.Title
	task.Description = input.Description
	task.Status = input.Status
	task.Expired = input.Expired
	task.ExpiredAt = input.ExpiredAt

	v := validator.New()

	if data.ValidateTask(v, task); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Tasks.Update(task)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"task": task}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
