package db

import (
	"context"
	"lunch_helper/util"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func RandomUser(id int32) User {
	return User{
		ID:                     id,
		LineID:                 util.RandomLineID(),
		Name:                   util.RandomName(),
		Picture:                util.RandomPicture(),
		GoogleMapsApiCallCount: 5,
		RoleID:                 2,
	}
}

func RandomUsers(count int) []User {
	users := make([]User, count)
	for i := 1; i <= count; i++ {
		users = append(users, RandomUser(int32(i)))
	}
	return users
}

func GenerateUserRows(users []User) *sqlmock.Rows {
	rows := sqlmock.
		NewRows([]string{"id", "line_id", "name", "picture", "google_maps_api_call_count", "role_id"})
	for _, user := range users {
		rows.AddRow(user.ID, user.LineID, user.Name, user.Picture, user.GoogleMapsApiCallCount, user.RoleID)
	}
	return rows
}

func TestQueries_CreateUser(t *testing.T) {
	user := RandomUser(1)
	expectUserRow := sqlmock.
		NewRows([]string{"id", "line_id", "name", "picture", "google_maps_api_call_count", "role_id"}).
		AddRow(user.ID, user.LineID, user.Name, user.Picture, user.GoogleMapsApiCallCount, user.RoleID)

	testMock.ExpectQuery(regexp.QuoteMeta(createUser)).
		WithArgs(user.LineID, user.Name, user.Picture).
		WillReturnRows(expectUserRow)

	ctx := context.TODO()
	result, err := testQueries.CreateUser(ctx, CreateUserParams{
		LineID:  user.LineID,
		Name:    user.Name,
		Picture: user.Picture,
	})
	assert.NoError(t, err)
	assert.Equal(t, user, result)

	if err := testMock.ExpectationsWereMet(); err != nil {
		t.Errorf("expected result was not achieved: %s", err)
	}
}

func TestQueries_GetUserByID(t *testing.T) {
	user := RandomUser(1)
	expectUserRow := sqlmock.
		NewRows([]string{"id", "line_id", "name", "picture", "google_maps_api_call_count", "role_id"}).
		AddRow(user.ID, user.LineID, user.Name, user.Picture, user.GoogleMapsApiCallCount, user.RoleID)

	testMock.ExpectQuery(regexp.QuoteMeta(getUserByID)).
		WithArgs(1).
		WillReturnRows(expectUserRow)

	ctx := context.TODO()
	result, err := testQueries.GetUserByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, user, result)

	if err := testMock.ExpectationsWereMet(); err != nil {
		t.Errorf("expected result was not achieved: %s", err)
	}
}

func TestQueries_GetUsers(t *testing.T) {
	existedUsers := RandomUsers(50)

	t.Run("Limit 10 Offset 0", func(t *testing.T) {
		currentUsers := existedUsers[:10]
		expectRows := GenerateUserRows(currentUsers)
		testMock.ExpectQuery(regexp.QuoteMeta(getUsers)).
			WithArgs(10, 0).
			WillReturnRows(expectRows)
		ctx := context.TODO()
		result, err := testQueries.GetUsers(ctx, GetUsersParams{
			Limit:  10,
			Offset: 0,
		})
		assert.NoError(t, err)
		assert.Equal(t, currentUsers, result)

		if err := testMock.ExpectationsWereMet(); err != nil {
			t.Errorf("expected result was not achieved: %s", err)
		}
	})

	t.Run("Limit 20 Offset 40", func(t *testing.T) {
		currentUsers := existedUsers[40-1 : 20+40]
		expectRows := GenerateUserRows(currentUsers)
		testMock.ExpectQuery(regexp.QuoteMeta(getUsers)).
			WithArgs(20, 40).
			WillReturnRows(expectRows)
		ctx := context.TODO()
		result, err := testQueries.GetUsers(ctx, GetUsersParams{
			Limit:  20,
			Offset: 40,
		})
		assert.NoError(t, err)
		assert.Equal(t, currentUsers, result)

		if err := testMock.ExpectationsWereMet(); err != nil {
			t.Errorf("expected result was not achieved: %s", err)
		}
	})
}
