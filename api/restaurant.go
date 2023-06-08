package api

import (
	"lunch_helper/util"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

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
