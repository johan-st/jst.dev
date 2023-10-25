package main

import (
	"net/http"
	"os"
	"time"

	"github.com/johan-st/jst.dev/pages"
	ai "github.com/johan-st/openAI"
)

func (srv *server) handleAiTranslation() http.HandlerFunc {
	// timing and logging
	l := srv.l.With("handler", "ApiTranslation")
	defer func(t time.Time) {
		l.Info(
			"ready",
			"time", time.Since(t),
		)
	}(time.Now())

	// setup
	var err error

	OPENAI_API_KEY := os.Getenv("OPENAI_API_KEY")
	if OPENAI_API_KEY == "" {
		l.Fatal("OPENAI_API_KEY not set")
	}

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		// timing and logging
		l := srv.l.With("handler", "ApiTranslation")
		defer func(t time.Time) {
			l.Debug(
				"responded",
				"time", time.Since(t),
			)
		}(time.Now())

		// get form data
		text := r.FormValue("text")
		lang := r.FormValue("target_lang")

		l.Debug("got form data", "formdata", r.Form)

		if text == "" || lang == "" {
			l.Info(
				"bad request",
				"text", text,
				"lang", lang,
			)
			srv.respCode(http.StatusBadRequest, w, r)
			return
		}

		// translate
		time.Sleep(4 * time.Second)
		translation := ai.Translation{Prompt: "test prompt", Choices: []string{text, lang, "en blåval"}}
		// translation, err := ai.Translate(OPENAI_API_KEY, lang, text)
		// if err != nil {
		// 	l.Error(
		// 		"translate",
		// 		"error", err,
		// 	)
		// 	srv.respCode(http.StatusInternalServerError, w, r)
		// 	return
		// }

		l.Debug("got translation", "translation", translation)

		err = pages.Translated(translation).Render(r.Context(), w)
		if err != nil {
			srv.respCode(http.StatusInternalServerError, w, r)
			l.Error(
				"render",
				"error", err,
			)
			return
		}

	}
}

func (srv *server) handleTestAI() http.HandlerFunc {
	// timing and logging
	l := srv.l.With("handler", "ApiTranslation")
	defer func(t time.Time) {
		l.Info(
			"ready",
			"time", time.Since(t),
		)
	}(time.Now())

	// setup

	OPENAI_API_KEY := os.Getenv("OPENAI_API_KEY")
	if OPENAI_API_KEY == "" {
		l.Fatal("OPENAI_API_KEY not set")
	}
	// chan, err := ai.TranslateStream(OPENAI_API_KEY, "Spanish", "snakar du norsk torsk eller är du glad att se mig (sexual inuendo)?")
	// for cha  
	// }
	// if err != nil {
	// 	l.Fatal("test stream", "error", err)
	// }

	// handler
	return func(w http.ResponseWriter, r *http.Request) {
		// timing and logging
		l := srv.l.With("handler", "ApiTranslation")
		defer func(t time.Time) {
			l.Debug(
				"responded",
				"time", time.Since(t),
			)
		}(time.Now())

		srv.respCode(http.StatusNotImplemented, w, r)
	}
}
