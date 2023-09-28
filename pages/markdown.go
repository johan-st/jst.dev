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

type PostMeta map[string]interface{}

type Post struct {
	Title string
	Body  []byte
	Path  string
	PostMeta
}

func (p Post) Render(ctx context.Context, w io.Writer) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return md.Convert(p.Body, w)
}

func FileToPost(file []byte, basePath string) (Post, error) {
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
	// make sure we have a start and end with a slash
	if basePath == "" {
		basePath = "/"
	}
	if basePath[len(basePath)-1] != '/' {
		basePath += "/"
	}
	if basePath[0] != '/' {
		basePath = "/" + basePath
	}

	document := md.Parser().Parse(text.NewReader(file))
	metaData := document.OwnerDocument().Meta()
	title, ok := metaData["title"]
	if !ok {
		return Post{}, fmt.Errorf("no title found on post. file starts %s", file[:60])
	}
	switch title.(type) {
	case string:
	default:
		return Post{}, fmt.Errorf("title is not a string. file starts %s", file[:60])
	}

	slug, ok := metaData["path"]
	if !ok {
		return Post{}, fmt.Errorf("no path found on post. file starts %s", file[:60])
	}
	switch slug.(type) {
	case string:
	default:
		return Post{}, fmt.Errorf("path is not a string. file starts %s", file[:60])
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

	return Post{
		Title:    title.(string), // we checked
		Body:     buf.Bytes(),
		Path:     basePath + slug.(string), // we checked
		PostMeta: metaData,
	}, nil
}
