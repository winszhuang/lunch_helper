package api

import (
	"lunch_helper/bot/carousel"
	"lunch_helper/bot/quickreply"
	"lunch_helper/constant"
	db "lunch_helper/db/sqlc"
	"lunch_helper/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

const (
	DefaultPageIndex             = 1
	MaximumNumberOfCarouselItems = 10
)

type SearchArgs struct {
	lat       float64
	lng       float64
	radius    int
	pageIndex int
}

func (s *Server) HandleSearchFirstPageRestaurants(c *gin.Context, event *linebot.Event) {
	userId := event.Source.UserID

	radius := s.messageCache.GetCurrentRadius(userId)
	uc, ok := s.messageCache.GetCurrentLocation(userId)
	if !ok {
		s.bot.SendTextWithQuickReplies(event.ReplyToken, "請先傳送位置資訊再做搜尋哦 ~", quickreply.QuickReplyLocation())
		return
	}

	searchArgs := &SearchArgs{
		lat:       uc.LatLng.Lat,
		lng:       uc.LatLng.Lng,
		pageIndex: DefaultPageIndex,
		radius:    radius,
	}

	s.searchSaveAndSend(c, event, searchArgs)
}

func (s *Server) HandleSearchNextPageRestaurants(c *gin.Context, event *linebot.Event) {
	args := util.ParseRegexQuery(event.Postback.Data, constant.LatLngPageIndex)
	if len(args) != 4 {
		s.bot.SendText(event.ReplyToken, "下一頁參數錯誤!!")
		return
	}

	lat, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		s.bot.SendText(event.ReplyToken, "解析Lat失敗")
		return
	}
	lng, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		s.bot.SendText(event.ReplyToken, "解析Lng失敗")
		return
	}
	radius, err := strconv.Atoi(args[2])
	if err != nil {
		s.bot.SendText(event.ReplyToken, "解析半徑失敗")
		return
	}
	pageIndex, err := strconv.Atoi(args[3])
	if err != nil {
		s.bot.SendText(event.ReplyToken, "解析頁數失敗")
		return
	}

	searchArgs := &SearchArgs{
		lat:       lat,
		lng:       lng,
		pageIndex: pageIndex,
		radius:    radius,
	}

	s.searchSaveAndSend(c, event, searchArgs)
}

func (s *Server) searchSaveAndSend(
	c *gin.Context,
	event *linebot.Event,
	args *SearchArgs,
) {
	list, errList := s.searchService.Search(
		args.lat,
		args.lng,
		args.radius,
		args.pageIndex,
		MaximumNumberOfCarouselItems,
	)
	if len(errList) > 0 {
		s.bot.SendText(event.ReplyToken, "搜尋有問題")
		for _, e := range errList {
			s.logService.Error(e)
		}
		return
	}

	if len(list) == 0 {
		s.bot.SendText(event.ReplyToken, "附近沒有店家")
		return
	}

	restaurantList := s.saveRestaurantsToDB(c, list)
	// #TODO 如果爬過就不再爬
	s.sendToCrawlerWork(restaurantList)
	s.sendRestaurantsWithCarousel(event, restaurantList, args)
}

func (s *Server) sendRestaurantsWithCarousel(event *linebot.Event, restaurantList []db.Restaurant, args *SearchArgs) {
	component := carousel.CreateCarouselWithNext(
		restaurantList,
		func(restaurant db.Restaurant) *linebot.BubbleContainer {
			return carousel.CreateRestaurantContainer(restaurant)
		},
		func() *linebot.BubbleContainer {
			if len(restaurantList) < MaximumNumberOfCarouselItems {
				return nil
			}
			return carousel.CreateRestaurantNextPageContainer(args.pageIndex+1, args.lat, args.lng, args.radius)
		},
	)
	s.bot.SendFlex(event.ReplyToken, "carousel", component)
}

func (s *Server) saveRestaurantsToDB(c *gin.Context, list []db.Restaurant) []db.Restaurant {
	restaurantList := []db.Restaurant{}
	for _, restaurant := range list {
		r, err := s.restaurantService.CreateRestaurant(c, db.CreateRestaurantParams{
			Name:             restaurant.Name,
			Rating:           restaurant.Rating,
			UserRatingsTotal: restaurant.UserRatingsTotal,
			Address:          restaurant.Address,
			GoogleMapPlaceID: restaurant.GoogleMapPlaceID,
			GoogleMapUrl:     restaurant.GoogleMapUrl,
			PhoneNumber:      restaurant.PhoneNumber,
			Image:            restaurant.Image,
		})
		if err != nil {
			s.logService.Errorf("CreateRestaurant error: %v", err)
		} else {
			restaurantList = append(restaurantList, r)
		}
	}
	return restaurantList
}

// 送去給爬蟲服務爬蟲
func (s *Server) sendToCrawlerWork(restaurantList []db.Restaurant) {
	for _, restaurant := range restaurantList {
		s.crawlerService.SendWork(restaurant)
	}
}
