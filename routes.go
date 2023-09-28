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

func defaultPageData() (pages.PageData, error) {
	cssInline, err := fs.ReadFile(embededFileSystem, "content/assets/inline.css")
	if err != nil {
		return pages.PageData{}, err
	}

	themeComponent, err := defaultTheme().Component()
	if err != nil {
		return pages.PageData{}, err
	}

	return pages.PageData{
		DocTitle: "dpj-ai",
		TopNav: []pages.Link{
			{Url: "/ai/translate", Text: "Translation", External: false},
			{Url: "/todos", Text: "ToDo", External: false},
			{Url: "/", Text: "Docs", External: false},
		},
		FooterLinks: []pages.Link{
			{Url: "https://www.dpj.se", Text: "Live site", External: true},
			{Url: "https://dpj.local", Text: "local site", External: true},
		},
		Metadata:    map[string]string{"Description": "desc", "Keywords": "e-comm docs", "Author": "dpj"},
		StyleInline: pages.Style(string(cssInline)),
		StyleTheme:  themeComponent,
	}, nil
}

func defaultTheme() pages.Theme {
	return pages.Theme{
		ColorPrimary:    "#f90",
		ColorSecondary:  "#fa3",
		ColorBackground: "#333",
		ColorText:       "#aaa",
		ColorBorder:     "#666",
		BorderRadius:    "1rem",
	}
}

type server struct {
	l      *log.Logger
	router *way.Router

	// persistence
	translations  []pages.Translation
	availableDocs []pages.Post
}

func newRouter(l *log.Logger) server {
	l = l.WithPrefix(l.GetPrefix() + ".router")

	// setup
	return server{
		l:             l,
		router:        way.NewRouter(),
		translations:  []pages.Translation{},
		availableDocs: []pages.Post{},
	}
}

// Register handlers for routes
func (srv *server) prepareRoutes() {

	// STATIC ASSETS
	srv.router.HandleFunc("GET", "/favicon.ico", srv.handleStaticFile("content/static/favicon.ico"))
	srv.router.HandleFunc("GET", "/static", srv.handleNotFound())
	srv.router.HandleFunc("GET", "/static/", srv.handleStaticDir("content/static", "/static/"))

	// GET
	srv.router.HandleFunc("GET", "/docs", srv.handleRedirect(http.StatusTemporaryRedirect, "/"))
	srv.router.HandleFunc("GET", "/docs/", srv.handleMarkdown("content/docs", "docs"))
	srv.router.HandleFunc("GET", "...", srv.handlePage()) // catch all

	// POST
	srv.router.HandleFunc("POST", "/ai/translate", srv.handleApiTranslationPost())
	srv.router.HandleFunc("POST", "/ai/stories", srv.handleAiStories())

	// 404
	srv.router.NotFound = srv.handleNotFound()
}

// HANDLERS

func (srv *server) handleAiStories() http.HandlerFunc {
	// timing and logging
	l := srv.l.
		WithPrefix(srv.l.GetPrefix() + ".AiStories")

	defer func(t time.Time) {
		l.Debug(
			logReady,
			logTimeSpent, time.Since(t),
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
func (srv *server) handleMarkdown(rootDir, basePath string) http.HandlerFunc {
	// timing and logging
	l := srv.l.
		WithPrefix(srv.l.GetPrefix()+".Markdown").
		With(
			logRootDir, rootDir,
			logBaseURL, basePath,
		)

	defer func(t time.Time) {
		l.Debug(
			logReady,
			logTimeSpent, time.Since(t),
		)
	}(time.Now())

	// setup

	pageData, err := defaultPageData()
	if err != nil {
		l.Fatal("Could not get default page data", logError, err)
	}
	markdownFS, err := fs.Sub(embededFileSystem, rootDir)
	if err != nil {
		l.Fatal("load filesystem", logError, err)
	}
	srv.availableDocs, err = getAvailablePosts(markdownFS, basePath)
	if err != nil {
		l.Fatal("Could not get available posts", logError, err)
	}

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		// time and log

		defer func(t time.Time) {
			l.Debug("serving page",
				logTimeSpent, time.Since(t),
				logReqPath, r.URL.Path,
			)
		}(time.Now())

		var content templ.Component

		for _, p := range srv.availableDocs {
			l.Debug("checking path", logFilePath, p.Path, logReqPath, r.URL.Path)
			if p.Path == r.URL.Path {
				content = pages.MarkdownPost(p)
				break
			}
		}
		if content == nil {
			srv.respCode(http.StatusNotFound, w, r)
			return
		}

		layout := pages.Layout(pageData, content)
		err = layout.Render(r.Context(), w)
		if err != nil {
			l.Error("Could not render template", logError, err)
			srv.respCode(http.StatusInternalServerError, w, r)
		}
	}
}

// handlePage serves a pages from templates.
func (srv *server) handlePage() http.HandlerFunc {
	// timing and logging
	l := srv.l.
		WithPrefix(srv.l.GetPrefix() + ".Page")

	defer func(t time.Time) {
		l.Debug(
			logReady,
			logTimeSpent, time.Since(t),
		)
	}(time.Now())

	// setup
	pageData, err := defaultPageData()
	if err != nil {
		l.Fatal("Could not get default page data", logError, err)
	}

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		// time and log
		defer func(t time.Time) {
			l.Debug("serving page",
				logTimeSpent, time.Since(t),
				"path", r.URL.Path,
			)
		}(time.Now())

		// uri

		var content templ.Component

		switch r.URL.Path {
		case "/":
			content = pages.Blog(&srv.availableDocs)
		case "/ai/translate":
			content = pages.OpenAI(srv.translations)
		default:
			file, err := fs.ReadFile(embededFileSystem, "content/docs/todo.md")
			if err != nil {
				l.Error("Could not read 'content/docs/todo.md'", logError, err)
			}
			post, err := pages.FileToPost(file, "")
			if err != nil {
				l.Error("Could not convert file to post", logError, err)
			}
			content = pages.MarkdownPost(post)
		}

		layout := pages.Layout(pageData, content)

		err = layout.Render(r.Context(), w)
		if err != nil {
			l.Error("Could not render template", logError, err)
			srv.respCode(http.StatusInternalServerError, w, r)
		}
	}
}

func (srv *server) handleStaticDir(rootDir, basePath string) http.HandlerFunc {
	// timing and logging
	l := srv.l.
		WithPrefix(srv.l.GetPrefix()+".StaticDir").
		With(
			logRootDir, rootDir,
			logBaseURL, basePath,
		)

	defer func(t time.Time) {
		l.Debug(
			logReady,
			logTimeSpent, time.Since(t),
		)
	}(time.Now())

	// setup
	subFs, err := fs.Sub(embededFileSystem, rootDir)
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
				logReqPath, r.URL.Path,
				logTimeSpent, time.Since(t),
			)
		}(time.Now())
		r.URL.Path = r.URL.Path[len(basePath)-1:]
		fileSrv.ServeHTTP(w, r)
	}
}

// handleStaticFile serves a predtermined static file.
func (srv *server) handleStaticFile(path string) http.HandlerFunc {
	// timing and logging
	l := srv.l.
		WithPrefix(srv.l.GetPrefix()+".StaticFile").
		With("file", path)

	defer func(t time.Time) {
		l.Debug(
			logReady,
			logTimeSpent, time.Since(t),
		)
	}(time.Now())

	// setup
	file, err := fs.ReadFile(embededFileSystem, path)
	if err != nil {
		l.Fatal(
			"load file",
			logFilePath, path,
			logError, err,
		)
	}

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(file)

		if err != nil {
			srv.respCode(http.StatusInternalServerError, w, r)
			l.Error(
				"serve file",
				logError, err,
			)
		}
	}
}

func (srv *server) handleRedirect(code int, url string) http.HandlerFunc {
	// timing and logging
	l := srv.l.
		WithPrefix(srv.l.GetPrefix()+".Redirect").
		With(
			"code", code,
			"to", url,
		)

	defer func(t time.Time) {
		l.Debug(
			logReady,
			logTimeSpent, time.Since(t),
		)
	}(time.Now())

	// setup
	if code < 300 || code > 399 {
		l.Fatal("invalid redirect code", "code", code, "valid range", "300-399")
	}

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		l.Debug(
			"redirecting",
			"from", r.URL.Path,
		)
		http.Redirect(w, r, url, code)
	}
}

func (srv *server) handleNotFound() http.HandlerFunc {
	// timing and logging
	l := srv.l.
		WithPrefix(srv.l.GetPrefix() + ".NotFound")

	defer func(t time.Time) {
		l.Debug(
			logReady,
			logTimeSpent, time.Since(t),
		)
	}(time.Now())

	// setup

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		l.Debug("responding", "referer", r.Header.Values("referer"))
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

func getAvailablePosts(filesystem fs.FS, basePath string) ([]pages.Post, error) {
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

		post, err := pages.FileToPost(file, basePath)
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
