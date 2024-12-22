package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"fmt"
	"time"

	"github.com/coder/websocket"
	"jst.dev/src/repo"
	"jst.dev/src/web"
)

func (s *Server) RegisterRoutes() http.Handler {
	navItems := navItems{
		{Text: "Home", Href: "/", Active: false},
		{Text: "Url", Href: "/url", Active: false},
	}
	mux := http.NewServeMux()

	// REGISTER ROUTES

	// URL Shortener
	// - handle short url redirects
	mux.HandleFunc("GET u.jst.dev/{shortCode}", s.handlerShortUrl())
	mux.HandleFunc("GET url.jst.dev/{shortCode}", s.handlerShortUrl())
	// - handle short url info
	mux.HandleFunc("GET /url/i/{shortCode}", s.handlerShortUrlInfo(navItems.at("/url")))
	// - handle short url new
	mux.HandleFunc("GET /url/", s.handlerShortUrlNew(navItems.at("/url")))
	// - handle short url internal redirect
	mux.HandleFunc("GET /url/{shortCode}", s.handlerRedirectWithUrlParam("/url/i/{shortCode}", []string{"shortCode"}, http.StatusMovedPermanently))

	// - handle root
	mux.HandleFunc("GET /", s.handlerRoot(navItems.at("/"), navItems.at("/404")))
	// mux.HandleFunc("GET /blog", s.handlerBlogIndex)
	// mux.HandleFunc("GET /blog/{slug}", s.handlerBlogPost)
	// mux.HandleFunc("GET /about", s.handlerAbout)
	// mux.HandleFunc("GET /health", s.healthHandler)
	// mux.HandleFunc("GET /websocket", s.websocketHandler)

	// API - Blog
	// mux.HandleFunc("GET /api/blog-post", s.handlerNotImplemented("list all posts"))
	// mux.HandleFunc("POST /api/blog-post", s.handlerNotImplemented("create a new post"))
	// mux.HandleFunc("GET /api/blog-post/{slug}", s.handlerNotImplemented("get post by slug"))
	// mux.HandleFunc("DELETE /api/blog-post/{slug}", s.handlerNotImplemented("delete post by slug"))

	// // API - Url Shortener
	mux.HandleFunc("POST /api/short-url", s.handlerShortUrlPost())

	// Serve static files
	fileServer := http.FileServer(http.FS(web.Files))
	mux.Handle("GET /assets/", fileServer)

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// s.debugRequest("corsMiddleware", r)

		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "fly.dev, jst.dev")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}
// HANDLERS
func (s *Server) handlerRoot(navItemsIndex, navItemsNotFound navItems) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// s.debugRequest("handlerRoot", r)

		if r.URL.Path != "/" {
			web.Layout(web.PageContext{
				Title: "404 Not Found",
				Meta: web.Meta{
					"description": "404 Not Found",
					"canonical":   r.URL.Path,
					"robots":      "noindex, nofollow",
				},
				TopNav: navItemsNotFound,
			}, web.NotFound()).Render(r.Context(), w)
			return
		}
		pageContext := web.PageContext{
			Title: "Home",
			Meta: web.Meta{
				"description": "Home",
			},
			TopNav: navItemsIndex,
			Scripts: []web.ScriptTag{
				{Src: "/assets/js/page/index.js", Async: true, Defer: true},
			},
		}
		web.Layout(pageContext, web.Index()).Render(r.Context(), w)
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	var status string

	err := s.RepoTurso.Health()
	if err != nil {
		status = "DOWN"
	} else {
		status = "OK"
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(status)); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	socket, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to open websocket", http.StatusInternalServerError)
		return
	}
	defer socket.Close(websocket.StatusGoingAway, "Server closing websocket")

	ctx := r.Context()
	socketCtx := socket.CloseRead(ctx)

	for {
		payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
		if err := socket.Write(socketCtx, websocket.MessageText, []byte(payload)); err != nil {
			log.Printf("Failed to write to socket: %v", err)
			break
		}
		time.Sleep(2 * time.Second)
	}
}

func (s *Server) handlerBlogIndex(w http.ResponseWriter, r *http.Request) {

	featuredPosts, err := s.RepoTurso.BlogPostsFeatured(5, 0)
	if err != nil {
		http.Error(w, "Failed to get featured posts", http.StatusInternalServerError)
		return
	}
	pageContext := web.PageContext{
		Title: "Blog",
		Meta: web.Meta{
			"description": "Blog",
		},
	}

	web.Layout(pageContext, web.Blog(featuredPosts)).Render(r.Context(), w)
}

func (s *Server) handlerBlogPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	post, err := s.RepoTurso.BlogPostBySlug(slug)
	if err != nil {
		http.Error(w, "Failed to get post", http.StatusInternalServerError)
		return
	}

	web.Layout(web.PageContext{
		Title: post.Title,
		Meta: web.Meta{
			"canonical": r.URL.Path,
		},
	}, web.BlogPost(post)).Render(r.Context(), w)
}

func (s *Server) handlerBlogPostEdit(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	post, err := s.RepoTurso.BlogPostBySlug(slug)
	if err != nil {
		http.Error(w, "Failed to get post", http.StatusInternalServerError)
		return
	}

	web.Layout(web.PageContext{
		Title: post.Title + " - Edit",
		Meta: web.Meta{
			"canonical": r.URL.Path,
		},
	}, web.BlogPostEdit(post)).Render(r.Context(), w)
}

func (s *Server) handlerBlogPostNew(w http.ResponseWriter, r *http.Request) {
	web.Layout(web.PageContext{
		Title: "New Post",
		Meta: web.Meta{
			"canonical": r.URL.Path,
		},
	}, web.BlogPostNew()).Render(r.Context(), w)
}

func (s *Server) handlerApiBlogPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	post, err := s.RepoTurso.BlogPostBySlug(slug)
	if err != nil {
		http.Error(w, "Failed to get post", http.StatusInternalServerError)
		return
	}
	post.Title = r.FormValue("title")
	post.Slug = r.FormValue("slug")
	post.Body = r.FormValue("body")

	http.Redirect(w, r, "/blog/"+post.Slug, http.StatusSeeOther)
}

func (s *Server) handlerApiBlogPostNew(w http.ResponseWriter, r *http.Request) {
	post := repo.BlogPost{
		Title: r.FormValue("title"),
		Slug:  r.FormValue("slug"),
		Body:  r.FormValue("body"),
	}
	// s.RepoTurso.CreateBlogPost(post)
	http.Redirect(w, r, "/blog/"+post.Slug+"/edit", http.StatusSeeOther)
}

func (s *Server) handlerAbout(w http.ResponseWriter, r *http.Request) {

	web.Layout(web.PageContext{
		Title: "About",
		Meta: web.Meta{
			"description": "About",
		},
	}, web.About()).Render(r.Context(), w)
}

func (s *Server) handlerNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info := web.DebugInfo{
			Pattern: r.URL.Path,
			Method:  r.Method,
			Host:    r.Host,
			Path:    r.URL.Path,
		}
		web.Layout(web.PageContext{
			Title: "Not Found",
			Meta: web.Meta{
				"canonical": r.URL.Path,
			},
		},
			web.NotFoundWithDebug(info),
		).Render(r.Context(), w)
	}
}

func (s *Server) handlerNotImplemented(message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not implemented: "+message, http.StatusNotImplemented)
	}
}

// handlerRedirectWithUrlParam redirects a request to a given url using the given param keys to replace the placeholders in the url.
// If the param key is not found in the url, it will panic.
// example:
// handlerRedirectWithUrlParam("/url/i/{shortCode}", []string{"shortCode"}, http.StatusMovedPermanently)
func (s *Server) handlerRedirectWithUrlParam(url string, paramKeys []string, status int) http.HandlerFunc {
	for _, paramKey := range paramKeys {
		if !strings.Contains(url, "{"+paramKey+"}") {
			// Yes! Actually panicing as there is no way to run this handler if the param key is not found in the target url
			panic("param key not found in url: " + paramKey)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		for _, paramKey := range paramKeys {
			url = strings.Replace(url, "{"+paramKey+"}", r.PathValue(paramKey), len(paramKeys))
		}
		http.Redirect(w, r, url, status)
	}
}

// RESPONSE WRITERS

func writeJSON(w http.ResponseWriter, v any) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

// HELPERS

type navItems []web.NavItem

// At returns a new navItems slice with the item at the given href set to active
func (ns navItems) at(href string) navItems {
	new := make(navItems, len(ns))
	for i, it := range ns {
		new[i] = it
		if it.Href == href {
			new[i].Active = true
		}
	}
	return new
}

func (s *Server) debugRequest(tag string, r *http.Request) {
	log.Printf("%s: %s %s %s, matchedPattern: %s", tag, r.Method, r.URL.Path, r.Host, r.Pattern)
}
