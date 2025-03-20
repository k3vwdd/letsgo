package main


import (
    "net/http"
	"github.com/go-chi/chi/v5"

)

func (app *application) routes() *chi.Mux {

    r := chi.NewRouter()

    fileserver := http.FileServer(http.Dir("./ui/static/"))
    r.Handle("GET /static/*", http.StripPrefix("/static", fileserver))

    r.Get("/", app.home)
    r.Get("/snippet/view/{id}", app.snippetView)
    r.Get("/snippet/create", app.snippetCreate)
    r.Post("/snippet/create", app.snippetCreatePost)

    return r
}
