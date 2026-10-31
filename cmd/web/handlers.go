package main

import (
	"errors"
	"net/http"
	"strings"

	//"strconv"

	"github.com/julienschmidt/httprouter"
	"goldleaf.ari.github.com/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "home.html", nil)
}

func (app *application) createPaste(w http.ResponseWriter, r *http.Request) {

	app.render(w, http.StatusOK, "create_paste.html", nil)
}

func (app *application) pasteView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())
	code := strings.ToUpper(params.ByName("code"))

	if len(code) != 6 {
		app.notFound(w)
		return
	}

	p, err := app.pastes.Read(code, "")
	switch {
	case errors.Is(err, models.ErrNoRecord):
		app.notFound(w)
	case errors.Is(err, models.ErrPasswordRequired):
	case err != nil:
		app.serverError(w, err)
	default:
		data := app.newTemplateData(r)
		data.Paste = p
		app.render(w, http.StatusOK, "read_paste.html", data)
	}

}
