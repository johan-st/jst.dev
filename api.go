package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/johan-st/jst.dev/pages"
	ai "github.com/sashabaranov/go-openai"
)

func (srv *server) handleApiTranslationPost() http.HandlerFunc {
	// timing and logging
	l := srv.l.
		WithPrefix(srv.l.GetPrefix() + ".ApiTranslationPost")

	defer func(t time.Time) {
		l.Debug(
			logReady,
			logTimeSpent, time.Since(t),
		)
	}(time.Now())

	// setup

	OPENAI_API_KEY := os.Getenv("OPENAI_API_KEY")
	if OPENAI_API_KEY == "" {
		l.Fatal("OPENAI_API_KEY not set")
	}

	client := ai.NewClient(OPENAI_API_KEY)

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		// timing and logging
		l := srv.l.With("handler", "ApiTranslation")
		defer func(t time.Time) {
			l.Debug(
				"respond",
				logTimeSpent, time.Since(t),
			)
		}(time.Now())

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

		if model != ai.GPT3Dot5Turbo && model != ai.GPT4 {
			srv.respCode(http.StatusBadRequest, w, r)
			l.Warn("recieved unhandled model field")
			return
		}

		currentTranslation := pages.Translation{Prompt: text}
		resp, err := client.CreateChatCompletion(
			r.Context(),
			ai.ChatCompletionRequest{
				Model: model,
				Messages: []ai.ChatCompletionMessage{
					{
						Role:    ai.ChatMessageRoleSystem,
						Content: "I am a translation service. Give me a text to translate.",
					},
					{
						Role:    ai.ChatMessageRoleUser,
						Content: fmt.Sprintf("Translate the following text into %s: %s", lang, text),
					},
				},
			},
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
		srv.translations = prepend(srv.translations, currentTranslation)

		// limit num of translations
		for len(srv.translations) > 10 {
			srv.translations = srv.translations[:10]
		}

		err = pages.Translated(srv.translations).Render(r.Context(), w)
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
