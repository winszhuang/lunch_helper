package service

import (
	"context"
	db "lunch_helper/db/sqlc"
)

type RestaurantService struct {
	dbStore db.Store
}

func NewRestaurantService(dbStore db.Store) *RestaurantService {
	return &RestaurantService{
		dbStore: dbStore,
	}
}

func (rs *RestaurantService) CreateRestaurant(ctx context.Context, arg db.CreateRestaurantParams) (db.Restaurant, error) {
	r, err := rs.GetRestaurantByGoogleMapPlaceId(ctx, arg.GoogleMapPlaceID)
	if err == nil {
		return r, nil
	}
	return rs.dbStore.CreateRestaurant(ctx, arg)
}

func (rs *RestaurantService) GetRestaurant(ctx context.Context, id int32) (db.Restaurant, error) {
	return rs.dbStore.GetRestaurant(ctx, id)
}

func (rs *RestaurantService) GetRestaurantByGoogleMapPlaceId(ctx context.Context, googleMapPlaceID string) (db.Restaurant, error) {
	return rs.dbStore.GetRestaurantByGoogleMapPlaceId(ctx, googleMapPlaceID)
}

func (rs *RestaurantService) UpdateMenuCrawled(ctx context.Context, arg db.UpdateMenuCrawledParams) error {
	return rs.dbStore.UpdateMenuCrawled(ctx, arg)
}
