package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/a-h/templ"
	log "github.com/charmbracelet/log"
	"github.com/johan-st/jst.dev/pages"
	"github.com/matryer/way"
)

//go:embed content
var embededFileSystem embed.FS
var (
	globalTitle = "dpj-docs"
	darkTheme   = pages.Theme{
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
		{Active: false, Url: "/", Text: "Docs", External: false},
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
	translations  []pages.Translation
	availableDocs []pages.Post
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
	srv.router.HandleFunc("GET", "/docs/", srv.handleMarkdown("content/docs"))
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

// handleMarkdown serves a pages from templates.
// NOTE: dirRoot can not end with a slash.
// NOTE: dirRoot is relative to the embeded filesystem.
func (srv *server) handleMarkdown(dirRoot string) http.HandlerFunc {
	// timing and logging
	l := srv.l.With("handler", "markdown")
	defer func(t time.Time) {
		l.Info(
			"ready",
			"time", time.Since(t),
		)
	}(time.Now())

	// setup
	// get base css styles
	styles, err := fs.ReadFile(embededFileSystem, "content/assets/inline.css")
	if err != nil {
		l.Fatal("Could not read main.css", "error", err)
	}

	baseStyles, err := pages.StyleTag(darkTheme, string(styles))
	if err != nil {
		l.Fatal("Could not create style tag", "error", err)
	}

	data := pages.MainData{
		DocTitle:      globalTitle,
		TopNav:        navLinks,
		FooterLinks:   footerLinks,
		Metadata:      metadata,
		ThemeStyleTag: baseStyles,
	}

	localFS, err := fs.Sub(embededFileSystem, dirRoot)
	if err != nil {
		l.Fatal("load filesystem", "error", err)
	}
	srv.availableDocs, err = getAvailablePosts(localFS)
	if err != nil {
		l.Fatal("Could not get available posts", "error", err)
	}

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		// time and log

		defer func(t time.Time) {
			l.Debug("serving page",
				"time", time.Since(t),
				"path", r.URL.Path,
			)
		}(time.Now())

		var content templ.Component

		for _, p := range srv.availableDocs {
			l.Debug("checking path", "path", p.Path, "url", r.URL.Path)
			if p.Path == r.URL.Path {
				content = pages.MarkdownPost(p)
				break
			}
		}
		if content == nil {
			srv.respCode(http.StatusNotFound, w, r)
			return
		}

		layout := pages.Layout(data, content)
		err = layout.Render(r.Context(), w)
		if err != nil {
			l.Error("Could not render template", "error", err)
			srv.respCode(http.StatusInternalServerError, w, r)
		}
	}
}

// handlePage serves a pages from templates.
func (srv *server) handlePage() http.HandlerFunc {
	// timing and logging
	l := srv.l.With("handler", "pages")
	defer func(t time.Time) {
		l.Info(
			"ready",
			"time", time.Since(t),
		)
	}(time.Now())

	// setup
	// get base css styles
	styles, err := fs.ReadFile(embededFileSystem, "content/assets/inline.css")
	if err != nil {
		l.Fatal("Could not read main.css", "error", err)
	}

	baseStyles, err := pages.StyleTag(darkTheme, string(styles))
	if err != nil {
		l.Fatal("Could not create style tag", "error", err)
	}

	data := pages.MainData{
		DocTitle:      globalTitle,
		TopNav:        navLinks,
		FooterLinks:   footerLinks,
		Metadata:      metadata,
		ThemeStyleTag: baseStyles,
	}

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

		var content templ.Component

		switch r.URL.Path {
		case "/":
			content = pages.Landing(&srv.availableDocs)
		case "/ai/translate":
			content = pages.OpenAI(srv.translations)
		default:
			file, err := fs.ReadFile(embededFileSystem, "content/docs/todo.md")
			if err != nil {
				l.Error("Could not read 'content/docs/todo.md'", "error", err)
			}
			post, err := pages.FileToPost(file, "")
			if err != nil {
				l.Error("Could not convert file to post", "error", err)
			}
			content = pages.MarkdownPost(post)
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
	file, err := fs.ReadFile(embededFileSystem, path)
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

func getAvailablePosts(filesystem fs.FS) ([]pages.Post, error) {
	var (
		paths []string
		posts []pages.Post
	)

	fs.WalkDir(filesystem, ".", addToList(&paths))
	for _, p := range paths {
		file, err := fs.ReadFile(filesystem, p)
		if err != nil {
			return nil, err
		}

		post, err := pages.FileToPost(file, "docs")
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func addToList(list *[]string) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !d.IsDir() {
			*list = append(*list, path)
		}
		return nil
	}
}
