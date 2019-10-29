package sumologic

import (
	"encoding/json"
	"fmt"
	"log"
)

//READ
func (s *Client) GetPermission(id string) (*ContentPermissions, error) {
	log.Println("####Begin GetPermission####")

	url := fmt.Sprintf("v2/content/%s/permissions", id)
	log.Printf("Permission read url: %s", url)

	//Execute the permission read request
	rawResponse, _, err := s.Get(url)

	//If there was an error, exit here and return it
	if err != nil {
		return nil, err
	}

	//Parse a Permission struct from the response
	var response ContentPermissionResponse
	err = json.Unmarshal(rawResponse, &response)

	//Exit here if there was an error parsing the json
	if err != nil {
		return nil, err
	}

	var contentPermission ContentPermissions
	contentPermission.Permissions = response.Explicit

	contentPermission.ID = id

	log.Println("####End GetPermission####")
	return &contentPermission, nil
}

//DELETE
func (s *Client) DeletePermissions(contentPermission ContentPermissions) error {
	log.Println("####Begin DeletePermissions####")

	log.Printf("Deleting Permissions: %s", contentPermission.Permissions)

	url := fmt.Sprintf("v2/content/%s/permissions/remove", contentPermission.ID)
	log.Printf("Permission delete url: %s", url)

	//since notifications are unused in the terraform implementation, populate dummy values here
	contentPermission.NotificationMessage = "Intentionally left blank."
	contentPermission.NotifyRecipients = false

	_, err := s.Put(url, contentPermission)

	log.Println("####End DeletePermissions####")
	return err
}

//CREATE
func (s *Client) AddPermissions(contentPermission ContentPermissions) (string, error) {
	log.Println("####Begin AddPermissions####")

	//since notifications are unused in the terraform implementation, populate dummy values here
	contentPermission.NotificationMessage = "Intentionally left blank."
	contentPermission.NotifyRecipients = false

	//construct the url
	url := fmt.Sprintf("v2/content/%s/permissions/add", contentPermission.ID)
	log.Printf("Create permission url: %s", url)

	//add the permission -- we don't need the response, so drop it
	_, err := s.Put(url, contentPermission)

	//Exit if there was an error during the request
	if err != nil {
		return "", err
	}

	log.Printf("New permission ID is: %s", contentPermission.ID)
	return contentPermission.ID, nil
}

//UPDATE

func (s *Client) UpdatePermission(update ContentPermissions) error {
	log.Println("####Begin permission update####")

	//Fetch the current list of permissions for this content object
	original, err := s.GetPermission(update.ID)

	//Error
	if err != nil {
		return err
	}

	log.Printf("original variable: %s", original.Permissions)
	log.Printf("update variable: %s", update.Permissions)

	add := permissionDifference(update.Permissions, original.Permissions)
	delete := permissionDifference(original.Permissions, update.Permissions)

	log.Printf("add: %s", add)
	log.Printf("delete: %s", delete)

	// make the calls to update the permission list
	var wrapper ContentPermissions
	wrapper.Permissions = add
	wrapper.ID = update.ID

	// Add
	if len(wrapper.Permissions) != 0 {
		_, err = s.AddPermissions(wrapper)
	}

	//Error
	if err != nil {
		return err
	}

	// Delete
	wrapper.Permissions = delete
	if len(wrapper.Permissions) != 0 {
		err = s.DeletePermissions(wrapper)
	}

	log.Println("####End permission update####")
	return err
}

func permissionDifference(a []Permission, b []Permission) []Permission {
	var diff []Permission
	m := make(map[Permission]bool)
	for _, permission := range b {
		m[permission] = true
	}

	for _, permission := range a {
		if _, ok := m[permission]; !ok {
			diff = append(diff, permission)
		}
	}
	return diff
}

type Permission struct {
	Name       string `json:"permissionName"`
	SourceType string `json:"sourceType"`
	SourceId   string `json:"sourceId"`
	ContentId  string `json:"contentId"`
}

type ContentPermissions struct {
	ID                  string
	NotifyRecipients    bool         `json:"notifyRecipients"`
	NotificationMessage string       `json:"notificationMessage"`
	Permissions         []Permission `json:"contentPermissionAssignments"`
}

type ContentPermissionResponse struct {
	Explicit []Permission `json:"explicitPermissions"`
	Implicit []Permission `json:"implicitPermissions"`
}
