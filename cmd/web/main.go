package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	addr      string
	staticDir string
	env       string
}

type Application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	var cfg Config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.env, "env", "production", "Environment (development|staging|production)")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	app := &Application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", neuteredFileSystem(fileServer)))

	mux.HandleFunc("/", app.Home)
	mux.HandleFunc("/snippet/view", app.SnippetView)
	mux.HandleFunc("/snippet/create", app.SnippetCreate)

	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	// The value returned from the flag.String() function is a pointer to the flag value, not the value itself,
	// so it is required to dereference the pointer before using it.
	infoLog.Printf("Starting server on %s", cfg.addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

//type neuteredFileSystem struct {
//	fs http.FileSystem
//}
//
//func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
//	f, err := nfs.fs.Open(path)
//	if err != nil {
//		return nil, err
//	}
//
//	s, err := f.Stat()
//	if s.IsDir() {
//		index := filepath.Join(path, "index.html")
//		if _, err := nfs.fs.Open(index); err != nil {
//			closeErr := f.Close()
//			if closeErr != nil {
//				return nil, closeErr
//			}
//			return nil, err
//		}
//	}
//	return f, nil
//}

func neuteredFileSystem(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
