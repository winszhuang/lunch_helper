package food_deliver

import (
	"fmt"
	"log"
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
		log.Printf("抓取%s的dishes。deliverUrl是: %s", googleMapUrl, deliverUrl)
		if err == nil && deliverUrl != "" {
			log.Printf("%s可以進來呼叫crawler.GetDishes搂", googleMapUrl)
			return crawler.GetDishes(deliverUrl)
		}
		if err != nil {
			log.Printf("parse source error for url %s: %v", googleMapUrl, err)
		}
	}
	return nil, fmt.Errorf("not found dishes for url %s", googleMapUrl)
}
