
-- name: CreateFood :one
INSERT INTO food (
    name,
    price,
    image,
    description,
    restaurant_id,
    edit_by
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateFood :exec
UPDATE food
SET name = $1, price = $2, image = $3, edit_by = $4, version = version + 1
WHERE id = $5;

-- name: DeleteFood :exec
DELETE FROM food
WHERE id = $1;

-- name: GetFoodsByRestaurantId :many
SELECT * FROM food WHERE restaurant_id = $1;
