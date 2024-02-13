package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/rlr524/snippetboxv2/internal/models"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

type Config struct {
	addr string
	//staticDir string
	env string
	dsn string
}

type Application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	cfg            Config
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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

	// Initialize a new decoder instance
	formDecoder := form.NewDecoder()

	// Initialize a new session manager and configure it to use the MySQL database
	// as the session store and set a lifetime of 12 hours.
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	app := &Application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Config struct to hold the non-default TLS settings
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	srv := &http.Server{
		Addr:         cfg.addr,
		ErrorLog:     errorLog,
		Handler:      app.Routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", cfg.addr)
	e := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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

// TODO: Change input elements in signup and create to button elements

//func neuteredFileSystem(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if strings.HasSuffix(r.URL.Path, "/") {
//			http.NotFound(w, r)
//			return
//		}
//		next.ServeHTTP(w, r)
//	})
//}
