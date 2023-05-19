package api

import (
	"fmt"
	"encoding/json"
	"net/http"
)

type MonitorStatus struct {
	Available bool `json:"available"`
  Type       string    `json:"type"`
  Result     float64   `json:"result"`
  LowerBound float64   `json:"lowerBound"`
  UpperBound float64   `json:"upperBound"`
  Predicted  float64   `json:"predicted"`
  Passed     bool      `json:"passed"`
  CreatedAt  string    `json:"createdAt"`
}

// {
//   "type":"decimal",
//   "result":2008.001,
//   "lowerBound":0.0,
//   "upperBound":803.8200528225457,
//   "predicted":358.9007652664038,
//   "passed":false,
//   "createdAt":"2023-05-18T00:01:51.552117Z"
// }
func (c *Client) GetMonitorStatus(monitor_id string) (*Monitor, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/monitors/status/%s", BaseUrl, monitor_id), nil)
	if err != nil {
		return nil, err
	}

  // Send request to the API
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

  // Parse the response
	monitor := Monitor{}
	err = json.Unmarshal(body, &monitor)
	if err != nil {
		return nil, err
	}
	return &monitor, nil
}

