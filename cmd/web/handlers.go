package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/k3vwdd/letsgo/internal/models"
	"github.com/k3vwdd/letsgo/internal/validator"
)

type snippetCreateForm struct {
    Title   string
    Content string
    Expires int
    FieldErrors map[string]string
    validator.Validator
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    snippets, err := app.snippets.Latest()
    if err != nil {
        app.serverError(w, r, err)
    }

    data := app.newTemplateData(r)
    data.Snippets = snippets

    app.render(w, r, http.StatusOK, "home.html", data )

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    snippet, err := app.snippets.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            http.NotFound(w, r)
        } else {
            app.serverError(w, r, err)
        }
        return
    }

    data := app.newTemplateData(r)
    data.Snippet = snippet

    app.render(w, r, http.StatusOK, "view.html", data)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    data := app.newTemplateData(r)
    data.Form = snippetCreateForm{
        Expires: 365,
    }

    app.render(w, r, http.StatusOK, "create.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
    }

    expires, err := strconv.Atoi(r.PostForm.Get("expires"))
    if err != nil {
        app.clientError(w, http.StatusBadRequest)
    }

    form := snippetCreateForm{
        Title: r.PostForm.Get("title"),
        Content: r.PostForm.Get("content"),
        Expires: expires,
    }

    form.CheckField(validator.NotBlank(form.Title), "title", "This field can't be blank")
    form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field can't have more then 100 characters")
    form.CheckField(validator.NotBlank(form.Content), "content", "This field can't be blank")
    form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must be equal to 1, 7, 365")

    if !form.Valid() {
        data := app.newTemplateData(r)
        data.Form = form
        app.render(w, r, http.StatusUnprocessableEntity, "create.html", data)
        return
    }

    id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Saved a new snippet...."))

}


