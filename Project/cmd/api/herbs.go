package main

import (
	"errors"
	"fmt"
	"gourmetspices.yerassyl.net/internal/data"
	"gourmetspices.yerassyl.net/internal/validator"
	"net/http"
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

	err = app.models.Herbs.Insert(herb)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/herbs/%d", herb.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"herb": herb}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showHerbHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	herb, err := app.models.Herbs.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"herb": herb}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateHerbHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	herb, err := app.models.Herbs.Get(id)
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
		Name         string     `json:"name"`
		Description  string     `json:"description"`
		Price        data.Price `json:"price"`
		CulinaryUses []string   `json:"culinary_uses"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	herb.Name = input.Name
	herb.Description = input.Description
	herb.Price = input.Price
	herb.CulinaryUses = input.CulinaryUses

	v := validator.New()
	if data.ValidateHerb(v, herb); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Herbs.Update(herb)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"herb": herb}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteHerbHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Herbs.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "herb successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
