package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Kwynto/gosession"
	"github.com/johan-st/jst.dev/pages"
	"github.com/sashabaranov/go-openai"
)

	func (srv *server) handleAiTranslationPost() http.HandlerFunc {
		// timing and logging
		l := srv.l.
			WithPrefix(srv.l.GetPrefix() + ".AiTranslationPost")

		defer func(t time.Time) {
			l.Info(
				logReady,
				logTimeSpent, time.Since(t),
			)
		}(time.Now())

		// setup

		OPENAI_API_KEY := os.Getenv("OPENAI_API_KEY")
		if OPENAI_API_KEY == "" {
			l.Fatal("OPENAI_API_KEY not set")
		}

		client := openai.NewClient(OPENAI_API_KEY)
		sessionHandler := newSessionHandler(l)

		// handler
		return func(w http.ResponseWriter, r *http.Request) {
			// timing and logging
			defer func(t time.Time) {
				l.Debug(
					"respond",
					logTimeSpent, time.Since(t),
				)
			}(time.Now())

			session := gosession.Start(&w, r)

			// get form data
			text := r.FormValue("text")
			lang := r.FormValue("target_lang")
			model := r.FormValue("model")

			if text == "" || lang == "" || model == "" {
				srv.respCode(http.StatusBadRequest, w, r)
				l.Warn("Bad request. field empty",
					"text", text,
					"lang", lang,
					"model", model,
				)
				return
			}

			if model != openai.GPT3Dot5Turbo && model != openai.GPT4 {
				srv.respCode(http.StatusBadRequest, w, r)
				l.Warn("recieved unhandled model field")
				return
			}

			currentTranslation := pages.Translation{Prompt: text}
			resp, err := client.CreateChatCompletion(
				r.Context(),
				openai.ChatCompletionRequest{
					Model: model,
					Messages: []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleSystem,
							Content: "I am a translation service. Give me a text to translate.",
						},
						{
							Role:    openai.ChatMessageRoleUser,
							Content: fmt.Sprintf("Translate the following text into %s: %s", lang, text),
						},
					},
				},
			)
			l.Info("openAI api response",
				"model", model,
				"tokens used", resp.Usage.TotalTokens,
			)

			if err != nil {
				srv.respCode(http.StatusInternalServerError, w, r)
				l.Error(
					"openAI api error",
					logError, err,
				)
				return
			}

			for _, c := range resp.Choices {
			currentTranslation.Choices = append(currentTranslation.Choices, c.Message.Content)
		}

		translations := sessionHandler.addTranslations(&session, currentTranslation)
		err = pages.Translated(translations).Render(r.Context(), w)
		if err != nil {
			srv.respCode(http.StatusInternalServerError, w, r)
			l.Error(
				"render 'translation' page",
				logError, err,
			)
			return
		}

	}
}

// HELPER

func prepend(slice []pages.Translation, item pages.Translation) []pages.Translation {
	// TODO:rework this solution
	slice = append(slice, pages.Translation{})
	copy(slice[1:], slice)
	slice[0] = item
	return slice

}
