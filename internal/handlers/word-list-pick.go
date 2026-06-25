package handlers

import (
	"fmt"
	"hert/gotest/internal/db"
	"hert/gotest/internal/html"
	"math/rand"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/labstack/echo/v5"
	"maragu.dev/gomponents"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

const urlWordListPick = "/word-list/:id/pick"

func HandleWordListPick(e *echo.Echo, queries *db.Queries) {
	e.GET(urlWordListPick, func(c *echo.Context) error {
		ctx := c.Request().Context()
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)

		if err != nil {
			return fmt.Errorf("Converting id param %s to int: %w", idParam, err)
		}

		wordListId := int32(id)

		wordList, err := queries.GetWordListToPick(ctx, wordListId)
		if err != nil {
			return fmt.Errorf("Getting word list %d for picking: %w", wordListId, err)
		}

		history, err := queries.GetWordListHistory(ctx, wordListId)

		if err != nil {
			return fmt.Errorf("Getting word list history %d for picking: %w", wordListId, err)
		}

		pageProps := &PagePropsWordListPick{
			WordList: wordList,
			History:  history,
		}
		return c.Render(http.StatusOK, "word-pick", pageWordListPick(pageProps))
	})

	e.POST(urlWordListPick, func(c *echo.Context) error {
		ctx := c.Request().Context()
		idValue := c.FormValue("id")
		id, err := strconv.Atoi(idValue)

		if err != nil {
			return fmt.Errorf("Parsing id param %s: %w", idValue, err)
		}

		wordList, err := queries.GetWordListToPick(ctx, int32(id))

		if err != nil {
			return fmt.Errorf("Fetching word list %s for picking: %w", idValue, err)
		}

		words := strings.Split(wordList.Words, ",")
		wordsPicked := strings.Split(wordList.WordsPicked, ",")

		if len(words) == len(wordsPicked) {
			err = queries.ClearWordsPickedFromList(ctx, int32(id))
			if err != nil {
				return fmt.Errorf("Clearing word list %s for picking: %w", idValue, err)
			}
			wordsPicked = []string{}
		}

		picked := make(map[string]struct{})
		wordsAvailable := []string{}
		for _, word := range wordsPicked {
			picked[word] = struct{}{}
		}
		for _, word := range words {
			if _, exists := picked[word]; !exists {
				wordsAvailable = append(wordsAvailable, word)
			}
		}
		shouldClearWordsPicked := len(wordsAvailable) == 1

		wordPicked := wordsAvailable[rand.Intn(len(wordsAvailable))]

		err = queries.CreateWordPicked(ctx, db.CreateWordPickedParams{WordListID: int32(id), Word: wordPicked})
		if err != nil {
			return fmt.Errorf("Error inserting picked word %s in list %s: %w", wordPicked, idValue, err)
		}

		if shouldClearWordsPicked {
			err = queries.ClearWordsPickedFromList(ctx, int32(id))
			if err != nil {
				return fmt.Errorf("Clearing word list %s for picking: %w", idValue, err)
			}

			wordList.WordsPicked = ""
		} else {
			wordList.WordsPicked = fmt.Sprintf("%s,%s", wordList.WordsPicked, wordPicked)
		}

		history, err := queries.GetWordListHistory(ctx, wordList.ID)

		if err != nil {
			return fmt.Errorf("Fetching history for word list %s for picking: %w", idValue, err)
		}

		pageProps := &PagePropsWordListPick{
			WordList:   wordList,
			History:    history,
			WordPicked: wordPicked,
		}
		return c.Render(http.StatusOK, "word-pick", pageWordListPick(pageProps))
	})
}

type PagePropsWordListPick struct {
	WordList   db.GetWordListToPickRow
	History    []db.GetWordListHistoryRow
	WordPicked string
}

func pageWordListPick(pageProps *PagePropsWordListPick) gomponents.Node {
	words := strings.Split(pageProps.WordList.Words, ",")
	wordsPicked := strings.Split(pageProps.WordList.WordsPicked, ",")
	id := int(pageProps.WordList.ID)

	return html.Page(
		html.PageProps{Title: "Pak er een", Description: "Pagina om een woord te pakken"},
		If(pageProps.WordPicked != "",
			Section(
				H2(Text("Gepakt")),
				P(Text(pageProps.WordPicked)),
			)),
		Section(
			H2(Text(pageProps.WordList.Title)),
			Ul(
				Map(words, func(word string) Node {
					return Li(Text(formatWordText(word, wordsPicked)))
				}),
			),
			Form(
				Method("post"),
				Input(Type("hidden"), Name("id"), Value(strconv.Itoa(id))),
				Button(Type("submit"), Text("Pak er een!")),
			),
			Section(
				H2(Text("Onlangs gekozen")),
				Ul(
					Map(pageProps.History, func(word db.GetWordListHistoryRow) Node {
						return Li(Text(fmt.Sprintf("%s - %s", word.PickedAt.Time.Format("2006-01-02 15:04:05"), word.Word)))
					}),
				),
			),
		),
	)
}

func formatWordText(word string, wordsPicked []string) string {
	if slices.Contains(wordsPicked, word) {
		return fmt.Sprintf("%s (uit)", word)
	}

	return word
}
