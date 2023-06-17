package food_deliver

import (
	"encoding/json"
	"fmt"
	"lunch_helper/food_deliver/model"
	"lunch_helper/util"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var UberEatsReg = regexp.MustCompile(`\"(https:\/\/www\.ubereats\.com[^"]+)\"`)

type UberEatsDishesCrawler struct{}

func (fp *UberEatsDishesCrawler) ParseSource(googleMapUrl string) (string, error) {
	source, err := util.FetchBytes(googleMapUrl)
	if err != nil {
		return "", err
	}

	ubereatsMatches := UberEatsReg.FindStringSubmatch(string(source))
	if len(ubereatsMatches) >= 2 {
		ubereatsURL, err := url.PathUnescape(ubereatsMatches[1])
		if err != nil {
			return "", fmt.Errorf("無法解析ubereatsURL %v", err)
		}
		return ubereatsURL, nil
	}

	return "", fmt.Errorf("找不到ubereatsURL %v", err)
}

func (fp *UberEatsDishesCrawler) GetDishes(uberEatsURL string) ([]model.Dish, error) {
	reader, err := util.Fetch(uberEatsURL)
	defer reader.Close()
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	scriptContent := doc.Find("#main-content").Children().First().Text()
	scriptContent = strings.TrimSpace(scriptContent)
	scriptContent = strings.Trim(scriptContent, "\"")
	if scriptContent == "" {
		return nil, fmt.Errorf("failed to find the #main-content element or its first child element, or the element is empty")
	}

	jsonData := &model.JsonData{}
	err = json.Unmarshal([]byte(scriptContent), jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %v", err)
	}

	var dishes []model.Dish
	for _, category := range jsonData.HasMenu.HasMenuSection {
		for _, menuItem := range category.HasMenuItem {
			dish := model.Dish{
				Name:        menuItem.Name,
				Description: menuItem.Description,
				Price:       menuItem.Offers.Price,
				Category:    category.Name,
			}
			dishes = append(dishes, dish)
		}
	}
	return dishes, nil
}
