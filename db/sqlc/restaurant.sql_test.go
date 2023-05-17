package db

import (
	"context"
	"database/sql"
	"lunch_helper/util"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func RandomRestaurant(id int32) Restaurant {
	return Restaurant{
		ID:               id,
		Name:             util.RandomName(),
		Rating:           util.RandomRating(),
		UserRatingsTotal: sql.NullInt32{Int32: util.RandomInt32(1, 1000), Valid: true},
		Address:          util.RandomChar(30),
		GoogleMapPlaceID: util.RandomChar(30),
		GoogleMapUrl:     util.RandomChar(30),
		PhoneNumber:      util.RandomPhoneNumber(),
	}
}

func TestQueries_CreateRestaurant(t *testing.T) {
	restaurant := RandomRestaurant(1)
	expectUserRow := sqlmock.
		NewRows([]string{"id", "name", "rating", "user_ratings_total", "address", "google_map_place_id", "google_map_url", "phone_number"}).
		AddRow(restaurant.ID, restaurant.Name, restaurant.Rating, restaurant.UserRatingsTotal, restaurant.Address, restaurant.GoogleMapPlaceID, restaurant.GoogleMapUrl, restaurant.PhoneNumber)

	testMock.ExpectQuery(regexp.QuoteMeta(createRestaurant)).
		WithArgs(restaurant.Name, restaurant.Rating, restaurant.UserRatingsTotal, restaurant.Address, restaurant.GoogleMapPlaceID, restaurant.GoogleMapUrl, restaurant.PhoneNumber).
		WillReturnRows(expectUserRow)

	ctx := context.TODO()
	result, err := testQueries.CreateRestaurant(ctx, CreateRestaurantParams{
		Name:             restaurant.Name,
		Rating:           restaurant.Rating,
		UserRatingsTotal: restaurant.UserRatingsTotal,
		Address:          restaurant.Address,
		GoogleMapPlaceID: restaurant.GoogleMapPlaceID,
		GoogleMapUrl:     restaurant.GoogleMapUrl,
		PhoneNumber:      restaurant.PhoneNumber,
	})
	assert.NoError(t, err)
	assert.Equal(t, restaurant, result)

	if err := testMock.ExpectationsWereMet(); err != nil {
		t.Errorf("expected result was not achieved: %s", err)
	}
}
