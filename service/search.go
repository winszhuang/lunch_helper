package service

import (
	"log"
	"lunch_helper/adapter"
	"lunch_helper/cache"
	"lunch_helper/constant"
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

func (s *SearchService) Search(searchArgs *constant.SearchArgs) ([]db.Restaurant, thirdparty.SearchError) {
	detailErrorList := []error{}
	isFetchLimited := false
	currentToken := s.nearByCache.GetLastPageToken(cache.LocationArgs{
		Lat:    searchArgs.Lat,
		Lng:    searchArgs.Lng,
		Radius: searchArgs.Radius,
		Text:   searchArgs.Text,
	})
	for {
		list, isEnough := s.nearByCache.GetRestaurantListByPagination(
			cache.LocationArgs{
				Lat:    searchArgs.Lat,
				Lng:    searchArgs.Lng,
				Radius: searchArgs.Radius,
				Text:   searchArgs.Text,
			},
			searchArgs.PageIndex,
			searchArgs.PageSize,
		)
		if isEnough || isFetchLimited {
			return list, thirdparty.SearchError{Err: nil, DetailErrors: detailErrorList}
		}

		log.Printf("第%d頁不夠繼續fetch!!", searchArgs.PageIndex)
		resp, nextPageToken, searchErr := s.nearBySearchOrTextSearch(searchArgs, currentToken)
		if searchErr.Err != nil {
			return []db.Restaurant{}, searchErr
		}

		// 為空表示後面已經沒有資料了
		if nextPageToken == "" {
			isFetchLimited = true
		}

		detailErrorList = append(detailErrorList, searchErr.DetailErrors...)

		// 加入cache清單
		s.nearByCache.Append(
			cache.LocationArgs{
				Lat:    searchArgs.Lat,
				Lng:    searchArgs.Lng,
				Radius: searchArgs.Radius,
				Text:   searchArgs.Text,
			},
			cache.NewPageDataOfPlaces(currentToken, nextPageToken, adapter.SearchResultsToRestaurants(resp, s.placeApi.GetApiKey())),
		)
		currentToken = nextPageToken
	}
}

func (s *SearchService) nearBySearchOrTextSearch(searchArgs *constant.SearchArgs, currentToken string) ([]thirdparty.SearchResult, string, thirdparty.SearchError) {
	if searchArgs.Text != "" {
		return s.placeApi.TextSearch(&maps.TextSearchRequest{
			Location: &maps.LatLng{
				Lat: searchArgs.Lat,
				Lng: searchArgs.Lng,
			},
			Radius:    uint(searchArgs.Radius),
			Type:      defaultSearchRequest.Type,
			Language:  defaultSearchRequest.Language,
			OpenNow:   defaultSearchRequest.OpenNow,
			PageToken: currentToken,
			Query:     searchArgs.Text,
		})
	}
	return s.placeApi.NearbySearch(&maps.NearbySearchRequest{
		Location: &maps.LatLng{
			Lat: searchArgs.Lat,
			Lng: searchArgs.Lng,
		},
		Radius:    uint(searchArgs.Radius),
		Type:      defaultSearchRequest.Type,
		Language:  defaultSearchRequest.Language,
		OpenNow:   defaultSearchRequest.OpenNow,
		PageToken: currentToken,
	})
}
