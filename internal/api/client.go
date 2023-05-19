/*
Package api provides a client for interacting with the Metaplane API.

url := "https://dev.api.metaplane.dev/v1/monitors/status/monitorId"

*/

package api

import (
  "io"
  "net/http"
  "errors"
  
  "github.com/hashicorp/go-retryablehttp"
)

const BaseUrl = "https://dev.api.metaplane.dev/v1"

type Client struct {
  ApiKey            string
	HTTPClient        *retryablehttp.Client
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("accept", "application/json")
  req.Header.Set("content-type", "application/json")
  req.Header.Set("Authorization", c.ApiKey)

	retryableReq, err := retryablehttp.FromRequest(req)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(retryableReq)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	switch statusCode := res.StatusCode; statusCode {
	case 200:
    body, err := io.ReadAll(res.Body)
    if err != nil {
      return nil, err
    }
    return body, err
	case 400:
    return nil, errors.New("Resource not found")
	case 401:
    return nil, errors.New("Unauthorized")
	case 403:
    return nil, errors.New("Forbidden")
	case 409:
    return nil, errors.New("Try again later")
	case 422:
    return nil, errors.New("Unprocessable entity or invalid request")
	case 500:
    return nil, errors.New("Invalid consent url")
	default:
    return nil, errors.New("Unknown error")
	}
}

func NewClient(apiKey *string) *Client {
  httpClient := retryablehttp.NewClient()
	c := Client{
		HTTPClient: httpClient,
		ApiKey: *apiKey,
	}

	return &c
}
