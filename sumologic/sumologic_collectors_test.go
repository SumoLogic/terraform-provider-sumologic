package sumologic

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

type mockHttpClient struct {
	response *http.Response
}

func (c *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return c.response, nil
}

func newTestClient(response *http.Response) *Client {
	httpClient := &mockHttpClient{
		response: response,
	}
	client := Client{
		AccessID:    "abcd",
		AccessKey:   "ef12",
		Environment: "us2",
		httpClient:  httpClient,
	}

	client.BaseURL, _ = url.Parse(endpoints[client.Environment])
	return &client
}

func TestGetCollectorFound(t *testing.T) {
	body := []byte(`{
		"collector": {
			"id": 111222333,
			"name": "collector1",
			"category": "foo/bar",
			"timeZone": "Etc/UTC",
			"links": [],
			"collectorType": "Hosted",
			"collectorVersion":"",
    		"lastSeenAlive":1542964933842,
    		"alive":true
		}
	}`)
	response := &http.Response{
		Status:     http.StatusText(200),
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}
	client := newTestClient(response)
	collector, err := client.GetCollector(1234)
	if err != nil {
		t.Fatalf("Expected GetCollector to succeed, received: %s", err)
	}
	var expectedCollector CollectorResponse
	json.Unmarshal(body, &expectedCollector)
	if !reflect.DeepEqual(*collector, expectedCollector.Collector) {
		t.Errorf("Expected GetCollector to return serialised body, instead of %v", *collector)
	}
}

func TestGetCollectorNotFound(t *testing.T) {
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
	collectorResponse, err := client.GetCollector(1234)
	if err != nil {
		t.Fatalf("Expected GetCollector to succeed, received: %s", err)
	}
	if collectorResponse != nil {
		t.Errorf("Expected GetCollector to return an empty response, instead of %v", *collectorResponse)
	}
}

func TestGetCollectorUnauthorized(t *testing.T) {
	body := []byte(`<html></html>`)
	response := &http.Response{
		Status:     http.StatusText(401),
		StatusCode: 401,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}
	client := newTestClient(response)
	_, err := client.GetCollector(1234)
	if err == nil {
		t.Error("Expected GetCollector to fail, but it succeeded")
	}
}
