package main

import (
	"fmt"
	"net/http"
	

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

	task := &data.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Expired:     input.Expired,
	}

	v := validator.New()

	if data.ValidateTask(v, task); !v.Valid() {
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

	// task := data.Task{
	// 	ID:          id,
	// 	Title:       "Laundry",
	// 	Description: "Laundry",
	// 	CreatedAt:   time.Now(),
	// 	Status:      "To-Do",
	// 	ExpiredAt:   time.Now().Local().Add(time.Hour * 12),
	// 	Expired:     false,
	// 	Version:     1,
	// }

	for _, v := range store {
		if v.ID == id {
			task = &data.Task{
				ID: v.ID,
				Title: v.Title,
				Description: v.Description,
				CreatedAt: v.CreatedAt,
				Status: v.Status,
				ExpiredAt: v.ExpiredAt,
				Expired: v.Expired,
				Version: v.Version,
			}

		}
	}



	err = app.writeJSON(w, http.StatusOK, envelope{"task": task}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)

	}
}
