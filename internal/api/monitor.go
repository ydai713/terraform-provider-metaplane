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
  "context"
  "fmt"
  "encoding/json"
  "net/http"
  "errors"
  "strings"
)

type Duration struct {
  Days              *int64              `json:"days,omitempty"`
  Hours             *int64              `json:"hours,omitempty"`
  Minutes           *int64              `json:"minutes,omitempty"`
}

type IncrementalClause struct {
  ColumnName        *string             `json:"columnName,omitempty"`
  Duration          *Duration          `json:"duration,omitempty"`
}

type Config struct {
  CustomSql         *string            `json:"customSql,omitempty"`
  IncrementalClause *IncrementalClause `json:"incrementalClause,omitempty"`
  CustomWhereClause *string            `json:"customWhereClause,omitempty"`
}

type Monitor struct {
  ID                string             `json:"id"`
  Type              string             `json:"type"`
  CronTab           string             `json:"cronTab"`
  IsEnabled         bool               `json:"isEnabled"`
  Config            *Config            `json:"config"`
  CreatedAt         string             `json:"createdAt"`
  UpdatedAt         string             `json:"updatedAt"`
  AbsolutePath      string             `json:"absolutePath"`
  ConnectionId      string             `json:"connectionId"`
  EntityType        string             `json:"entityType"`
}

type NewMonitor struct {
  ConnectionId        string           `json:"connectionId"`
  Type                string           `json:"type"`
  EntityType          string           `json:"entityType"`
  CronTab             string           `json:"cronTab"`
  AbsolutePath        string           `json:"absolutePathString"`
  Config              Config           `json:"config,omitempty"`
}

type Monitors struct {
  Data                []Monitor        `json:"data"`
}

type UpdateMonitor struct {
  MonitorId           string
  CronTab             string           `json:"cronTab,omitempty"`
  IsEnabled           bool             `json:"isEnabled,omitempty"`
  Config              Config           `json:"config,omitempty"`
}

func (c *Client) GetMonitor(monitorId string) (*Monitor, error) {
  req, err := http.NewRequest("GET", fmt.Sprintf("%s/monitors/%s", BaseUrl, monitorId), nil)
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

func (c *Client) CreateMonitor(ctx context.Context, newMonitor NewMonitor) (*Monitor, error) {
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

      return c.UpdateMonitor(ctx, updateMonitor)
    }
   	return nil, err
  }
  
  monitor := Monitor{}
  err = json.Unmarshal(body, &monitor)
  if err != nil {
   	return nil, err
  }

  return &monitor, nil
}

func (c *Client) UpdateMonitor(ctx context.Context, updateMonitor UpdateMonitor) (*Monitor, error) {
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
   	return nil, errors.New(string(rb))
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
      return monitor.ID, nil
  	}
  }
  
  return "", errors.New("Monitor is not found")
}
