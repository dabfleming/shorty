package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dabfleming/shorty/internal/datastore"
	"github.com/dabfleming/shorty/internal/slugs"
	"github.com/ua-parser/uap-go/uaparser"
)

const defaultSlugLength = 7

// Server models our http server
type Server struct {
	mux    *http.ServeMux
	ds     datastore.Datastore
	parser *uaparser.Parser
}

// New returns a new server
func New(ds datastore.Datastore, parser *uaparser.Parser) (Server, error) {
	s := Server{
		ds:     ds,
		parser: parser,
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
	fmt.Fprint(w, `<!DOCTYPE html>
		<html>
		<head><title>Shorty</title></head>
		<body>
		<h2>Create Short Link</h2>
		<form method="post" action="/new">
		Full URL: <input type="text" name="url" value="https://" /><br />
		Requested short url (optional): http://SERVER_NAME/<input type="text" name="slug" /><br />
		<input type="submit" />
		</form>
		<h2><a href="/info/">View Link Stats</a></h2>
		<h2>Debug/Testing Links</h2>
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

	if url.URL == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Track the visit
	ua := r.Header.Get("User-Agent")
	client := s.parser.Parse(ua)
	ip := r.RemoteAddr
	if idx := strings.Index(ip, ":"); idx != -1 {
		ip = ip[0:idx]
	}
	err = s.ds.TrackHit(ctx, url.ID, client, ip)
	if err != nil {
		log.Printf("Error tracking hit: %v", err)
	}

	w.Header().Set("Location", url.URL)
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
		// Generate a random slug
		// TODO Cope better with collisions
		slug = slugs.Random(defaultSlugLength)
	}

	// Request to save to DB
	log.Printf("Requested '%v' link to '%v'.", slug, url) // TODO Remove dubug output
	err = s.ds.SaveNewURL(ctx, slug, url)
	if err != nil && strings.HasPrefix(err.Error(), "Error 1062") {
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
	fmt.Fprintf(w, `<!DOCTYPE html>
		<html>
		<head><title>Shorty</title></head>
		<body>
		<h2>Link Created:</h2>
		<a href="/%v">/%v</a> now links to: <pre>%v</pre>
		</body></html>`, slug, slug, url)
}

func (s *Server) infoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slug := strings.TrimPrefix(r.URL.Path, "/info/")

	if len(slug) > 0 {
		s.infoDetailHandler(w, r)
		return
	}

	vs, err := s.ds.GetVisitCounts(ctx)
	if err != nil {
		log.Printf("Error getting visitor counts: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, `<!DOCTYPE html>
		<html>
		<head><title>Shorty</title></head>
		<body>
		<h2>Visits:</h2>
		<table border="2">
		<tr><th>Short URL</th><th>Full URL</th><th>Visit Count</th></tr>
		`)
	for _, v := range vs {
		fmt.Fprintf(w, `<tr><td><a href="/info/%v">%v</a></td><td>%v</td><td>%v</td></tr>`, v.Slug, v.Slug, v.URL, v.Count)
	}
	fmt.Fprint(w, `</table>
		</body>
		</html>
		`)
}

func (s *Server) infoDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slug := strings.TrimPrefix(r.URL.Path, "/info/")

	log.Printf("Lookup on slug: %v", slug)
	url, visits, err := s.ds.GetVisits(ctx, slug)
	if err != nil {
		log.Printf("Error getting visits: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, `<!DOCTYPE html>
		<html>
		<head><title>Shorty</title></head>
		<body>
		<h2>Visits to /%v</h2>
		<pre>%v</pre>
		<p>Visit detail, most recent first.</p>
		<table border="2">
		<tr><th>Device</th><th>OS</th><th>Browser</th><th>IP</th><th>Time</th></tr>
		`, url.Slug, url.URL)
	for _, v := range visits {
		fmt.Fprintf(w, `<tr><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td></tr>`, v.Device, v.OS, v.Browser, v.IP, v.Time)
	}
	fmt.Fprint(w, `</table>
		</body>
		</html>
		`)
}
