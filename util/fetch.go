package util

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func Fetch(sourceUrl string) ([]byte, error) {
	proxyURL, _ := url.Parse("http://180.183.120.117:8080")
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	req, err := http.NewRequest("GET", sourceUrl, nil)
	if err != nil {
		return nil, err
	}

	// #TODO 這邊繼續增加
	req.Header.Set("Cookie", "1P_JAR=2023-06-16-06; AEC=AUEFqZcH4PWA3VyJGdRuh4dITMGcrIk1hKCd5Wk-SMOgFsnxKHiV2Dkg6to; NID=511=IW_WIa5Ls8Q9eCtEts7ZAMAs64RNqELrBW6_oxbYpkgkiMwBQx1wOC88AjYQ_kh3ma-vlRrmOmyVjlGJ2eEYTKB9TX7sAAGYPYAYuvf-RFGbHDnoIBTfB_jkmS4S4qKltj6SRn8U1m-YkzOvf4XvQtdsUzg-isXd_9FrWno3bUE")
	req.Header.Set("Accept", "*/*")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-TW,zh;q=0.9")
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("User-Agent", "PostmanRuntime/7.32.3")
	req.Header.Set("Connection", "keep-alive")
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	log.Println("---")
	log.Println(response.Header.Get("Content-Type"))
	log.Println("---")
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func FetchBody(sourceUrl string) (io.ReadCloser, error) {
	response, err := http.Get(sourceUrl)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}
