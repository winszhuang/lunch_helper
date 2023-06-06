package service

import (
	"log"
	"lunch_helper/adapter"
	"lunch_helper/cache"
	db "lunch_helper/db/sqlc"
	"lunch_helper/thirdparty"

	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

var defaultSearchRequest = &maps.NearbySearchRequest{
	Type:     maps.PlaceTypeRestaurant,
	Language: "zh-TW",
	OpenNow:  true,
}

type SearchService struct {
	nearByCache       *cache.NearByRestaurantCache
	placeApi          thirdparty.PlaceApi
	crawlerService    *CrawlerService
	restaurantService *RestaurantService
	workChan          chan thirdparty.SearchResult
}

const WORKER_CHAN_SIZE = 10

func NewSearchService(
	nearByCache *cache.NearByRestaurantCache,
	placeApi thirdparty.PlaceApi,
	crawlerService *CrawlerService,
	restaurantService *RestaurantService,
	workerCount int,
) *SearchService {
	service := &SearchService{
		nearByCache:       nearByCache,
		placeApi:          placeApi,
		crawlerService:    crawlerService,
		restaurantService: restaurantService,
		workChan:          make(chan thirdparty.SearchResult, WORKER_CHAN_SIZE),
	}

	for i := 0; i < workerCount; i++ {
		go service.doWork()
	}

	return service
}

func (s *SearchService) doWork() {
	ctx := context.Background()
	apiKey := s.placeApi.GetApiKey()

	for w := range s.workChan {
		// 1. 確認資料庫是否有該餐廳資訊，沒有就註冊
		// 2. 確認餐廳是否有被爬取過餐點資訊，有的話就跳過，沒有的話就爬取
		// 3. 爬取後更新餐點到資料庫
		restaurant, err := s.restaurantService.GetRestaurantByGoogleMapPlaceId(ctx, w.Data.PlaceID)
		if err != nil {
			item := adapter.SearchResultToRestaurant(w, apiKey)
			restaurant, err = s.restaurantService.CreateRestaurant(ctx, db.CreateRestaurantParams{
				Name:             item.Name,
				Rating:           item.Rating,
				UserRatingsTotal: item.UserRatingsTotal,
				Address:          item.Address,
				GoogleMapPlaceID: item.GoogleMapPlaceID,
				GoogleMapUrl:     item.GoogleMapUrl,
				PhoneNumber:      item.PhoneNumber,
				Image:            item.Image,
			})
			if err != nil {
				log.Printf("Create Restaurant %s error: %v", item.Name, err)
			}
		}

		// 沒有爬取過就爬取
		if !restaurant.MenuCrawled {
			// 加入爬蟲代辦清單，並且更新dishes至資料庫
			s.crawlerService.SendWork(restaurant.GoogleMapUrl)

			if err = s.restaurantService.UpdateMenuCrawled(ctx, db.UpdateMenuCrawledParams{
				ID:          restaurant.ID,
				MenuCrawled: true,
			}); err != nil {
				log.Printf("Update MenuCrawled error: %v", err)
			}

		}
	}
}

func (s *SearchService) sendSearchDataToWorker(data []thirdparty.SearchResult) {
	for _, d := range data {
		go func(singleDeliverData thirdparty.SearchResult) {
			s.workChan <- singleDeliverData
		}(d)
	}
}

func (s *SearchService) Search(lat, lng float64, radius, pageIndex, pageSize int) ([]db.Restaurant, error) {
	currentToken := s.nearByCache.GetLastPageToken(cache.LocationArgs{
		Lat:    lat,
		Lng:    lng,
		Radius: radius,
	})
	for {
		list, isEnough := s.nearByCache.GetRestaurantListByPagination(
			cache.LocationArgs{
				Lat:    lat,
				Lng:    lng,
				Radius: radius,
			},
			pageIndex,
			pageSize,
		)
		if isEnough {
			return list, nil
		}

		// 資料不夠，繼續fetch
		log.Println("資料不夠，繼續fetch")
		resp, nextPageToken, err := s.placeApi.NearbySearch(&maps.NearbySearchRequest{
			Location: &maps.LatLng{
				Lat: lat,
				Lng: lng,
			},
			Radius:    uint(radius),
			Type:      defaultSearchRequest.Type,
			Language:  defaultSearchRequest.Language,
			OpenNow:   defaultSearchRequest.OpenNow,
			PageToken: currentToken,
		})
		if err != nil {
			return nil, err
		}

		// 獲取到的店家資訊丟給worker來爬蟲
		s.sendSearchDataToWorker(resp)

		// 加入cache清單
		s.nearByCache.Append(
			cache.LocationArgs{
				Lat:    lat,
				Lng:    lng,
				Radius: radius,
			},
			cache.NewPageDataOfPlaces(currentToken, nextPageToken, adapter.SearchResultsToRestaurants(resp, s.placeApi.GetApiKey())),
		)
		currentToken = nextPageToken
	}
}
