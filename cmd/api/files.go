package main

import (
	"net/http"
)

func (app *application) createFileHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Data string `json:"data"`
		Name string `json:"name"`
	}

	err := app.createFile([]byte(input.Data), "./internal/files/"+input.Name)
	if err != nil {
		app.logger.PrintError(err, nil)
	}
	headers := make(http.Header)
	headers.Set("Location", "")

	err = app.writeJSON(w, http.StatusCreated, envelope{"data": input.Data}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) readFileHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}

	data, err := app.readFile("./internal/files/" + input.Name)
	if err != nil {
		app.logger.PrintError(err, nil)
	}
	headers := make(http.Header)
	headers.Set("Location", "")

	err = app.writeJSON(w, http.StatusCreated, envelope{"data": data}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
