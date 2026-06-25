-- name: GetWordListToPick :one
SELECT wl.id, wl.title, wl.words, (COALESCE(STRING_AGG(wp.word, ','), ''))::text as words_picked
FROM app.word_lists wl
LEFT JOIN app.words_picked wp ON wl.id = wp.word_list_id
WHERE wl.id = $1
GROUP BY wl.id, wl.title, wl.words;

-- name: GetWordListHistory :many
SELECT word, picked_at
FROM app.words_picked_history
WHERE word_list_id = $1
ORDER BY picked_at DESC
LIMIT 14;

-- name: CreateWordPicked :exec
INSERT INTO app.words_picked (word_list_id, word)
VALUES ($1, $2);

-- name: ClearWordsPickedFromList :exec
DELETE FROM app.words_picked
WHERE word_list_id = $1;