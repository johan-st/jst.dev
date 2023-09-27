// Code generated by templ@v0.2.334 DO NOT EDIT.

package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "fmt"

import ai "github.com/johan-st/openAI"

func OpenAI() templ.Component {
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
		_, err = templBuffer.WriteString("<div class=\"translation__form-container\"><form class=\"translation__form\" hx-post=\"/ai/translate\" hx-target=\"#translation\" hx-indicator=\"#spinner\"><h3>")
		if err != nil {
			return err
		}
		var_2 := `Text to translate:`
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h3><label for=\"text\">")
		if err != nil {
			return err
		}
		var_3 := `prompt:`
		_, err = templBuffer.WriteString(var_3)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><textarea id=\"text\" name=\"text\" rows=\"4\" cols=\"50\"></textarea><div class=\"translation__form-input-container\"><label for=\"target_lang\">")
		if err != nil {
			return err
		}
		var_4 := `Target Language:`
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label><select id=\"target_lang\" name=\"target_lang\"><option value=\"Swedish\" selected>")
		if err != nil {
			return err
		}
		var_5 := `Swedish`
		_, err = templBuffer.WriteString(var_5)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option><option value=\"English\">")
		if err != nil {
			return err
		}
		var_6 := `English`
		_, err = templBuffer.WriteString(var_6)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option><option value=\"Danish\">")
		if err != nil {
			return err
		}
		var_7 := `Danish`
		_, err = templBuffer.WriteString(var_7)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option><option value=\"Finnish\">")
		if err != nil {
			return err
		}
		var_8 := `Finnish`
		_, err = templBuffer.WriteString(var_8)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option><option value=\"French\">")
		if err != nil {
			return err
		}
		var_9 := `French`
		_, err = templBuffer.WriteString(var_9)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option><option value=\"Dutch\">")
		if err != nil {
			return err
		}
		var_10 := `Dutch`
		_, err = templBuffer.WriteString(var_10)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option><option value=\"German\">")
		if err != nil {
			return err
		}
		var_11 := `German`
		_, err = templBuffer.WriteString(var_11)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option></select><input class=\"button\" type=\"submit\" value=\"Translate\" hx-disabled-elt=\"this\"></div></form></div><div id=\"translation\">")
		if err != nil {
			return err
		}
		err = Translated(ai.Translation{}).Render(ctx, templBuffer)
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

func Translated(tran ai.Translation) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_12 := templ.GetChildren(ctx)
		if var_12 == nil {
			var_12 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		err = spinner("spinner").Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<div class=\"translation__results\"><h3>")
		if err != nil {
			return err
		}
		var_13 := `Translated Text:`
		_, err = templBuffer.WriteString(var_13)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h3>")
		if err != nil {
			return err
		}
		if len(tran.Choices) == 0 {
			_, err = templBuffer.WriteString("<p>")
			if err != nil {
				return err
			}
			var_14 := `nothing yet...`
			_, err = templBuffer.WriteString(var_14)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</p>")
			if err != nil {
				return err
			}
		}
		for i, c := range tran.Choices {
			_, err = templBuffer.WriteString("<h4>")
			if err != nil {
				return err
			}
			var_15 := `Choice `
			_, err = templBuffer.WriteString(var_15)
			if err != nil {
				return err
			}
			var var_16 string = fmt.Sprintf("%d", i+1)
			_, err = templBuffer.WriteString(templ.EscapeString(var_16))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h4> <p>")
			if err != nil {
				return err
			}
			var var_17 string = c
			_, err = templBuffer.WriteString(templ.EscapeString(var_17))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</p>")
			if err != nil {
				return err
			}
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
