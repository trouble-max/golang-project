package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/herbs", app.listHerbsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/herbs", app.createHerbHandler)
	router.HandlerFunc(http.MethodGet, "/v1/herbs/:id", app.showHerbHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/herbs/:id", app.updateHerbHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/herbs/:id", app.deleteHerbHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	return app.recoverPanic(app.rateLimit(router))
}
