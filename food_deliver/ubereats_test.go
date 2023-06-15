package food_deliver

import (
	"encoding/json"
	"io/ioutil"
	"lunch_helper/food_deliver/model"
	"testing"

	"github.com/stretchr/testify/require"
)

type Test struct {
	Context string `json:"@context"`
	Type    string `json:"type"`
}

func TestUberEatsDishesCrawler_ParseSource(t *testing.T) {
	crawler := &UberEatsDishesCrawler{}
	url, err := crawler.ParseSource("https://maps.google.com/?cid=11815448699177978476")
	require.NoError(t, err)

	expectUrl := `https://www.ubereats.com/tw/store/大煙囪烤肉便當/PI7R2L-ZS_eLIQ7Y0Dndaw?utm_campaign\\u003dplace-action-link\\u0026utm_medium\\u003dorganic\\u0026utm_source\\u003dgoogle\`
	require.Equal(t, expectUrl, url)
}

func TestUberEatsDishesCrawler_GetDishes(t *testing.T) {
	source, err := ioutil.ReadFile("ubereats_test_data.json")
	require.NoError(t, err)

	var expectJson []model.Dish
	err = json.Unmarshal(source, &expectJson)
	require.NoError(t, err)

	const pageUrl = "https://www.ubereats.com/tw/store/%E5%BF%A0%E5%91%88%E6%BA%AB%E5%B7%9E%E5%A4%A7%E9%A4%9B%E9%A3%A9/MIE0wXk4TYOxf5a_4ZJD4g?utm_campaign=place-action-link&utm_medium=organic&utm_source=google"
	crawler := &UberEatsDishesCrawler{}
	got, err := crawler.GetDishes(pageUrl)
	require.NoError(t, err)
	require.Equal(t, expectJson, got)
}
