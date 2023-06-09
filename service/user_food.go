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

func (uf *UserFoodService) GetByFoodId(ctx context.Context, arg db.GetUserFoodByFoodIdParams) (db.GetUserFoodByFoodIdRow, error) {
	return uf.dbStore.GetUserFoodByFoodId(ctx, arg)
}

func (uf *UserFoodService) List(ctx context.Context, args db.GetUserFoodsParams) ([]db.GetUserFoodsRow, error) {
	return uf.dbStore.GetUserFoods(ctx, args)
}

func (uf *UserFoodService) ListAll(ctx context.Context, userID int32) ([]db.GetAllUserFoodsRow, error) {
	return uf.dbStore.GetAllUserFoods(ctx, userID)
}

func (uf *UserFoodService) Delete(ctx context.Context, args db.DeleteUserFoodParams) error {
	return uf.dbStore.DeleteUserFood(ctx, args)
}
