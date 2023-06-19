package service

import (
	"encoding/json"
	"io/ioutil"
	"lunch_helper/cache"
	"lunch_helper/constant"
	"lunch_helper/thirdparty"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var testSearchArgs = &constant.SearchArgs{
	Lat:       24.172746869207582,
	Lng:       120.67216126454902,
	Radius:    100,
	PageIndex: 1,
	PageSize:  10,
	Text:      "牛排",
}

func TestSearchService_Search(t *testing.T) {
	// 讀取假資料，讓TextSearch都回傳此
	result := initFakeSearchResultData(t)
	placeApi := MockPlaceApi{}
	// 預計就只有13筆資料，不管帶什麼參數指定什麼nextPageToken，都只有13筆
	placeApi.On("TextSearch", mock.Anything).Return(result, "", thirdparty.SearchError{Err: nil, DetailErrors: []error{}})
	placeApi.On("GetApiKey", mock.Anything).Return("your google api key")

	searchService := NewSearchService(
		&cache.NearByRestaurantCache{},
		&placeApi,
	)

	// 測試獲取第一頁
	list, searchErr := searchService.Search(testSearchArgs)
	require.NoError(t, searchErr.Err)
	require.Len(t, list, 10)

	// 測試獲取第二頁
	nextPageSearchArgs := testSearchArgs
	nextPageSearchArgs.PageIndex = 2
	nextList, nextSearchErr := searchService.Search(nextPageSearchArgs)
	require.NoError(t, nextSearchErr.Err)
	require.Len(t, nextList, 3)

	// 測試獲取第三頁
	thirdPageSearchArgs := testSearchArgs
	thirdPageSearchArgs.PageIndex = 3
	thirdList, thirdSearchErr := searchService.Search(thirdPageSearchArgs)
	require.NoError(t, thirdSearchErr.Err)
	require.Len(t, thirdList, 0)
}

func initFakeSearchResultData(t *testing.T) []thirdparty.SearchResult {
	source, err := ioutil.ReadFile("text_search_data.json")
	require.NoError(t, err)

	var result []thirdparty.SearchResult
	err = json.Unmarshal(source, &result)
	require.NoError(t, err)

	require.Len(t, result, 13)

	return result
}
