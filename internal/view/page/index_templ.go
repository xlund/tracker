// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.663
package page

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "github.com/xlund/tracker/internal/domain"

func Index(users []domain.User) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div><h1>Chess Games Tracker</h1><p>The Chess Games Tracker is a tool designed to simplify the way you keep track of your chess games. It is built for both seasoned players and newcomers who want a straightforward and effective way to manage their chess activities.</p><section><h2>Features</h2><hr><h3>Live Updates</h3><p>Get real-time updates on scores, pairings, and standings. Stay informed with accurate and timely information throughout your tournament.</p><h3>Simple and Fast</h3><p>Our focus is always simplicity and ease of use. From this comes the benefit of being fast by being simple.</p><h3>Private</h3><p>We don't keep track of you or what you do. Privacy is at the core of our buisness</p></section></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
