package util

import (
	"io"
	"io/ioutil"
	"net/http"
)

func Fetch(url string) (io.ReadCloser, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-TW,zh;q=0.9")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("Connection", "keep-alive")

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

func FetchBytes(url string) ([]byte, error) {
	reader, err := Fetch(url)
	if err != nil {
		return nil, err
	}

	defer reader.Close()

	return ioutil.ReadAll(reader)
}
