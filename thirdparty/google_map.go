package thirdparty

import (
	"context"
	"fmt"
	"sync"

	"googlemaps.github.io/maps"
)

type PlaceApi interface {
	NearbySearch(nearbySearchRequest *maps.NearbySearchRequest) ([]SearchResult, string, SearchError)
	TextSearch(query *maps.TextSearchRequest) ([]SearchResult, string, SearchError)
	GetApiKey() string
}

type GoogleMapPlaceApi struct {
	client *maps.Client
	apiKey string
}

type SearchResult struct {
	Data   maps.PlacesSearchResult
	Detail maps.PlaceDetailsResult
}

type SearchError struct {
	Err          error
	DetailErrors []error
}

var defaultFields = []maps.PlaceDetailsFieldMask{
	"url",
	"formatted_phone_number",
	"formatted_address",
}

func NewGoogleMapPlaceApi(apiKey string) (GoogleMapPlaceApi, error) {
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return GoogleMapPlaceApi{}, err
	}
	return GoogleMapPlaceApi{client, apiKey}, nil
}

func (m *GoogleMapPlaceApi) NearbySearch(nearbySearchRequest *maps.NearbySearchRequest) ([]SearchResult, string, SearchError) {
	resp, err := m.client.NearbySearch(context.Background(), nearbySearchRequest)
	if err != nil {
		return nil, "", SearchError{Err: err}
	}

	results, errs := m.appendDetail(resp.Results, nearbySearchRequest.Language)
	return results, resp.NextPageToken, SearchError{Err: nil, DetailErrors: errs}
}

func (m *GoogleMapPlaceApi) TextSearch(query *maps.TextSearchRequest) ([]SearchResult, string, SearchError) {
	resp, err := m.client.TextSearch(context.Background(), query)
	if err != nil {
		return nil, "", SearchError{Err: err}
	}

	results, errs := m.appendDetail(resp.Results, query.Language)
	return results, resp.NextPageToken, SearchError{Err: nil, DetailErrors: errs}
}

func (m *GoogleMapPlaceApi) GetApiKey() string {
	return m.apiKey
}

func (m *GoogleMapPlaceApi) appendDetail(list []maps.PlacesSearchResult, language string) ([]SearchResult, []error) {
	var wg sync.WaitGroup
	results := make([]SearchResult, len(list))
	errorList := []error{}
	for i, result := range list {
		wg.Add(1)
		go func(i int, result maps.PlacesSearchResult, errList []error) {
			detailResp, err := m.client.PlaceDetails(context.Background(), &maps.PlaceDetailsRequest{
				PlaceID: result.PlaceID,
				// https://developers.google.com/maps/documentation/places/web-service/details?hl=zh-tw#fields
				// 只保留需要的欄位，不然很花開銷
				Fields:   defaultFields,
				Language: language,
			})
			if err != nil {
				errList = append(errList, fmt.Errorf("Failed to get place details for %s %s: %v, placeId is %s", result.Name, result.Vicinity, err, result.PlaceID))
				detailResp = maps.PlaceDetailsResult{}
			}

			results[i] = SearchResult{
				Data:   result,
				Detail: detailResp,
			}
			wg.Done()
		}(i, result, errorList)
	}

	wg.Wait()
	return results, errorList
}
