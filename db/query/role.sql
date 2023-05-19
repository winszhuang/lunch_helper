-- name: CreateRole :one
INSERT INTO "role" (name, description)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateRole :one
UPDATE "role"
SET name = $2, description = $3
WHERE id = $1
RETURNING *;

-- name: GetRoles :many
SELECT *
FROM "role";