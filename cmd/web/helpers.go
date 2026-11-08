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
	// TODO: Replace the numeric Print argument with proper formatted or structured logging.
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

	// TODO: Execute the template into a buffer first; only write the status and body after execution succeeds.
	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		// TODO: Add CSRF token generation here when CSRF middleware is enabled.
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
		// TODO: Render a password-entry form and preserve the paste code.
	case errors.Is(err, models.ErrInvalidCredentials):
		// TODO: Re-render the password-entry form with an invalid-password message.
	case err != nil:
		app.serverError(w, err)
	default:
		data := app.newTemplateData(r)
		data.Paste = p
		app.render(w, http.StatusOK, "read_paste.html", data)
	}
}
