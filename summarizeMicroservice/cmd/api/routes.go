package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/v1/summaries", app.insertSummaryHandler)
	router.HandlerFunc(http.MethodGet, "/v1/summaries/latest", app.getLatestSummaryHandler)
	router.HandlerFunc(http.MethodGet, "/v1/test/chatgpt", app.testChatgptHandler)
	return app.recoverPanic(router)
}
