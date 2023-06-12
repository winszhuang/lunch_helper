package api

import (
	"fmt"
	"lunch_helper/bot/carousel"
	"lunch_helper/bot/quickreply"
	db "lunch_helper/db/sqlc"
	"net/url"
	"strconv"
	"strings"

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
	query := strings.Split(event.Postback.Data, "?")[1]
	values, err := url.ParseQuery(query)
	if err != nil {
		s.logService.Errorf("parse query params error: %v", err)
		s.bot.SendText(event.ReplyToken, "下一頁參數錯誤!!")
		return
	}

	latStr := values.Get("lat")
	lngStr := values.Get("lng")
	radiusStr := values.Get("radius")
	pageIndexStr := values.Get("pageIndex")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		s.bot.SendText(event.ReplyToken, "解析Lat失敗")
		return
	}
	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		s.bot.SendText(event.ReplyToken, "解析Lng失敗")
		return
	}
	radius, err := strconv.Atoi(radiusStr)
	if err != nil {
		s.bot.SendText(event.ReplyToken, "解析半徑失敗")
		return
	}
	pageIndex, err := strconv.Atoi(pageIndexStr)
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
	list, searchErr := s.searchService.Search(
		args.lat,
		args.lng,
		args.radius,
		args.pageIndex,
		MaximumNumberOfCarouselItems,
	)
	if searchErr.Err != nil {
		s.logService.Error(searchErr.Err)
		s.bot.SendText(event.ReplyToken, "搜尋有問題")
		return
	}

	// 紀錄呼叫map detail api是否有異常
	for _, e := range searchErr.DetailErrors {
		s.logService.Error(e)
	}

	if len(list) == 0 {
		s.bot.SendText(event.ReplyToken, "附近沒有店家")
		return
	}

	restaurantList := s.saveRestaurantsToDB(c, list)
	s.sendRestaurantsWithCarousel(event, restaurantList, args)

	// send to crawl
	for _, restaurant := range restaurantList {
		if !restaurant.MenuCrawled {
			go func(r db.Restaurant) {
				<-s.crawlerService.SendWork(r)
				// 確定做完才更新"已爬蟲"
				s.restaurantService.UpdateMenuCrawled(c, db.UpdateMenuCrawledParams{
					MenuCrawled: true,
					ID:          r.ID,
				})
			}(restaurant)
		}
	}
}

func (s *Server) sendRestaurantsWithCarousel(event *linebot.Event, restaurantList []db.Restaurant, args *SearchArgs) {
	component := carousel.CreateCarouselWithNext(
		restaurantList,
		func(restaurant db.Restaurant) *linebot.BubbleContainer {
			return carousel.CreateRestaurantCarouselItem(restaurant, func(r db.Restaurant) []linebot.FlexComponent {
				return carousel.PostBackContentsWithShowMenuAndLikeAndViewOnMap(r)
			})
		},
		func() *linebot.BubbleContainer {
			if len(restaurantList) < MaximumNumberOfCarouselItems {
				return nil
			}
			nextData := fmt.Sprintf(
				"/searchnext?lat=%s&lng=%s&radius=%d&pageIndex=%d",
				strconv.FormatFloat(args.lat, 'f', 6, 64),
				strconv.FormatFloat(args.lng, 'f', 6, 64),
				args.radius,
				args.pageIndex+1,
			)
			return carousel.CreateNextPageContainer(nextData)
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
