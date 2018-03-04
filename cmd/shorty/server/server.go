package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dabfleming/shorty/internal/datastore"
)

var (
	urls map[string]string
)

// Server models our http server
type Server struct {
	mux *http.ServeMux
	ds  datastore.Datastore
}

// New returns a new server
func New(ds datastore.Datastore) (Server, error) {
	urls = map[string]string{
		"foo": "https://www.google.ca",
		"bar": "https://twitter.com",
	}

	s := Server{
		ds: ds,
	}

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
		<a href="/goog">goog</a><br />
		<a href="/twitter">twitter</a><br />
		<a href="/fb">fb</a><br />
		<a href="/foo">foo (not found)</a><br />
		`)
}

func (s *Server) forwardHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	url, err := s.ds.GetURLBySlug(ctx, r.URL.Path[1:])
	if err != nil {
		log.Printf("Error looking up url: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if url == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Server) infoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "TODO: Implement info tools")
}
