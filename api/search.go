package api

import (
	"fmt"
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

func (s *Server) SearchFirstPageRestaurants(c *gin.Context, event *linebot.Event) {
	userId := event.Source.UserID

	radius := s.messageCache.GetCurrentRadius(userId)
	uc, ok := s.messageCache.GetCurrentLocation(userId)
	if !ok {
		s.bot.SendTextWithQuickReplies(event.ReplyToken, "請先傳送位置資訊再做搜尋哦 ~", quickreply.QuickReplyLocation())
		return
	}

	list, err := s.searchService.Search(
		uc.LatLng.Lat,
		uc.LatLng.Lng,
		radius,
		DefaultPageIndex,
		MaximumNumberOfCarouselItems,
	)
	if err != nil {
		msg := fmt.Sprintf("搜尋有問題: %v", err)
		s.bot.SendText(event.ReplyToken, msg)
		return
	}

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
}

func (s *Server) SearchNextPageRestaurants(c *gin.Context, event *linebot.Event) {
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

	list, err := s.searchService.Search(
		lat,
		lng,
		radius,
		pageIndex,
		MaximumNumberOfCarouselItems,
	)
	if err != nil {
		msg := fmt.Sprintf("搜尋有問題: %v", err)
		s.bot.SendText(event.ReplyToken, msg)
		return
	}

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
}
