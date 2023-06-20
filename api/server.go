package api

import (
	"log"
	"lunch_helper/bot"
	"lunch_helper/bot/quickreply"
	"lunch_helper/cache"
	"lunch_helper/constant"
	"lunch_helper/food_deliver"
	"lunch_helper/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"googlemaps.github.io/maps"
)

type Server struct {
	router                *gin.Engine
	bot                   bot.BotClient
	messageCache          *cache.MessageCache
	nearByCache           *cache.NearByRestaurantCache
	searchService         *service.SearchService
	userService           *service.UserService
	userFoodService       *service.UserFoodService
	restaurantService     *service.RestaurantService
	userRestaurantService *service.UserRestaurantService
	foodService           *service.FoodService
	logService            *service.LogService
	foodDeliverApi        *food_deliver.FoodDeliverApi
}

func NewServer(
	bot bot.BotClient,
	messageCache *cache.MessageCache,
	nearByCache *cache.NearByRestaurantCache,
	searchService *service.SearchService,
	userService *service.UserService,
	userFoodService *service.UserFoodService,
	restaurantService *service.RestaurantService,
	userRestaurantService *service.UserRestaurantService,
	foodService *service.FoodService,
	logService *service.LogService,
	foodDeliverApi *food_deliver.FoodDeliverApi,
) *Server {
	server := &Server{
		bot:                   bot,
		messageCache:          messageCache,
		nearByCache:           nearByCache,
		searchService:         searchService,
		userService:           userService,
		userFoodService:       userFoodService,
		restaurantService:     restaurantService,
		userRestaurantService: userRestaurantService,
		foodService:           foodService,
		logService:            logService,
		foodDeliverApi:        foodDeliverApi,
	}
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "welcome to lunch helper ~"})
	})

	router.POST("/callback", func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeFollow {
				log.Printf("user %s EventTypeFollow", event.Source.UserID)
				server.RegisterUser(c, event)
				return
			}

			if event.Type == linebot.EventTypeUnfollow {
				log.Printf("user %s EventTypeUnfollow", event.Source.UserID)
				return
			}

			// check user is register
			if !isUserRegister(server, event, c) {
				server.bot.SendText(event.ReplyToken, "尚未註冊會員!! 請先註冊再使用功能。若已再頻道中，請先封鎖頻道再重新解封加入")
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			if event.Type == linebot.EventTypePostback {
				handlePostbackEvent(server, event, c)
				return
			}

			if event.Type == linebot.EventTypeMessage {
				switch event.Message.(type) {
				case *linebot.LocationMessage:
					handleLocationEvent(server, event)
					return
				case *linebot.TextMessage:
					handleTextEvent(server, event, c)
					return
				}
			}
		}

		c.Status(http.StatusOK)
	})

	server.router = router
	return server
}

func (server *Server) Start(port string) error {
	return server.router.Run(":" + port)
}

func handleLocationEvent(server *Server, event *linebot.Event) {
	messageData := event.Message.(*linebot.LocationMessage)
	server.messageCache.UpdateLocation(event.Source.UserID, &cache.UserLocation{
		ID:      messageData.ID,
		LatLng:  maps.LatLng{Lat: messageData.Latitude, Lng: messageData.Longitude},
		Title:   messageData.Title,
		Address: messageData.Address,
	})
	server.bot.SendText(event.ReplyToken, "更新定位資訊成功，可以開始搜尋")
}

func handleTextEvent(server *Server, event *linebot.Event, c *gin.Context) {
	messageData := event.Message.(*linebot.TextMessage)
	server.logService.Debugf("message %s", messageData.Text)

	userId := event.Source.UserID
	message := strings.TrimSpace(messageData.Text)

	switch message {
	case string(constant.Search):
		server.HandleSearchFirstPageRestaurants(c, event)
	case string(constant.SearchLocation):
		server.bot.SendText(event.ReplyToken, "請傳送定位資訊")
	case string(constant.SearchText):
		server.messageCache.SetMode(userId, constant.SearchText)
		server.bot.SendText(event.ReplyToken, `請輸入搜尋關鍵字，輸入"無"或者"n"將會清除關鍵字`)
	case string(constant.SearchRadius):
		server.messageCache.SetMode(userId, constant.SearchRadius)
		server.bot.SendTextWithQuickReplies(event.ReplyToken, "請選擇半徑(單位公尺)", quickreply.QuickReplyRadiusOptions())
	case string(constant.FavoriteRestaurants):
		server.HandleShowFirstPageUserRestaurants(c, event)
	case string(constant.FavoriteFoods):
		server.HandleShowFirstPageUserFoods(c, event)
	case string(constant.PickRestaurant):
		server.HandlePickingRestaurant(c, event)
	case string(constant.Trends):
		server.bot.SendText(event.ReplyToken, "尚未開放，敬請期待")
	case string(constant.NotificationSetting):
		server.bot.SendText(event.ReplyToken, "尚未開放，敬請期待")
		// #TODO 增加user_notification table
		// #TODO 修改成flex message讓使用者可以對不同item(時間)做新增編輯刪除
	default:
		switch server.messageCache.GetUserMode(userId) {
		case constant.SearchText:
			if message == "n" || message == "無" {
				message = ""
			}
			server.messageCache.UpdateSearchText(userId, message)
			server.bot.SendText(event.ReplyToken, "更新搜尋關鍵字成功")
			server.messageCache.SetMode(userId, "")
		case constant.SearchRadius:
			num, err := strconv.Atoi(message)
			if err != nil {
				server.bot.SendTextWithQuickReplies(event.ReplyToken, "請正確輸入數字再提交!!", quickreply.QuickReplyRadiusOptions())
				return
			}
			if num < 0 || num > 2000 {
				server.bot.SendTextWithQuickReplies(event.ReplyToken, "目前僅提供0 ~ 2000公尺距離範圍，請重新嘗試", quickreply.QuickReplyRadiusOptions())
				return
			}
			server.messageCache.UpdateSearchRadius(userId, num)
			server.bot.SendText(event.ReplyToken, "更新半徑成功")
			server.messageCache.SetMode(userId, "")
		}
	}
}

func handlePostbackEvent(server *Server, event *linebot.Event, c *gin.Context) {
	userId := event.Source.UserID
	server.logService.Debugf("current postback data: %s", event.Postback.Data)

	// no params postback data
	switch event.Postback.Data {
	case string(constant.Search):
		server.HandleSearchFirstPageRestaurants(c, event)
	case string(constant.SearchLocation):
		server.bot.SendTextWithQuickReplies(event.ReplyToken, "請傳送定位資訊", quickreply.QuickReplyLocation())
	case string(constant.SearchText):
		server.messageCache.SetMode(userId, constant.SearchText)
		server.bot.SendText(event.ReplyToken, `請輸入搜尋關鍵字，輸入"無"或者"n"將會清除關鍵字`)
	case string(constant.SearchRadius):
		server.messageCache.SetMode(userId, constant.SearchRadius)
		server.bot.SendTextWithQuickReplies(event.ReplyToken, "請選擇半徑(單位公尺)", quickreply.QuickReplyRadiusOptions())
	case string(constant.SearchAI):
		// #TODO aiMode
		server.bot.SendText(event.ReplyToken, "尚未開放，敬請期待")
	case string(constant.FavoriteRestaurants):
		server.HandleShowFirstPageUserRestaurants(c, event)
	case string(constant.FavoriteFoods):
		server.HandleShowFirstPageUserFoods(c, event)
	case string(constant.PickRestaurant):
		server.HandlePickingRestaurant(c, event)
	case string(constant.NotificationSetting):
		server.bot.SendText(event.ReplyToken, "尚未開放，敬請期待")
		// #TODO 增加user_notification table
		// #TODO 修改成flex message讓使用者可以對不同item(時間)做新增編輯刪除
	case string(constant.Trends):
		server.bot.SendText(event.ReplyToken, "尚未開放，敬請期待")
	}

	// postback data with params
	switch {
	case strings.Contains(event.Postback.Data, "/searchnext"):
		server.HandleSearchNextPageRestaurants(c, event)
	case strings.Contains(event.Postback.Data, "/restaurantmenu"):
		server.HandleGetFoods(c, event)
	case strings.Contains(event.Postback.Data, "/showfood"):
		server.HandleShowFood(c, event)
	case strings.Contains(event.Postback.Data, "/userlikefood"):
		server.HandleLikeFood(c, event)
	case strings.Contains(event.Postback.Data, "/userlikerestaurant"):
		server.HandleLikeRestaurant(c, event)
	case strings.Contains(event.Postback.Data, "/showuserlikerestaurantnext"):
		server.HandleShowNextPageUserRestaurants(c, event)
	case strings.Contains(event.Postback.Data, "/showuserlikefoodnext"):
		server.HandleShowNextPageUserFoods(c, event)
	case strings.Contains(event.Postback.Data, "/userunlikefood"):
		server.HandleUnlikeFood(c, event)
	case strings.Contains(event.Postback.Data, "/userunlikerestaurant"):
		server.HandleUnLikeRestaurant(c, event)
	}
}

func isUserRegister(s *Server, event *linebot.Event, c *gin.Context) bool {
	_, err := s.userService.GetUserByLineID(c, event.Source.UserID)
	return err == nil
}
