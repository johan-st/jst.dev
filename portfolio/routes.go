package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/matryer/way"
	"gitlab.com/golang-commonmark/markdown"
)

type handler struct {
	l      *log.Logger
	fs     fs.FS
	router *way.Router
}

// Register handlers for routes
func (h *handler) routes() {

	// STATIC ASSETS
	h.router.HandleFunc("GET", "/favicon.ico", h.handleFavicon())
	h.router.HandleFunc("GET", "/assets/", h.handleAssets())

	// basic pages
	h.router.HandleFunc("GET", "/:page", h.handlePage())

	// 405
	h.router.HandleFunc("*", "*", h.handleNotAllowed())

	// 404
	h.router.NotFound = h.handleNotFound()
}

// PAGE type and data
var (
	baseFiles = []string{
		"template/layout/base.gohtml",
		"template/layout/header.gohtml",
		"template/layout/nav.gohtml",
	}
	baseCSS = []string{"/assets/style.css"}
	baseJS  = []string{}
	// baseJS    = []string{"https://cdn.tailwindcss.com"}
	baseMeta = map[string]string{
		"description": "Portfolio and playground for Johan Strand",
		"author":      "Johan Strand",
		"keywords":    "portfolio, johan-st, Johan Strand, projects, blog, images, full stack, software, developer, web-dev, web developer, golang, go, javascript, react, reactjs, nextjs, nodejs, typescript, ts, html, css, sass, scss, tailwindcss, tailwind, postgres, sql, mongodb, nosql, docker, kubernetes, k8s, aws, amazon web services, cloud, cloud computing, serverless, lambda, api, rest, graphql, jamstack, server side rendering, ssr, static site generator, ssg, web development, webdev, web development, webdev, fullstack, full stack, full-stack, fullstack developer, full stack developer, full-stack developer, fullstack dev, full stack dev, full-stack dev, fullstack development, full stack development, full-stack development, fullstack web development, full stack web development, full-stack web development, fullstack webdev, full stack webdev, full-stack webdev, fullstack web dev, full stack web dev, full-stack web dev",
	}
)

type pageDataGetter func(*http.Request) (any, error)

type page struct {
	file       string
	linkText   string
	path       string
	tmplParsed *template.Template

	Title    string
	Meta     map[string]string
	CSS      []string
	JS       []string
	NavLinks map[string]string

	// empty string means no markdown
	markdownFile string
	// rendered markdown
	Markdown string

	PageData any
	pageDataGetter
}

type pageDataAdmin struct {
	Message string
	User    string
	Error   string
}

func getDataAdmin(req *http.Request) (any, error) {
	return pageDataAdmin{
		Message: "hello world",
		User:    "master Johan",
	}, nil
}

// HANDLERS

func (h *handler) handlePage() http.HandlerFunc {
	// pages
	pages := []page{
		{
			file:     "index.gohtml",
			linkText: "Home",
			path:     "",

			Title: "Home | jst.dev",
			Meta:  baseMeta,
			CSS:   baseCSS,
			JS:    baseJS,
		},
		{
			file:     "admin.gohtml",
			linkText: "Administration",
			path:     "admin",

			Title: "Admin | jst.dev",
			Meta:  baseMeta,
			CSS:   baseCSS,
			JS:    baseJS,

			pageDataGetter: getDataAdmin,
		},
		{
			file:         "about.gohtml",
			linkText:     "About",
			path:         "about",
			markdownFile: "about.md",

			Title: "About | jst.dev",
			Meta:  baseMeta,
			CSS:   baseCSS,
			JS:    baseJS,
		},
	}

	// setup
	l := h.l.With("handler", "handlePage")

	defer func(t time.Time) {
		l.Info("templates parsed and ready to be served", "time", time.Since(t))
	}(time.Now())

	tmplBase, err := template.ParseFS(h.fs, baseFiles...)
	if err != nil {
		l.Fatal("Could not parse base template", "error", err)
	}
	err = buildTemplates(h.fs, tmplBase, &pages)
	if err != nil {
		l.Fatal("Could not build templates", "error", err)
	}
	err = renderMarkdown(h.fs, &pages)
	if err != nil {
		l.Fatal("Could not render markdown", "error", err)
	}
	// make and add nav-links
	links := make(map[string]string)
	for _, p := range pages {
		links[p.path] = p.linkText
	}
	for i := range pages {
		pages[i].NavLinks = links
	}

	err = testExecution(pages)
	if err != nil {
		l.Fatal("Trial execution error:", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		requestedPage := way.Param(r.Context(), "page")

		// DEBUG

		for _, p := range pages {
			if requestedPage == p.path {
				l.Debug(
					"serving page",
					"page", requestedPage,
					"path", r.URL.Path,
				)

				if p.pageDataGetter != nil {
					var pData any
					pData, err := p.pageDataGetter(r)
					if err != nil {
						l.Warn("Failed to get page data")
						pData = nil
					}
					p.PageData = pData

				}
				err = p.tmplParsed.Execute(w, p)
				if err != nil {
					l.Error("Could not execute template", "error", err)
					h.respondError(w, r, "internal server error", http.StatusInternalServerError)
				}
				return
			}
		}

		l.Debug(
			"serving page 404 not found",
			"page", requestedPage,
			"path", r.URL.Path,
		)
		h.handleNotFound()(w, r)

	}
}
func (h *handler) handleAssets() http.HandlerFunc {
	// setup
	l := h.l.With("handler", "handleAssets")
	return func(w http.ResponseWriter, r *http.Request) {
		reqFile := strings.TrimPrefix(r.URL.Path, "/assets/")
		if reqFile == "" {
			h.respondError(w, r, "not found", http.StatusNotFound)
			return
		}

		file, err := h.fs.Open("assets/" + reqFile)
		if err != nil {
			l.Debug("could not open asset", "file", reqFile, "error", err)
			h.respondError(w, r, "not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		bytes, err := io.ReadAll(file)
		if err != nil {
			l.Error("could not read asset", "file", reqFile, "error", err)
			h.respondError(w, r, "internal server error", http.StatusInternalServerError)
			return
		}

		mimeType := mime.TypeByExtension(path.Ext(reqFile))
		l.Debug("serving asset", "file", reqFile, "Content-Type", mimeType)
		w.Header().Add("Content-Type", mimeType)

		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	}

}

// handleFavicon serves the favicon.ico.
func (h *handler) handleFavicon() http.HandlerFunc {
	// setup
	l := h.l.With("handler", "handleFavicon")
	file, err := h.fs.Open("assets/favicon.ico")
	if err != nil {
		l.Fatal("could not open favicon", "error", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		l.Fatal("could not read favicon", "error", err)
	}

	mimeType := mime.TypeByExtension(path.Ext("favicon.ico"))

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		l.Debug("serving favicon", "Content-Type", mimeType)

		w.Header().Add("Content-Type", mimeType)
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	}
}

func (h *handler) handleNotFound() http.HandlerFunc {
	// setup
	page404 := page{
		file:     "404.gohtml",
		linkText: "404",
		path:     "/404",

		Title: "404 | jst.dev",
		Meta:  baseMeta,

		CSS: baseCSS,
		JS:  baseJS,

		PageData: nil,
	}

	l := h.l.With("handler", "handleNotFound")
	tmpl, err := template.ParseFS(h.fs, append(baseFiles, "template/page/"+page404.file)...)
	if err != nil {
		l.Fatal("Could not parse 404 template", "error", err)
	}

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		l.Debug("handling request", "path", r.URL.Path)

		w.WriteHeader(http.StatusNotFound)
		err = tmpl.Execute(w, page404)
		if err != nil {
			l.Error("Could not execute template", "error", err)
			h.respondError(w, r, "internal server error", http.StatusInternalServerError)
		}
	}
}
func (h *handler) handleNotAllowed() http.HandlerFunc {
	// setup
	l := h.l.With("handler", "handleNotAllowed")
	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		l.Info("Method not allowed", "path", r.URL.Path)
		h.respondCode(w, r, http.StatusMethodNotAllowed)
	}
}

// RESPONDERS

func (h *handler) respondJson(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func (h *handler) respondCode(w http.ResponseWriter, r *http.Request, code int) {
	w.WriteHeader(code)
}

// respondError sends out a respons containing an error. This helper function is meant to be generic enough to serve most needs to communicate errors to the users
func (h *handler) respondError(w http.ResponseWriter, r *http.Request, msg string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "<html><h1>%d</h1><pre>%s</pre></html>", statusCode, msg)
}

// OTHER ESSENTIALS

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// t := time.Now()

	h.router.ServeHTTP(w, r)

	// h.l.Print(t.UTC().Local(),
	// 	"method", r.Method,
	// 	"url", r.Host+r.URL.String(),
	// 	"remote", r.RemoteAddr,
	// 	"user-agent", r.UserAgent(),
	// 	"time elapsed", time.Since(t))
}

// Template helpers

func buildTemplates(fs fs.FS, base *template.Template, pages *[]page) error {

	for i, p := range *pages {
		baseClone, err := base.Clone()
		if err != nil {
			return err
		}
		tmpl, err := baseClone.ParseFS(fs, "template/page/"+p.file)
		if err != nil {
			return err
		}
		(*pages)[i].tmplParsed = tmpl
	}
	return nil
}

func renderMarkdown(fs fs.FS, pages *[]page) error {
	md := markdown.New(markdown.XHTMLOutput(true))

	for i, p := range *pages {
		if p.markdownFile == "" {
			continue
		}
		file, err := fs.Open("template/page/" + p.markdownFile)
		if err != nil {
			return err
		}
		defer file.Close()

		bytes, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		(*pages)[i].Markdown = md.RenderToString(bytes)
	}

	return nil
}

func testExecution(pages []page) error {
	for _, p := range pages {
		err := p.tmplParsed.Execute(io.Discard, p)
		if err != nil {
			return err
		}
	}
	return nil
}
