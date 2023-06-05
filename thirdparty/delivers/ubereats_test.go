package thirdparty

import (
	"encoding/json"
	"io/ioutil"
	"lunch_helper/thirdparty/model"
	"testing"

	"github.com/stretchr/testify/require"
)

type Test struct {
	Context string `json:"@context"`
	Type    string `json:"type"`
}

func TestGetDishesByUberEats(t *testing.T) {
	source, err := ioutil.ReadFile("ubereats_test_data.json")
	require.NoError(t, err)

	var expectJson []model.Dish
	err = json.Unmarshal(source, &expectJson)
	require.NoError(t, err)

	const pageUrl = "https://www.ubereats.com/tw/store/%E5%BF%A0%E5%91%88%E6%BA%AB%E5%B7%9E%E5%A4%A7%E9%A4%9B%E9%A3%A9/MIE0wXk4TYOxf5a_4ZJD4g?utm_campaign=place-action-link&utm_medium=organic&utm_source=google"
	got, err := GetDishesByUberEats(pageUrl)
	require.NoError(t, err)
	require.Equal(t, expectJson, got)
}
