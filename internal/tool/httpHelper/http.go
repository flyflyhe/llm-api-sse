package httpHelper

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

func PostJson(url string, data interface{}, headers map[string]string) (*http.Response, error) {
	client := http.Client{Timeout: time.Second * 10}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}

	if headers != nil {
		for k, v := range headers {
			request.Header.Set(k, v)
		}
	}

	request.Header.Set("Content-Type", "application/json")

	return client.Do(request)
}
