-- name: CreateUser :one
INSERT INTO "user" (
    line_id,
    name,
    picture
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetUserByID :one
SELECT *
FROM "user"
WHERE id = $1;

-- name: GetUsers :many
SELECT *
FROM "user"
LIMIT $1 OFFSET $2;
