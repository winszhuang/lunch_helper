package api

import (
	"lunch_helper/bot"
	db "lunch_helper/db/sqlc"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Server struct {
	store  db.Store
	router *gin.Engine
	bot    *bot.BotClient
}

func NewServer(store db.Store, bot *bot.BotClient) *Server {
	server := &Server{store: store, bot: bot}
	router := gin.Default()

	router.POST("/callback", func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		for _, event := range events {
			switch event.Type {
			case linebot.EventTypeJoin:
				server.RegisterUser(c, event)
			case linebot.EventTypeMessage:
				switch messageData := event.Message.(type) {
				case *linebot.LocationMessage:
					c.JSON(http.StatusOK, gin.H{})
				case *linebot.TextMessage:
					message := strings.TrimSpace(messageData.Text)
					bot.SendText(event.ReplyToken, message)
					// if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					// 	c.AbortWithStatus(http.StatusInternalServerError)
					// 	return
					// }
				}
			}
		}

		c.Status(http.StatusOK)
	})

	server.router = router
	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Start(port string) error {
	return server.router.Run(":" + port)
}
