package pages

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	goldmark_meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

var (
	md   goldmark.Markdown
	once sync.Once
)

type PageMeta map[string]interface{}

type Page struct {
	Title string
	Body  []byte
	Path  string
	PageMeta
}

func (p Page) Render(ctx context.Context, w io.Writer) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return md.Convert(p.Body, w)
}

func MdToPage(file []byte, baseUrl string) (Page, error) {
	// setup markdown parser
	once.Do(func() {
		md = goldmark.New(
			goldmark.WithExtensions(
				highlighting.NewHighlighting(
					highlighting.WithStyle("native"),
				),
				extension.GFM,
				goldmark_meta.New(
					goldmark_meta.WithStoresInDocument(),
					// goldmark_meta.WithTable(),
				),
				extension.Table,
			),
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
		)
	})
	// make sure we start and end with a slash
	if baseUrl == "" {
		baseUrl = "/"
	}
	if baseUrl[len(baseUrl)-1] != '/' {
		baseUrl += "/"
	}
	if baseUrl[0] != '/' {
		baseUrl = "/" + baseUrl
	}

	document := md.Parser().Parse(text.NewReader(file))
	metaData := document.OwnerDocument().Meta()
	title, ok := metaData["title"]
	if !ok {
		return Page{}, fmt.Errorf("no title found on post. file starts %s", file[:100])
	}
	switch title.(type) {
	case string:
	default:
		return Page{}, fmt.Errorf("title is not a string. file starts %s", file[:100])
	}

	slug, ok := metaData["path"]
	if !ok {
		return Page{}, fmt.Errorf("no path found on post. file starts %s", file[:100])
	}
	switch slug.(type) {
	case string:
	default:
		return Page{}, fmt.Errorf("path is not a string. file starts %s", file[:100])
	}

	// make sure we DON'T have a trailing slash
	if slug.(string) == "/" {
		slug = ""
	}
	if slug.(string)[len(slug.(string))-1] == '/' {
		slug = slug.(string)[:len(slug.(string))-1]
	}

	buf := &bytes.Buffer{}
	md.Convert(file, buf)

	return Page{
		Title:    title.(string), // we checked
		Body:     buf.Bytes(),
		Path:     baseUrl + slug.(string), // we checked
		PageMeta: metaData,
	}, nil
}
