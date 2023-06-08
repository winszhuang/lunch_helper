package service

import (
	"context"
	db "lunch_helper/db/sqlc"
)

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

func (fs *FoodService) GetFood(ctx context.Context, id int32) (db.Food, error) {
	return fs.dbStore.GetFood(ctx, id)
}

func (fs *FoodService) GetFoods(ctx context.Context, restaurantID int32) ([]db.Food, error) {
	return fs.dbStore.GetFoodsByRestaurantId(ctx, restaurantID)
}
