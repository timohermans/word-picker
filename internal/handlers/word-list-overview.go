package handlers

import (
	"hert/gotest/internal/db"
	"hert/gotest/internal/html"
	"strings"

	"fmt"
	"strconv"

	"github.com/labstack/echo/v5"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

const (
	urlWordListOverview = "/"
)

func HandleWordListOverview(e *echo.Echo, queries *db.Queries) {
	e.GET(urlWordListOverview, func(c *echo.Context) error {
		ctx := c.Request().Context()
		wordLists, err := queries.ListWordLists(ctx)

		if err != nil {
			c.Redirect(302, "/error")
		}

		return c.Render(200, "home", pageWordListOverview(wordLists))
	})
}

func pageWordListOverview(wordLists []db.AppWordList) Node {
	return html.Page(html.PageProps{
		Title:       "Picker",
		Description: "Page where you can create lists and select words. The selected word cannot be reselected until all words have been seleced once",
	},
		H2(Text("Picker")),
		Div(A(Href("/word-list/add"), Text("Add new list"))),

		If(len(wordLists) == 0, Div(Text("No lists yet."))),
		If(len(wordLists) > 0,
			Ul(
				Map(wordLists, func(wordList db.AppWordList) Node {
					id := strconv.Itoa(int(wordList.ID))
					pickUrl := strings.Replace(urlWordListPick, ":id", id, 1)
					listText := fmt.Sprintf("%s (%s)", wordList.Title, strconv.Itoa(len(wordList.GetWords())))

					return Li(
						A(Href(pickUrl), Text(listText)),
					)
				}),
			)),
	)
}
