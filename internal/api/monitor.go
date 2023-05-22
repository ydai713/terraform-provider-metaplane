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
  "strings"
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

type NewMonitor struct {
  ConnectionId    string   `json:"connectionId"`
  Type            string   `json:"type"`
  EntityType      string   `json:"entityType"`
  CronTab         string   `json:"cronTab"`
  AbsolutePath    string   `json:"absolutePathString"`
}

type UpdateMonitor struct {
  MonitorId       string
  CronTab         string   `json:"cronTab"`
  IsEnabled       bool     `json:"isEnabled"`
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

func (c *Client) CreateMonitor(newMonitor NewMonitor) (*Monitor, error) {
  rb, err := json.Marshal(newMonitor)
  if err != nil {
 	  return nil, err
  }
  
  req, err := http.NewRequest("POST", fmt.Sprintf("%s/monitors", BaseUrl), strings.NewReader(string(rb)))
  if err != nil {
   	return nil, err
  }
  
  body, err := c.doRequest(req)
  if err != nil {
    if strings.Contains(err.Error(), "already exists") {
      // get the monitor id that has the same type and absolute path
      monitorId, err := c.getDuplicatedMonitor(newMonitor.ConnectionId, newMonitor.AbsolutePath, newMonitor.Type)
      if err != nil {
        return nil, err
      }

      // update isEnabled to true
      updateMonitor := UpdateMonitor {
        MonitorId: monitorId,
        IsEnabled: true,
        CronTab:   newMonitor.CronTab,
      }

      return c.UpdateMonitor(updateMonitor)
    }
   	return nil, err
  }
  
  monitor := Monitor{}
  err = json.Unmarshal(body, &monitor)
  if err != nil {
   	return nil, err
  }

  monitor.ConnectionId = newMonitor.ConnectionId
  return &monitor, nil
}

func (c *Client) UpdateMonitor(updateMonitor UpdateMonitor) (*Monitor, error) {
  monitorId := updateMonitor.MonitorId

  rb, err := json.Marshal(updateMonitor)
  if err != nil {
   	return nil, errors.New("1")
  }
  
  req, err := http.NewRequest("POST", fmt.Sprintf("%s/monitors/%s", BaseUrl, monitorId), strings.NewReader(string(rb)))
  if err != nil {
   	return nil, errors.New("2")
  }
  
  body, err := c.doRequest(req)
  if err != nil {
   	return nil, errors.New("3")
  }
  
  monitor := Monitor{}
  err = json.Unmarshal(body, &monitor)
  if err != nil {
   	return nil, errors.New("4")
  }

  return &monitor, nil
}

func (c *Client) getDuplicatedMonitor(connectionId string, absolutePath string, monitorType string) (string, error) {
  req, err := http.NewRequest("GET", fmt.Sprintf("%s/monitors/connection/%s?includeDisabled=true", BaseUrl, connectionId), nil)
  if err != nil {
  	return "", err
  }
  
  // Send request to the API
  body, err := c.doRequest(req)
  if err != nil {
  	return "", err
  }
  
  // Parse the response
  monitors := Monitors{}
  err = json.Unmarshal(body, &monitors)
  if err != nil {
  	return "", err
  }
  
  for _, monitor := range monitors.Data {
  	if strings.ToUpper(monitor.Type) == strings.ToUpper(monitorType) && strings.ToUpper(monitor.AbsolutePath) == strings.ToUpper(absolutePath) {
      return monitor.MonitorId, nil
  	}
  }
  
  return "", errors.New("Monitor is not found")
}
