package adapter

import db "lunch_helper/db/sqlc"

func UserRestaurantRowsToRestaurants(urs []db.GetUserRestaurantsRow) []db.Restaurant {
	result := []db.Restaurant{}
	for _, ur := range urs {
		result = append(result, UserRestaurantRowToRestaurant(ur))
	}
	return result
}

func UserRestaurantRowToRestaurant(ur db.GetUserRestaurantsRow) db.Restaurant {
	return db.Restaurant{
		ID:               ur.ID,
		Name:             ur.Name,
		Rating:           ur.Rating,
		UserRatingsTotal: ur.UserRatingsTotal,
		Address:          ur.Address,
		GoogleMapPlaceID: ur.GoogleMapPlaceID,
		GoogleMapUrl:     ur.GoogleMapUrl,
		PhoneNumber:      ur.PhoneNumber,
		Image:            ur.Image,
		MenuCrawled:      ur.MenuCrawled,
	}
}
