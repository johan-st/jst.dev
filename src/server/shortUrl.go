package server

import (
	"net/http"
	"time"

	"jst.dev/src/web"
)

func (s *Server) handlerShortUrl(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")

	urlRedirect, err := s.RepoTurso.GetShortUrlRedirect(shortCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, urlRedirect, http.StatusTemporaryRedirect)
}

func (s *Server) handlerShortUrlNew(w http.ResponseWriter, r *http.Request) {
	shortUrls, err := s.RepoTurso.GetShortUrls()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	web.LayoutMinimal(
		web.PageContext{
			Title: "New Short URL",
		},
		web.NewShortUrl(shortUrls),
	).Render(r.Context(), w)
}

func (s *Server) handlerShortUrlPost(w http.ResponseWriter, r *http.Request) {
	shortUrl, err := s.RepoTurso.UrlShortenerInsert(r.FormValue("url"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/i/"+shortUrl.ShortCode, http.StatusSeeOther)
}

func (s *Server) handlerShortUrlInfo(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")

	shortUrl, err := s.RepoTurso.GetShortUrl(shortCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	unixSeconds := int64(shortUrl.Id.Time()) / 1000
	unixNano := (int64(shortUrl.Id.Time()) % 1000) * 10
	web.LayoutMinimal(
		web.PageContext{
			Title: "Short URL Info",
		},
		web.ShortUrlInfo(shortUrl, time.Unix(unixSeconds, unixNano).Format(time.RFC3339)),
	).Render(r.Context(), w)
}
