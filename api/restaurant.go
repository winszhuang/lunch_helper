package api

import (
	"lunch_helper/adapter"
	"lunch_helper/bot/carousel"
	db "lunch_helper/db/sqlc"
	"lunch_helper/util"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type ListArgs struct {
	limit  int
	offset int
}

func (s *Server) HandleLikeRestaurant(c *gin.Context, event *linebot.Event) {
	restaurantId, err := util.ParseId("userlikerestaurant", event.Postback.Data)
	if err != nil {
		s.logService.Errorf("failed to parse restaurant id: %v, error: %s", restaurantId, err)
		s.bot.SendText(event.ReplyToken, "取得餐點失敗")
		return
	}

	userLineId := event.Source.UserID
	user, err := s.userService.GetUserByLineID(c, userLineId)
	if err != nil {
		s.logService.Errorf("failed to get user id: %v", err)
		s.bot.SendText(event.ReplyToken, "取得使用者資訊失敗")
		return
	}

	_, err = s.userRestaurantService.Create(c, user.ID, int32(restaurantId))
	if err != nil {
		s.logService.Errorf("failed to create user food: %v", err)
		s.bot.SendText(event.ReplyToken, "加入使用者收藏店家失敗")
		return
	}

	s.bot.SendText(event.ReplyToken, "成功加入收藏店家")
}

func (s *Server) HandleShowUserRestaurant(c *gin.Context, event *linebot.Event) {
	userLineId := event.Source.UserID
	user, err := s.userService.GetUserByLineID(c, userLineId)
	if err != nil {
		s.logService.Errorf("failed to get user id: %v", err)
		s.bot.SendText(event.ReplyToken, "取得使用者資訊失敗")
		return
	}

	restaurantList, err := s.userRestaurantService.List(c, db.GetUserRestaurantsParams{
		UserID: user.ID,
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		s.logService.Errorf("failed to get user restaurant: %v", err)
		s.bot.SendText(event.ReplyToken, "取得使用者收藏餐廳失敗")
		return
	}

	s.sendUserRestaurantsWithCarousel(
		event,
		adapter.UserRestaurantRowsToRestaurants(restaurantList),
		&ListArgs{
			limit:  10,
			offset: 0,
		},
	)
}

func (s *Server) sendUserRestaurantsWithCarousel(event *linebot.Event, restaurantList []db.Restaurant, args *ListArgs) {
	component := carousel.CreateCarouselWithNext(
		restaurantList,
		func(restaurant db.Restaurant) *linebot.BubbleContainer {
			return carousel.CreateRestaurantContainer(restaurant)
		},
		func() *linebot.BubbleContainer {
			// #TODO 補上下一頁item
			return nil
			// if len(restaurantList) < MaximumNumberOfCarouselItems {
			// 	return nil
			// }
			// return carousel.CreateRestaurantNextPageContainer(args.pageIndex+1, args.lat, args.lng, args.radius)
		},
	)
	s.bot.SendFlex(event.ReplyToken, "carousel", component)
}
