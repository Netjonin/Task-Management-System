package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// previously returns *httprouter.Router
func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	// it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/tasks", app.requireActivatedUser(app.listTasksHandler))
	router.HandlerFunc(http.MethodPost, "/v1/tasks", app.requireActivatedUser(app.createTaskHandler))
	router.HandlerFunc(http.MethodGet, "/v1/tasks/:id", app.requireActivatedUser(app.showTaskHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/tasks/:id", app.requireActivatedUser(app.updateTaskHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/tasks/:id", app.requireActivatedUser(app.deleteTaskHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
