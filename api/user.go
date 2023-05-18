package api

import (
	"fmt"
	db "lunch_helper/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (s *Server) RegisterUser(c *gin.Context, event *linebot.Event) {
	userId := event.Source.UserID

	userProfile, err := s.bot.GetProfile(userId).Do()
	if err != nil {
		s.bot.SendText(event.ReplyToken, "取得使用者資料失敗!! 使用者未能註冊服務!")
		return
	}

	arg := db.CreateUserParams{
		LineID:  userProfile.UserID,
		Name:    userProfile.DisplayName,
		Picture: userProfile.PictureURL,
	}

	_, err = s.store.CreateUser(c, arg)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			s.bot.SendText(event.ReplyToken, fmt.Sprintf("資料庫新建使用者失敗:  %s", pgErr.Code.Name()))
			return
		}
		s.bot.SendText(event.ReplyToken, "新建使用者失敗!!")
		return
	}
}
