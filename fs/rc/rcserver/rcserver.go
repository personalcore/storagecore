// Package rcserver implements the HTTP endpoint to serve the remote control
package rcserver

import (
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/personalcore/storagecore/fs"
	"github.com/personalcore/storagecore/fs/cache"
	"github.com/personalcore/storagecore/fs/config"
	"github.com/personalcore/storagecore/fs/list"
	"github.com/personalcore/storagecore/fs/rc"
	"github.com/personalcore/storagecore/fs/rc/jobs"
	libhttp "github.com/personalcore/storagecore/lib/http"
	"github.com/personalcore/storagecore/lib/http/serve"
)

// Start the remote control server if configured
//
// If the server wasn't configured the *Server returned may be nil
func Start(ctx context.Context, opt *rc.Options) (*Server, error) {
	jobs.SetOpt(opt) // set the defaults for jobs
	if opt.Enabled {
		// Serve on the DefaultServeMux so can have global registrations appear
		s, err := newServer(ctx, opt, http.DefaultServeMux)
		if err != nil {
			return nil, err
		}
		return s, s.Serve()
	}
	return nil, nil
}

// Server contains everything to run the rc server
type Server struct {
	ctx            context.Context // for global config
	server         *libhttp.Server
	files          http.Handler
	pluginsHandler http.Handler
	opt            *rc.Options
}

func newServer(ctx context.Context, opt *rc.Options, mux *http.ServeMux) (*Server, error) {
	fileHandler := http.Handler(nil)
	pluginsHandler := http.Handler(nil)
	// Add some more mime types which are often missing
	_ = mime.AddExtensionType(".wasm", "application/wasm")
	_ = mime.AddExtensionType(".js", "application/javascript")

	// File handling
	if opt.Files != "" {
		fs.Logf(nil, "Serving files from %q", opt.Files)
		fileHandler = http.FileServer(http.Dir(opt.Files))
	}

	s := &Server{
		ctx:            ctx,
		opt:            opt,
		files:          fileHandler,
		pluginsHandler: pluginsHandler,
	}

	var err error
	s.server, err = libhttp.NewServer(ctx,
		libhttp.WithConfig(opt.HTTP),
		libhttp.WithAuth(opt.Auth),
		libhttp.WithTemplate(opt.Template),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init server: %w", err)
	}

	router := s.server.Router()
	router.Use(
		middleware.SetHeader("Accept-Ranges", "bytes"),
		middleware.SetHeader("Server", "rclone/"+fs.Version),
	)

	// // Add the debug handler which is installed in the default mux
	// router.Handle("/debug/pprof/*", mux)

	// FIXME split these up into individual functions
	router.Get("/*", s.handler)
	router.Head("/*", s.handler)
	router.Post("/*", s.handler)
	router.Options("/*", s.handler)

	return s, nil
}

// Serve runs the http server in the background.
//
// Use s.Close() and s.Wait() to shutdown server
func (s *Server) Serve() error {
	s.server.Serve()
	return nil
}

// writeError writes a formatted error to the output
func writeError(path string, in rc.Params, w http.ResponseWriter, err error, status int) {
	fs.Errorf(nil, "rc: %q: error: %v", path, err)
	params, status := rc.Error(path, in, err, status)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err = rc.WriteJSON(w, params)
	if err != nil {
		// can't return the error at this point
		fs.Errorf(nil, "rc: writeError: failed to write JSON output from %#v: %v", in, err)
	}
}

// handler reads incoming requests and dispatches them
func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimLeft(r.URL.Path, "/")

	switch r.Method {
	case "POST":
		s.handlePost(w, r, path)
	case "OPTIONS":
		s.handleOptions(w, r, path)
	case "GET", "HEAD":
		s.handleGet(w, r, path)
	default:
		writeError(path, nil, w, fmt.Errorf("method %q not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}
}

func (s *Server) handlePost(w http.ResponseWriter, r *http.Request, path string) {
	ctx := r.Context()
	contentType := r.Header.Get("Content-Type")

	var (
		contentTypeMediaType string
		contentTypeParams    map[string]string
	)
	if contentType != "" {
		var err error
		contentTypeMediaType, contentTypeParams, err = mime.ParseMediaType(contentType)
		if err != nil {
			writeError(path, nil, w, fmt.Errorf("failed to parse Content-Type: %w", err), http.StatusBadRequest)
			return
		}
	}

	values := r.URL.Query()
	if contentTypeMediaType == "application/x-www-form-urlencoded" {
		// Parse the POST and URL parameters into r.Form, for others r.Form will be empty value
		err := r.ParseForm()
		if err != nil {
			writeError(path, nil, w, fmt.Errorf("failed to parse form/URL parameters: %w", err), http.StatusBadRequest)
			return
		}
		values = r.Form
	}

	// Read the POST and URL parameters into in
	in := make(rc.Params)
	for k, vs := range values {
		if len(vs) > 0 {
			in[k] = vs[len(vs)-1]
		}
	}

	// Parse a JSON blob from the input
	if contentTypeMediaType == "application/json" {
		// Check the charset is utf-8 or unset
		if charset, ok := contentTypeParams["charset"]; ok && !strings.EqualFold(charset, "utf-8") {
			writeError(path, in, w, fmt.Errorf("unsupported charset %q for JSON input", charset), http.StatusBadRequest)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&in)
		if err != nil {
			writeError(path, in, w, fmt.Errorf("failed to read input JSON: %w", err), http.StatusBadRequest)
			return
		}
	}
	// Find the call
	call := rc.Calls.Get(path)
	if call == nil {
		writeError(path, in, w, fmt.Errorf("couldn't find method %q", path), http.StatusNotFound)
		return
	}

	// Check to see if it requires authorisation
	if !s.opt.NoAuth && call.AuthRequired && !s.server.UsingAuth() {
		writeError(path, in, w, fmt.Errorf("authentication must be set up on the rc server to use %q or the --rc-no-auth flag must be in use", path), http.StatusForbidden)
		return
	}

	inOrig := in.Copy()

	if call.NeedsRequest {
		// Add the request to RC
		in["_request"] = r
	}

	if call.NeedsResponse {
		in["_response"] = w
	}

	fs.Debugf(nil, "rc: %q: with parameters %+v", path, in)
	job, out, err := jobs.NewJob(ctx, call.Fn, in)
	if job != nil {
		w.Header().Add("x-rclone-jobid", fmt.Sprintf("%d", job.ID))
	}
	if err != nil {
		writeError(path, inOrig, w, err, http.StatusInternalServerError)
		return
	}
	if out == nil {
		out = make(rc.Params)
	}

	fs.Debugf(nil, "rc: %q: reply %+v: %v", path, out, err)
	w.Header().Set("Content-Type", "application/json")
	err = rc.WriteJSON(w, out)
	if err != nil {
		// can't return the error at this point - but have a go anyway
		writeError(path, inOrig, w, err, http.StatusInternalServerError)
		fs.Errorf(nil, "rc: handlePost: failed to write JSON output: %v", err)
	}
}

func (s *Server) handleOptions(w http.ResponseWriter, r *http.Request, path string) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) serveRoot(w http.ResponseWriter, r *http.Request) {
	remoteNames := config.GetRemoteNames()
	sort.Strings(remoteNames)
	directory := serve.NewDirectory("", s.server.HTMLTemplate())
	directory.Name = "List of all rclone remotes."
	q := url.Values{}
	for _, remoteName := range remoteNames {
		q.Set("fs", remoteName)
		directory.AddHTMLEntry("["+remoteName+":]", true, -1, time.Time{})
	}
	sortParm := r.URL.Query().Get("sort")
	orderParm := r.URL.Query().Get("order")
	directory.ProcessQueryParams(sortParm, orderParm)

	directory.Serve(w, r)
}

func (s *Server) serveRemote(w http.ResponseWriter, r *http.Request, path string, fsName string) {
	f, err := cache.Get(s.ctx, fsName)
	if err != nil {
		writeError(path, nil, w, fmt.Errorf("failed to make Fs: %w", err), http.StatusInternalServerError)
		return
	}
	if path == "" || strings.HasSuffix(path, "/") {
		path = strings.Trim(path, "/")
		entries, err := list.DirSorted(r.Context(), f, false, path)
		if err != nil {
			writeError(path, nil, w, fmt.Errorf("failed to list directory: %w", err), http.StatusInternalServerError)
			return
		}
		// Make the entries for display
		directory := serve.NewDirectory(path, s.server.HTMLTemplate())
		for _, entry := range entries {
			_, isDir := entry.(fs.Directory)
			var modTime time.Time
			if !s.opt.ServeNoModTime {
				modTime = entry.ModTime(r.Context())
			}
			directory.AddHTMLEntry(entry.Remote(), isDir, entry.Size(), modTime)
		}
		sortParm := r.URL.Query().Get("sort")
		orderParm := r.URL.Query().Get("order")
		directory.ProcessQueryParams(sortParm, orderParm)

		directory.Serve(w, r)
	} else {
		path = strings.Trim(path, "/")
		o, err := f.NewObject(r.Context(), path)
		if err != nil {
			writeError(path, nil, w, fmt.Errorf("failed to find object: %w", err), http.StatusInternalServerError)
			return
		}
		serve.Object(w, r, o)
	}
}

// Match URLS of the form [fs]/remote
var fsMatch = regexp.MustCompile(`^\[(.*?)\](.*)$`)

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request, path string) {
	// Look to see if this has an fs in the path
	fsMatchResult := fsMatch.FindStringSubmatch(path)

	switch {
	case fsMatchResult != nil && s.opt.Serve:
		// Serve /[fs]/remote files
		s.serveRemote(w, r, fsMatchResult[2], fsMatchResult[1])
		return
	case path == "metrics" && s.opt.EnableMetrics:
		promHandlerFunc(w, r)
		return
	case path == "*" && s.opt.Serve:
		// Serve /* as the remote listing
		s.serveRoot(w, r)
		return
	case s.files != nil:
		// Serve the files
		r.URL.Path = "/" + path
		s.files.ServeHTTP(w, r)
		return
	case path == "" && s.opt.Serve:
		// Serve the root as a remote listing
		s.serveRoot(w, r)
		return
	}
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// Wait blocks while the server is serving requests
func (s *Server) Wait() {
	s.server.Wait()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	return s.server.Shutdown()
}
