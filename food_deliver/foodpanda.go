package food_deliver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"lunch_helper/food_deliver/model"
	"lunch_helper/util"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var FoodPandaReg = regexp.MustCompile(`\"(https:\/\/www\.foodpanda\.com\.tw[^"]+)\"`)

const FoodPandaApiUrl = "https://tw.fd-api.com/api/v5/vendors/%s?include=menus,bundles,multiple_discounts&language_id=6&opening_type=delivery&basket_currency=TWD&show_pro_deals=true"

type FoodPandaDishesCrawler struct{}

func test(url string) {
	has := strings.Contains(url, "https://www.foodpanda.com.tw/restaurant")
	if has {
		log.Println("foodpanda網址有再裡面阿!!")
	} else {
		log.Println("完蛋... foodpanda網址沒再裡面")
	}

	hasF := strings.Contains(url, "foodpanda")
	if hasF {
		log.Println("有foodpanda字眼阿!!")

	} else {
		log.Println("完蛋...沒有foodpanda字眼")

	}
}

func (fp *FoodPandaDishesCrawler) ParseSource(googleMapUrl string) (string, error) {
	source, err := util.Fetch(googleMapUrl)
	if err != nil {
		return "", err
	}

	log.Println(source)
	err = ioutil.WriteFile("templates/test.html", source, 0644)
	if err != nil {
		log.Println("寫入失敗")
	}

	test(string(source))
	foodpandaMatches := FoodPandaReg.FindStringSubmatch(string(source))
	log.Println("-----foodpandaMatches")
	for _, v := range foodpandaMatches {
		log.Println(v)
	}
	log.Println("-----")
	if len(foodpandaMatches) >= 2 {
		foodpandaURL, err := url.PathUnescape(foodpandaMatches[1])
		if err != nil {
			return "", fmt.Errorf("無法解析foodpandaURL %v", err)
		}
		return foodpandaURL, nil
	}

	return "", fmt.Errorf("找不到foodpandaURL %v", err)
}

func (fp *FoodPandaDishesCrawler) GetDishes(foodPandaURL string) ([]model.Dish, error) {
	restaurantID := getIDByUrl(foodPandaURL)
	apiUrl := fmt.Sprintf(FoodPandaApiUrl, restaurantID)
	source, err := util.Fetch(apiUrl)
	if err != nil {
		return nil, err
	}

	var foodPandaResponse model.Source
	err = json.Unmarshal(source, &foodPandaResponse)
	if err != nil {
		return nil, err
	}

	categories := foodPandaResponse.Data.Menus[0].MenuCategories
	dishes := make([]model.Dish, 0)
	for _, category := range categories {
		categoryName := category.Name
		for _, product := range category.Products {
			name := product.Name
			description := product.Description
			image := ""
			if len(product.Images) > 0 {
				image = product.Images[0].ImageURL
			}
			price := ""
			if len(product.ProductVariations) > 0 {
				price = strconv.FormatFloat(product.ProductVariations[0].Price, 'f', -1, 64)
			}
			dish := model.Dish{Name: name, Description: description, Image: image, Price: price, Category: categoryName}
			dishes = append(dishes, dish)
		}
	}

	return dishes, nil
}

func getIDByUrl(url string) string {
	split := strings.Split(url, "restaurant/")
	return strings.Split(split[1], "/")[0]
}
