-- name: ExistsWordList :one
SELECT EXISTS (
	SELECT 1
	FROM app.word_lists
	WHERE title = $1
);

-- name: CreateWordList :one
INSERT INTO app.word_lists (title, words)
VALUES ($1, $2)
RETURNING id;