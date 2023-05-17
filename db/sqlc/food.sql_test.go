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

func RandomFood(id int32) Food {
	return Food{
		ID:           id,
		Name:         util.RandomName(),
		Price:        string(rune(util.RandomInt(80, 400))),
		Image:        sql.NullString{String: util.RandomPicture(), Valid: true},
		Description:  sql.NullString{String: util.RandomString(10), Valid: true},
		RestaurantID: util.RandomInt32(1, 100),
		EditBy:       sql.NullInt32{Int32: util.RandomInt32(1, 100), Valid: true},
	}
}

func RandomFoods(count int) []Food {
	foods := make([]Food, count)
	for i := 1; i <= count; i++ {
		foods = append(foods, RandomFood(int32(i)))
	}
	return foods
}

func GenerateFoodRows(foods []Food) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id", "name", "price", "image", "description", "restaurant_id", "version", "edit_by"})
	for _, food := range foods {
		rows.AddRow(food.ID, food.Name, food.Price, food.Image, food.Description, food.RestaurantID, food.Version, food.EditBy)
	}
	return rows
}

func TestQueries_CreateFood(t *testing.T) {
	food := RandomFood(1)
	expectFoodRow := sqlmock.
		NewRows([]string{"id", "name", "price", "image", "description", "restaurant_id", "version", "edit_by"}).
		AddRow(food.ID, food.Name, food.Price, food.Image, food.Description, food.RestaurantID, food.Version, food.EditBy)

	testMock.ExpectQuery(regexp.QuoteMeta(createFood)).
		WithArgs(food.Name, food.Price, food.Image, food.Description, food.RestaurantID, food.EditBy).
		WillReturnRows(expectFoodRow)

	ctx := context.TODO()
	result, err := testQueries.CreateFood(ctx, CreateFoodParams{
		Name:         food.Name,
		Price:        food.Price,
		Image:        sql.NullString{String: food.Image.String, Valid: true},
		Description:  sql.NullString{String: food.Description.String, Valid: true},
		RestaurantID: food.RestaurantID,
		EditBy:       sql.NullInt32{Int32: food.EditBy.Int32, Valid: true},
	})
	assert.NoError(t, err)
	assert.Equal(t, food, result)

	if err := testMock.ExpectationsWereMet(); err != nil {
		t.Errorf("expected result was not achieved: %s", err)
	}
}

func TestQueries_DeleteFood(t *testing.T) {
	testMock.ExpectExec(regexp.QuoteMeta(deleteFood)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.TODO()
	err := testQueries.DeleteFood(ctx, 1)
	assert.NoError(t, err)

	if err := testMock.ExpectationsWereMet(); err != nil {
		t.Errorf("expected result was not achieved: %s", err)
	}
}

func TestQueries_UpdateFood(t *testing.T) {
	testMock.ExpectExec(regexp.QuoteMeta(updateFood)).
		WithArgs("New Name", "New Price", sql.NullString{String: "new_image.jpg", Valid: true}, sql.NullInt32{Int32: 2, Valid: true}, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	ctx := context.TODO()
	err := testQueries.UpdateFood(ctx, UpdateFoodParams{
		Name:   "New Name",
		Price:  "New Price",
		Image:  sql.NullString{String: "new_image.jpg", Valid: true},
		EditBy: sql.NullInt32{Int32: 2, Valid: true},
		ID:     1,
	})
	assert.NoError(t, err)

	if err := testMock.ExpectationsWereMet(); err != nil {
		t.Errorf("expected result was not achieved: %s", err)
	}
}
