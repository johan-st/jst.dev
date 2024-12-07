package web

import (
	"log"
	"net/http"

	"fmt"

	"jst.dev/internal/repo"
)

var messagesRepo *repo.MessageRepo

func HandlerApiMessages() http.HandlerFunc {
	messagesRepo := repo.NewMessageRepo()
	handlerPost := handlerPost(messagesRepo)
	handlerGet := handlerGet(messagesRepo)

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlerPost.ServeHTTP(w, r)
		case http.MethodGet:
			handlerGet.ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	}
}

func handlerPost(messagesRepo *repo.MessageRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}

		msg := r.FormValue("message")
		if msg == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		id := messagesRepo.MessageAdd(msg)
		component := Msg(id.String(), msg)
		err = component.Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Fatalf("Error rendering in HelloWebHandler: %e", err)
		}
	}
}

func handlerGet(messagesRepo *repo.MessageRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		messages := messagesRepo.MessagesListAll()
		fmt.Printf("%v\n", messages)
		for _, msg := range messages {
			component := Msg(msg.ID.String(), msg.Message)
			err := component.Render(r.Context(), w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log.Fatalf("Error rendering in HelloWebHandler: %e", err)
			}
		}
	}
}
