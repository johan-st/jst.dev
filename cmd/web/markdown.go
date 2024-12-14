package web

import (
	"bytes"
	"fmt"
	"log"

	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
)

func ComponentMarkdown(md string) templ.Component {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		log.Fatalf("failed to convert markdown to HTML: %v", err)
	}
	fmt.Println(buf.String())
	return templ.Raw(buf.String())
}
