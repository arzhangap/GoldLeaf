package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"goldleaf.ari.github.com/ui"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		app.notFound(writer)
	})

	fileServer := http.FileServer(http.FS(ui.File))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	router.Handler(http.MethodGet, "/", http.HandlerFunc(app.home))
	router.Handler(http.MethodGet, "/paste/create", http.HandlerFunc(app.createPaste))
	router.Handler(http.MethodPost, "/paste/create", http.HandlerFunc(app.createPastePost))

	router.Handler(http.MethodGet, "/paste/view/:code", http.HandlerFunc(app.pasteView))
	router.Handler(http.MethodPost, "/paste/view", http.HandlerFunc(app.pasteViewPost))

	return router
}
