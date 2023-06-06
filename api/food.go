package api

import (
	"log"
	"lunch_helper/bot/flex"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (s *Server) HandleGetFoods(c *gin.Context, event *linebot.Event) {
	log.Printf("event.Postback.Data: %s", event.Postback.Data)
	restaurantId := strings.Split(event.Postback.Data, "/restaurantmenu=")[1]
	id, err := strconv.Atoi(restaurantId)
	if err != nil {
		log.Printf("failed to parse restaurant id: %v", err)
	}

	restaurant, err := s.restaurantService.GetRestaurant(c, int32(id))
	if err != nil {
		log.Printf("failed to get restaurant: %v", err)
		return
	}

	foods, err := s.foodService.GetFoods(c, int32(id))
	if err != nil {
		s.bot.SendText(event.ReplyToken, "取得菜單失敗")
		log.Printf("failed to get foods: %v", err)
		return
	}

	if len(foods) > 0 {
		container := flex.CreateFoodListContainer(foods, restaurant)
		s.bot.SendFlex(event.ReplyToken, "菜單", &container)
		return
	}

	// 沒有菜單的情況
	if restaurant.GoogleMapUrl == "" {
		s.bot.SendText(event.ReplyToken, "未在google上找到相關菜單")
		return
	}
	s.crawlerService.SendPriorityWork(restaurant)
	s.bot.SendText(event.ReplyToken, "爬取菜單中，請稍後再試")

	// #TODO 等待直到該任務爬蟲完
	// s.crawlerService.CheckWorkDone(restaurant.GoogleMapUrl)
	// container := flex.CreateFoodListContainer(foods, restaurant)
	// s.bot.SendFlex(event.ReplyToken, "菜單", &container)
}
