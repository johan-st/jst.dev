// Code generated by templ@v0.2.334 DO NOT EDIT.

package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "fmt"

func Content(p Page) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_1 := templ.GetChildren(ctx)
		if var_1 == nil {
			var_1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<h1>")
		if err != nil {
			return err
		}
		var var_2 string = p.Title
		_, err = templBuffer.WriteString(templ.EscapeString(var_2))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1><div class=\"content\" id=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString("page_" + p.Path))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\">")
		if err != nil {
			return err
		}
		err = p.Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func Blog(pages *[]Page) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_3 := templ.GetChildren(ctx)
		if var_3 == nil {
			var_3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<h1>")
		if err != nil {
			return err
		}
		var_4 := `Docs`
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1><div>")
		if err != nil {
			return err
		}
		err = listPages(pages).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func Blog404(pages *[]Page) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_5 := templ.GetChildren(ctx)
		if var_5 == nil {
			var_5 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<h1>")
		if err != nil {
			return err
		}
		var_6 := `404`
		_, err = templBuffer.WriteString(var_6)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1><p>")
		if err != nil {
			return err
		}
		var_7 := `This article does not exist..  Here are some alternatives`
		_, err = templBuffer.WriteString(var_7)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p>")
		if err != nil {
			return err
		}
		err = listPages(pages).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func listPages(pages *[]Page) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_8 := templ.GetChildren(ctx)
		if var_8 == nil {
			var_8 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, err = templBuffer.WriteString("<style>")
		if err != nil {
			return err
		}
		var_9 := `
        .page-list {
            list-style: none;
            padding: 0;
        }
        .page-list li:last-child {
            border: none;
        }
        .page-list li h3 {
            font-size: 1.2em;
            font-weight: bold;
            color: var(--clr-primary);
            margin: 0;
        }
        .page-list li a {
            margin-bottom: 1em;
            padding: 1em;
            background: var(--clr-background-alt);
            border-radius: var(--border-radius);
            border: 1px solid transparent;
            text-decoration: none;
            display: block;
            color: var(--clr-text);
            // transition: all 0.2s ease-in-out;
        }
        .page-list li a:hover {
            // box-shadow: var(--shadow-glow);
            // box-shadow: var(--shadow);
            border: 1px solid var(--clr-secondary);
            // transform: scale(1.01);
        }
        .page-list li a:hover>h3 {
            text-decoration: underline;
            text-decoration-color: var(--clr-secondary);
        }
        .page-list li .date {
            margin-top: 0;
            font-size: 0.8em;
        }
        .page-list .meta-key {
            margin-top: 1rem;
            font-weight: 100;
            color: var(--clr-text-muted);
            display:block;
        }

    `
		_, err = templBuffer.WriteString(var_9)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</style><ul hx-boost=\"true\" class=\"page-list\">")
		if err != nil {
			return err
		}
		for _, p := range *pages {
			_, err = templBuffer.WriteString("<li><a href=\"")
			if err != nil {
				return err
			}
			var var_10 templ.SafeURL = templ.SafeURL(p.Path)
			_, err = templBuffer.WriteString(templ.EscapeString(string(var_10)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\"><h3>")
			if err != nil {
				return err
			}
			var var_11 string = p.Title
			_, err = templBuffer.WriteString(templ.EscapeString(var_11))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h3>")
			if err != nil {
				return err
			}
			if v, ok := p.PageMeta["date"]; ok {
				_, err = templBuffer.WriteString("<div class=\"date\">")
				if err != nil {
					return err
				}
				var var_12 string = fmt.Sprintf("%v", v)
				_, err = templBuffer.WriteString(templ.EscapeString(var_12))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div>")
				if err != nil {
					return err
				}
			}
			if v, ok := p.PageMeta["description"]; ok {
				_, err = templBuffer.WriteString("<div><span class=\"meta-key\">")
				if err != nil {
					return err
				}
				var_13 := `description:`
				_, err = templBuffer.WriteString(var_13)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</span>")
				if err != nil {
					return err
				}
				var var_14 string = fmt.Sprintf("%v", v)
				_, err = templBuffer.WriteString(templ.EscapeString(var_14))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div>")
				if err != nil {
					return err
				}
			}
			if v, ok := p.PageMeta["excerpt"]; ok {
				_, err = templBuffer.WriteString("<div><span class=\"meta-key\">")
				if err != nil {
					return err
				}
				var_15 := `excerpt:`
				_, err = templBuffer.WriteString(var_15)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</span>")
				if err != nil {
					return err
				}
				var var_16 string = fmt.Sprintf("%v", v)
				_, err = templBuffer.WriteString(templ.EscapeString(var_16))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div>")
				if err != nil {
					return err
				}
			}
			if v, ok := p.PageMeta["embedding"]; ok {
				_, err = templBuffer.WriteString("<div><span class=\"meta-key\">")
				if err != nil {
					return err
				}
				var_17 := `embedding:`
				_, err = templBuffer.WriteString(var_17)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</span>")
				if err != nil {
					return err
				}
				var var_18 string = fmt.Sprintf("%v", v)
				_, err = templBuffer.WriteString(templ.EscapeString(var_18))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div>")
				if err != nil {
					return err
				}
			}
			if v, ok := p.PageMeta["tags"]; ok {
				_, err = templBuffer.WriteString("<div><span class=\"meta-key\">")
				if err != nil {
					return err
				}
				var_19 := `tags:`
				_, err = templBuffer.WriteString(var_19)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</span>")
				if err != nil {
					return err
				}
				var var_20 string = fmt.Sprintf("%v", v)
				_, err = templBuffer.WriteString(templ.EscapeString(var_20))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div>")
				if err != nil {
					return err
				}
			}
			if v, ok := p.PageMeta["keywords"]; ok {
				_, err = templBuffer.WriteString("<div><span class=\"meta-key\">")
				if err != nil {
					return err
				}
				var_21 := `keywords:`
				_, err = templBuffer.WriteString(var_21)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</span>")
				if err != nil {
					return err
				}
				var var_22 string = fmt.Sprintf("%v", v)
				_, err = templBuffer.WriteString(templ.EscapeString(var_22))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</a></li>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</ul>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
