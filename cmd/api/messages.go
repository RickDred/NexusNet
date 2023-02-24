package main

import (
	"NexusNet/internal/data"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) showDirectMessagesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	direct, err := app.models.Directs.Get(id)

	user := app.contextGetUser(r)
	if user.ID != direct.User1 && user.ID != direct.User2 {
		app.noAccessRights(w, r)
		return
	}

	messages, err := app.models.Messages.GetAllFromDirect(int(direct.ID))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"messages": messages}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) writeMessageHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Content string `json:"content"`
	}

	user := app.contextGetUser(r)

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	// if there is error with decoding, we are sending corresponding message
	err = app.readJSON(w, r, &input) //non-nil pointer as the target decode destination
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
	}

	message := &data.Message{
		Content:  input.Content,
		DirectId: id,
		SenderId: user.ID,
	}

	err = app.models.Messages.Insert(message)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/messages/%d", message.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"message": message}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
