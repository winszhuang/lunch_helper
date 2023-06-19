package service

import (
	"lunch_helper/thirdparty"

	"github.com/stretchr/testify/mock"
	"googlemaps.github.io/maps"
)

// MockPlaceApi 是 PlaceApi 的 Mock 物件
type MockPlaceApi struct {
	mock.Mock
}

func (m *MockPlaceApi) NearbySearch(nearbySearchRequest *maps.NearbySearchRequest) ([]thirdparty.SearchResult, string, thirdparty.SearchError) {
	// TODO: 根據需要進行測試邏輯的實現
	args := m.Called(nearbySearchRequest)
	return args.Get(0).([]thirdparty.SearchResult), args.String(1), args.Get(2).(thirdparty.SearchError)
}

// TextSearch 是 MockPlaceApi 的 TextSearch 方法實現
func (m *MockPlaceApi) TextSearch(query *maps.TextSearchRequest) ([]thirdparty.SearchResult, string, thirdparty.SearchError) {
	// TODO: 根據需要進行測試邏輯的實現
	args := m.Called(query)
	return args.Get(0).([]thirdparty.SearchResult), args.String(1), args.Get(2).(thirdparty.SearchError)
}

// GetApiKey 是 MockPlaceApi 的 GetApiKey 方法實現
func (m *MockPlaceApi) GetApiKey() string {
	// TODO: 根據需要進行測試邏輯的實現
	args := m.Called()
	return args.String(0)
}
