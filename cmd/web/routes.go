package main


import (
    "net/http"
	"github.com/go-chi/chi/v5"
    "github.com/justinas/alice"

)

func (app *application) routes() http.Handler {

    r := chi.NewRouter()
    dynamic := alice.New(app.sessionManager.LoadAndSave)

    fileserver := http.FileServer(http.Dir("./ui/static/"))
    r.Handle("GET /static/*", http.StripPrefix("/static", fileserver))
    r.Handle("GET /", dynamic.ThenFunc(app.home))
    r.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
    r.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
    r.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))


    standardMW := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

    return standardMW.Then(r)

}
