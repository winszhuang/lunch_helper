package service

import (
	"log"
	"lunch_helper/adapter"
	"lunch_helper/cache"
	db "lunch_helper/db/sqlc"
	"lunch_helper/thirdparty"

	"googlemaps.github.io/maps"
)

var defaultSearchRequest = &maps.NearbySearchRequest{
	Type:     maps.PlaceTypeRestaurant,
	Language: "zh-TW",
	OpenNow:  true,
}

type SearchService struct {
	nearByCache *cache.NearByRestaurantCache
	placeApi    thirdparty.PlaceApi
}

const WORKER_CHAN_SIZE = 10

func NewSearchService(
	nearByCache *cache.NearByRestaurantCache,
	placeApi thirdparty.PlaceApi,
) *SearchService {
	return &SearchService{
		nearByCache: nearByCache,
		placeApi:    placeApi,
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
