package qdb

import (
	"net/http"
)

func callGet(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return httpClient.Do(req)
}

func callGetWithAuth(url string, username string, password string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if username != "" {
		req.SetBasicAuth(username, password)
	}
	if err != nil {
		return nil, err
	}
	return httpClient.Do(req)
}
