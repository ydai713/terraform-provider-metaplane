/*
Implement CRUD for Monitor

GetMonitor: requires monitor_id and connection_id.
  Metaplane does not have a GET method specifically for monitors. Instead, use
  the GET method for the API (list for connection), which requires the
  connection_id. In the response, use the monitor_id to get the monitor.
CreateMonitor: built-in
UpdateMonitor: built-in
DeleteMonitor: there is no delete method for monitors. Instead, use the UPDATE method
  to set the monitor to inactive.
*/
package api

import (
  "fmt"
  "encoding/json"
  "net/http"
  "errors"
)

type Monitor struct {
  MonitorId       string   `json:"id"`
  ConnectionId    string
  Type            string   `json:"type"`
  CronTab         string   `json:"cronTab"`
  IsEnabled       bool     `json:"isEnabled"`
  UpdatedAt       string   `json:"updatedAt"`
  AbsolutePath    string   `json:"absolutePath"`
  CreatedAt       string   `json:"createdAt"`
}

type Monitors struct {
  Data           []Monitor `json:"data"`
}

func (c *Client) GetMonitor(connectionId string, monitorId string) (*Monitor, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/monitors/connection/%s", BaseUrl, connectionId), nil)
	if err != nil {
		return nil, err
	}

  // Send request to the API
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

  // Parse the response
	monitors := Monitors{}
	err = json.Unmarshal(body, &monitors)
	if err != nil {
		return nil, err
	}

	for _, monitor := range monitors.Data {
		if monitor.MonitorId == monitorId {
      monitor.ConnectionId = connectionId
      return &monitor, nil
		}
	}

  return nil, errors.New("Monitor is not found")
}

