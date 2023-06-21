package api

import (
	"database/sql"
	"fmt"
	"lunch_helper/bot/carousel"
	"lunch_helper/bot/quickreply"
	"lunch_helper/constant"
	db "lunch_helper/db/sqlc"
	"lunch_helper/food_deliver/model"
	"lunch_helper/service"
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

func (s *Server) HandleSearchFirstPageRestaurants(c *gin.Context, event *linebot.Event) {
	userId := event.Source.UserID

	text := s.messageCache.GetCurrentSearchText(userId)
	radius := s.messageCache.GetCurrentRadius(userId)
	uc, ok := s.messageCache.GetCurrentLocation(userId)
	if !ok {
		s.bot.SendTextWithQuickReplies(event.ReplyToken, "請先傳送位置資訊再做搜尋哦 ~", quickreply.QuickReplyLocation())
		return
	}

	searchArgs := &constant.SearchArgs{
		Lat:       uc.LatLng.Lat,
		Lng:       uc.LatLng.Lng,
		Radius:    radius,
		Text:      text,
		PageIndex: DefaultPageIndex,
		PageSize:  MaximumNumberOfCarouselItems,
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
	textStr := values.Get("text")

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

	searchArgs := &constant.SearchArgs{
		Lat:       lat,
		Lng:       lng,
		Radius:    radius,
		PageIndex: pageIndex,
		PageSize:  MaximumNumberOfCarouselItems,
		Text:      textStr,
	}

	s.searchSaveAndSend(c, event, searchArgs)
}

func (s *Server) searchSaveAndSend(
	c *gin.Context,
	event *linebot.Event,
	args *constant.SearchArgs,
) {
	list, searchErr := s.searchService.Search(args)
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
			go s.handleMenuCrawl(restaurant, c)
		}
	}
}

func (s *Server) handleMenuCrawl(r db.Restaurant, c *gin.Context) {
	defer func() {
		if err := s.restaurantService.UpdateMenuCrawled(c, db.UpdateMenuCrawledParams{
			MenuCrawled: true,
			ID:          r.ID,
		}); err != nil {
			s.logService.Errorf("update menu crawled error: %v, restaurant name is %s, restaurant id is %d", err, r.Name, r.ID)
		}
	}()

	fetchInfo, err := s.foodDeliverApi.CheckFoodDeliverFromGoogleMap(r.GoogleMapUrl)
	if err != nil {
		s.logService.Debugf("no food deliver link from %s, restaurant name is %s, restaurant id is %d", r.GoogleMapUrl, r.Name, r.ID)
		return
	}

	result := <-s.taskService.SendRateLimitTask(func() service.Result {
		dishes, err := s.foodDeliverApi.GetDishes(fetchInfo)
		return service.Result{Data: dishes, Err: err}
	})
	if result.Err != nil {
		s.logService.Errorf("get dishes from google map error: %v, restaurant name is %s, restaurant id is %d", err, r.Name, r.ID)
	} else {
		s.logService.Debugf("get dishes from google map success!!, restaurant name is %s, restaurant id is %d", r.Name, r.ID)
		if _, errList := s.foodService.CreateFoodsByDishes(c, service.CreateFoodsByDishesParams{
			RestaurantID: r.ID,
			Dishes:       result.Data.([]model.Dish),
			EditBy:       sql.NullInt32{Valid: false},
		}); len(errList) > 0 {
			s.logService.Errorf("create foods by dishes error: %v, restaurant name is %s, restaurant id is %d", errList, r.Name, r.ID)
		}
	}
}

func (s *Server) sendRestaurantsWithCarousel(event *linebot.Event, restaurantList []db.Restaurant, args *constant.SearchArgs) {
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
				"/searchnext?lat=%s&lng=%s&radius=%d&pageIndex=%d&text=%s",
				strconv.FormatFloat(args.Lat, 'f', 6, 64),
				strconv.FormatFloat(args.Lng, 'f', 6, 64),
				args.Radius,
				args.PageIndex+1,
				args.Text,
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
