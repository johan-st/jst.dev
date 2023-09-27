package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/a-h/templ"
	log "github.com/charmbracelet/log"
	"github.com/johan-st/jst.dev/pages"
	"github.com/matryer/way"
)

//go:embed content
var embededFileSystem embed.FS
var (
	darkTheme = pages.Theme{
		ColorPrimary:    "#f90",
		ColorSecondary:  "#fa3",
		ColorBackground: "#333",
		ColorText:       "#aaa",
		ColorBorder:     "#666",
		BorderRadius:    "1rem",
	}

	metadata = map[string]string{
		"Description": "desc",
		"Keywords":    "e-comm docs",
		"Author":      "dpj",
	}

	// nav and footer links
	navLinks = []pages.Link{
		{Active: false, Url: "/ai/translate", Text: "Translation", External: false},
		{Active: false, Url: "/todos", Text: "ToDo", External: false},
	}
	footerLinks = []pages.Link{
		{Active: false, Url: "https://www.dpj.se", Text: "Live site", External: true},
		{Active: false, Url: "https://dpj.local", Text: "local site", External: true},
	}
)

type server struct {
	l      *log.Logger
	router *way.Router

	// persistence
	translations []pages.Translation
}

func newRouter(l *log.Logger) server {
	return server{
		l:      l,
		router: way.NewRouter(),
	}
}

// Register handlers for routes
func (srv *server) prepareRoutes() {

	// STATIC ASSETS
	srv.router.HandleFunc("GET", "/favicon.ico", srv.handleStaticFile("content/static/favicon.ico"))
	srv.router.HandleFunc("GET", "/static", srv.handleNotFound())
	srv.router.HandleFunc("GET", "/static/", srv.handleStaticDir("/static/", "content/static"))

	// GET
	srv.router.HandleFunc("GET", "...", srv.handlePage())
	
	// POST
	srv.router.HandleFunc("POST", "/ai", srv.handleAiTranslationPost())
	srv.router.HandleFunc("POST", "/ai/translate", srv.handleAiTranslationPost())
	srv.router.HandleFunc("POST", "/ai/stories", srv.handleAiStories())
	// h.router.HandleFunc("GET", "/admin/:page", h.handleAdminTempl())
	// h.router.HandleFunc("GET", "/admin/images/:id", h.handleAdminImage())

	// 404
	srv.router.NotFound = srv.handleNotFound()
}

// HANDLERS

func (srv *server) handleAiStories() http.HandlerFunc {
	// timing and logging
	l := srv.l.With("handler", "ApiAiStories")
	defer func(t time.Time) {
		l.Debug(
			"ready and waiting...",
			"time", time.Since(t),
		)
	}(time.Now())

	// setup

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		srv.respCode(http.StatusNotImplemented, w, r)
	}
}

// handlePage serves a pages from templates.
func (srv *server) handlePage() http.HandlerFunc {
	// timing and logging
	l := srv.l.With("handler", "root-pages")
	defer func(t time.Time) {
		l.Info(
			"ready",
			"time", time.Since(t),
		)
	}(time.Now())

	// setup
	// get base css styles
	styles, err := os.ReadFile("pages/assets/inline.css")
	if err != nil {
		l.Fatal("Could not read main.css", "error", err)
	}

	baseStyles, err := pages.StyleTag(darkTheme, string(styles))
	if err != nil {
		l.Fatal("Could not create style tag", "error", err)
	}

	data := pages.MainData{
		DocTitle:      "local_test",
		TopNav:        navLinks,
		FooterLinks:   footerLinks,
		Metadata:      metadata,
		ThemeStyleTag: baseStyles,
	}

	availablePosts := srv.getAvailablePosts("content/blog")

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		// time and log
		defer func(t time.Time) {
			l.Debug("serving page",
			"time", time.Since(t),
			"path", r.URL.Path,
		)
		}(time.Now())

		// uri
		reqUri , err := url.ParseRequestURI(r.URL.Path)
		normalizedPath := reqUri.EscapedPath()
		if err != nil {
			l.Error("Could not normalize path. using raw path", "error", err)
			normalizedPath = r.URL.Path
		}
		if normalizedPath != r.URL.Path {
			l.Warn("Normalized path", "from", r.URL.Path, "to", normalizedPath)
		}

		var content templ.Component

		switch normalizedPath {
		case "":
			content = pages.Landing(availablePosts)
		case "/ai/translate":
			content = pages.OpenAI(srv.translations)
		default:
			file, err := os.ReadFile("_docs/thoughts.md")
			if err != nil {
				l.Error("Could not read _docs/thoughts.md", "error", err)
			}

			content = pages.MarkdownFile(file)
		}

		layout := pages.Layout(data, content)

		err = layout.Render(r.Context(), w)
		if err != nil {
			l.Error("Could not render template", "error", err)
			srv.respCode(http.StatusInternalServerError, w, r)
		}
	}
}

func (srv *server) handleStaticDir(urlRoot, dirRoot string) http.HandlerFunc {
	// timing and logging
	l := srv.l.With("handler", "StaticDir")
	defer func(t time.Time) {
		l.Info(
			"ready",
			"urlRoot", urlRoot,
			"dirRoot", dirRoot,
			"time", time.Since(t),
		)
	}(time.Now())

	// setup
	subFs, err := fs.Sub(embededFileSystem, dirRoot)
	if err != nil {
		l.Fatal(
			"load filesystem",
			"err", err,
		)
	}
	fileSrv := http.FileServer(http.FS(subFs))

	// handler
	// return http.FileServer(http.FS(staticFS))
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(t time.Time) {
			l.Info(
				"serve file",
				"path", r.URL.Path,
				"time", time.Since(t),
			)
		}(time.Now())
		r.URL.Path = r.URL.Path[len(urlRoot)-1:]
		fileSrv.ServeHTTP(w, r)
	}
}

// handleStaticFile serves a predtermined static file.
func (srv *server) handleStaticFile(path string) http.HandlerFunc {
	// timing and logging
	l := srv.l.With("handler", "StaticFile")
	defer func(t time.Time) {
		l.Info(
			"ready",
			"time", time.Since(t),
		)
	}(time.Now())

	// setup
	file, err := os.ReadFile(path)
	if err != nil {
		l.Fatal(
			"load file",
			"path", path)
	}

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(file)

		if err != nil {
			srv.respCode(http.StatusInternalServerError, w, r)
			l.Error(
				"serve file",
				"path", path,
			)
		}
	}
}

// handleStaticFile serves a predtermined static file.
func (srv *server) handleNotFound() http.HandlerFunc {
	// timing and logging
	l := srv.l.With("handler", "NotFound")
	defer func(t time.Time) {
		l.Info(
			"ready",
			"time", time.Since(t),
		)
	}(time.Now())

	// setup

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		srv.respCode(http.StatusNotFound, w, r)
	}

}

// RESPONDERS

func (srv *server) respCode(code int, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf("%d - %s", code, http.StatusText(code))))
}

// http.ServeHttp

func (srv *server) Handler() http.Handler {
	return srv.router
}

// HELPERS

func (srv *server) getAvailablePosts(dir string) []pages.Post {
	// get all files in dir
	_, err := os.ReadDir(dir)
	if err != nil {
		srv.l.Error("read dir", "err", err)
	}

	return []pages.Post{
		{Title: "test for the one", Slug: "testslug"},
		{Title: "test 2 the moon", Slug: "testslug-2"},
	}
}
