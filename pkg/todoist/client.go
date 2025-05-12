package todoist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	APIKey  string
	BaseURL string
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: "https://api.todoist.com/api/v1",
	}
}

func (c *Client) request(method, endpoint string, body any, query any) (*http.Response, error) {
	requestURL, err := url.Parse(c.BaseURL + endpoint)
	if err != nil {
		return nil, err
	}

	if query != nil {
		requestURL, err = addQueryParams(requestURL, query)
		if err != nil {
			return nil, err
		}
	}

	var reqBody *bytes.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(jsonData)
	} else {
		reqBody = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(method, requestURL.String(), reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}

func addQueryParams(requestURL *url.URL, query any) (*url.URL, error) {
	q := requestURL.Query()
	queryData, err := json.Marshal(query)
	if err != nil {
		fmt.Println("Error marshalling query data:", err)
		return nil, err
	}

	var queryMap map[string]any
	err = json.Unmarshal(queryData, &queryMap)
	if err != nil {
		fmt.Println("Error unmarshalling query data:", err)
		return nil, err
	}

	for key, value := range queryMap {
		q.Set(key, fmt.Sprintf("%v", value))
	}

	requestURL.RawQuery = q.Encode()
	return requestURL, nil
}
