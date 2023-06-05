package thirdparty

import (
	"encoding/json"
	"fmt"
	"lunch_helper/thirdparty/model"
	"net/http"
	"strconv"
	"strings"
)

const (
	apiURL = "https://tw.fd-api.com/api/v5/vendors/%s?include=menus,bundles,multiple_discounts&language_id=6&opening_type=delivery&basket_currency=TWD&show_pro_deals=true"
)

func GetDishesByFoodPanda(restaurantPageURL string) ([]model.Dish, error) {
	restaurantID := getIDByUrl(restaurantPageURL)
	return fetchDishes(restaurantID)
}

func fetchDishes(restaurantID string) ([]model.Dish, error) {
	url := fmt.Sprintf(apiURL, restaurantID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var source model.Source
	err = json.NewDecoder(resp.Body).Decode(&source)
	if err != nil {
		return nil, err
	}

	categories := source.Data.Menus[0].MenuCategories

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
