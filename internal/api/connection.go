/*
Implement CRUD for Connection

GetConnection: requires connection_id
  Metaplane does not have a GET method specifically for connections. Instead,
  use the GET method for the API (list all). In the response, use the 
  name to get the desired connection.
*/
package api

import (
	"fmt"
	"encoding/json"
	"net/http"
  "errors"
)

type Connection struct {
	Name            string     `json:"name"`
	ConnectionId    string     `json:"id"`
  Type            string     `json:"type"`
  IsEnabled       bool       `json:"isEnabled"`
  UpdatedAt       string     `json:"updatedAt"`
  CreatedAt       string     `json:"createdAt"`
  Status          string     `json:"status"`
}

func (c *Client) GetConnection(name string) (*Connection, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/connections", BaseUrl), nil)
	if err != nil {
		return nil, err
	}

  // Send request to the API
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

  // Parse the response
  var connections []Connection
	err = json.Unmarshal(body, &connections)
	if err != nil {
		return nil, err
	}

	for _, connection := range connections {
		if connection.Name == name {
      return &connection, nil
		}
	}

  return nil, errors.New("Connection is not found")
}

