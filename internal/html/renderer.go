package html

import (
	"io"

	"github.com/labstack/echo/v5"
	"maragu.dev/gomponents"
)

type GomponentRendererRender struct{}

func (renderer *GomponentRendererRender) Render(c *echo.Context, w io.Writer, templateName string, data any) error {
	c.Logger().Info("Rendering gomcomponent with name", "name", templateName)
	node, ok := data.(gomponents.Node)

	if !ok {
		return &HtmlError{Message: "This renderer only supports gomponents"}
	}

	err := node.Render(w)
	return err
}
