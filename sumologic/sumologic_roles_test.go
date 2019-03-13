package sumologic

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestGetRoleFound(t *testing.T) {
	body := []byte(`{
		"id": "1112223334445556",
		"name": "role1",
		"description": "description",
		"filterPredicate": "filter",
		"users": [
		  "AAABBBCCCDDDEEEF"
		],
		"capabilities": [
		  "viewCollectors"
		]
	}`)
	response := &http.Response{
		Status:     http.StatusText(200),
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}
	client := newTestClient(response)
	role, err := client.GetRole("1112223334445556")
	if err != nil {
		t.Fatalf("Expected GetRole to succeed, received: %s", err)
	}
	var expectedRole Role
	json.Unmarshal(body, &expectedRole)
	if !reflect.DeepEqual(*role, expectedRole) {
		t.Errorf("Expected GetRole to return serialised body, instead of %v", *role)
	}
}

func TestGetRoleNotFound(t *testing.T) {
	body := []byte(`{
		"status": 404,
		"code" : "roles.role.invalid",
		"message" : "The specified role ID is invalid."
	`)
	response := &http.Response{
		Status:     http.StatusText(404),
		StatusCode: 404,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}
	client := newTestClient(response)
	roleResponse, err := client.GetRole("1234")
	if err != nil {
		t.Fatalf("Expected GetRole to succeed, received: %s", err)
	}
	if roleResponse != nil {
		t.Errorf("Expected GetRole to return an empty response, instead of %v", *roleResponse)
	}
}

func TestGetRoleUnauthorized(t *testing.T) {
	body := []byte(`<html></html>`)
	response := &http.Response{
		Status:     http.StatusText(401),
		StatusCode: 401,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}
	client := newTestClient(response)
	_, err := client.GetRole("1234")
	if err == nil {
		t.Error("Expected GetRole to fail, but it succeeded")
	}
}
