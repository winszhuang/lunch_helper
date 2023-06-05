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
	googleMapLinkChan chan string
	deliverLinkChan   chan string
	foodService       FoodService
}

type CrawFood struct {
	Name        string
	Price       string
	Image       string
	Description string
	Category    string
}

type CrawlRequest struct {
	name string
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
		deliverLinkSpider: deliverLinkSpider,
		foodDeliverApi:    foodDeliverApi,
		foodService:       foodService,
		googleMapLinkChan: make(chan string),
		deliverLinkChan:   make(chan string, MAX_COUNT_OF_DELIVER_LINK_CHAN),
	}

	for i := 0; i < MAX_COUNT_OF_DO_FETCH_DISHES_WORKER; i++ {
		// 單純fetch資料，使用多個goroutine處理
		go service.fetchDishes()
	}
	// 使用selenium爬蟲，並免負荷太大，只開一個goroutine處理
	go service.crawlFoodDeliverLinks()

	return service
}

func (s *CrawlerService) SendWork(googleMapRestaurantUrl string) {
	go func() {
		s.googleMapLinkChan <- googleMapRestaurantUrl
	}()
}

func (s *CrawlerService) crawlFoodDeliverLinks() {
	for url := range s.googleMapLinkChan {
		// 從 google map 網站上的店家頁面爬取合作外送平台的url
		foodDeliverLink, err := s.deliverLinkSpider.ScrapeDeliverLink(url)
		if err == nil {
			go func(url string) {
				s.deliverLinkChan <- url
			}(foodDeliverLink)
		} else {
			log.Printf("url %s crawl error: %v", url, err)
		}
	}
}

func (s *CrawlerService) fetchDishes() {
	ctx := context.Background()
	for deliverLink := range s.deliverLinkChan {
		dishes, err := s.foodDeliverApi.GetDishes(deliverLink)
		if err != nil {
			log.Printf("url %s dishes fetch error: %v", deliverLink, err)
			continue
		}

		errList := s.saveFoods(ctx, dishes)
		if len(errList) > 0 {
			for _, err := range errList {
				log.Printf("url %s dishes save error: %v", deliverLink, err)
			}
		}
	}
}

func (s *CrawlerService) saveFoods(ctx context.Context, dishes []model.Dish) []error {
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
			RestaurantID: 1,
			EditBy:       sql.NullInt32{Int32: 0, Valid: false},
		}); err != nil {
			errList = append(errList, err)
		}
	}
	return errList
}
