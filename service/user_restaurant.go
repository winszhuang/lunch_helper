package service

import (
	"context"
	db "lunch_helper/db/sqlc"
)

type UserRestaurantService struct {
	dbStore db.Store
}

func NewUserRestaurantService(dbStore db.Store) *UserRestaurantService {
	return &UserRestaurantService{
		dbStore: dbStore,
	}
}

func (uf *UserRestaurantService) Create(ctx context.Context, userId int32, restaurantId int32) (db.UserRestaurant, error) {
	return uf.dbStore.CreateUserRestaurant(ctx, db.CreateUserRestaurantParams{
		UserID:       userId,
		RestaurantID: restaurantId,
	})
}
