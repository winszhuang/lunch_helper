package spider

import (
	"fmt"
	"log"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type DeliverLinkSpider interface {
	ScrapeDeliverLink(url string) chan Response
	Quit()
}

type GoogleDeliverLinkSpider struct {
	service        *selenium.Service
	WebDriver      selenium.WebDriver
	doing          chan bool
	reconnectChan  chan bool
	disconnectChan chan bool
	isHeadless     bool
}

// 抓取外送平台店家網址
func NewGoogleDeliverLinkSpider(chromeDriverPath string, isHeadless bool) (*GoogleDeliverLinkSpider, error) {
	service, err := selenium.NewChromeDriverService(chromeDriverPath, 4444)
	if err != nil {
		return nil, err
	}

	wd, err := NewWebDriver(isHeadless)
	if err != nil {
		return nil, err
	}

	d := &GoogleDeliverLinkSpider{
		WebDriver:      wd,
		service:        service,
		doing:          make(chan bool, 1),
		reconnectChan:  make(chan bool, 1),
		disconnectChan: make(chan bool, 1),
		isHeadless:     isHeadless,
	}

	go d.doCheckConnected()

	return d, nil
}

func (d *GoogleDeliverLinkSpider) doCheckConnected() {
	for {
		if !d.isConnected() {
			d.disconnectChan <- true
			err := d.RestartWebDriver()
			if err != nil {
				log.Fatalf("重啟不能，沒救了: %v", err)
			} else {
				d.reconnectChan <- true
			}
		}
		time.Sleep(time.Second * 1)
	}
}

func NewWebDriver(isHeadless bool) (selenium.WebDriver, error) {
	caps := selenium.Capabilities{"browserName": "chrome"}

	capabilities := chrome.Capabilities{
		Args: []string{
			"window-size=1920x1080",
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"disable-gpu",
		},
	}
	if isHeadless {
		capabilities.Args = append(capabilities.Args, "--headless")
	}
	caps.AddChrome(capabilities)

	chromeOptions := map[string]interface{}{
		"args": []string{
			"--excludeSwitches=enable-automation",
			"--Referer=https://www.google.com.tw/",
			"--Sec-Ch-Ua=Google Chrome;v=113, Chromium;v=113, Not-A.Brand;v=24",
			"--Accept-Language=zh-TW,zh;q=0.9,en-US;q=0.8,en;q=0.7",
			"--Accept-Encoding=gzip, deflate, br",
			"--Sec-Ch-Ua-Platform=Windows",
			"--Accept=text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			"--Cache-Control=max-age=0",
			"--User-Agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
		},
	}
	caps["chromeOptions"] = chromeOptions

	return selenium.NewRemote(caps, "")
}

func (d *GoogleDeliverLinkSpider) RestartWebDriver() error {
	wd, err := NewWebDriver(d.isHeadless)
	if err != nil {
		return err
	}
	d.WebDriver = wd
	return nil
}

func (d *GoogleDeliverLinkSpider) isConnected() bool {
	_, err := d.WebDriver.CurrentWindowHandle()
	return err == nil
}

func (d *GoogleDeliverLinkSpider) Quit() {
	d.WebDriver.Quit()
}

func (d *GoogleDeliverLinkSpider) findDeliverLink() (string, error) {
	foodpandaUrl, err := d.WebDriver.FindElement(selenium.ByCSSSelector, "[aria-label='foodpanda.com.tw']")
	if err != nil {
		return "", err
	}
	if foodpandaUrl != nil {
		return foodpandaUrl.GetAttribute("href")
	}

	ubereatsUrl, err := d.WebDriver.FindElement(selenium.ByCSSSelector, "[aria-label='ubereats.com']")
	if err != nil {
		return "", err
	}
	if ubereatsUrl != nil {
		return ubereatsUrl.GetAttribute("href")
	}
	return "", nil
}

func (d *GoogleDeliverLinkSpider) WorkAfterReconnect(callback func()) {
	<-d.reconnectChan
	callback()
}

type Response struct {
	ResultLink  string
	ShouldRetry bool
	Err         error
}

type Cache struct {
	ResultLink string
	Err        error
}

func (d *GoogleDeliverLinkSpider) ScrapeDeliverLink(url string) chan Response {
	// 確保一次只能有一個抓取動作
	d.doing <- true

	cache := make(chan Cache)
	response := make(chan Response)

	go func() {
		defer func() {
			// 不能關，會報調
			// close(cache)
			<-d.doing
		}()
		// 判斷
		// #TODO 考慮搶先任務(當前任務中斷)
		select {
		// 有提前斷線就走這
		case <-d.disconnectChan:
			log.Printf("執行%s的時候斷線搂", url)
			log.Printf("需要重新執行%s的爬蟲", url)
			d.WorkAfterReconnect(func() {
				response <- Response{ShouldRetry: true}
			})
		case data := <-cache:
			response <- Response{ResultLink: data.ResultLink, ShouldRetry: false, Err: data.Err}
		}
	}()

	// 實際爬蟲
	go func() {
		defer func() {
			if r := recover(); r != nil {
				err := r.(error)
				inputToCache(cache, Cache{ResultLink: "", Err: err})
			}
		}()
		err := d.WebDriver.Get(url)
		if err != nil {
			panic(err)
		}

		var orderEl selenium.WebElement
		getOrderButtonCondition := func(wd selenium.WebDriver) (bool, error) {
			orderEl, err = wd.FindElement(selenium.ByCSSSelector, "[aria-label='預訂']")
			if err != nil {
				return false, nil
			}

			// #NOTICE 這段很重要!! 沒有加上這個會等不到幾秒就抱錯
			enabled, err := orderEl.IsEnabled()
			if !enabled {
				return false, nil
			}

			return orderEl != nil, err
		}

		err = d.WebDriver.WaitWithTimeout(getOrderButtonCondition, time.Second*4)
		if err != nil {
			panic(err)
		}

		if orderEl == nil {
			panic(fmt.Errorf("找不到預訂按鈕"))
		}

		// 只有一個合作店家的情況
		href, _ := orderEl.GetAttribute("href")
		if href != "" {
			inputToCache(cache, Cache{ResultLink: href, Err: nil})
			return
		}

		orderEl.Click()
		selectDeliversCondition := func(wd selenium.WebDriver) (bool, error) {
			el, err := wd.FindElement(selenium.ByCSSSelector, "[aria-label='選擇服務供應商']")
			if err != nil {
				return false, nil
			}

			// #NOTICE 這段很重要!! 沒有加上這個會等不到幾秒就抱錯
			enabled, _ := el.IsEnabled()
			if enabled {
				return true, nil
			} else {
				return false, nil
			}
		}

		err = d.WebDriver.WaitWithTimeout(selectDeliversCondition, 5*time.Second)
		if err != nil {
			panic(err)
		}

		url, err = d.findDeliverLink()
		if err != nil || url == "" {
			panic(fmt.Errorf("找不到該合作店家"))
		}

		inputToCache(cache, Cache{ResultLink: url, Err: nil})
	}()

	return response
}

func inputToCache(cache chan<- Cache, data Cache) {
	select {
	case cache <- data:
	default:
		log.Printf("cache channel is full")
	}
}
