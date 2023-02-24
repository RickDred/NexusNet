package main

import (
	"NexusNet/internal/data"
	"fmt"
	"net/http"
)

func (app *application) createDirectHandler(w http.ResponseWriter, r *http.Request) {
	user2, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)

	direct := &data.Direct{
		User1: user.ID,
		User2: user2,
	}

	err = app.models.Directs.Insert(direct)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/direct/%d", direct.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"direct": direct}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showDirectsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	directs, err := app.models.Directs.GetAllFromUser(int(user.ID))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"directs": directs}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
