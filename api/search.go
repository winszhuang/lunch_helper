package api

import (
	"log"
	"lunch_helper/adapter"
	"lunch_helper/bot/carousel"
	"lunch_helper/bot/quickreply"
	"lunch_helper/cache"
	"lunch_helper/constant"
	db "lunch_helper/db/sqlc"
	"lunch_helper/util"
	"strconv"

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

const defaultToken = ""

func (s *Server) SearchRestaurants(c *gin.Context, event *linebot.Event) {
	userId := event.Source.UserID
	radius := s.messageCache.GetCurrentRadius(userId)
	uc, ok := s.messageCache.GetCurrentLocation(userId)
	if !ok {
		s.bot.SendTextWithQuickReplies(event.ReplyToken, "請先傳送位置資訊再做搜尋哦 ~", quickreply.QuickReplyLocation())
		return
	}

	currentToken := defaultToken
	for {
		list, isEnough := s.nearByCache.GetRestaurantListByPagination(
			cache.LocationArgs{
				Lat:    uc.LatLng.Lat,
				Lng:    uc.LatLng.Lng,
				Radius: radius,
			},
			DefaultPageIndex,
			MaximumNumberOfCarouselItems,
		)
		if isEnough {
			component := carousel.CreateCarouselWithNext(
				list,
				func(restaurant db.Restaurant) *linebot.BubbleContainer {
					return carousel.CreateRestaurantContainer(restaurant)
				},
				DefaultPageIndex+1,
				uc.LatLng.Lat,
				uc.LatLng.Lng,
				radius,
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
			PageToken: currentToken,
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
			cache.NewPageDataOfPlaces(currentToken, pageToken, adapter.SearchResultToRestaurant(resp, s.placeApi.GetApiKey())),
		)
		currentToken = pageToken
	}
}

func (s *Server) SearchNextPageRestaurants(c *gin.Context, event *linebot.Event) {
	args := util.ParseRegexQuery(event.Postback.Data, constant.LatLngPageIndex)
	lat, err := strconv.ParseFloat(args[0], 64)
	lng, err := strconv.ParseFloat(args[1], 64)
	radius, err := strconv.Atoi(args[2])
	pageIndex, err := strconv.Atoi(args[3])
	if err != nil {
		s.bot.SendText(event.ReplyToken, "轉換錯誤")
		return
	}
	if len(args) != 4 {
		s.bot.SendText(event.ReplyToken, "下一頁參數錯誤!!")
		return
	}

	currentToken := defaultToken
	for {
		list, isEnough := s.nearByCache.GetRestaurantListByPagination(
			cache.LocationArgs{
				Lat:    lat,
				Lng:    lng,
				Radius: radius,
			},
			pageIndex,
			MaximumNumberOfCarouselItems,
		)
		if isEnough {
			component := carousel.CreateCarouselWithNext(
				list,
				func(restaurant db.Restaurant) *linebot.BubbleContainer {
					return carousel.CreateRestaurantContainer(restaurant)
				},
				pageIndex+1,
				lat,
				lng,
				radius,
			)
			s.bot.SendFlex(event.ReplyToken, "carousel", component)
			return
		}

		// 資料不夠，繼續fetch
		log.Println("資料不夠，繼續fetch")
		currentToken = s.nearByCache.GetLastPageToken(cache.LocationArgs{
			Lat:    lat,
			Lng:    lng,
			Radius: radius,
		})
		resp, pageToken, err := s.placeApi.NearbySearch(&maps.NearbySearchRequest{
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
			log.Printf("NearbySearch api error: %v", err)
			s.bot.SendText(event.ReplyToken, "取得附近店家資訊失敗!!")
			return
		}
		s.nearByCache.Append(
			cache.LocationArgs{
				Lat:    lat,
				Lng:    lng,
				Radius: radius,
			},
			cache.NewPageDataOfPlaces(currentToken, pageToken, adapter.SearchResultToRestaurant(resp, s.placeApi.GetApiKey())),
		)
		currentToken = pageToken
	}
}
