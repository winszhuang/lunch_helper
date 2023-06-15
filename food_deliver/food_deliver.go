package food_deliver

import (
	"fmt"
	"lunch_helper/food_deliver/model"
)

type DishesCrawler interface {
	ParseSource(string) (string, error)
	GetDishes(string) ([]model.Dish, error)
}

var crawlerList = []DishesCrawler{
	&FoodPandaDishesCrawler{},
	&UberEatsDishesCrawler{},
}

type FoodDeliverApi struct{}

func NewFoodDeliverApi() *FoodDeliverApi {
	return &FoodDeliverApi{}
}

func (f *FoodDeliverApi) GetDishesFromGoogleMap(googleMapUrl string) ([]model.Dish, error) {
	for _, crawler := range crawlerList {
		deliverUrl, err := crawler.ParseSource(googleMapUrl)
		if err == nil && deliverUrl != "" {
			return crawler.GetDishes(deliverUrl)
		}
	}
	return nil, fmt.Errorf("not found dishes for url %s", googleMapUrl)
}
