package sumologic

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// GetConnection returns connection information for given id
func (s *Client) GetConnection(id string) (*Connection, error) {
	log.Println("#### Begin GetConnection ####")

	url := fmt.Sprintf("v1/connections/%s", id)
	log.Printf("connection read url: %s", url)

	rawConnection, _, err := s.Get(url, false)
	if err != nil {
		log.Printf("SSAIN: The err: %s", err.Error())
		if strings.Contains(err.Error(), "Connection with given ID does not exist.") {
			return nil, nil
		}
		return nil, err
	}

	if rawConnection == nil {
		return nil, nil
	}

	var connection Connection
	err = json.Unmarshal(rawConnection, &connection)
	if err != nil {
		return nil, err
	}

	log.Println("#### End GetConnection ####")
	return &connection, nil
}

// DeleteConnection deletes connection linked to id
func (s *Client) DeleteConnection(id string, connectionType string) error {
	log.Println("#### Begin DeleteConnection ####")

	log.Printf("Deleting connection Id: %s, Type %s", id, connectionType)

	url := fmt.Sprintf("v1/connections/%s?type=%s", id, connectionType)
	log.Printf("connection delete url: %s", url)

	// Execute the connection delete request
	_, err := s.Delete(url)
	log.Println("#### End DeleteConnection ####")
	return err
}

// CreateConnection creates connection with given params
func (s *Client) CreateConnection(connection Connection) (string, error) {
	log.Println("#### Begin CreateConnection ####")

	url := "v1/connections"
	log.Printf("Create connection url: %s", url)

	connection.Type = convertConToDef(connection.Type)
	responseData, err := s.Post(url, connection, false)
	if err != nil {
		return "", err
	}

	var connectionResponse Connection
	err = json.Unmarshal(responseData, &connectionResponse)
	if err != nil {
		return "", err
	}

	log.Printf("New connection ID is: %s", connectionResponse.ID)
	return connectionResponse.ID, nil
}

// UpdateConnection updates connection with given params
func (s *Client) UpdateConnection(connection Connection) error {
	log.Println("#### Begin connection update ####")

	url := fmt.Sprintf("v1/connections/%s", connection.ID)
	log.Printf("Update connection job status url: %s", url)

	connection.Type = convertConToDef(connection.Type)
	_, err := s.Put(url, connection, false)

	log.Println("#### End connection update ####")
	return err
}

// Post and Put request use *Definitions as types. The API returns *Connections as types.
func convertConToDef(in string) string {
	switch in {
	case "WebhookConnection":
		return "WebhookDefinition"
	}
	return ""
}
