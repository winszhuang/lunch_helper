package util

import (
	"net/url"
	"strings"
)

func BindUrl(apiBaseUrl, endPoint string) (string, error) {
	u, err := url.Parse(apiBaseUrl)
	if err != nil {
		return "", err
	}

	baseURL := u.Scheme + "://" + u.Host
	endPoint = strings.TrimLeft(endPoint, "/")
	baseURL += "/" + endPoint

	return baseURL, nil
}
