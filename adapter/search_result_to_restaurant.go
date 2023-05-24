package adapter

import (
	"database/sql"
	db "lunch_helper/db/sqlc"
	"lunch_helper/thirdparty"
	"lunch_helper/util"

	"github.com/shopspring/decimal"
)

func SearchResultToRestaurant(result []thirdparty.SearchResult, apiKey string) []db.Restaurant {
	var restaurants []db.Restaurant
	for _, r := range result {
		var photo sql.NullString
		if len(r.Data.Photos) > 0 {
			reference := r.Data.Photos[0].PhotoReference
			photo = sql.NullString{String: util.GetGoogleImageUrl(reference, apiKey), Valid: true}
		} else {
			photo = sql.NullString{String: "", Valid: false}
		}
		restaurants = append(restaurants, db.Restaurant{
			// id in this place is not important
			ID:               0,
			Name:             r.Data.Name,
			Rating:           decimal.NewFromFloat32(r.Data.Rating),
			UserRatingsTotal: sql.NullInt32{Int32: int32(r.Data.UserRatingsTotal), Valid: true},
			Address:          r.Data.Vicinity,
			GoogleMapPlaceID: r.Data.PlaceID,
			GoogleMapUrl:     r.Detail.URL,
			PhoneNumber:      r.Detail.FormattedPhoneNumber,
			Image:            photo,
		})
	}
	return restaurants
}
