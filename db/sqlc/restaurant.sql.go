// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: restaurant.sql

package db

import (
	"context"
	"database/sql"

	"github.com/shopspring/decimal"
)

const createRestaurant = `-- name: CreateRestaurant :one
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
RETURNING id, name, rating, user_ratings_total, address, google_map_place_id, google_map_url, phone_number
`

type CreateRestaurantParams struct {
	Name             string          `json:"name"`
	Rating           decimal.Decimal `json:"rating"`
	UserRatingsTotal sql.NullInt32   `json:"user_ratings_total"`
	Address          string          `json:"address"`
	GoogleMapPlaceID string          `json:"google_map_place_id"`
	GoogleMapUrl     string          `json:"google_map_url"`
	PhoneNumber      string          `json:"phone_number"`
}

func (q *Queries) CreateRestaurant(ctx context.Context, arg CreateRestaurantParams) (Restaurant, error) {
	row := q.db.QueryRowContext(ctx, createRestaurant,
		arg.Name,
		arg.Rating,
		arg.UserRatingsTotal,
		arg.Address,
		arg.GoogleMapPlaceID,
		arg.GoogleMapUrl,
		arg.PhoneNumber,
	)
	var i Restaurant
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Rating,
		&i.UserRatingsTotal,
		&i.Address,
		&i.GoogleMapPlaceID,
		&i.GoogleMapUrl,
		&i.PhoneNumber,
	)
	return i, err
}
