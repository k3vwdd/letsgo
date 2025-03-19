package main

import (
	"log"
	"net/http"

	//"github.com/k3vwdd/letsgo/internal/utils"

	"github.com/go-chi/chi/v5"
)

func main() {

    r := chi.NewRouter()

    fileserver := http.FileServer(http.Dir("./ui/static/"))
    r.Handle("GET /static/*", http.StripPrefix("/static", fileserver))

    r.Get("/", home)
    r.Get("/snippet/view/{id}", snippetView)
    r.Get("/snippet/create", snippetCreate)
    r.Post("/snippet/create", snippetCreatePost)

    log.Print("server started on port :4000")
    err := http.ListenAndServe(":4000", r)
    if err != nil {
        log.Fatal(err)
    }

}

