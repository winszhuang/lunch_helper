package service

import (
	"context"
	db "lunch_helper/db/sqlc"
)

type UserFoodService struct {
	dbStore db.Store
}

func NewUserFoodService(dbStore db.Store) *UserFoodService {
	return &UserFoodService{
		dbStore: dbStore,
	}
}

func (uf *UserFoodService) Create(ctx context.Context, userId int32, foodId int32) (db.UserFood, error) {
	return uf.dbStore.CreateUserFood(ctx, db.CreateUserFoodParams{
		UserID: userId,
		FoodID: foodId,
	})
}

func (uf *UserFoodService) List(ctx context.Context, args db.GetUserFoodsParams) ([]db.GetUserFoodsRow, error) {
	return uf.dbStore.GetUserFoods(ctx, args)
}
