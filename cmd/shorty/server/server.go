package server

import (
	"fmt"
	"log"
	"net/http"
)

var (
	urls map[string]string
)

// Server models our http server
type Server struct {
	mux *http.ServeMux
}

// New returns a new server
func New() (Server, error) {
	urls = map[string]string{
		"foo": "https://www.google.ca",
		"bar": "https://twitter.com",
	}

	s := Server{}

	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/info/", s.logMiddleware(s.infoHandler))
	s.mux.HandleFunc("/", s.logMiddleware(s.routerHandler))

	return s, nil
}

// Go runs the server
func (s *Server) Go() {
	log.Fatal(http.ListenAndServe(":8080", s.mux))
}

func (s *Server) logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v\n", r.Method, r.URL.Path)
		next(w, r)
	}
}

// routerHandler routes to the correct handler for short urls or the root page
func (s *Server) routerHandler(w http.ResponseWriter, r *http.Request) {
	// If this is not a request for /, assume it's a short URL and forward
	if len(r.URL.Path) > 1 {
		s.forwardHandler(w, r)
		return
	}

	// request for /, serve up some links for testing
	fmt.Fprint(w, `
		<a href="/info/">info/</a><br />
		<a href="/foo">foo</a><br />
		<a href="/bar">bar</a><br />
		<a href="/baz">baz</a><br />
		`)
}

func (s *Server) forwardHandler(w http.ResponseWriter, r *http.Request) {
	url := s.lookupShortURL(r.URL.Path[1:])
	if url == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Server) lookupShortURL(short string) string {
	for key, value := range urls {
		if key == short {
			return value
		}
	}
	return ""
}

func (s *Server) infoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "TODO: Implement info tools")
}
