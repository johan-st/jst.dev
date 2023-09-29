package main

import (
	"github.com/Kwynto/gosession"
	"github.com/charmbracelet/log"
	"github.com/johan-st/jst.dev/pages"
)

type sessionHandler struct {
	l *log.Logger
}

// type session struct {
// 	translations []pages.Translation
// }

func newSessionHandler(l *log.Logger) sessionHandler {
	// setup
	l = l.WithPrefix(l.GetPrefix() + ".sessionsHandler")

	// return
	return sessionHandler{l}
}

func (s *sessionHandler) getTranslations(id *gosession.SessionId) []pages.Translation {
	translations := []pages.Translation{}
	translationsMaybe := id.Get("translations")

	switch trans := translationsMaybe.(type) {
	case []pages.Translation:
		s.l.Debug("restore translations", "num", len(trans))
		translations = trans
	case nil:
		s.l.Debug("translations not found in session, setting a new empty one")
		id.Set("translations", []pages.Translation{})
	default:

		s.l.Warn("translations wrong type, setting a new empty one", "translations", translationsMaybe)
		id.Set("translations", []pages.Translation{})
	}

	return translations
}

func (s *sessionHandler) addTranslations(id *gosession.SessionId, tran pages.Translation) []pages.Translation {
	translations := s.getTranslations(id)
	translations = prepend(translations, tran)

	// limit num of translations in session
	for len(translations) > 10 {
		translations = translations[:10]
	}
	id.Set("translations", translations)

	return translations
}
