package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/matryer/way"
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
	// srv.router.HandleFunc("GET", "/robots.txt", srv.handleRobots())
	h.router.HandleFunc("GET", "/assets/", h.handleAssets())

	h.router.HandleFunc("GET", "/:page", h.handlePage())
	h.router.HandleFunc("*", "*", h.handleNotAllowed())
	// 404
	h.router.NotFound = h.handleNotFound()
}

// PAGE type and data
var (
	baseFiles = []string{"template/layout/base.gohtml", "template/layout/header.gohtml"}
	baseCSS   = []string{"/assets/style.css"}
	baseJS    = []string{}
	// baseJS    = []string{"https://cdn.tailwindcss.com"}
	baseMeta = map[string]string{
		"description": "Portfolio and playground for Johan Strand",
		"author":      "Johan Strand",
		"keywords":    "portfolio, johan-st, Johan Strand, projects, blog, images, full stack, software, developer, web-dev, web developer, golang, go, javascript, react, reactjs, nextjs, nodejs, typescript, ts, html, css, sass, scss, tailwindcss, tailwind, postgres, sql, mongodb, nosql, docker, kubernetes, k8s, aws, amazon web services, cloud, cloud computing, serverless, lambda, api, rest, graphql, jamstack, server side rendering, ssr, static site generator, ssg, web development, webdev, web development, webdev, fullstack, full stack, full-stack, fullstack developer, full stack developer, full-stack developer, fullstack dev, full stack dev, full-stack dev, fullstack development, full stack development, full-stack development, fullstack web development, full stack web development, full-stack web development, fullstack webdev, full stack webdev, full-stack webdev, fullstack web dev, full stack web dev, full-stack web dev",
	}
)

type page struct {
	file     string
	linkText string
	path     string

	Title string
	Meta  map[string]string
	CSS   []string
	JS    []string

	PageData any
}

type adminData struct {
	Message string
	User    string
	Error   string
}

// HANDLERS
func (h *handler) handlePage() http.HandlerFunc {
	// pages

	var (
		pageIndex = page{
			file:     "index.gohtml",
			linkText: "Home",
			path:     "/",

			Title: "Home | jst.dev",
			Meta:  baseMeta,
			CSS:   baseCSS,
			JS:    baseJS,

			PageData: nil,
		}

		pageAdmin = page{
			file:     "admin.gohtml",
			linkText: "Admin",
			path:     "/admin",

			Title: "Admin | jst.dev",
			Meta:  baseMeta,
			CSS:   baseCSS,
			JS:    baseJS,

			PageData: adminData{},
		}
	)

	// setup
	l := h.l.With("handler", "handlePage")
	defer func(t time.Time) {
		l.Info("teplates parsed and ready to be served", "time", time.Since(t))
	}(time.Now())

	tmplBase, err := template.ParseFS(h.fs, baseFiles...)
	if err != nil {
		l.Fatal("Could not parse base template", "error", err)
	}

	// INDEX PAGE
	tmplIndex, err := tmplBase.Clone()
	if err != nil {
		l.Fatal("Could not clone base template", "error", err)
	}
	tmplIndex, err = tmplIndex.ParseFS(h.fs, "template/page/"+pageIndex.file)
	if err != nil {
		l.Fatal("Could not parse index template", "error", err)
	}

	// ADMIN PAGE
	tmplAdmin, err := tmplBase.Clone()
	if err != nil {
		l.Fatal("Could not clone base template", "error", err)
	}
	tmplAdmin, err = tmplAdmin.ParseFS(h.fs, "template/page/"+pageAdmin.file)
	if err != nil {
		l.Fatal("Could not parse admin template", "error", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		l.Debug("handling request", "path", r.URL.Path)
		requestedPage := way.Param(r.Context(), "page")
		switch requestedPage {
		case "":
			l.Debug("serving page", "page", "index")
			err = tmplIndex.Execute(w, pageIndex)
			if err != nil {
				l.Error("Could not execute template", "error", err)
				h.respondError(w, r, "internal server error", http.StatusInternalServerError)
			}
		case "admin":
			l.Debug("serving page", "page", "admin")

			pageAdmin.PageData = adminData{
				Message: "Hello, admin!",
				Error:   "This is an error message",
				User:    "Johan",
			}

			err = tmplAdmin.Execute(w, pageAdmin)
			if err != nil {
				l.Error("Could not execute template", "error", err)
				h.respondError(w, r, "internal server error", http.StatusInternalServerError)
			}

		default:
			l.Debug("serving page", "page", "default")
			h.handleNotFound()(w, r)
		}

	}
}

func (h *handler) handleAssets() http.HandlerFunc {
	// setup
	l := h.l.With("handler", "handleAssets")
	return func(w http.ResponseWriter, r *http.Request) {
		l.Debug("handling request", "path", r.URL.Path)
		file := strings.TrimPrefix(r.URL.Path, "/assets/")
		if file == "" {
			h.respondError(w, r, "not found", http.StatusNotFound)
			return
		}

		f, err := staticFS.ReadFile("assets/" + file)
		if err != nil {
			l.Debug("could not find asset", "file", file, "error", err)
			h.respondError(w, r, "not found", http.StatusNotFound)
			return
		}

		mimeType := mime.TypeByExtension(path.Ext(file))
		l.Debug("serving asset", "file", file, "Content-Type", mimeType)
		w.Header().Add("Content-Type", mimeType)

		w.WriteHeader(http.StatusOK)
		w.Write(f)
	}

}

// handleFavicon serves the favicon.ico.
func (h *handler) handleFavicon() http.HandlerFunc {
	// setup
	l := h.l.With("handler", "handleFavicon")

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		l.Debug("handling request", "path", r.URL.Path)
		http.ServeFile(w, r, "assets/favicon.ico")
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

		l.Printf("%#v", page404)

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
