package todoist

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type Client struct {
    APIKey  string
    BaseURL string
}

func NewClient(apiKey string) *Client {
    return &Client{
        APIKey:  apiKey,
        BaseURL: "https://api.todoist.com/rest/v2",
    }
}

func (c *Client) request(method, endpoint string, body any) (*http.Response, error) {
    url := c.BaseURL + endpoint
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

    req, err := http.NewRequest(method, url, reqBody)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bearer "+c.APIKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    return client.Do(req)
}

