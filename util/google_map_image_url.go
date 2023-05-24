package util

import "fmt"

func GetGoogleImageUrl(reference string, apiKey string) string {
	return fmt.Sprintf("https://maps.googleapis.com/maps/api/place/photo?maxwidth=400&photoreference=%s&key=%s", reference, apiKey)
}
