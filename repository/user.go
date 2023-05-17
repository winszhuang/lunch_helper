package repository

import (
	"context"
	db "lunch_helper/db/sqlc"
)

type userRepository struct {
	dbStore *db.SQLStore
}

func NewUserRepository(dbStore *db.SQLStore) *userRepository {
	return &userRepository{
		dbStore: dbStore,
	}
}

func (ur *userRepository) Create(c context.Context, params db.CreateUserParams) (db.User, error) {
	return ur.dbStore.CreateUser(c, params)
}

func (ur *userRepository) GetByID(c context.Context, id int32) (db.User, error) {
	return ur.dbStore.GetUserByID(c, id)
}

func (ur *userRepository) GetByLineID(c context.Context, lineID string) (db.User, error) {
	return ur.dbStore.GetUserByLineID(c, lineID)
}

func (ur *userRepository) List(c context.Context, params db.GetUsersParams) ([]db.User, error) {
	return ur.dbStore.GetUsers(c, params)
}
