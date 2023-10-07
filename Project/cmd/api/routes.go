package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {

	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/herbs", app.createHerbHandler)
	router.HandlerFunc(http.MethodGet, "/v1/herbs/:id", app.showHerbHandler)

	return router
}
