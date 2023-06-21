package food_deliver

import (
	"fmt"
	"lunch_helper/food_deliver/model"
)

type DishesCrawler interface {
	ParseSource(string) (string, error)
	GetDishes(string) ([]model.Dish, error)
	GetDeliverName() FoodDeliverName
}

var crawlerList = []DishesCrawler{
	NewFoodPandaDishesCrawler(),
	NewUberEatsDishesCrawler(),
}

type FoodDeliverApi struct {
}

type FetchInfo struct {
	DeliverName FoodDeliverName
	FetchLink   string
}

func NewFoodDeliverApi() *FoodDeliverApi {
	return &FoodDeliverApi{}
}

func (f *FoodDeliverApi) getCrawler(deliverName FoodDeliverName) (DishesCrawler, error) {
	for _, crawler := range crawlerList {
		if crawler.GetDeliverName() == deliverName {
			return crawler, nil
		}
	}
	return nil, fmt.Errorf("not found crawler %s", deliverName)
}

func (f *FoodDeliverApi) CheckFoodDeliverFromGoogleMap(googleMapUrl string) (*FetchInfo, error) {
	for _, crawler := range crawlerList {
		deliverUrl, err := crawler.ParseSource(googleMapUrl)
		if err == nil && deliverUrl != "" {
			return &FetchInfo{
				DeliverName: crawler.GetDeliverName(),
				FetchLink:   deliverUrl,
			}, nil
		}
	}
	return nil, fmt.Errorf("not found deliver for url %s", googleMapUrl)
}

func (f *FoodDeliverApi) GetDishes(fetchInfo *FetchInfo) ([]model.Dish, error) {
	link := fetchInfo.FetchLink
	crawler, err := f.getCrawler(fetchInfo.DeliverName)
	if err != nil {
		return nil, err
	}

	return crawler.GetDishes(link)
}
