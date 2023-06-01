-- name: CreateRestaurant :one
INSERT INTO restaurant (
    name,
    rating,
    user_ratings_total,
    address,
    google_map_place_id,
    google_map_url,
    phone_number,
    image
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetRestaurant :one
SELECT * FROM Restaurant WHERE id = $1;

-- name: GetRestaurantByGoogleMapPlaceId :one
SELECT * FROM Restaurant WHERE google_map_place_id = $1;

-- name: UpdateMenuCrawled :exec
UPDATE Restaurant SET menu_crawled = $1 WHERE id = $2;
