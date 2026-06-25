-- name: IsReady :one
SELECT EXISTS(
    SELECT 1
);