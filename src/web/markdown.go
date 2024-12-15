package web

import (
	"bytes"
	"log"

	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
)

func ComponentMarkdown(md string) templ.Component {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		log.Fatalf("failed to convert markdown to HTML: %v", err)
	}
	return templ.Raw(buf.String())
}
