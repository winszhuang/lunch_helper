package food_deliver

import (
	"encoding/json"
	"io/ioutil"
	"lunch_helper/food_deliver/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFoodPandaDishesCrawler_ParseSource(t *testing.T) {
	crawler := &FoodPandaDishesCrawler{}

	url, err := crawler.ParseSource("https://maps.google.com/?cid=9498528971100991227")
	require.NoError(t, err)

	containText := "https://www.foodpanda.com.tw/restaurant/g1hl/"
	require.Contains(t, url, containText)
}

func TestFoodPandaDishesCrawler_GetDishes(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		source, err := ioutil.ReadFile("foodpanda_test_data.json")
		require.NoError(t, err)

		var expectJson []model.Dish
		err = json.Unmarshal(source, &expectJson)
		require.NoError(t, err)

		crawler := &FoodPandaDishesCrawler{}
		got, err := crawler.GetDishes("https://www.foodpanda.com.tw/restaurant/u0gl/ming-zhi-gao-xian-xian-chao-guan?utm_source=google&utm_medium=organic&utm_campaign=google_reserve_place_order_action")
		require.NoError(t, err)
		require.Equal(t, expectJson, got)
	})

	t.Run("not a foodpanda url", func(t *testing.T) {
		crawler := &FoodPandaDishesCrawler{}
		got, err := crawler.GetDishes("https://www.ubereats.com/tw/store/大煙囪烤肉便當/PI7R2L-ZS_eLIQ7Y0Dndaw?utm_campaign\\\\u003dplace-action-link\\\\u0026utm_medium\\\\u003dorganic\\\\u0026utm_source\\\\u003dgoogle\\")
		require.Error(t, err)
		require.Nil(t, got)
	})
}
