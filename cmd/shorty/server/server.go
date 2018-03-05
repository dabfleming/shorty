package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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
	s.mux.HandleFunc("/new", s.logMiddleware(s.newLinkHandler))
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
		<html>
		<head><title>Shorty</title></head>
		<body>
		<h2>Create Short Link</h2>
		<form method="post" action="/new">
		Full URL: <input type="text" name="url" value="https://" /><br />
		Requested short url (optional): http://SERVER_NAME/<input type="text" name="slug" /><br />
		<input type="submit" />
		</form>
		<h2>Debug Links</h2>
		<a href="/info/">info/</a><br />
		<a href="/goog">goog</a><br />
		<a href="/twitter">twitter</a><br />
		<a href="/fb">fb</a><br />
		<a href="/foo">foo (not found)</a><br />
		</body>
		</html>
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

func (s *Server) newLinkHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Check method
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse form
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error parsing form: %v", err)
		return
	}

	// Checks on URL
	url := r.PostForm.Get("url")
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Must include a URL.")
		return
	}
	if strings.HasPrefix(url, "http://") == false && strings.HasPrefix(url, "https://") == false {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "URL must begin with 'http://' or 'https://'.")
		return
	}

	// Check for requested slug
	slug := r.PostForm.Get("slug")
	if slug == "" {
		// TODO Generate a random slug
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprint(w, "TODO: Generate random slug")
		return
	}

	// Request to save to DB
	log.Printf("Requested '%v' link to '%v'.", slug, url) // TODO Remove dubug output
	err = s.ds.SaveNewURL(ctx, slug, url)
	if strings.HasPrefix(err.Error(), "Error 1062") {
		// Duplicate slug
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "Error, the short url '%v' is already in use.", slug)
		return
	}
	if err != nil {
		log.Printf("Error saving new short url: %T, %v", err, err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `Link created: <a href="/%v">/%v</a> now links to <pre>%v</pre>`, slug, slug, url)
}

func (s *Server) infoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "TODO: Implement info tools")
}
