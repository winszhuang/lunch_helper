package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"lunch_helper/api"
	"lunch_helper/bot"
	"lunch_helper/bot/richmenu"
	"lunch_helper/cache"
	"lunch_helper/config"
	"lunch_helper/constant"
	db "lunch_helper/db/sqlc"
	"lunch_helper/food_deliver"
	"lunch_helper/service"
	"lunch_helper/thirdparty"
	"lunch_helper/worker"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	// prevent the cloud server from sleeping by pinging the web server every 13 minutes
	go worker.PingWebServerEveryMinutes(config.ApiUrl, 13)

	// load db
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	// run db migration
	err = runDBMigration(config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	// load linebot client
	bc, err := bot.NewBotClient(config.LineBotChannelSecret, config.LineBotChannelAccessToken)
	if err != nil {
		log.Fatalf("init linebot client error: %v", err)
	}
	// set linebot webhook url
	if err = bc.SetWebHookUrl(config.ApiUrl, config.LineBotEndpoint); err != nil {
		log.Fatalf("setting linebot webhook url error: %v", err)
	}

	// clean richmenu
	if errList := bc.ResetRichMenu(); len(errList) > 0 {
		errString := ""
		for _, err := range errList {
			errString += err.Error() + "\n"
		}
		log.Fatalf(errString)
	}

	// setup richmenu
	if err = bc.SetupRichMenu(
		bot.RichMenuRequest{
			AliasName:      constant.RichMenuSearch,
			ImagePath:      "richmenu-search.png",
			RichMenuStruct: richmenu.CreateSearchRichMenu(),
		},
		bot.RichMenuRequest{
			AliasName:      constant.RichMenuUser,
			ImagePath:      "richmenu-user.png",
			RichMenuStruct: richmenu.CreateUserRichMenu(),
		},
	); err != nil {
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

	// init db store
	store := db.NewStore(conn)

	// init food deliver api
	foodDeliverApi := food_deliver.NewFoodDeliverApi()

	// init service
	logService := service.NewLogService("log/debug.txt", "log/error.txt")
	defer logService.Sync()

	userService := service.NewUserService(store)
	userFoodService := service.NewUserFoodService(store)
	restaurantService := service.NewRestaurantService(store)
	userRestaurantService := service.NewUserRestaurantService(store)
	foodService := service.NewFoodService(store)
	searchService := service.NewSearchService(nearByCache, &placeApi)

	// init api server
	server := api.NewServer(
		bc,
		messageCache,
		nearByCache,
		searchService,
		userService,
		userFoodService,
		restaurantService,
		userRestaurantService,
		foodService,
		logService,
		foodDeliverApi,
	)

	server.Start(port)
}

func runDBMigration(dbSource string) error {
	migration, err := migrate.New("file://db/migration", dbSource)
	if err != nil {
		return fmt.Errorf("cannot create new migrate instance: %w", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrate up: %w", err)
	}

	fmt.Println("db migrated successfully")
	return nil
}
