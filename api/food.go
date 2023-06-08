package api

import (
	"lunch_helper/bot/flex"
	db "lunch_helper/db/sqlc"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (s *Server) HandleGetFoods(c *gin.Context, event *linebot.Event) {
	s.logService.Debugf("event.Postback.Data: %s", event.Postback.Data)
	restaurantId := strings.Split(event.Postback.Data, "/restaurantmenu=")[1]
	id, err := strconv.Atoi(restaurantId)
	if err != nil {
		s.logService.Errorf("failed to parse restaurant id: %v", err)
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

	if len(foods) > 0 {
		container := flex.CreateFoodListContainer(foods, restaurant)
		s.bot.SendFlex(event.ReplyToken, "菜單", &container)
		return
	}

	// 沒有菜單的情況
	if restaurant.GoogleMapUrl == "" {
		s.logService.Errorf("restaurant %s has no google map url. google map id is %s", restaurant.Name, restaurant.GoogleMapPlaceID)
		s.bot.SendText(event.ReplyToken, "未在google上找到相關菜單")
		return
	}

	if restaurant.MenuCrawled {
		s.bot.SendText(event.ReplyToken, "網路上爬不到菜單哦")
		return
	}

	s.bot.SendText(event.ReplyToken, "暫時沒菜單，請等我爬取(等個3秒以上)")

	// 等待爬蟲完該任務再執行後續處理
	crawlSuccess := <-s.crawlerService.SendPriorityWork(restaurant)
	if crawlSuccess {
		foods, err := s.foodService.GetFoods(c, int32(id))
		if err != nil {
			s.bot.PushText(event.Source.UserID, "取得菜單失敗")
			return
		}
		container := flex.CreateFoodListContainer(foods, restaurant)
		s.bot.PushFlex(event.Source.UserID, "菜單", &container)
	} else {
		s.bot.PushText(event.Source.UserID, "爬不到菜單哦")
	}
	s.restaurantService.UpdateMenuCrawled(c, db.UpdateMenuCrawledParams{
		ID:          restaurant.ID,
		MenuCrawled: true,
	})
}
