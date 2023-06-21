package api

import (
	"database/sql"
	"fmt"
	"lunch_helper/adapter"
	"lunch_helper/bot/carousel"
	"lunch_helper/bot/flex"
	db "lunch_helper/db/sqlc"
	"lunch_helper/food_deliver/model"
	"lunch_helper/service"
	"lunch_helper/util"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (s *Server) HandleGetFoods(c *gin.Context, event *linebot.Event) {
	id, err := util.ParseId("restaurantmenu", event.Postback.Data)
	if err != nil {
		s.logService.Errorf("failed to parse restaurant id: %v, error: %s", id, err)
		s.bot.SendText(event.ReplyToken, "取得餐點失敗")
		return
	}

	restaurant, err := s.restaurantService.GetRestaurant(c, int32(id))
	if err != nil {
		s.logService.Errorf("failed to get restaurant: %v", err)
		return
	}

	foods, err := s.foodService.GetFoods(c, int32(id))
	if err != nil {
		s.bot.SendText(event.ReplyToken, "取得菜單失敗")
		s.logService.Errorf("failed to get foods: %v", err)
		return
	}

	hasMenu := len(foods) > 0
	if hasMenu {
		container := flex.CreateFoodListContainer(foods, restaurant)
		s.bot.SendFlex(event.ReplyToken, "菜單", &container)
	} else {
		s.handleNoMenuCase(c, event, restaurant)
	}
}

// 處理沒有菜單的情況
func (s *Server) handleNoMenuCase(c *gin.Context, event *linebot.Event, r db.Restaurant) {
	if r.GoogleMapUrl == "" {
		s.logService.Errorf("restaurant %s has no google map url. google map id is %s", r.Name, r.GoogleMapPlaceID)
		s.bot.SendText(event.ReplyToken, "未在google上找到相關菜單")
		return
	}

	if r.MenuCrawled {
		s.bot.SendText(event.ReplyToken, "網路上爬不到菜單哦")
		return
	}

	s.bot.SendText(event.ReplyToken, "尚未爬取該菜單，請稍等")
	fetchInfo, err := s.foodDeliverApi.CheckFoodDeliverFromGoogleMap(r.GoogleMapUrl)
	if err != nil {
		s.logService.Debugf("no food deliver link from %s, restaurant name is %s, restaurant id is %d", r.GoogleMapUrl, r.Name, r.ID)
		s.bot.PushText(event.Source.UserID, "網路上爬不到菜單哦")
		return
	}

	result := <-s.taskService.SendPriorityTask(func() service.Result {
		dishes, err := s.foodDeliverApi.GetDishes(fetchInfo)
		return service.Result{Data: dishes, Err: err}
	})
	if result.Err != nil {
		s.logService.Errorf("get dishes from google map error: %v, restaurant name is %s, restaurant id is %d", err, r.Name, r.ID)
		s.bot.PushText(event.Source.UserID, "爬取菜單發生錯誤")
	} else {
		s.logService.Debugf("get dishes from google map success!!, restaurant name is %s, restaurant id is %d", r.Name, r.ID)
		foods, errList := s.foodService.CreateFoodsByDishes(c, service.CreateFoodsByDishesParams{
			Dishes:       result.Data.([]model.Dish),
			RestaurantID: r.ID,
			EditBy:       sql.NullInt32{Valid: false},
		})
		if len(errList) > 0 {
			s.logService.Errorf("failed to create foods: %v", errList)
			s.bot.PushText(event.Source.UserID, "創建菜單異常")
			return
		}
		container := flex.CreateFoodListContainer(foods, r)
		altText := fmt.Sprintf("菜單: %s", r.Name)
		s.bot.PushFlex(event.ReplyToken, altText, &container)
	}
}

func (s *Server) HandleShowFood(c *gin.Context, event *linebot.Event) {
	id, err := util.ParseId("showfood", event.Postback.Data)
	if err != nil {
		s.logService.Errorf("failed to parse food id: %v", err)
		s.bot.SendText(event.ReplyToken, "取得餐點失敗(parse error)")
		return
	}

	food, err := s.foodService.GetFood(c, int32(id))
	if err != nil {
		s.logService.Errorf("failed to get food: %v", err)
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

	userFood, err := s.userFoodService.GetByFoodId(c, db.GetUserFoodByFoodIdParams{
		UserID: user.ID,
		FoodID: food.ID,
	})

	s.logService.Debugf("user food: %v", userFood)
	isUserFood := err == nil
	container := flex.CreateFoodItem(food, isUserFood)
	s.bot.SendFlex(event.ReplyToken, food.Name, &container)
}

func (s *Server) HandleLikeFood(c *gin.Context, event *linebot.Event) {
	foodId, err := util.ParseId("userlikefood", event.Postback.Data)
	if err != nil {
		s.logService.Errorf("failed to parse food: %v, data: %s", err, event.Postback.Data)
		s.bot.SendText(event.ReplyToken, "取得餐點失敗(parse error)")
		return
	}

	userLineId := event.Source.UserID
	user, err := s.userService.GetUserByLineID(c, userLineId)
	if err != nil {
		s.logService.Errorf("failed to get user id: %v", err)
		s.bot.SendText(event.ReplyToken, "取得使用者資訊失敗")
		return
	}

	_, err = s.userFoodService.Create(c, user.ID, int32(foodId))
	if err != nil {
		s.logService.Errorf("failed to create user food: %v", err)
		s.bot.SendText(event.ReplyToken, "加入使用者收藏餐點失敗")
		return
	}

	s.bot.SendText(event.ReplyToken, "成功加入收藏餐點")
}

func (s *Server) HandleUnlikeFood(c *gin.Context, event *linebot.Event) {
	foodId, err := util.ParseId("userunlikefood", event.Postback.Data)
	if err != nil {
		s.logService.Errorf("failed to parse food: %v, data: %s", err, event.Postback.Data)
		s.bot.SendText(event.ReplyToken, "取得餐點失敗(parse error)")
		return
	}

	userLineId := event.Source.UserID
	user, err := s.userService.GetUserByLineID(c, userLineId)
	if err != nil {
		s.logService.Errorf("failed to get user id: %v", err)
		s.bot.SendText(event.ReplyToken, "取得使用者資訊失敗")
		return
	}

	if err = s.userFoodService.Delete(c, db.DeleteUserFoodParams{
		UserID: user.ID,
		FoodID: int32(foodId),
	}); err != nil {
		s.logService.Errorf("failed to delete user food: %v", err)
		s.bot.SendText(event.ReplyToken, "取消收藏餐點失敗")
		return
	}

	s.bot.SendText(event.ReplyToken, "成功取消收藏餐點")
}

func (s *Server) HandleShowFirstPageUserFoods(c *gin.Context, event *linebot.Event) {
	userLineId := event.Source.UserID
	user, err := s.userService.GetUserByLineID(c, userLineId)
	if err != nil {
		s.logService.Errorf("failed to get user id: %v", err)
		s.bot.SendText(event.ReplyToken, "取得使用者資訊失敗")
		return
	}

	listArgs := &ListArgs{PageIndex: 1, PageSize: 10}
	foodList, err := s.userFoodService.List(c, db.GetUserFoodsParams{
		UserID: user.ID,
		Limit:  int32(listArgs.PageSize),
		Offset: int32((listArgs.PageIndex - 1) * 10),
	})
	if err != nil {
		s.logService.Errorf("failed to get user foods: %v", err)
		s.bot.SendText(event.ReplyToken, "取得使用者收藏餐點失敗")
		return
	}

	if len(foodList) == 0 {
		s.bot.SendText(event.ReplyToken, "使用者沒有任何收藏餐點哦，請先搜尋並收藏餐點後再使用此功能")
	} else {
		s.sendUserFoodsWithCarousel(
			event,
			adapter.UserFoodRowsToFoods(foodList),
			&ListArgs{PageIndex: listArgs.PageIndex + 1, PageSize: listArgs.PageSize},
		)
	}

}

func (s *Server) HandleShowNextPageUserFoods(c *gin.Context, event *linebot.Event) {
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

	foodList, err := s.userFoodService.List(c, db.GetUserFoodsParams{
		UserID: user.ID,
		Limit:  int32(pageSize),
		Offset: int32((pageIndex - 1) * pageSize),
	})
	if err != nil {
		s.logService.Errorf("failed to get user foods: %v", err)
		s.bot.SendText(event.ReplyToken, "取得使用者收藏餐點失敗")
		return
	}

	if len(foodList) == 0 {
		s.bot.SendText(event.ReplyToken, "使用者沒有任何收藏餐點哦，請先搜尋並收藏餐點後再使用此功能")
	} else {
		s.sendUserFoodsWithCarousel(
			event,
			adapter.UserFoodRowsToFoods(foodList),
			&ListArgs{PageIndex: pageIndex + 1, PageSize: pageSize},
		)
	}

}

func (s *Server) sendUserFoodsWithCarousel(event *linebot.Event, foodList []db.Food, nextListArgs *ListArgs) {
	component := carousel.CreateCarouselWithNext(
		foodList,
		func(food db.Food) *linebot.BubbleContainer {
			return carousel.CreateFoodCarouselItem(food)
		},
		func() *linebot.BubbleContainer {
			if len(foodList) < MaximumNumberOfCarouselItems {
				return nil
			}
			nextData := fmt.Sprintf(
				"/showuserlikefoodnext?pageIndex=%d&pageSize=%d",
				nextListArgs.PageIndex,
				nextListArgs.PageSize,
			)
			return carousel.CreateNextPageContainer(nextData)
		},
	)
	s.bot.SendFlex(event.ReplyToken, "carousel", component)
}
