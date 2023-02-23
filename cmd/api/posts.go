package main

import (
	"NexusNet/internal/data"
	"NexusNet/internal/validator"
	"fmt"
	"net/http"
)

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	//Declare an anonymous struct to hold the information that we expect to be in the
	// HTTP request body (note that the field names and types in the struct are a subset
	// of the Movie struct that we created earlier). This struct will be our *target
	// decode destination*.
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		AuthorID    int64  `json:"author_id"`
	}

	// if there is error with decoding, we are sending corresponding message
	err := app.readJSON(w, r, &input) //non-nil pointer as the target decode destination
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
	}

	post := &data.Post{
		Title:       input.Title,
		Description: input.Description,
		AuthorID:    input.AuthorID,
	}

	err = app.models.Posts.Insert(post)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/posts/%d", post.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"post": post}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listPostsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title      string
		AuthorID   int
		AuthorName string
		Page       int
		PageSize   int
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()
	input.Title = app.readString(qs, "title", "")
	input.AuthorID = app.readInt(qs, "author_id", 0, v)
	input.AuthorName = app.readString(qs, "author_name", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.SortSafelist = []string{"title", "id", "-title", "-id", "author_id", "-author_id", "author_name", "-author_name"}

	input.Filters.Sort = app.readString(qs, "sort", "id")

	posts, _, err := app.models.Posts.GetAll(input.Title, input.AuthorID, input.AuthorName, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send a JSON response containing the movie data.
	err = app.writeJSON(w, http.StatusOK, envelope{"posts": posts}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
