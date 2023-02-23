package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
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

	router.HandlerFunc(http.MethodGet, "/posts", app.listPostsHandler)
	router.HandlerFunc(http.MethodPost, "/posts/create", app.requireActivatedUser(app.createPostHandler))
	router.HandlerFunc(http.MethodPatch, "/posts/:id", app.requireActivatedUser(app.updatePostHandler))
	router.HandlerFunc(http.MethodDelete, "/posts/:id", app.requireActivatedUser(app.deletePostHandler))
	router.HandlerFunc(http.MethodGet, "/posts/:id", app.showPostHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}

//import (
//	"net/http"
//
//	"github.com/julienschmidt/httprouter"
//)
//
//func (app *application) routes() http.Handler {
//	// Initialize a new httprouter router instance.
//	router := httprouter.New()
//	router.NotFound = http.HandlerFunc(app.notFoundResponse)
//	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
//
//	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
//	// movie routes here
//	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requireActivatedUser(app.createPostHandler))
//	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
//	router.HandlerFunc(http.MethodPut, "/v1/movies/:id", app.updateMovieHandler)
//	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovieHandler)
//
//	// user routes here
//	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
//	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
//
//	// Add the route for the POST /v1/tokens/authentication endpoint.
//	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
//
//	// Return the http-router instance.
//	// wrapping the router with rateLimiter() middleware to limit requests' frequency
//	return app.recoverPanic(app.rateLimit(router))
//}
