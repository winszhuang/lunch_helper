package thirdparty

import (
	"encoding/json"
	"fmt"
	"lunch_helper/thirdparty/model"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetDishesByUberEats(restaurantPageURL string) ([]model.Dish, error) {
	var dishes []model.Dish

	res, err := http.Get(restaurantPageURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
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
