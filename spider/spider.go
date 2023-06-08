package spider

import (
	"fmt"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type DeliverLinkSpider interface {
	ScrapeDeliverLink(url string) (string, error)
	Quit()
}

type GoogleDeliverLinkSpider struct {
	service   *selenium.Service
	WebDriver selenium.WebDriver
	doing     chan bool
}

// 抓取外送平台店家網址
func NewGoogleDeliverLinkSpider(chromeDriverPath string) (*GoogleDeliverLinkSpider, error) {
	service, err := selenium.NewChromeDriverService(chromeDriverPath, 4444)
	if err != nil {
		return nil, err
	}

	caps := selenium.Capabilities{"browserName": "chrome"}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"window-size=1920x1080",
		"--no-sandbox",
		"--disable-dev-shm-usage",
		"disable-gpu",
		// "--headless",  // comment out this line to see the browser
	}})

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

	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		return nil, err
	}

	d := &GoogleDeliverLinkSpider{
		WebDriver: wd,
		service:   service,
		doing:     make(chan bool, 1),
	}
	return d, nil
}

func (d *GoogleDeliverLinkSpider) Quit() {
	d.service.Stop()
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

func (d *GoogleDeliverLinkSpider) ScrapeDeliverLink(url string) (string, error) {
	// 確保一次只能有一個抓取動作
	d.doing <- true
	defer func() {
		<-d.doing
	}()

	err := d.WebDriver.Get(url)
	if err != nil {
		return "", err
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

	err = d.WebDriver.WaitWithTimeout(getOrderButtonCondition, 3500)
	if err != nil {
		return "", err
	}

	if orderEl == nil {
		return "", fmt.Errorf("找不到預訂按鈕")
	}

	// 只有一個合作店家的情況
	href, _ := orderEl.GetAttribute("href")
	if href != "" {
		return href, nil
	}

	orderEl.Click()
	selectDeliversCondition := func(wd selenium.WebDriver) (bool, error) {
		el, err := wd.FindElement(selenium.ByCSSSelector, "[aria-label='選擇服務供應商']")
		if err != nil {
			return false, nil
		}

		// #NOTICE 這段很重要!! 沒有加上這個會等不到幾秒就抱錯
		enabled, err := el.IsEnabled()
		if enabled {
			return true, nil
		} else {
			return false, nil
		}
	}

	err = d.WebDriver.WaitWithTimeout(selectDeliversCondition, 5*time.Second)
	if err != nil {
		return "", err
	}

	url, err = d.findDeliverLink()
	if err != nil || url == "" {
		return "", fmt.Errorf("找不到該合作店家")
	}

	return url, nil
}
