package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tasks", app.createTaskHandler)
	router.HandlerFunc(http.MethodGet, "/v1/tasks/:id", app.showTaskHandler)
	return router
}
