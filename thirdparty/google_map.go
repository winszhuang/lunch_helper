package thirdparty

import (
	"context"
	"log"
	"sync"

	"googlemaps.github.io/maps"
)

type PlaceApi interface {
	NearbySearch(nearbySearchRequest *maps.NearbySearchRequest) ([]SearchResult, string, error)
	TextSearch(query *maps.TextSearchRequest) ([]SearchResult, string, error)
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

func NewGoogleMapPlaceApi(apiKey string) (GoogleMapPlaceApi, error) {
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return GoogleMapPlaceApi{}, err
	}
	return GoogleMapPlaceApi{client, apiKey}, nil
}

func (m *GoogleMapPlaceApi) NearbySearch(nearbySearchRequest *maps.NearbySearchRequest) ([]SearchResult, string, error) {
	resp, err := m.client.NearbySearch(context.Background(), nearbySearchRequest)
	if err != nil {
		return nil, "", err
	}

	results := m.appendDetail(resp.Results)
	return results, resp.NextPageToken, nil
}

func (m *GoogleMapPlaceApi) TextSearch(query *maps.TextSearchRequest) ([]SearchResult, string, error) {
	resp, err := m.client.TextSearch(context.Background(), query)
	if err != nil {
		return nil, "", err
	}

	results := m.appendDetail(resp.Results)
	return results, resp.NextPageToken, nil
}

func (m *GoogleMapPlaceApi) GetApiKey() string {
	return m.apiKey
}

func (m *GoogleMapPlaceApi) appendDetail(list []maps.PlacesSearchResult) []SearchResult {
	var wg sync.WaitGroup
	results := make([]SearchResult, len(list))
	for i, result := range list {
		wg.Add(1)
		go func(i int, result maps.PlacesSearchResult) {
			detailResp, err := m.client.PlaceDetails(context.Background(), &maps.PlaceDetailsRequest{
				PlaceID: result.PlaceID,
			})
			if err != nil {
				log.Printf("Failed to get place details for %s %s: %v, placeId is %s", result.Name, result.Vicinity, err, result.PlaceID)
				detailResp = maps.PlaceDetailsResult{}
			}

			results[i] = SearchResult{
				Data:   result,
				Detail: detailResp,
			}
			wg.Done()
		}(i, result)
	}

	wg.Wait()
	return results
}
