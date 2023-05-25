package main

import (
	"database/sql"
	"flag"
	"log"
	"lunch_helper/api"
	"lunch_helper/bot"
	"lunch_helper/bot/richmenu"
	"lunch_helper/cache"
	"lunch_helper/config"
	db "lunch_helper/db/sqlc"
	"lunch_helper/service"
	"lunch_helper/thirdparty"
)

var (
	appEnv string
	port   string
)

func main() {
	// parse flags when run main
	flag.StringVar(&appEnv, "APP_ENV", "dev", "current environment")
	flag.StringVar(&port, "PORT", "8080", "current port")
	flag.Parse()

	// load config
	config, err := config.New(".", appEnv)
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// load db
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	// load linebot client
	bc, err := bot.NewBotClient(config.LineBotChannelSecret, config.LineBotChannelAccessToken)
	if err != nil {
		log.Fatalf("init linebot client error: %v", err)
	}
	// set linebot webhook url
	err = bc.SetWebHookUrl(config.ApiUrl, config.LineBotEndpoint)
	if err != nil {
		log.Fatalf("setting linebot webhook url error: %v", err)
	}
	// setup richmenu
	err = bc.SetupRichMenu(richmenu.CreateSearchRichMenu(), "richmenu.png")
	if err != nil {
		log.Fatalf("setup richmenu error: %v", err)
	}

	// init user input cache and nearby place cache
	messageCache := cache.NewMessageCache()
	nearByCache := cache.NewNearByRestaurantCache()

	// init place api
	placeApi, err := thirdparty.NewGoogleMapPlaceApi(config.GoogleMapApiKey)
	if err != nil {
		log.Fatalf("init google map api error: %v", err)
	}

	// init service
	searchService := service.NewSearchService(nearByCache, &placeApi)

	// run server
	store := db.NewStore(conn)
	server := api.NewServer(
		store,
		bc,
		messageCache,
		nearByCache,
		searchService,
	)

	server.Start(port)
}
