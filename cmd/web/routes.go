package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// example, vanilla router + middleware
func (app *application) oldRoutes() http.Handler {
	mux := http.NewServeMux()
	httprouter.New()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}

// example, minimal lightweight quality of life router + middleware using httprouter + alice
func (app *application) routes() http.Handler {
	router := httprouter.New()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	router.Handler(http.MethodGet, "/static/*filePath", http.StripPrefix("static", fileServer))
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/view", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	// middleware
	baseMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return baseMiddleware.Then(router)
}
