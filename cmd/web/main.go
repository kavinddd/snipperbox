package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"snippetbox.kavinddd.net/internal/models"
)

type application struct {
	errLog   *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

func main() {

	port := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:snippetbox@tcp(localhost:3306)/snippetbox?parseTime=true", "Database server name")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stdin, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(dsn)

	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errLog:   errLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
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
