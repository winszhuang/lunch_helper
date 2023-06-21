package service

import (
	"context"
	"database/sql"
	db "lunch_helper/db/sqlc"
	"lunch_helper/food_deliver/model"
	"lunch_helper/util"
)

type CreateFoodByDishParams struct {
	Dish         model.Dish
	RestaurantID int32
	EditBy       sql.NullInt32
}

type CreateFoodsByDishesParams struct {
	Dishes       []model.Dish
	RestaurantID int32
	EditBy       sql.NullInt32
}

type FoodService struct {
	dbStore db.Store
}

func NewFoodService(dbStore db.Store) *FoodService {
	return &FoodService{
		dbStore: dbStore,
	}
}

func (fs *FoodService) CreateFood(ctx context.Context, arg db.CreateFoodParams) (db.Food, error) {
	list, err := fs.GetFoods(ctx, arg.RestaurantID)
	if err == nil {
		for _, food := range list {
			if food.Name == arg.Name {
				return food, nil
			}
		}
	}

	return fs.dbStore.CreateFood(ctx, arg)
}

func (fs *FoodService) CreateFoodByDish(ctx context.Context, arg CreateFoodByDishParams) (db.Food, error) {
	dish := arg.Dish
	params := db.CreateFoodParams{
		Name:         dish.Name,
		Price:        dish.Price,
		Image:        util.CheckNullString(dish.Image),
		Description:  util.CheckNullString(dish.Description),
		RestaurantID: arg.RestaurantID,
		EditBy:       arg.EditBy,
	}
	return fs.CreateFood(ctx, params)
}

func (fs *FoodService) CreateFoodsByDishes(ctx context.Context, arg CreateFoodsByDishesParams) ([]db.Food, []error) {
	var foods []db.Food
	var errList []error
	for _, dish := range arg.Dishes {
		if food, err := fs.CreateFoodByDish(ctx, CreateFoodByDishParams{
			Dish:         dish,
			RestaurantID: arg.RestaurantID,
			EditBy:       arg.EditBy,
		}); err != nil {
			errList = append(errList, err)
		} else {
			foods = append(foods, food)
		}
	}
	return foods, errList
}

func (fs *FoodService) GetFood(ctx context.Context, id int32) (db.Food, error) {
	return fs.dbStore.GetFood(ctx, id)
}

func (fs *FoodService) GetFoods(ctx context.Context, restaurantID int32) ([]db.Food, error) {
	return fs.dbStore.GetFoodsByRestaurantId(ctx, restaurantID)
}
