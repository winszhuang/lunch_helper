package service

import (
	"context"
	"database/sql"
	"log"
	db "lunch_helper/db/sqlc"
	"lunch_helper/spider"
	"lunch_helper/thirdparty"
	"lunch_helper/thirdparty/model"
)

type CrawlerService struct {
	deliverLinkSpider spider.DeliverLinkSpider
	foodDeliverApi    thirdparty.FoodDeliverApi
	googleMapLinkChan chan ChanData
	// 急件處理
	priorityGoogleMapLinkChan chan ChanData
	deliverLinkChan           chan ChanData
	foodService               FoodService
}

type ChanData struct {
	GoogleMapRestaurantUrl   string
	FoodDeliverRestaurantUrl string
	RestaurantId             int32
}

const (
	MAX_COUNT_OF_DELIVER_LINK_CHAN      = 100
	MAX_COUNT_OF_DO_FETCH_DISHES_WORKER = 10
)

func NewCrawlerService(
	deliverLinkSpider spider.DeliverLinkSpider,
	foodDeliverApi thirdparty.FoodDeliverApi,
	foodService FoodService,
) *CrawlerService {
	service := &CrawlerService{
		deliverLinkSpider:         deliverLinkSpider,
		foodDeliverApi:            foodDeliverApi,
		foodService:               foodService,
		googleMapLinkChan:         make(chan ChanData),
		priorityGoogleMapLinkChan: make(chan ChanData),
		deliverLinkChan:           make(chan ChanData, MAX_COUNT_OF_DELIVER_LINK_CHAN),
	}

	for i := 0; i < MAX_COUNT_OF_DO_FETCH_DISHES_WORKER; i++ {
		// 單純fetch資料，使用多個goroutine處理
		go service.fetchDishes()
	}
	// 使用selenium爬蟲，並免負荷太大，只開一個goroutine處理
	go service.doCrawl()

	return service
}

func (s *CrawlerService) SendWork(restaurantData db.Restaurant) {
	go func() {
		s.googleMapLinkChan <- ChanData{
			GoogleMapRestaurantUrl: restaurantData.GoogleMapUrl,
			RestaurantId:           restaurantData.ID,
		}
	}()
}

// 發送優先任務，worker中其他工作會先暫緩，優先處理這個
func (s *CrawlerService) SendPriorityWork(restaurantData db.Restaurant) {
	go func() {
		s.priorityGoogleMapLinkChan <- ChanData{
			GoogleMapRestaurantUrl: restaurantData.GoogleMapUrl,
			RestaurantId:           restaurantData.ID,
		}
	}()
}

// #TODO 可等待任務完成
// func (s *CrawlerService) CheckWorkDone(googleMapRestaurantUrl string) {

// }

func (s *CrawlerService) doCrawl() {
	for {
		select {
		// 急件處理
		case chanData := <-s.priorityGoogleMapLinkChan:
			s.crawlFoodDeliverLinks(chanData)
		default:
			chanData := <-s.googleMapLinkChan
			s.crawlFoodDeliverLinks(chanData)
		}
	}
}

func (s *CrawlerService) crawlFoodDeliverLinks(chanData ChanData) {
	// 從 google map 網站上的店家頁面爬取合作外送平台的url
	foodDeliverLink, err := s.deliverLinkSpider.ScrapeDeliverLink(chanData.GoogleMapRestaurantUrl)
	if err == nil {
		chanData.FoodDeliverRestaurantUrl = foodDeliverLink
		go func(chanData ChanData) {
			s.deliverLinkChan <- chanData
		}(chanData)
	} else {
		log.Printf("url %s crawl error: %v", chanData.GoogleMapRestaurantUrl, err)
	}
}

func (s *CrawlerService) fetchDishes() {
	ctx := context.Background()
	for chanData := range s.deliverLinkChan {
		dishes, err := s.foodDeliverApi.GetDishes(chanData.FoodDeliverRestaurantUrl)
		if err != nil {
			log.Printf("url %s dishes fetch error: %v", chanData.FoodDeliverRestaurantUrl, err)
			continue
		}

		errList := s.saveFoods(ctx, dishes, chanData.RestaurantId)
		if len(errList) > 0 {
			for _, err := range errList {
				log.Printf("url %s dishes save error: %v", chanData.FoodDeliverRestaurantUrl, err)
			}
		}
	}
}

func (s *CrawlerService) saveFoods(ctx context.Context, dishes []model.Dish, restaurantId int32) []error {
	// dishes儲存至資料庫
	errList := []error{}
	for _, dish := range dishes {
		var image sql.NullString
		var description sql.NullString
		if dish.Image == "" {
			image = sql.NullString{String: "", Valid: false}
		} else {
			image = sql.NullString{String: dish.Image, Valid: true}
		}
		if dish.Description == "" {
			description = sql.NullString{String: "", Valid: false}
		} else {
			description = sql.NullString{String: dish.Description, Valid: true}
		}
		if _, err := s.foodService.CreateFood(ctx, db.CreateFoodParams{
			Name:         dish.Name,
			Price:        dish.Price,
			Image:        image,
			Description:  description,
			RestaurantID: restaurantId,
			EditBy:       sql.NullInt32{Int32: 0, Valid: false},
		}); err != nil {
			errList = append(errList, err)
		}
	}
	return errList
}
