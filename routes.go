package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/Kwynto/gosession"
	"github.com/a-h/templ"
	log "github.com/charmbracelet/log"
	"github.com/johan-st/jst.dev/pages"
)

//go:embed content
var embededFileSystem embed.FS

type server struct {
	l              *log.Logger
	router         *http.ServeMux
	availablePosts []pages.BlogPost
	defaultData    pages.Data
}

// Register handlers for routes
func (srv *server) prepareRoutes() {

	// STATIC ASSETS
	srv.router.HandleFunc("GET /favicon.ico", srv.handleStaticFile("content/assets/favicon.ico"))
	srv.router.HandleFunc("GET /public/", srv.handlePublicDir("content/public", "/public/"))

	// AI
	// srv.router.HandleFunc("GET /ai", srv.handleRedirect(http.StatusTemporaryRedirect, "/ai/translate"))
	// srv.router.HandleFunc("GET /ai/audio", srv.handleNotImplemented())
	// srv.router.HandleFunc("GET /ai/chat", srv.handleNotImplemented())
	// srv.router.HandleFunc("GET /ai/content-filter", srv.handleNotImplemented())
	// srv.router.HandleFunc("GET /ai/stories", srv.handleNotImplemented())
	// srv.router.HandleFunc("GET /ai/tutor", srv.handleNotImplemented())
	// srv.router.HandleFunc("GET /ai/translate", srv.handleAiTranslation())
	// AI POST
	// srv.router.HandleFunc("POST /ai/translate", srv.handleAiTranslationPost())
	// srv.router.HandleFunc("POST /ai/stories", srv.handleAiStories())

	// DOCS
	srv.router.HandleFunc("GET /blog", srv.handleBlogIndex())
	srv.router.HandleFunc("GET /blog/", srv.handleBlog())

	// LANDING
	srv.router.HandleFunc("GET /", srv.handleRedirect(http.StatusTemporaryRedirect, "/404"))

	// PAGES
	srv.router.HandleFunc("GET /about", srv.handleMarkdownFile("content/pages/about.md"))
	srv.router.HandleFunc("GET /contact", srv.handleMarkdownFile("content/pages/contact.md"))

	// 404
	srv.router.HandleFunc("GET /404", srv.handleNotFound())
}

// HANDLERS

func (srv *server) handleMarkdownFile(path string) http.HandlerFunc {
	l := srv.l.
		WithPrefix(srv.l.GetPrefix() + ".MarkdownFile")

	defer func(t time.Time) {
		l.Info(
			logReady,
			logTimeSpent, time.Since(t),
			"file", path,
		)
	}(time.Now())

	// setup
	file, err := fs.ReadFile(embededFileSystem, path)
	if err != nil {
		l.Fatal(
			"load file",
			"error", err,
		)
	}

	mdPage, err := pages.MdToPage(file, "")
	if err != nil {
		l.Fatal(
			"convert markdown",
			"error", err,
		)
	}

	content := pages.PageContent(mdPage)

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		l.Debug(
			"responding",
			logReqPath, r.URL.Path,
		)
		err = pages.Layout(srv.defaultData, content).Render(r.Context(), w)
		if err != nil {
			l.Error("Could not render template", logError, err)
			srv.respCode(http.StatusInternalServerError, w, r)
		}
	}
}

// handleBlog serves a pages from templates.
// NOTE: dirRoot can not end with a slash.
// NOTE: dirRoot is relative to the embeded filesystem.
func (srv *server) handleBlog() http.HandlerFunc {
	l := srv.l.
		WithPrefix(srv.l.GetPrefix() + ".Docs")

	defer func(t time.Time) {
		l.Info(
			logReady,
			logTimeSpent, time.Since(t),
		)
	}(time.Now())
	// setup

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		// time and log

		defer func(t time.Time) {
			l.Debug("responding",
				logTimeSpent, time.Since(t),
				logReqPath, r.URL.Path,
			)
		}(time.Now())

		var content templ.Component

		for _, p := range srv.availablePosts {
			if p.Path == r.URL.Path {
				content = pages.BlogContent(p)
				l.Debug("document found",
					logFilePath, p.Path,
					logReqPath, r.URL.Path)

			}
		}
		if content == nil { // content not found
			l.Warn("no such document",
				logReqPath, r.URL.Path,
				"referer", r.Header.Values("referer"),
			)
			content = pages.Blog404(&srv.availablePosts)
		}

		err := pages.Layout(srv.defaultData, content).Render(r.Context(), w)
		if err != nil {
			l.Error("Could not render template", logError, err)
			srv.respCode(http.StatusInternalServerError, w, r)
		}
	}
}

// handlePage serves a pages from templates.
func (srv *server) handleAiTranslation() http.HandlerFunc {
	// timing and logging
	l := srv.l.
		WithPrefix(srv.l.GetPrefix() + ".PageAiTranslation")

	defer func(t time.Time) {
		l.Info(
			logReady,
			logTimeSpent, time.Since(t),
		)
	}(time.Now())

	// setup

	sessionHandler := newSessionHandler(l)
	pageData, err := defaultPageData()
	if err != nil {
		l.Fatal("Could not get default page data", logError, err)
	}

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			content templ.Component
			session gosession.SessionId
		)
		// time and log
		defer func(t time.Time) {
			l.Debug("serving page",
				logTimeSpent, time.Since(t),
				"path", r.URL.Path,
			)
		}(time.Now())

		session = gosession.Start(&w, r)
		trans := sessionHandler.getTranslations(&session)
		content = pages.OpenAiTranslate(trans)

		layout := pages.Layout(pageData, content)
		err = layout.Render(r.Context(), w)
		if err != nil {
			l.Error("Could not render template", logError, err)
			srv.respCode(http.StatusInternalServerError, w, r)
		}
	}
}

// handlePage serves a pages from templates.
func (srv *server) handleBlogIndex() http.HandlerFunc {
	// timing and logging
	l := srv.l.
		WithPrefix(srv.l.GetPrefix() + ".DocsIndex")

	defer func(t time.Time) {
		l.Info(
			logReady,
			logTimeSpent, time.Since(t),
		)
	}(time.Now())

	// setup

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		// time and log
		defer func(t time.Time) {
			l.Debug("serving page",
				logTimeSpent, time.Since(t),
				"path", r.URL.Path,
			)
		}(time.Now())

		content := pages.BlogIndex(&srv.availablePosts)

		layout := pages.Layout(srv.defaultData, content)
		err := layout.Render(r.Context(), w)
		if err != nil {
			l.Error("Could not render template", logError, err)
			srv.respCode(http.StatusInternalServerError, w, r)
		}
	}
}

func (srv *server) handlePublicDir(rootDir, basePath string) http.HandlerFunc {
	// timing and logging
	l := srv.l.
		WithPrefix(srv.l.GetPrefix()+".PublicDir").
		With(
			logRootDir, rootDir,
			logBaseURL, basePath,
		)

	defer func(t time.Time) {
		l.Info(
			logReady,
			logTimeSpent, time.Since(t),
		)
	}(time.Now())

	// setup
	PublicFs, err := fs.Sub(embededFileSystem, rootDir)
	if err != nil {
		l.Fatal(
			"load filesystem",
			"err", err,
		)
	}
	fileSrv := http.FileServer(http.FS(PublicFs))

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(t time.Time) {
			l.Debug(
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
		l.Info(
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
	return srv.subHandleRedirect(srv.l, code, url)
}

func (srv *server) handleNotImplemented() http.HandlerFunc {
	l := srv.l.WithPrefix(
		srv.l.GetPrefix() + ".NotImplemented",
	)
	l.Error("handler not implemented")
	return func(w http.ResponseWriter, r *http.Request) {
		l.Error("Not implemented",
			"referer", r.Header.Values("referer"),
			logReqPath, r.URL.Path,
		)
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Not implemented:\nmesssage: "))
	}
}

func (srv *server) handleNotFound() http.HandlerFunc {
	return srv.subHandleNotFound(srv.l)
}

// SUB HANDLERS
// sub handlers are used by other handlers

func (srv *server) subHandleRedirect(l *log.Logger, code int, url string) http.HandlerFunc {
	l = l.
		WithPrefix(srv.l.GetPrefix()+".Redirect").
		With(
			"code", code,
			"to", url,
		)
	defer func(t time.Time) {
		l.Info(
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

func (srv *server) subHandleNotFound(l *log.Logger) http.HandlerFunc {
	// timing and logging
	l = l.
		WithPrefix(srv.l.GetPrefix() + ".NotFound")

	defer func(t time.Time) {
		l.Info(
			logReady,
			logTimeSpent, time.Since(t),
		)
	}(time.Now())

	// setup
	pageData, err := defaultPageData()
	layout := pages.Layout(pageData, pages.NotFound())

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		l.Warn("responding", "referer", r.Header.Values("referer"), logReqPath, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		err = layout.Render(r.Context(), w)
		if err != nil {
			l.Error("Could not render template", logError, err)
			srv.respCode(http.StatusInternalServerError, w, r)
		}
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

// MIDDLEWARE

// func (srv *server) withSession(next http.HandlerFunc) http.HandlerFunc {
// 	// setup

// 	// middleware
// 	return func(w http.ResponseWriter, r *http.Request) {

// 	next(w, r)
// 	}
// }

// HELPERS

func FsToPosts(filesystem fs.FS, basePathForMatching string) ([]pages.BlogPost, error) {
	var (
		paths   []string
		mdPages []pages.BlogPost
	)

	fs.WalkDir(filesystem, ".", addToList(&paths))
	for _, p := range paths {
		file, err := fs.ReadFile(filesystem, p)
		if err != nil {
			return nil, err
		}

		mdPage, err := pages.MdToPage(file, basePathForMatching)
		if err != nil {
			return nil, err
		}
		mdPages = append(mdPages, mdPage)
	}
	return mdPages, nil
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

func newRouter(l *log.Logger) server {
	var posts []pages.BlogPost

	// setup
	tempFs, err := fs.Sub(embededFileSystem, "content/blog")
	if err != nil {
		l.Fatal("load filesystem", logError, err)
	}
	posts, err = FsToPosts(tempFs, "blog")
	if err != nil {
		l.Fatal("Could not get available posts", logError, err)
	}
	if len(posts) == 0 {
		l.Error("No blogposts found", logError, err)
	}
	if l.GetLevel() == log.DebugLevel {
		for _, p := range posts {
			l.Debug("Blogpost loaded",
				"title", p.Path,
				"path", p.Path,
				"bytes", len(p.Body),
			)
		}
	}

	pageData, err := defaultPageData()
	if err != nil {
		l.Fatal("Could not get default page data", logError, err)
	}

	return server{
		l:              l,
		router:         http.NewServeMux(),
		availablePosts: posts,
		defaultData:    pageData,
	}
}

// DEFAULTS

func defaultPageData() (pages.Data, error) {
	cssInline, err := fs.ReadFile(embededFileSystem, "content/assets/inline.css")
	if err != nil {
		return pages.Data{}, err
	}

	themeComponent, err := defaultTheme().Component()
	if err != nil {
		return pages.Data{}, err
	}

	return pages.Data{
		DocTitle: "dpj-ai",
		TopNav: []pages.Link{
			{Url: "/ai", Text: "AI", External: false},
			{Url: "/about", Text: "About", External: false},
			{Url: "/contact", Text: "Contact", External: false},
			{Url: "/blog", Text: "Blog", External: false},
		},
		FooterLinks: []pages.Link{
			{Url: "https://github.com/johan-st", Text: "GitHub", External: true},
			{Url: "mailto:jst@jst.dev", Text: "jst@jst.dev", External: true},
		},
		Metadata:    map[string]string{"Description": "desc", "Keywords": "e-comm docs", "Author": "Johan Strand"},
		StyleInline: pages.Style(string(cssInline)),
		StyleTheme:  themeComponent,
	}, nil
}

func defaultTheme() pages.Theme {
	return pages.Theme{
		ColorPrimary:       "#f90",
		ColorSecondary:     "#fc0",
		ColorBackground:    "#333",
		ColorBackgroundAlt: "#202020",
		ColorText:          "#aaa",
		ColorTextMuted:     "#666",
		ColorBorder:        "#666",
		BorderRadius:       "1rem",
	}
}
