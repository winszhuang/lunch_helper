package thirdparty

import (
	"errors"
	thirdparty "lunch_helper/thirdparty/delivers"
	"lunch_helper/thirdparty/model"
	"strings"
)

type DELIVER_NAME string

const (
	UberEats  DELIVER_NAME = "ubereats"
	FoodPanda DELIVER_NAME = "foodpanda"
)

type FoodDeliverApi interface {
	GetDishes(url string) ([]model.Dish, error)
}

type CommonFoodDeliverApi struct{}

func NewCommonFoodDeliverApi() *CommonFoodDeliverApi {
	return &CommonFoodDeliverApi{}
}

func (c *CommonFoodDeliverApi) GetDishes(url string) ([]model.Dish, error) {
	d, err := checkDeliver(url)
	if err != nil {
		return nil, err
	}

	switch d {
	case UberEats:
		return thirdparty.GetDishesByUberEats(string(d))
	case FoodPanda:
		return thirdparty.GetDishesByFoodPanda(string(d))
	default:
		return nil, errors.New("not already implement this get dishes")
	}
}

func checkDeliver(url string) (DELIVER_NAME, error) {
	if strings.Contains(url, "ubereats.com") {
		return UberEats, nil
	}
	if strings.Contains(url, "foodpanda.com") {
		return FoodPanda, nil
	}
	return "", errors.New("no supported deliver")
}
