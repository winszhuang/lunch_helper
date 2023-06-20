// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: user_restaurant.sql

package db

import (
	"context"
	"database/sql"

	"github.com/shopspring/decimal"
)

const createUserRestaurant = `-- name: CreateUserRestaurant :one
INSERT INTO user_restaurant (
    user_id,
    restaurant_id
) VALUES (
    $1, $2
)
RETURNING user_id, restaurant_id
`

type CreateUserRestaurantParams struct {
	UserID       int32 `json:"user_id"`
	RestaurantID int32 `json:"restaurant_id"`
}

func (q *Queries) CreateUserRestaurant(ctx context.Context, arg CreateUserRestaurantParams) (UserRestaurant, error) {
	row := q.db.QueryRowContext(ctx, createUserRestaurant, arg.UserID, arg.RestaurantID)
	var i UserRestaurant
	err := row.Scan(&i.UserID, &i.RestaurantID)
	return i, err
}

const deleteUserRestaurant = `-- name: DeleteUserRestaurant :exec
DELETE FROM user_restaurant
WHERE user_id = $1 AND restaurant_id = $2
`

type DeleteUserRestaurantParams struct {
	UserID       int32 `json:"user_id"`
	RestaurantID int32 `json:"restaurant_id"`
}

func (q *Queries) DeleteUserRestaurant(ctx context.Context, arg DeleteUserRestaurantParams) error {
	_, err := q.db.ExecContext(ctx, deleteUserRestaurant, arg.UserID, arg.RestaurantID)
	return err
}

const getAllUserRestaurants = `-- name: GetAllUserRestaurants :many
SELECT id, name, rating, user_ratings_total, address, google_map_place_id, google_map_url, phone_number, image, menu_crawled, user_id, restaurant_id
FROM restaurant
JOIN user_restaurant ON user_restaurant.restaurant_id = restaurant.id
WHERE user_restaurant.user_id = $1
`

type GetAllUserRestaurantsRow struct {
	ID               int32           `json:"id"`
	Name             string          `json:"name"`
	Rating           decimal.Decimal `json:"rating"`
	UserRatingsTotal sql.NullInt32   `json:"user_ratings_total"`
	Address          string          `json:"address"`
	GoogleMapPlaceID string          `json:"google_map_place_id"`
	GoogleMapUrl     string          `json:"google_map_url"`
	PhoneNumber      string          `json:"phone_number"`
	Image            sql.NullString  `json:"image"`
	MenuCrawled      bool            `json:"menu_crawled"`
	UserID           int32           `json:"user_id"`
	RestaurantID     int32           `json:"restaurant_id"`
}

func (q *Queries) GetAllUserRestaurants(ctx context.Context, userID int32) ([]GetAllUserRestaurantsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllUserRestaurants, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllUserRestaurantsRow
	for rows.Next() {
		var i GetAllUserRestaurantsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Rating,
			&i.UserRatingsTotal,
			&i.Address,
			&i.GoogleMapPlaceID,
			&i.GoogleMapUrl,
			&i.PhoneNumber,
			&i.Image,
			&i.MenuCrawled,
			&i.UserID,
			&i.RestaurantID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserRestaurants = `-- name: GetUserRestaurants :many
SELECT id, name, rating, user_ratings_total, address, google_map_place_id, google_map_url, phone_number, image, menu_crawled, user_id, restaurant_id
FROM restaurant
JOIN user_restaurant ON user_restaurant.restaurant_id = restaurant.id
WHERE user_restaurant.user_id = $1
LIMIT $2 OFFSET $3
`

type GetUserRestaurantsParams struct {
	UserID int32 `json:"user_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetUserRestaurantsRow struct {
	ID               int32           `json:"id"`
	Name             string          `json:"name"`
	Rating           decimal.Decimal `json:"rating"`
	UserRatingsTotal sql.NullInt32   `json:"user_ratings_total"`
	Address          string          `json:"address"`
	GoogleMapPlaceID string          `json:"google_map_place_id"`
	GoogleMapUrl     string          `json:"google_map_url"`
	PhoneNumber      string          `json:"phone_number"`
	Image            sql.NullString  `json:"image"`
	MenuCrawled      bool            `json:"menu_crawled"`
	UserID           int32           `json:"user_id"`
	RestaurantID     int32           `json:"restaurant_id"`
}

func (q *Queries) GetUserRestaurants(ctx context.Context, arg GetUserRestaurantsParams) ([]GetUserRestaurantsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserRestaurants, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserRestaurantsRow
	for rows.Next() {
		var i GetUserRestaurantsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Rating,
			&i.UserRatingsTotal,
			&i.Address,
			&i.GoogleMapPlaceID,
			&i.GoogleMapUrl,
			&i.PhoneNumber,
			&i.Image,
			&i.MenuCrawled,
			&i.UserID,
			&i.RestaurantID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
