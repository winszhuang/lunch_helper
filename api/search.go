package api

import (
	"log"
	"lunch_helper/adapter"
	"lunch_helper/bot/carousel"
	"lunch_helper/bot/quickreply"
	"lunch_helper/cache"
	db "lunch_helper/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"googlemaps.github.io/maps"
)

const (
	DefaultPageIndex             = 1
	MaximumNumberOfCarouselItems = 10
)

var defaultSearchRequest = &maps.NearbySearchRequest{
	Radius:   500,
	Type:     maps.PlaceTypeRestaurant,
	Language: "zh-TW",
	OpenNow:  true,
}

// steps:
// 1. 查看
func (s *Server) SearchRestaurants(c *gin.Context, event *linebot.Event) {
	userId := event.Source.UserID
	radius := s.messageCache.GetCurrentRadius(userId)
	uc, ok := s.messageCache.GetCurrentLocation(userId)
	if !ok {
		s.bot.SendTextWithQuickReplies(event.ReplyToken, "請先傳送位置資訊再做搜尋哦 ~", quickreply.QuickReplyLocation())
		return
	}

	for {
		list, noMiss := s.nearByCache.GetRestaurantListByPagination(
			cache.LocationArgs{
				Lat:    uc.LatLng.Lat,
				Lng:    uc.LatLng.Lng,
				Radius: radius,
			},
			DefaultPageIndex,
			MaximumNumberOfCarouselItems,
		)
		if noMiss {
			component := carousel.CreateCarouselWithNext(
				list,
				func(restaurant db.Restaurant) *linebot.BubbleContainer {
					return carousel.CreateRestaurantContainer(restaurant)
				},
				DefaultPageIndex+1,
				uc.LatLng.Lat,
				uc.LatLng.Lng,
			)
			s.bot.SendFlex(event.ReplyToken, "carousel", component)
			return
		}

		// 資料不夠，繼續fetch
		log.Println("資料不夠，繼續fetch")
		resp, pageToken, err := s.placeApi.NearbySearch(&maps.NearbySearchRequest{
			Location: &maps.LatLng{
				Lat: uc.LatLng.Lat,
				Lng: uc.LatLng.Lng,
			},
			Radius:    defaultSearchRequest.Radius,
			Type:      defaultSearchRequest.Type,
			Language:  defaultSearchRequest.Language,
			OpenNow:   defaultSearchRequest.OpenNow,
			PageToken: "",
		})
		if err != nil {
			log.Printf("NearbySearch api error: %v", err)
			s.bot.SendText(event.ReplyToken, "取得附近店家資訊失敗!!")
			return
		}
		s.nearByCache.Append(
			cache.LocationArgs{
				Lat:    uc.LatLng.Lat,
				Lng:    uc.LatLng.Lng,
				Radius: int(defaultSearchRequest.Radius),
			},
			cache.NewPageDataOfPlaces("", pageToken, adapter.SearchResultToRestaurant(resp, s.placeApi.GetApiKey())),
		)
	}
}
