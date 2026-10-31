package main

import (
	"net/http"
	//"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "home.html", nil)
}

func (app *application) createPaste(w http.ResponseWriter, r *http.Request) {

	app.render(w, http.StatusOK, "create_paste.html", nil)
}

func (app *application) pasteView(w http.ResponseWriter, r *http.Request) {

	app.render(w, http.StatusOK, "read_paste.html", nil)
}
