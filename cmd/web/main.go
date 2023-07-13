package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/rlr524/snippetboxv2/internal/models"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Config struct {
	addr string
	//staticDir string
	env string
	dsn string
}

type Application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	cfg           Config
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	var cfg Config

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbPass := os.Getenv("DB_PASS")
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	//flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static/", "Path to static assets")
	flag.StringVar(&cfg.env, "env", "production", "Environment (development|staging|production)")
	flag.StringVar(&cfg.dsn, "dsn", fmt.Sprintf("web:%s@tcp(lancer:3306)/snippetbox?parseTime=true", dbPass),
		"MySQL data source")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	db, err := openDB(cfg.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			errorLog.Print(err)
		}
	}(db)

	// Init a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &Application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  app.Routes(),
	}

	infoLog.Printf("Starting server on %s", cfg.addr)
	e := srv.ListenAndServe()
	errorLog.Fatal(e)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func NeuteredFileSystem(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
