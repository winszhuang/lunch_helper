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

// url example :
// - https://www.ubereats.com/tw/store/%E5%BA%B7%E5%AF%B6%E8%97%A5%E7%87%89%E6%8E%92%E9%AA%A8-%E7%BE%8E%E5%BE%B7%E7%B8%BD%E5%BA%97/-NQAJwmsRJGIo-7WUwHfHQ?utm_campaign=place-action-link&utm_medium=organic&utm_source=google
// - https://www.foodpanda.com.tw/restaurant/zb5n/mu-he-tang-tai-zhong-bei-ping-dian-ri-shi-jing-gai-fan-wu-long-mian-la-mian-zhuan-mai-dian?utm_source=google&utm_medium=organic&utm_campaign=google_reserve_place_order_action
func (c *CommonFoodDeliverApi) GetDishes(url string) ([]model.Dish, error) {
	d, err := checkDeliver(url)
	if err != nil {
		return nil, err
	}

	switch d {
	case UberEats:
		return thirdparty.GetDishesByUberEats(url)
	case FoodPanda:
		return thirdparty.GetDishesByFoodPanda(url)
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
