package adapter

import db "lunch_helper/db/sqlc"

func UserFoodRowsToFoods(urs []db.GetUserFoodsRow) []db.Food {
	result := []db.Food{}
	for _, ur := range urs {
		result = append(result, UserFoodRowToFood(ur))
	}
	return result
}

func UserFoodRowToFood(ur db.GetUserFoodsRow) db.Food {
	return db.Food{
		ID:           ur.ID,
		Name:         ur.Name,
		Price:        ur.Price,
		Image:        ur.Image,
		Description:  ur.Description,
		RestaurantID: ur.RestaurantID,
		Version:      ur.Version,
		EditBy:       ur.EditBy,
	}
}
