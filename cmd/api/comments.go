package main

import (
	"NexusNet/internal/data"
	"errors"
	"fmt"
	"net/http"
)

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Content string `json:"content"`
	}

	postId, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}

	// if there is error with decoding, we are sending corresponding message
	err = app.readJSON(w, r, &input) //non-nil pointer as the target decode destination
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
	}

	comment := &data.Comment{
		AuthorID: user.ID,
		PostID:   postId,
		Content:  input.Content,
	}

	err = app.models.Comments.Insert(comment)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/stories/%d", comment.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"comment": comment}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listPostCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postId, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	comments, err := app.models.Comments.GetAllFromPost(int(postId))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"comments": comments}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Comments.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "comment successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Retrieve the movie record as normal.
	comment, err := app.models.Comments.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Use pointers for the Title, Year and Runtime fields.
	var input struct {
		Content *string `json:"content"`
	}

	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}
	if comment.AuthorID != user.ID {
		app.noAccessRights(w, r)
		return
	}

	// Decode the JSON as normal.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Content != nil {
		comment.Content = *input.Content
	}

	err = app.models.Comments.Update(comment)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"comment": comment}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
