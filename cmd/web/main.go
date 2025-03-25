package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/k3vwdd/letsgo/internal/models"
)

type application struct {
    logger          *slog.Logger
    snippets        *models.SnippetModel
    templateCache   map[string]*template.Template
    formDecoder     *form.Decoder
    sessionManager  *scs.SessionManager
}

func openDB(dsn string) (*sql.DB, error) {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        db.Close()
        return nil, err
    }

    return db, nil
}

func main() {

    addr := flag.String("addr", ":4000", "HTTP Network address")
    dsn := flag.String("dsn", "web:mysql!@/snippetbox?parseTime=true", "MySQL data source name")
    flag.Parse()

    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    db, err := openDB(*dsn)
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }
    defer db.Close()

    templateCache, err := newTemplateCache()
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }

    formDecoder := form.NewDecoder()

    sessionManger := scs.New()
    sessionManger.Store = mysqlstore.New(db)
    sessionManger.Lifetime = 12 * time.Hour

    app := &application{
        logger: logger,
        snippets: &models.SnippetModel{DB: db},
        templateCache: templateCache,
        formDecoder:   formDecoder,
        sessionManager: sessionManger,
    }

    serv := &http.Server{
        Addr: *addr,
        Handler: app.routes(),
        ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
    }

    logger.Info("starting server", "addr", *addr)

    err = serv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
    logger.Error(err.Error())
    os.Exit(1)

}

