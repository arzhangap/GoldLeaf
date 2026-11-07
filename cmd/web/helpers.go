package main

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"goldleaf.ari.github.com/internal/models"
)

func (app *application) notFound(w http.ResponseWriter) {
	app.render(w, http.StatusNotFound, "not_found.html", nil)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Print(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		//Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		//CSRFToken:       nosurf.Token(r),
	}
}

func (app *application) showPaste(r *http.Request, w http.ResponseWriter, code, password string) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if len(code) != 6 {
		app.notFound(w)
		return
	}

	p, err := app.pastes.Read(code, password)
	switch {
	case errors.Is(err, models.ErrNoRecord):
		app.notFound(w)
	case errors.Is(err, models.ErrPasswordRequired):
	case errors.Is(err, models.ErrInvalidCredentials):
	case err != nil:
		app.serverError(w, err)
	default:
		data := app.newTemplateData(r)
		data.Paste = p
		app.render(w, http.StatusOK, "read_paste.html", data)
	}
}
