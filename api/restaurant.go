package api

import (
	"fmt"
	"lunch_helper/adapter"
	"lunch_helper/bot/carousel"
	"lunch_helper/bot/flex"
	db "lunch_helper/db/sqlc"
	"lunch_helper/util"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type ListArgs struct {
	PageIndex int
	PageSize  int
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

	restaurant, err := s.restaurantService.GetRestaurant(c, int32(restaurantId))
	if err != nil {
		s.logService.Errorf("failed to get restaurant: %v", err)
		s.bot.SendText(event.ReplyToken, "取得店家資訊失敗")
		return
	}
	msg := fmt.Sprintf("-%s-成功加入收藏店家", restaurant.Name)
	s.bot.SendText(event.ReplyToken, msg)
}

func (s *Server) HandleUnLikeRestaurant(c *gin.Context, event *linebot.Event) {
	restaurantId, err := util.ParseId("userunlikerestaurant", event.Postback.Data)
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

	if err = s.userRestaurantService.Delete(c, db.DeleteUserRestaurantParams{
		UserID:       user.ID,
		RestaurantID: int32(restaurantId),
	}); err != nil {
		s.logService.Errorf("failed to delete user food: %v", err)
		s.bot.SendText(event.ReplyToken, "取消收藏收藏店家失敗")
		return
	}

	// #TODO 需要補上店家名稱
	// msg := fmt.Sprintf("-%s-成功加入收藏店家")
	s.bot.SendText(event.ReplyToken, "成功取消收藏店家")
}

func (s *Server) HandleShowFirstPageUserRestaurants(c *gin.Context, event *linebot.Event) {
	userLineId := event.Source.UserID
	user, err := s.userService.GetUserByLineID(c, userLineId)
	if err != nil {
		s.logService.Errorf("failed to get user id: %v", err)
		s.bot.SendText(event.ReplyToken, "取得使用者資訊失敗")
		return
	}

	listArgs := &ListArgs{PageIndex: 1, PageSize: 10}
	restaurantList, err := s.userRestaurantService.List(c, db.GetUserRestaurantsParams{
		UserID: user.ID,
		Limit:  int32(listArgs.PageSize),
		Offset: int32((listArgs.PageIndex - 1) * 10),
	})
	if err != nil {
		s.logService.Errorf("failed to get user restaurant: %v", err)
		s.bot.SendText(event.ReplyToken, "取得使用者收藏餐廳失敗")
		return
	}

	if len(restaurantList) == 0 {
		s.bot.SendText(event.ReplyToken, "使用者沒有任何收藏店家哦，請先搜尋並收藏店家後再使用此功能")
	} else {
		s.sendUserRestaurantsWithCarousel(
			event,
			adapter.UserRestaurantRowsToRestaurants(restaurantList),
			&ListArgs{PageIndex: listArgs.PageIndex + 1, PageSize: listArgs.PageSize},
		)
	}
}

func (s *Server) HandleShowNextPageUserRestaurants(c *gin.Context, event *linebot.Event) {
	userLineId := event.Source.UserID
	user, err := s.userService.GetUserByLineID(c, userLineId)
	if err != nil {
		s.logService.Errorf("failed to get user id: %v", err)
		s.bot.SendText(event.ReplyToken, "取得使用者資訊失敗")
		return
	}

	query := strings.Split(event.Postback.Data, "?")[1]
	values, err := url.ParseQuery(query)
	if err != nil {
		s.logService.Errorf("parse query params error: %v", err)
		s.bot.SendText(event.ReplyToken, "下一頁參數錯誤!!")
		return
	}

	pageIndexStr := values.Get("pageIndex")
	pageSizeStr := values.Get("pageSize")

	pageIndex, err := strconv.Atoi(pageIndexStr)
	if err != nil {
		s.bot.SendText(event.ReplyToken, "解析pageIndex失敗")
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		s.bot.SendText(event.ReplyToken, "解析pageSize失敗")
		return
	}

	restaurantList, err := s.userRestaurantService.List(c, db.GetUserRestaurantsParams{
		UserID: user.ID,
		Limit:  int32(pageSize),
		Offset: int32((pageIndex - 1) * pageSize),
	})
	if err != nil {
		s.logService.Errorf("failed to get user restaurant: %v", err)
		s.bot.SendText(event.ReplyToken, "取得使用者收藏餐廳失敗")
		return
	}

	if len(restaurantList) == 0 {
		s.bot.SendText(event.ReplyToken, "使用者沒有任何收藏店家哦，請先搜尋並收藏店家後再使用此功能")
	} else {
		s.sendUserRestaurantsWithCarousel(
			event,
			adapter.UserRestaurantRowsToRestaurants(restaurantList),
			&ListArgs{PageIndex: pageIndex + 1, PageSize: pageSize},
		)
	}
}

func (s *Server) HandlePickingRestaurant(c *gin.Context, event *linebot.Event) {
	user, err := s.userService.GetUserByLineID(c, event.Source.UserID)
	if err != nil {
		s.logService.Errorf("failed to get user id: %v", err)
		s.bot.SendText(event.ReplyToken, "取得使用者資訊失敗")
		return
	}

	userRestaurants, err := s.userRestaurantService.ListAll(c, user.ID)
	if err != nil {
		s.logService.Errorf("failed to get user restaurants: %v", err)
		s.bot.SendText(event.ReplyToken, "取得使用者收藏店家失敗")
		return
	}

	if len(userRestaurants) == 0 {
		s.bot.SendText(event.ReplyToken, "使用者沒有任何收藏店家哦，請先搜尋並收藏店家後再使用此功能")
		return
	}

	s.sendRandomUserRestaurantWithFlexItem(event, userRestaurants)
}

func (s *Server) sendUserRestaurantsWithCarousel(event *linebot.Event, restaurantList []db.Restaurant, nextListArgs *ListArgs) {
	component := carousel.CreateCarouselWithNext(
		restaurantList,
		func(restaurant db.Restaurant) *linebot.BubbleContainer {
			return carousel.CreateRestaurantCarouselItem(restaurant, func(r db.Restaurant) []linebot.FlexComponent {
				return carousel.PostBackContentsWithShowMenuAndUnLikeAndViewOnMap(r)
			})
		},
		func() *linebot.BubbleContainer {
			if len(restaurantList) < MaximumNumberOfCarouselItems {
				return nil
			}
			nextData := fmt.Sprintf(
				"/showuserlikerestaurantnext?pageIndex=%d&pageSize=%d",
				nextListArgs.PageIndex,
				nextListArgs.PageSize,
			)
			return carousel.CreateNextPageContainer(nextData)
		},
	)
	s.bot.SendFlex(event.ReplyToken, "carousel", component)
}

func (s *Server) sendRandomUserRestaurantWithFlexItem(event *linebot.Event, userRestaurants []db.GetAllUserRestaurantsRow) {
	seed := time.Now().UnixNano()
	random := rand.New(rand.NewSource(seed))

	randomIndex := random.Intn(len(userRestaurants))
	randomRestaurant := userRestaurants[randomIndex]

	isUserFood := true
	container := flex.CreateRestaurantItem(adapter.AllUserRestaurantRowToRestaurant(randomRestaurant), isUserFood)
	s.bot.SendFlex(event.ReplyToken, randomRestaurant.Name, &container)
}
