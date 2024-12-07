package web

import (
	"log"
	"net/http"

	"jst.dev/internal/repo"
)

func HandlerBlogWeb(repo *repo.TursoRepo) func(w http.ResponseWriter, r *http.Request) {
	base := DataBase{
		Title: "Blog - jst.dev",
		Meta:  Meta{"description": "Blog"},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := repo.BlogPostsFeatured(4, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatalf("Error getting featured posts: %e", err)
		}
		if len(posts) < 2 {
			http.Error(w, "Not enough posts", http.StatusInternalServerError)
			log.Fatalf("Not enough posts: %d", len(posts))
			// TODO: handle this by showing the first post, expanded.
		}
		component := Blog(base, posts[0], posts[1:])
		err = component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Fatalf("Error rendering in BlogWebHandler: %e", err)
		}
	}
}
