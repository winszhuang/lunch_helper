// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: food.sql

package db

import (
	"context"
	"database/sql"
)

const createFood = `-- name: CreateFood :one
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
RETURNING id, name, price, image, description, restaurant_id, version, edit_by
`

type CreateFoodParams struct {
	Name         sql.NullString `json:"name"`
	Price        sql.NullString `json:"price"`
	Image        sql.NullString `json:"image"`
	Description  sql.NullString `json:"description"`
	RestaurantID sql.NullInt32  `json:"restaurant_id"`
	EditBy       sql.NullInt32  `json:"edit_by"`
}

func (q *Queries) CreateFood(ctx context.Context, arg CreateFoodParams) (Food, error) {
	row := q.db.QueryRowContext(ctx, createFood,
		arg.Name,
		arg.Price,
		arg.Image,
		arg.Description,
		arg.RestaurantID,
		arg.EditBy,
	)
	var i Food
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Price,
		&i.Image,
		&i.Description,
		&i.RestaurantID,
		&i.Version,
		&i.EditBy,
	)
	return i, err
}

const deleteFood = `-- name: DeleteFood :exec
DELETE FROM food
WHERE id = $1
`

func (q *Queries) DeleteFood(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteFood, id)
	return err
}

const updateFood = `-- name: UpdateFood :exec
UPDATE food
SET name = $1, price = $2, image = $3, edit_by = $4, version = version + 1
WHERE id = $5
`

type UpdateFoodParams struct {
	Name   sql.NullString `json:"name"`
	Price  sql.NullString `json:"price"`
	Image  sql.NullString `json:"image"`
	EditBy sql.NullInt32  `json:"edit_by"`
	ID     int32          `json:"id"`
}

func (q *Queries) UpdateFood(ctx context.Context, arg UpdateFoodParams) error {
	_, err := q.db.ExecContext(ctx, updateFood,
		arg.Name,
		arg.Price,
		arg.Image,
		arg.EditBy,
		arg.ID,
	)
	return err
}
