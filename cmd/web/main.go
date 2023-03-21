package main

import (
	"flag"
	"github.com/rlr524/snippetboxv2/handlers"
	"log"
	"net/http"
	"strings"
)

func main() {
	type config struct {
		addr      string
		staticDir string
	}

	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", neuteredFileSystem(fileServer)))

	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/snippet/view", handlers.SnippetView)
	mux.HandleFunc("/snippet/create", handlers.SnippetCreate)

	// The value returned from the flag.String() function is a pointer to the flag value, not the value itself,
	// so it is required to dereference the pointer before using it.
	log.Printf("Starting server on %s", cfg.addr)
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
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
