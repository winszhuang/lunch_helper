package service

import (
	"context"
	db "lunch_helper/db/sqlc"
)

type UserService struct {
	dbStore db.Store
}

func NewUserService(dbStore db.Store) *UserService {
	return &UserService{
		dbStore: dbStore,
	}
}

func (us *UserService) GetUserByLineID(ctx context.Context, lineId string) (db.User, error) {
	return us.dbStore.GetUserByLineID(ctx, lineId)
}

func (us *UserService) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	return us.dbStore.CreateUser(ctx, arg)
}
