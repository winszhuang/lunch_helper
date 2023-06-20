// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package db

import (
	"context"
)

type Querier interface {
	CreateFeedback(ctx context.Context, arg CreateFeedbackParams) (Feedback, error)
	CreateFood(ctx context.Context, arg CreateFoodParams) (Food, error)
	CreateOperateRecord(ctx context.Context, arg CreateOperateRecordParams) (OperateRecord, error)
	CreateRestaurant(ctx context.Context, arg CreateRestaurantParams) (Restaurant, error)
	CreateRole(ctx context.Context, arg CreateRoleParams) (Role, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateUserFood(ctx context.Context, arg CreateUserFoodParams) (UserFood, error)
	CreateUserRestaurant(ctx context.Context, arg CreateUserRestaurantParams) (UserRestaurant, error)
	DeleteFood(ctx context.Context, id int32) error
	DeleteUserFood(ctx context.Context, arg DeleteUserFoodParams) error
	DeleteUserRestaurant(ctx context.Context, arg DeleteUserRestaurantParams) error
	GetAllUserFoods(ctx context.Context, userID int32) ([]GetAllUserFoodsRow, error)
	GetAllUserRestaurants(ctx context.Context, userID int32) ([]GetAllUserRestaurantsRow, error)
	GetFeedback(ctx context.Context, arg GetFeedbackParams) ([]Feedback, error)
	GetFeedbackByDateRange(ctx context.Context, arg GetFeedbackByDateRangeParams) ([]Feedback, error)
	GetFeedbackByStatus(ctx context.Context) ([]Feedback, error)
	GetFood(ctx context.Context, id int32) (Food, error)
	GetFoodsByRestaurantId(ctx context.Context, restaurantID int32) ([]Food, error)
	GetOperateRecords(ctx context.Context, arg GetOperateRecordsParams) ([]OperateRecord, error)
	GetOperateRecordsByDateRange(ctx context.Context, arg GetOperateRecordsByDateRangeParams) ([]OperateRecord, error)
	GetOperateRecordsByUserID(ctx context.Context, arg GetOperateRecordsByUserIDParams) ([]OperateRecord, error)
	GetRestaurant(ctx context.Context, id int32) (Restaurant, error)
	GetRestaurantByGoogleMapPlaceId(ctx context.Context, googleMapPlaceID string) (Restaurant, error)
	GetRoles(ctx context.Context) ([]Role, error)
	GetUserByID(ctx context.Context, id int32) (User, error)
	GetUserByLineID(ctx context.Context, lineID string) (User, error)
	GetUserFoodByFoodId(ctx context.Context, arg GetUserFoodByFoodIdParams) (GetUserFoodByFoodIdRow, error)
	GetUserFoods(ctx context.Context, arg GetUserFoodsParams) ([]GetUserFoodsRow, error)
	GetUserRestaurants(ctx context.Context, arg GetUserRestaurantsParams) ([]GetUserRestaurantsRow, error)
	GetUsers(ctx context.Context, arg GetUsersParams) ([]User, error)
	UpdateFood(ctx context.Context, arg UpdateFoodParams) error
	UpdateMenuCrawled(ctx context.Context, arg UpdateMenuCrawledParams) error
	UpdateRole(ctx context.Context, arg UpdateRoleParams) (Role, error)
}

var _ Querier = (*Queries)(nil)
