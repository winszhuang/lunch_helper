
-- name: CreateUserRestaurant :one
INSERT INTO user_restaurant (
    user_id,
    restaurant_id
) VALUES (
    $1, $2
)
RETURNING *;

-- name: DeleteUserRestaurant :exec
DELETE FROM user_restaurant
WHERE user_id = $1 AND restaurant_id = $2;

-- name: GetUserRestaurants :many
SELECT *
FROM restaurant
JOIN user_restaurant ON user_restaurant.restaurant_id = restaurant.id
WHERE user_restaurant.user_id = $1
LIMIT $2 OFFSET $3;
