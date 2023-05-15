-- name: CreateRestaurant :one
INSERT INTO restaurant (
    name,
    rating,
    user_ratings_total,
    address,
    google_map_place_id,
    google_map_url,
    phone_number
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;
