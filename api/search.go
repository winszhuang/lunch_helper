package api

import (
	"log"
	"lunch_helper/adapter"
	"lunch_helper/bot/carousel"
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

func (s *Server) SearchRestaurantByLocation(c *gin.Context, event *linebot.Event) {
	messageData, ok := event.Message.(*linebot.LocationMessage)
	if !ok {
		s.bot.SendText(event.ReplyToken, "定位資訊有問題!!請從新發送")
		return
	}

	for {
		list, noMiss := s.nearByCache.GetRestaurantListByPagination(
			cache.LocationArgs{
				Lat:    messageData.Latitude,
				Lng:    messageData.Longitude,
				Radius: int(defaultSearchRequest.Radius),
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
				messageData.Latitude,
				messageData.Longitude,
			)
			s.bot.SendFlex(event.ReplyToken, "carousel", component)
			return
		}

		// 資料不夠，繼續fetch
		resp, pageToken, err := s.placeApi.NearbySearch(&maps.NearbySearchRequest{
			Location: &maps.LatLng{
				Lat: messageData.Latitude,
				Lng: messageData.Longitude,
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
				Lat:    messageData.Latitude,
				Lng:    messageData.Longitude,
				Radius: int(defaultSearchRequest.Radius),
			},
			adapter.SearchResultToRestaurant(resp),
			pageToken,
		)
	}
}
