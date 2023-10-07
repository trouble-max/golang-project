package main

import (
	"fmt"
	"gourmetspices.yerassyl.net/internal/validator"
	"net/http"
	"time"

	"gourmetspices.yerassyl.net/internal/data"
)

func (app *application) createHerbHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name         string     `json:"name"`
		Description  string     `json:"description"`
		Price        data.Price `json:"price"`
		CulinaryUses []string   `json:"culinary_uses"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	herb := &data.Herb{
		Name:         input.Name,
		Description:  input.Description,
		Price:        input.Price,
		CulinaryUses: input.CulinaryUses,
	}

	v := validator.New()

	if data.ValidateHerb(v, herb); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
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
