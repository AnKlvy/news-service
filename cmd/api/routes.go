package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	//router.HandlerFunc(http.MethodGet, "/v1/news", app.listNewsHandler)
	//router.HandlerFunc(http.MethodPost, "/v1/news", app.createNewsHandler)
	//router.HandlerFunc(http.MethodGet, "/v1/news/:id", app.showNewsHandler)
	//router.HandlerFunc(http.MethodPatch, "/v1/news/:id", app.updateNewsHandler)
	//router.HandlerFunc(http.MethodDelete, "/v1/news/:id", app.deleteNewsHandler)

	// Оборачиваем роутер в middleware rateLimit().
	return app.recoverPanic(app.rateLimit(router))
}
