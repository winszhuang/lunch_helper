package adapter

import (
	"database/sql"
	db "lunch_helper/db/sqlc"
	"lunch_helper/thirdparty"
	"lunch_helper/util"

	"github.com/shopspring/decimal"
)

func SearchResultToRestaurant(r thirdparty.SearchResult, apiKey string) db.Restaurant {
	var photo sql.NullString
	if len(r.Data.Photos) > 0 {
		reference := r.Data.Photos[0].PhotoReference
		photo = sql.NullString{String: util.GetGoogleImageUrl(reference, apiKey), Valid: true}
	} else {
		photo = sql.NullString{String: "", Valid: false}
	}

	return db.Restaurant{
		// id in this place is not important
		ID:               0,
		Name:             r.Data.Name,
		Rating:           decimal.NewFromFloat32(r.Data.Rating),
		UserRatingsTotal: sql.NullInt32{Int32: int32(r.Data.UserRatingsTotal), Valid: true},
		Address:          r.Detail.FormattedAddress,
		GoogleMapPlaceID: r.Data.PlaceID,
		GoogleMapUrl:     r.Detail.URL,
		PhoneNumber:      r.Detail.FormattedPhoneNumber,
		Image:            photo,
	}
}

func SearchResultsToRestaurants(result []thirdparty.SearchResult, apiKey string) []db.Restaurant {
	var restaurants []db.Restaurant

	for _, r := range result {
		restaurants = append(restaurants, SearchResultToRestaurant(r, apiKey))
	}
	return restaurants
}
