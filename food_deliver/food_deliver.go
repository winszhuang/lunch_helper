package food_deliver

import (
	"fmt"
	"lunch_helper/food_deliver/model"
	"time"
)

type DishesCrawler interface {
	ParseSource(string) (string, error)
	GetDishes(string) ([]model.Dish, error)
}

var crawlerList = []DishesCrawler{
	&FoodPandaDishesCrawler{},
	&UberEatsDishesCrawler{},
}

// 限速器
var limiter = time.Tick(200 * time.Millisecond)

type FoodDeliverApi struct{}

func NewFoodDeliverApi() *FoodDeliverApi {
	return &FoodDeliverApi{}
}

func (f *FoodDeliverApi) GetDishesFromGoogleMap(googleMapUrl string) ([]model.Dish, error) {
	for _, crawler := range crawlerList {
		deliverUrl, err := crawler.ParseSource(googleMapUrl)
		if err == nil && deliverUrl != "" {
			// 限速
			<-limiter
			return crawler.GetDishes(deliverUrl)
		}
	}
	return nil, fmt.Errorf("not found dishes for url %s", googleMapUrl)
}
