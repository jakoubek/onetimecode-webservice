package main

import (
	"github.com/go-chi/chi/v5"

	"net/http"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(app.recoverPanic)
	router.Use(app.rateLimit)
	router.Use(app.enableCORS)

	router.NotFound(http.HandlerFunc(app.notFoundResponse))
	router.MethodNotAllowed(http.HandlerFunc(app.methodNotAllowedResponse))

	router.Get("/", app.indexHandler())
	router.Get("/healthz", app.healthcheckHandler)

	router.Group(func(router chi.Router) {
		router.Use(app.logRequests)

		router.Get("/number", app.numberHandler)
		router.Get("/alphanumeric", app.alphanumericHandler)
		router.Get("/ksuid", app.ksuidHandler)
		router.Get("/uuid", app.uuidHandler)
		router.Get("/dice", app.diceHandler)
		router.Get("/coin", app.coinHandler)
	})

	return router
}
