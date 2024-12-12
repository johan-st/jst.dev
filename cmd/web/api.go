package web

import (
	"encoding/json"
	"net/http"

	"jst.dev/internal/repo"
)

var messagesRepo *repo.MessageRepo

func HandlerApiMessages() http.HandlerFunc {
	handlerGet := handlerGet()

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlerGet.ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handlerGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		messages := messagesRepo.MessagesListAll()
		json.NewEncoder(w).Encode(messages)
	}
}
