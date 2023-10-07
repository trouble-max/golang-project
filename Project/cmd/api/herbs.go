package main

import (
	"fmt"
	"net/http"
)

func (app *application) createHerbHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new herb")
}

func (app *application) showHerbHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "show the details of herb %d\n", id)
}
