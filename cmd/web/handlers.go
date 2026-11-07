package main

import (
	"fmt"
	"net/http"
	"strings"

	//"strconv"

	"github.com/julienschmidt/httprouter"
	"goldleaf.ari.github.com/internal/validator"
)

type pasteForm struct {
	title   string
	body    string
	expires string
	validator.Validator
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "home.html", nil)
}

func (app *application) createPaste(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = pasteForm{expires: "1h"}
	app.render(w, http.StatusUnprocessableEntity, "create_paste.html", data)
}

func (app *application) createPastePost(w http.ResponseWriter, r *http.Request) {
	var form pasteForm
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.title = strings.TrimSpace(r.PostForm.Get("title"))
	form.body = strings.TrimSpace(r.PostForm.Get("body"))
	form.expires = strings.TrimSpace(r.PostForm.Get("expires"))

	form.CheckField(validator.NotBlank(form.title), "title", "عنوان نمی تواند خالی باشد")
	form.CheckField(validator.NotBlank(form.body), "body", "عنوان نمی تواند خالی باشد")
	form.CheckField(validator.MaxChars(form.title, 100), "title", "عنوان بیش از حد مجاز طولانی است.")
	form.CheckField(validator.MaxChars(form.body, 500), "body", "عنوان بیش از حد مجاز طولانی است.")
	form.CheckField(validator.PermittedValue(form.expires, "1h", "1d", "7d"), "expires", "تاریخ انقضا نادرست است")

	if !form.Validator.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create_paste.html", data)
		return
	}

	code, err := app.pastes.Insert(form.title, form.body, "", form.expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/paste/view/%s", code), http.StatusSeeOther)
}

func (app *application) pasteView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	code := params.ByName("code")

	app.showPaste(r, w, code, "")
}

func (app *application) pasteViewPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	code := r.PostForm.Get("code")

	app.showPaste(r, w, code, "")
}
