package server

import (
	"net/http"

	"jst.dev/src/web"
)

func (s *Server) handlerShortUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortCode := r.PathValue("shortCode")

		urlRedirect, err := s.RepoTurso.GetShortUrlRedirect(shortCode)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
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
				Title: "New Short URL",
				TopNav: navItems,
			},
			web.ShortUrl(shortUrls, web.NewShortUrlResult{}),
		).Render(r.Context(), w)
	}
}

func (s *Server) handlerShortUrlPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortUrl, err := s.RepoTurso.UrlShortenerInsert(r.FormValue("url"))
		var location string
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
	for _, it := range navItems {
		if it.Href == "/url" {
			it.Active = true
		}
	}
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

