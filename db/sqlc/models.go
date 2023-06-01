// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package db

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type Feedback struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	FoodID    int32     `json:"food_id"`
	EditBy    int32     `json:"edit_by"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
}

type Food struct {
	ID           int32          `json:"id"`
	Name         string         `json:"name"`
	Price        string         `json:"price"`
	Image        sql.NullString `json:"image"`
	Description  sql.NullString `json:"description"`
	RestaurantID int32          `json:"restaurant_id"`
	Version      int16          `json:"version"`
	EditBy       sql.NullInt32  `json:"edit_by"`
}

type OperateRecord struct {
	ID              int32          `json:"id"`
	UserID          int32          `json:"user_id"`
	FoodID          int32          `json:"food_id"`
	Before          sql.NullString `json:"before"`
	After           sql.NullString `json:"after"`
	UpdateAt        time.Time      `json:"update_at"`
	OperateCategory int16          `json:"operate_category"`
}

type Restaurant struct {
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
}

type Role struct {
	ID          int32          `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
}

type User struct {
	ID                     int32  `json:"id"`
	LineID                 string `json:"line_id"`
	Name                   string `json:"name"`
	Picture                string `json:"picture"`
	GoogleMapsApiCallCount int16  `json:"google_maps_api_call_count"`
	RoleID                 int32  `json:"role_id"`
}

type UserFood struct {
	UserID int32 `json:"user_id"`
	FoodID int32 `json:"food_id"`
}

type UserRestaurant struct {
	UserID       int32 `json:"user_id"`
	RestaurantID int32 `json:"restaurant_id"`
}
