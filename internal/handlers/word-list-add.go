package handlers

import (
	"hert/gotest/internal/db"
	h "hert/gotest/internal/html"
	"strconv"
	"strings"

	"github.com/labstack/echo/v5"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type WordListCommand struct {
	Title string `form:"title"`
	Words string `form:"words"`
}

type PagePropsWordListAdd struct {
	Message string
}

const (
	urlWordListAdd = "/word-list/add"
)

func HandleWordListAdd(e *echo.Echo, queries *db.Queries) {
	url := urlWordListAdd
	e.GET(url, func(c *echo.Context) error {
		return c.Render(200, "word-list-add", pageWordListAdd(&PagePropsWordListAdd{}))
	})

	e.POST(url, func(c *echo.Context) error {
		ctx := c.Request().Context()

		var wordList WordListCommand
		props := &PagePropsWordListAdd{}
		if err := c.Bind(&wordList); err != nil {
			c.Logger().Warn(err.Error())
			props.Message = "Something went wrong binding the form to the code."
		}

		exists, err := queries.ExistsWordList(ctx, wordList.Title)

		if err != nil {
			c.Logger().Error(err.Error())
			return c.Render(200, "word-list-add", pageWordListAdd(props))
		}

		if exists {
			props.Message = "A list with this title already exists."
			return c.Render(200, "word-list-add", pageWordListAdd(props))
		}

		id, err := queries.CreateWordList(ctx, db.CreateWordListParams{
			Title: wordList.Title,
			Words: strings.ReplaceAll(wordList.Words, "\r\n", ","),
		})

		if err != nil {
			props.Message = "Something went wrong adding the list."
			c.Logger().Error(err.Error())
			return c.Render(200, "word-list-add", pageWordListAdd(props))
		}

		url := strings.Replace(urlWordListPick, ":id", strconv.Itoa(int(id)), 1)
		return c.Redirect(302, url)
	})
}

func pageWordListAdd(props *PagePropsWordListAdd) Node {
	return h.Page(
		h.PageProps{Title: "Add new list", Description: "Add a new list of words."},
		H2(Text("Voeg nieuwe woordenlijst toe")),
		If(props.Message != "", Div(Class("error"), Text(props.Message))),
		Form(
			Method("POST"),
			Div(
				Label(For("title"), Text("Title")),
				Input(ID("title"), Name("title")),
			),
			Div(
				Label(For("words"), Text("Words")),
				Textarea(ID("words"), Name("words"), Rows("4"), Placeholder("Word 1\nWord 2")),
			),
			Div(Button(Type("submit"), Text("Add"))),
		),
	)
}
