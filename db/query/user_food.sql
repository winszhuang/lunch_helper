-- name: CreateUserFood :one
INSERT INTO user_food (
    user_id,
    food_id
) VALUES (
    $1, $2
)
RETURNING *;

-- name: DeleteUserFood :exec
DELETE FROM user_food
WHERE user_id = $1 AND food_id = $2;
