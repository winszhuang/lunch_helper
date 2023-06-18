package worker

import (
	"fmt"
	"net/http"
	"time"
)

// #NOTICE render.com每15分鐘休眠，需要定時ping來喚醒
func PingWebServerEveryMinutes(serverUrl string, minutes int) {
	ticker := time.NewTicker(time.Duration(minutes) * time.Minute)

	for {
		select {
		case <-ticker.C:
			pingWebServer(serverUrl)
		}
	}
}

func pingWebServer(serverUrl string) {
	_, err := http.Get(serverUrl)
	if err != nil {
		fmt.Println("Error pinging web server:", err)
	} else {
		fmt.Println("Web server ping successful")
	}
}
