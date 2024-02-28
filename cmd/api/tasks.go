package main

import (
	"fmt"
	"net/http"
	"time"

	"TMS.netjonin.net/internal/data"
	"TMS.netjonin.net/internal/validator"
)

var task *data.Task

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
	num := len(store) + 1
	task := &data.Task{
		ID:          int64(num),
		Title:       input.Title,
		Description: input.Description,
		CreatedAt:   time.Now(),
		Status:      input.Status,
		ExpiredAt:   time.Now().Add(time.Duration(4)),
		Expired:     input.Expired,
		Version:     1,
	}

	v := validator.New()

	if data.ValidateTask(v, task); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	store[num] = *task
	fmt.Fprintf(w, "%+v\n", len(store))

}

func (app *application) showTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	for _, v := range store {
		if v.ID == id {
			task = &data.Task{
				ID:          v.ID,
				Title:       v.Title,
				Description: v.Description,
				CreatedAt:   v.CreatedAt,
				Status:      v.Status,
				ExpiredAt:   v.ExpiredAt,
				Expired:     v.Expired,
				Version:     v.Version,
			}

		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"task": task}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)

	}
}
