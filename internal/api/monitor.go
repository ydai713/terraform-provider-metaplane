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
	MonitorId       string   `tfsdk:"monitor_id"      json:"id"`
	ConnectionId    string   `tfsdk:"connection_id"`      
  Type            string   `tfsdk:"type"            json:"type"`
  CronTab         string   `tfsdk:"cron_tab"        json:"cronTab"`
  IsEnabled       bool     `tfsdk:"is_enabled"      json:"isEnabled"`
  UpdatedAt       string   `tfsdk:"updated_at"      json:"updatedAt"`
  AbsolutePath    string   `tfsdk:"absolute_path"   json:"absolutePath"`
  CreatedAt       string   `tfsdk:"created_at"      json:"createdAt"`
}

type Monitors struct {
  Data         []Monitor `json:"data"`
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

