/*
Package api provides a client for interacting with the Metaplane API.

url := "https://dev.api.metaplane.dev/v1/monitors/status/monitorId"

*/

package api

import (
  "io"
  "net/http"
  "encoding/json"
  "errors"
  
  "github.com/hashicorp/go-retryablehttp"
)

const BaseUrl = "https://dev.api.metaplane.dev/v1"

type Client struct {
  ApiKey            string
  HTTPClient        *retryablehttp.Client
}

type ErrorResponse struct {
  StatusCode      int         `json:"statusCode"`
  ErrorMessage    string      `json:"errorMessage"`
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

  body, err := io.ReadAll(res.Body)
  if err != nil {
    return nil, err
  }

	statusCode := res.StatusCode
  if statusCode == 200 {
    return body, err
  }

  error_response := ErrorResponse{}
  err = json.Unmarshal(body, &error_response)
  if err != nil {
    return nil, errors.New(string(body))
  }
  return nil, errors.New(error_response.ErrorMessage)
}

func NewClient(apiKey *string) *Client {
  httpClient := retryablehttp.NewClient()
  c := Client{
  	HTTPClient: httpClient,
  	ApiKey: *apiKey,
  }
  return &c
}
