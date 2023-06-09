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

-- name: GetUserFoods :many
SELECT *
FROM food
JOIN user_food ON user_food.food_id = food.id
WHERE user_food.user_id = $1
LIMIT $2 OFFSET $3;
