package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type PageProps struct {
	Title       string
	Description string
}

func Page(props PageProps, children ...Node) Node {
	return HTML5(HTML5Props{
		Title:       props.Title,
		Description: props.Description,
		Language:    "en",
		Head:        []Node{},
		// Head: []Node{
		// 	Link(Rel("stylesheet"), Href(appCSSPath)),
		// 	Script(Src(htmxJSPath), Defer()),
		// 	Script(Src(appJSPath), Defer()),
		// },
		Body: []Node{
			Header(
				A(Text("Lijsten"), Href("/")),
			),
			Div(
				Group(children),
			),
		},
	})

}
