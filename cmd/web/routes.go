package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Setting up file server
	fileServer := http.FileServer(http.Dir(app.config.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// setting up routing and request handling
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippedView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return app.logRequest(
		secureHeaders(
			mux))
}
