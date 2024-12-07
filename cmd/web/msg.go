package web

import (
	"net/http"

	"jst.dev/internal/repo"
)

func HandlerMsg(messagesRepo *repo.MessageRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		messages := messagesRepo.MessagesListAll()
		component := MsgList(messages)
		if err := component.Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
