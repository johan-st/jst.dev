package server

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"jst.dev/src/repo"
	"jst.dev/src/web"
)

func (s *Server) handlerShortUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortCode := r.PathValue("shortCode")

		urlRedirect, err := s.RepoTurso.GetShortUrlRedirect(shortCode)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Redirect(w, r, "https://jst.dev/url/", http.StatusSeeOther)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		http.Redirect(w, r, urlRedirect, http.StatusTemporaryRedirect)
	}
}

func (s *Server) handlerShortUrlNew(navItems navItems) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortUrls, err := s.RepoTurso.GetShortUrls()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		web.Layout(
			web.PageContext{
				Title:  "New Short URL",
				TopNav: navItems,
			},
			web.ShortUrl(shortUrls),
		).Render(r.Context(), w)
	}
}

func (s *Server) handlerShortUrlPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			parsedUrl *url.URL
			shortUrl  repo.UrlShort
			err       error
			location  string
		)

		parsedUrl, err = validateAndParseUrl(r.FormValue("url"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		shortUrl, err = s.RepoTurso.UrlShortenerInsert(parsedUrl)
		if r.Host == "url.jst.dev" {
			location = "https://url.jst.dev/i/" + shortUrl.ShortCode
		} else {
			location = "/url/i/" + shortUrl.ShortCode
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, location, http.StatusSeeOther)
	}
}

func (s *Server) handlerShortUrlInfo(navItems navItems) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortCode := r.PathValue("shortCode")

		shortUrl, err := s.RepoTurso.GetShortUrl(shortCode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		web.Layout(
			web.PageContext{
				Title:  "Short URL Info",
				TopNav: navItems,
			},
			web.ShortUrlInfo(&shortUrl),
		).Render(r.Context(), w)
	}
}

func (s *Server) handlerShortUrlRedirectToInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortCode := r.PathValue("shortCode")
		http.Redirect(w, r, "/url/i/"+shortCode, http.StatusMovedPermanently)
	}
}

func validateAndParseUrl(u string) (*url.URL, error) {
	parsedUrl, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	// If there's no scheme but there is a path, the host might be in the path
	if parsedUrl.Scheme == "" && parsedUrl.Host == "" && parsedUrl.Path != "" {
		// Split on first "/" to separate host from path
		parts := strings.SplitN(parsedUrl.Path, "/", 2)
		parsedUrl.Host = parts[0]
		if len(parts) > 1 {
			parsedUrl.Path = "/" + parts[1]
		} else {
			parsedUrl.Path = ""
		}
	}

	// Set default scheme if missing
	if parsedUrl.Scheme == "" {
		parsedUrl.Scheme = "https"
	}

	if parsedUrl.Host == "" {
		return nil, errors.New("host is required")
	}
	return parsedUrl, nil
}
