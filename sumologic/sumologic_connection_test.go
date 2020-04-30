package sumologic

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestGetConnectionFound(t *testing.T) {
	body := []byte(`{
	"id": "1738",
  "type": "WebhookConnection",
  "name": "sumologic connection test",
  "description": "integration",
  "url": "https://connection.endpoint",
  "headers": [
    {
      "name": "X-Header",
      "value": "Some-Header"
    }
  ],
  "defaultPayload": "{\"client\":\"Sumo Logic\"}",
  "webhookType": "Webhook"
}`)
	response := &http.Response{
		Status:     http.StatusText(200),
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}
	client := newTestClient(response)
	connection, err := client.GetConnection("1738")
	if err != nil {
		t.Fatalf("Expected GetConnection to succeed, received: %s", err)
	}
	var expectedConnection Connection
	json.Unmarshal(body, &expectedConnection)
	if !reflect.DeepEqual(*connection, expectedConnection) {
		t.Errorf("Expected GetConnection to return serialised body, instead of %v", *connection)
	}
}

func TestGetConnectionNotFound(t *testing.T) {
	body := []byte(`{
		"status": 404,
		"code" : "collectors.collector.invalid",
		"message" : "The specified collector ID is invalid."
	`)
	response := &http.Response{
		Status:     http.StatusText(404),
		StatusCode: 404,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}
	client := newTestClient(response)
	collectorResponse, err := client.GetConnection("1738")
	if err != nil {
		t.Fatalf("Expected GetConnection to succeed, received: %s", err)
	}
	if collectorResponse != nil {
		t.Errorf("Expected GetConnection to return an empty response, instead of %v", *collectorResponse)
	}
}

func TestGetConnectionUnauthorized(t *testing.T) {
	body := []byte(`<html></html>`)
	response := &http.Response{
		Status:     http.StatusText(401),
		StatusCode: 401,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}
	client := newTestClient(response)
	_, err := client.GetConnection("1738")
	if err == nil {
		t.Error("Expected GetConnection to fail, but it succeeded")
	}
}
