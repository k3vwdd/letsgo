package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	//"github.com/k3vwdd/letsgo/internal/utils"
)

func home(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Server", "Go")

    files := []string{
        "./ui/html/base.html",
        "./ui/html/partials/nav.html",
        "./ui/html/pages/home.html",
    }

    ts, err := template.ParseFiles(files...)
    if err != nil {
        log.Fatal(err.Error())
        http.Error(w, "Internal Server error", http.StatusInternalServerError)
        return
    }

    err = ts.ExecuteTemplate(w, "base", nil)
    if err != nil {
        log.Fatal(err.Error())
        http.Error(w, "Internal Server error", http.StatusInternalServerError)
    }
}

func snippetView(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil || id < 1 {
        http.NotFound(w, r)
        return
    }

    fmt.Fprintf(w, "Display a specific snippet with id %d....", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Display a form for creating a new snippet...."))
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Saved a new snippet...."))

}


