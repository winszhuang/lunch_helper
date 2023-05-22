package api

import (
	"log"
	"lunch_helper/bot"
	"lunch_helper/cache"
	"lunch_helper/constant"
	db "lunch_helper/db/sqlc"
	"lunch_helper/thirdparty"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"googlemaps.github.io/maps"
)

type Server struct {
	store        db.Store
	router       *gin.Engine
	bot          bot.BotClient
	placeApi     thirdparty.PlaceApi
	messageCache *cache.MessageCache
	nearByCache  *cache.NearByRestaurantCache
}

func NewServer(
	store db.Store,
	bot bot.BotClient,
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
			userId := event.Source.UserID

			switch event.Type {
			case linebot.EventTypeFollow:
				log.Printf("user %s EventTypeFollow", event.Source.UserID)
				server.RegisterUser(c, event)
			case linebot.EventTypeUnfollow:
				log.Printf("user %s EventTypeUnfollow", event.Source.UserID)
			case linebot.EventTypePostback:
				switch event.Postback.Data {
				case string(constant.Search):
					server.SearchRestaurants(c, event)
				case string(constant.SearchLocation):
					server.bot.SendText(event.ReplyToken, "請傳送定位資訊")
				case string(constant.SearchText):
					server.messageCache.SetMode(userId, constant.SearchText)
					server.bot.SendText(event.ReplyToken, "請輸入搜尋關鍵字")
				case string(constant.SearchRadius):
					server.messageCache.SetMode(userId, constant.SearchRadius)
					// #TODO 修改成flex message讓使用者只能點選500、1000、100等規格
					server.bot.SendText(event.ReplyToken, "請輸入半徑(單位公尺)")
				case string(constant.UserOption):
					// #TODO change richmenu to user option menu
				case string(constant.FavoriteRestaurants):
					// #TODO server.ListFavoriteRestaurants api
				case string(constant.FavoriteFoods):
					// #TODO server.ListFavoriteFoods api
				case string(constant.PickRestaurant):
					// #TODO server.PickRestaurant api
				case string(constant.NotificationSetting):
					// #TODO 增加user_notification table
					// #TODO 修改成flex message讓使用者可以對不同item(時間)做新增編輯刪除
				case string(constant.SearchOption):
					// #TODO change richmenu to search option menu
				}
			case linebot.EventTypeMessage:
				switch messageData := event.Message.(type) {
				case *linebot.LocationMessage:
					server.messageCache.UpdateLocation(userId, &cache.UserLocation{
						ID:      messageData.ID,
						LatLng:  maps.LatLng{messageData.Latitude, messageData.Longitude},
						Title:   messageData.Title,
						Address: messageData.Address,
					})
					server.bot.SendText(event.ReplyToken, "更新定位資訊成功，可以開始搜尋")
				case *linebot.TextMessage:
					message := strings.TrimSpace(messageData.Text)
					switch message {
					case string(constant.Search):
						server.SearchRestaurants(c, event)
					case string(constant.SearchLocation):
						server.bot.SendText(event.ReplyToken, "請傳送定位資訊")
					case string(constant.SearchText):
						server.messageCache.SetMode(userId, constant.SearchText)
						server.bot.SendText(event.ReplyToken, "請輸入搜尋關鍵字")
					case string(constant.SearchRadius):
						server.messageCache.SetMode(userId, constant.SearchRadius)
						// #TODO 修改成flex message讓使用者只能點選500、1000、100等規格
						server.bot.SendText(event.ReplyToken, "請輸入半徑(單位公尺)")
					case string(constant.UserOption):
						// #TODO change richmenu to user option menu
					case string(constant.FavoriteRestaurants):
						// #TODO server.ListFavoriteRestaurants api
					case string(constant.FavoriteFoods):
						// #TODO server.ListFavoriteFoods api
					case string(constant.PickRestaurant):
						// #TODO server.PickRestaurant api
					case string(constant.NotificationSetting):
						// #TODO 增加user_notification table
						// #TODO 修改成flex message讓使用者可以對不同item(時間)做新增編輯刪除
					case string(constant.SearchOption):
						// #TODO change richmenu to search option menu
					default:
						switch server.messageCache.GetUserMode(userId) {
						case constant.SearchText:
							server.messageCache.UpdateSearchText(userId, message)
							server.bot.SendText(event.ReplyToken, "更新搜尋關鍵字成功")
							server.messageCache.SetMode(userId, "")
						case constant.SearchRadius:
							num, err := strconv.Atoi(message)
							if err != nil {
								server.bot.SendText(event.ReplyToken, "請正確輸入數字再提交!!")
								return
							}
							server.messageCache.UpdateSearchRadius(userId, num)
							server.bot.SendText(event.ReplyToken, "更新半徑成功")
							server.messageCache.SetMode(userId, "")
						}
						bot.SendText(event.ReplyToken, message)
						log.Println("-----------------------")
					}
				}
			}
		}

		c.Status(http.StatusOK)
	})

	server.router = router
	return server
}

func handleTextMessage() {

}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Start(port string) error {
	return server.router.Run(":" + port)
}
