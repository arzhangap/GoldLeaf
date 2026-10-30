package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "goland:pass@/goldleaf?parseTime=true", "MySQL data source name")

	// custom logs
	logInfo := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	logErr := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// open a db
	db, err := openDB(*dsn)
	if err != nil {
		logErr.Fatal(err)
		return
	}
	defer db.Close()

	//get template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logErr.Fatal(err)
	}

	app := application{
		infoLog:       logInfo,
		errorLog:      logErr,
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	log.Print("Starting server on :4000")
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	log.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(3 * time.Minute)

	return db, err
}
