package thirdparty

import (
	"io/ioutil"
	"lunch_helper/thirdparty/model"
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/require"
)

func TestGetDishesByFoodPanda(t *testing.T) {
	source, err := ioutil.ReadFile("foodpanda_test_data.json")
	require.NoError(t, err)

	var expectJson []model.Dish
	err = json.Unmarshal(source, &expectJson)
	require.NoError(t, err)

	got, err := GetDishesByFoodPanda("https://www.foodpanda.com.tw/restaurant/u0gl/ming-zhi-gao-xian-xian-chao-guan?utm_source=google&utm_medium=organic&utm_campaign=google_reserve_place_order_action")
	require.NoError(t, err)
	require.Equal(t, expectJson, got)
}
