package main

import (
	"database/sql"
	"flag"
	"log"
	"lunch_helper/api"
	"lunch_helper/bot"
	"lunch_helper/config"
	db "lunch_helper/db/sqlc"
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
		log.Fatal("cannot load config:", err)
	}

	// load db
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// load linebot client
	bc, err := bot.NewBotClient(config.LineBotChannelSecret, config.LineBotChannelAccessToken)
	if err != nil {
		log.Fatal("init linebot client error: ", err)
	}
	// set linebot webhook url
	err = bc.SetWebHookUrl(config.ApiUrl, config.LineBotEndpoint)
	if err != nil {
		log.Fatal("setting linebot webhook url error: ", err)
	}

	// run server
	store := db.NewStore(conn)
	server := api.NewServer(store, bc)

	server.Start(port)
}
