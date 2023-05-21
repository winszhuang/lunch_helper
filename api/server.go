package api

import (
	"log"
	"lunch_helper/bot"
	"lunch_helper/cache"
	db "lunch_helper/db/sqlc"
	"lunch_helper/thirdparty"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Server struct {
	store        db.Store
	router       *gin.Engine
	bot          *bot.BotClient
	placeApi     thirdparty.PlaceApi
	messageCache *cache.MessageCache
	nearByCache  *cache.NearByRestaurantCache
}

func NewServer(
	store db.Store,
	bot *bot.BotClient,
	placeApi thirdparty.PlaceApi,
	messageCache *cache.MessageCache,
	nearByCache *cache.NearByRestaurantCache,
) *Server {
	server := &Server{store: store, bot: bot, placeApi: placeApi, messageCache: messageCache, nearByCache: nearByCache}
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.POST("/callback", func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		for _, event := range events {
			switch event.Type {
			case linebot.EventTypeFollow:
				log.Printf("user %s EventTypeFollow", event.Source.UserID)
				server.RegisterUser(c, event)
			case linebot.EventTypeUnfollow:
				log.Printf("user %s EventTypeUnfollow", event.Source.UserID)
			case linebot.EventTypeMessage:
				switch messageData := event.Message.(type) {
				case *linebot.LocationMessage:
					server.SearchRestaurantByLocation(c, event)
				case *linebot.TextMessage:
					message := strings.TrimSpace(messageData.Text)
					bot.SendText(event.ReplyToken, message)
					log.Println("-----------------------")
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
