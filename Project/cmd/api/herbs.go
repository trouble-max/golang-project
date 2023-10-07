package main

import (
	"fmt"
	"gourmetspices.yerassyl.net/internal/data"
	"net/http"
	"time"
)

func (app *application) createHerbHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new herb")
}

func (app *application) showHerbHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	herb := data.Herb{
		ID:           id,
		CreatedAt:    time.Time{},
		Name:         "Acacia Powder",
		Description:  "Acacia spp. is found across sub-Saharan Africa and is known for its gum resin. Acacia gum, or gum Arabic, is harvested, dried, then milled into a fine powder. Acacia powder is a popular addition to a multitude of products as it acts to bind formulas and stabilize emulsions.",
		Price:        8.25,
		CulinaryUses: []string{"emulsifier", "stabilizer", "thickener"},
		Version:      1,
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"herb": herb}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
