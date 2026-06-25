package db

import (
	"strings"
)

func (wordList *AppWordList) GetWords() []string {
	return strings.Split(wordList.Words, ",")
}
