package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"snippetbox.kavinddd.net/internal/models"
)

type application struct {
	errLog         *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {

	port := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:snippetbox@tcp(localhost:3306)/snippetbox?parseTime=true", "MYSQL data source name")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stdin, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(dsn)

	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()

	if err != nil {
		errLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		errLog:         errLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	mux := app.routes()

	infoLog.Printf("Server is running on port %s", *port)

	server := &http.Server{
		Addr:     *port,
		Handler:  mux,
		ErrorLog: errLog,
	}

	err = server.ListenAndServe()
	errLog.Fatal(err)
}

func openDB(dsn *string) (*sql.DB, error) {
	db, err := sql.Open("mysql", *dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
