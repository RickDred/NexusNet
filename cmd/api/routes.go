package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method. Note that http.MethodGet and
	// http.MethodPost are constants which equate to the strings "GET" and "POST"
	// respectively.
	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/signup", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/login", app.authenticationUserHandler)
	router.HandlerFunc(http.MethodGet, "/users/:id", app.showUserHandler)

	router.HandlerFunc(http.MethodGet, "/", app.listPostsHandler)
	router.HandlerFunc(http.MethodPost, "/posts/create", app.requireActivatedUser(app.createPostHandler))
	router.HandlerFunc(http.MethodPatch, "/posts/:id", app.requireActivatedUser(app.updatePostHandler))
	router.HandlerFunc(http.MethodDelete, "/posts/:id", app.requireActivatedUser(app.deletePostHandler))
	router.HandlerFunc(http.MethodGet, "/posts/:id", app.showPostHandler)

	router.HandlerFunc(http.MethodGet, "/stories/:id", app.requireActivatedUser(app.showStoryHandler))
	router.HandlerFunc(http.MethodGet, "/users/:id/stories", app.requireActivatedUser(app.listUserStoriesHandler))
	router.HandlerFunc(http.MethodDelete, "/stories/:id", app.requireActivatedUser(app.deleteStoryHandler))
	router.HandlerFunc(http.MethodPost, "/stories", app.requireActivatedUser(app.createStoryHandler))

	router.HandlerFunc(http.MethodPost, "/post/:id/comments/create", app.requireActivatedUser(app.createCommentHandler))
	router.HandlerFunc(http.MethodPatch, "/post/:post_id/comments/:id", app.requireActivatedUser(app.updateCommentHandler))
	router.HandlerFunc(http.MethodGet, "/post/:id/comments", app.requireActivatedUser(app.listPostCommentsHandler))

	router.HandlerFunc(http.MethodGet, "/direct/:id", app.requireActivatedUser(app.showDirectMessagesHandler))
	router.HandlerFunc(http.MethodGet, "/direct", app.requireActivatedUser(app.showDirectsHandler))
	router.HandlerFunc(http.MethodPost, "/direct/:id", app.requireActivatedUser(app.writeMessageHandler))
	router.HandlerFunc(http.MethodPost, "/directs/:id", app.requireActivatedUser(app.createDirectHandler))

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}

// 2MICW67VDEHGMOQS5TEUX5MCZU
