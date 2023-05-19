package api

import (
  "fmt"
  "encoding/json"
  "net/http"
)

type MonitorStatus struct {
  Available    bool      `json:"available"`
  Type         string    `json:"type"`
  Result       float64   `json:"result"`
  LowerBound   float64   `json:"lowerBound"`
  UpperBound   float64   `json:"upperBound"`
  Predicted    float64   `json:"predicted"`
  Passed       bool      `json:"passed"`
  CreatedAt    string    `json:"createdAt"`
}

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

