package sumologic

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"
)

const publishedPipelineJSON = `{
  "id": "00000000HX546",
  "customerId": "00000000HX123",
  "name": "prod-log-pipeline",
  "description": "Handles prod log parsing.",
  "pipelineType": "route_based",
  "isEnabled": true,
  "routeExpression": "_sourceCategory=prod/*",
  "order": 2,
  "version": 7,
  "createdAt": "2026-09-01T10:00:00Z",
  "modifiedAt": "2026-09-12T11:30:00Z",
  "createdByUserId": "00000000HX546",
  "modifiedByUserId": "00000000HX999",
  "state": "published",
  "nodes": [
    {
      "id": "NODE1",
      "name": "Source",
      "nodeType": "source",
      "outputs": [{"target": "Triage"}]
    },
    {
      "id": "NODE2",
      "name": "Triage",
      "nodeType": "router",
      "outputs": [
        {"target": "Security Processing", "condition": "event_category=security", "order": 1},
        {"target": "Log Normalization"}
      ]
    },
    {
      "id": "NODE3",
      "name": "Log Normalization",
      "nodeType": "processing_group",
      "filterExpression": "env=prod",
      "isEnabled": true,
      "processors": [
        {
          "id": "PROC1",
          "name": "json-parser-1",
          "processorType": "JSON_PARSE",
          "order": 1,
          "isEnabled": true,
          "config": {"fieldContainingJson": "_raw"}
        },
        {
          "id": "PROC2",
          "name": "add-env",
          "processorType": "ADD_FIELDS",
          "order": 2,
          "isEnabled": false,
          "config": {"fields": [{"fieldName": "env", "value": "prod"}]}
        }
      ],
      "outputs": [{"target": "Sumo Logic"}]
    },
    {
      "id": "NODE4",
      "name": "Sumo Logic",
      "nodeType": "destination"
    }
  ]
}`

func TestDataPipelineUnmarshalPublished(t *testing.T) {
	var pipeline DataPipeline
	if err := json.Unmarshal([]byte(publishedPipelineJSON), &pipeline); err != nil {
		t.Fatalf("unmarshalling published pipeline: %v", err)
	}

	if pipeline.ID != "00000000HX546" {
		t.Errorf("ID = %q, want %q", pipeline.ID, "00000000HX546")
	}
	if pipeline.Version != 7 {
		t.Errorf("Version = %d, want 7", pipeline.Version)
	}
	if pipeline.Order != 2 {
		t.Errorf("Order = %d, want 2", pipeline.Order)
	}
	if pipeline.State != "published" {
		t.Errorf("State = %q, want %q", pipeline.State, "published")
	}
	if got := len(pipeline.Nodes); got != 4 {
		t.Fatalf("len(Nodes) = %d, want 4", got)
	}
	if pipeline.Nodes[3].Outputs != nil {
		t.Errorf("destination node Outputs = %v, want nil", pipeline.Nodes[3].Outputs)
	}
}

func TestFlattenDataPipelineNodes(t *testing.T) {
	var pipeline DataPipeline
	if err := json.Unmarshal([]byte(publishedPipelineJSON), &pipeline); err != nil {
		t.Fatalf("unmarshalling published pipeline: %v", err)
	}

	flattened := flattenDataPipelineNodes(pipeline.Nodes)
	if got := len(flattened); got != 4 {
		t.Fatalf("len(flattened) = %d, want 4", got)
	}

	wantNames := []string{"Source", "Triage", "Log Normalization", "Sumo Logic"}
	for i, want := range wantNames {
		node := flattened[i].(map[string]interface{})
		if got := node["name"].(string); got != want {
			t.Errorf("node[%d][name] = %q, want %q", i, got, want)
		}
	}

	source := flattened[0].(map[string]interface{})
	if got := source["node_type"].(string); got != "source" {
		t.Errorf("source node_type = %q, want %q", got, "source")
	}
	if got := source["processor"].([]interface{}); len(got) != 0 {
		t.Errorf("source processor = %v, want empty", got)
	}

	group := flattened[2].(map[string]interface{})
	if got := group["filter_expression"].(string); got != "env=prod" {
		t.Errorf("group filter_expression = %q, want %q", got, "env=prod")
	}
	if got := group["is_enabled"].(bool); !got {
		t.Error("group is_enabled = false, want true")
	}
	if got := group["processor"].([]interface{}); len(got) != 2 {
		t.Fatalf("len(group processor) = %d, want 2", len(got))
	}

	destination := flattened[3].(map[string]interface{})
	if got := destination["output"].([]interface{}); len(got) != 0 {
		t.Errorf("destination output = %v, want empty", got)
	}
	if got := destination["is_enabled"].(bool); got {
		t.Error("destination is_enabled = true, want false")
	}
}

func TestFlattenDataPipelineProcessorsPreservesConfigJSON(t *testing.T) {
	var pipeline DataPipeline
	if err := json.Unmarshal([]byte(publishedPipelineJSON), &pipeline); err != nil {
		t.Fatalf("unmarshalling published pipeline: %v", err)
	}

	processors := flattenDataPipelineProcessors(pipeline.Nodes[2].Processors)
	if got := len(processors); got != 2 {
		t.Fatalf("len(processors) = %d, want 2", got)
	}

	first := processors[0].(map[string]interface{})
	if got := first["processor_type"].(string); got != "JSON_PARSE" {
		t.Errorf("processor_type = %q, want %q", got, "JSON_PARSE")
	}
	if got := first["order"].(int); got != 1 {
		t.Errorf("order = %d, want 1", got)
	}
	if got := first["is_enabled"].(bool); !got {
		t.Error("is_enabled = false, want true")
	}

	var config map[string]interface{}
	if err := json.Unmarshal([]byte(first["config"].(string)), &config); err != nil {
		t.Fatalf("config is not valid JSON: %v", err)
	}
	if got := config["fieldContainingJson"]; got != "_raw" {
		t.Errorf("config[fieldContainingJson] = %v, want %q", got, "_raw")
	}

	second := processors[1].(map[string]interface{})
	if got := second["is_enabled"].(bool); got {
		t.Error("second processor is_enabled = true, want false")
	}
	var nested map[string]interface{}
	if err := json.Unmarshal([]byte(second["config"].(string)), &nested); err != nil {
		t.Fatalf("nested config is not valid JSON: %v", err)
	}
	if _, ok := nested["fields"].([]interface{}); !ok {
		t.Errorf("config[fields] = %v, want a JSON array", nested["fields"])
	}
}

func TestFlattenDataPipelineOutputs(t *testing.T) {
	var pipeline DataPipeline
	if err := json.Unmarshal([]byte(publishedPipelineJSON), &pipeline); err != nil {
		t.Fatalf("unmarshalling published pipeline: %v", err)
	}

	outputs := flattenDataPipelineOutputs(pipeline.Nodes[1].Outputs)
	want := []interface{}{
		map[string]interface{}{
			"target":    "Security Processing",
			"condition": "event_category=security",
			"order":     1,
		},
		map[string]interface{}{
			"target":    "Log Normalization",
			"condition": "",
			"order":     0,
		},
	}

	if !reflect.DeepEqual(outputs, want) {
		t.Errorf("outputs = %#v\nwant %#v", outputs, want)
	}
}

func TestFlattenDataPipelineNodesDraftWithoutIDs(t *testing.T) {
	const draftPipelineJSON = `{
      "id": "00000000HX777",
      "customerId": "00000000HX123",
      "name": "draft-pipeline",
      "pipelineType": "route_based",
      "isEnabled": false,
      "order": 5,
      "version": 1,
      "createdAt": "2026-09-14T09:00:00Z",
      "modifiedAt": "2026-09-14T09:00:00Z",
      "createdByUserId": "00000000HX546",
      "modifiedByUserId": "00000000HX546",
      "state": "draft",
      "nodes": [
        {
          "name": "Log Normalization",
          "nodeType": "processing_group",
          "processors": [
            {
              "name": "regex-1",
              "processorType": "REGEX_PARSE",
              "order": 1,
              "isEnabled": true,
              "config": {"expression": "(?<code>\\d+)"}
            }
          ],
          "outputs": [{"target": "Sumo Logic"}]
        }
      ]
    }`

	var pipeline DataPipeline
	if err := json.Unmarshal([]byte(draftPipelineJSON), &pipeline); err != nil {
		t.Fatalf("unmarshalling draft pipeline: %v", err)
	}

	if pipeline.State != "draft" {
		t.Errorf("State = %q, want %q", pipeline.State, "draft")
	}

	flattened := flattenDataPipelineNodes(pipeline.Nodes)
	if got := len(flattened); got != 1 {
		t.Fatalf("len(flattened) = %d, want 1", got)
	}

	node := flattened[0].(map[string]interface{})
	if got := node["id"].(string); got != "" {
		t.Errorf("draft node id = %q, want empty", got)
	}
	if got := node["name"].(string); got != "Log Normalization" {
		t.Errorf("draft node name = %q, want %q", got, "Log Normalization")
	}
	if got := node["is_enabled"].(bool); got {
		t.Error("draft node is_enabled = true, want false")
	}

	processors := node["processor"].([]interface{})
	if got := len(processors); got != 1 {
		t.Fatalf("len(processors) = %d, want 1", got)
	}
	processor := processors[0].(map[string]interface{})
	if got := processor["id"].(string); got != "" {
		t.Errorf("draft processor id = %q, want empty", got)
	}
	if got := processor["processor_type"].(string); got != "REGEX_PARSE" {
		t.Errorf("draft processor_type = %q, want %q", got, "REGEX_PARSE")
	}
}

func TestFlattenDataPipelineNilSlices(t *testing.T) {
	if got := flattenDataPipelineNodes(nil); got == nil || len(got) != 0 {
		t.Errorf("flattenDataPipelineNodes(nil) = %#v, want empty non-nil", got)
	}
	if got := flattenDataPipelineProcessors(nil); got == nil || len(got) != 0 {
		t.Errorf("flattenDataPipelineProcessors(nil) = %#v, want empty non-nil", got)
	}
	if got := flattenDataPipelineOutputs(nil); got == nil || len(got) != 0 {
		t.Errorf("flattenDataPipelineOutputs(nil) = %#v, want empty non-nil", got)
	}
}

func TestGetDataPipelineFound(t *testing.T) {
	response := &http.Response{
		Status:     http.StatusText(200),
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(publishedPipelineJSON))),
	}
	client := newTestClient(response)

	pipeline, err := client.GetDataPipeline("00000000HX546")
	if err != nil {
		t.Fatalf("Expected GetDataPipeline to succeed, received: %s", err)
	}
	if pipeline == nil {
		t.Fatal("Expected GetDataPipeline to return a pipeline, got nil")
	}
	if pipeline.ID != "00000000HX546" {
		t.Errorf("ID = %q, want %q", pipeline.ID, "00000000HX546")
	}
	if got := len(pipeline.Nodes); got != 4 {
		t.Errorf("len(Nodes) = %d, want 4", got)
	}
}

func TestGetDataPipelineNotFound(t *testing.T) {
	body := []byte(`{
		"status": 404,
		"code": "pipelines.pipeline.notfound",
		"message": "The specified pipeline ID is invalid."
	}`)
	response := &http.Response{
		Status:     http.StatusText(404),
		StatusCode: 404,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}
	client := newTestClient(response)

	pipeline, err := client.GetDataPipeline("00000000HX999")
	if err != nil {
		t.Fatalf("Expected GetDataPipeline to succeed, received: %s", err)
	}
	if pipeline != nil {
		t.Errorf("Expected GetDataPipeline to return nil, instead got %v", *pipeline)
	}
}

func TestGetDataPipelineUnauthorized(t *testing.T) {
	body := []byte(`<html></html>`)
	response := &http.Response{
		Status:     http.StatusText(401),
		StatusCode: 401,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}
	client := newTestClient(response)

	_, err := client.GetDataPipeline("00000000HX546")
	if err == nil {
		t.Error("Expected GetDataPipeline to fail, but it succeeded")
	}
}

func TestFindDataPipelineByNameNotFound(t *testing.T) {
	body := []byte(`{
		"pipelines": [
			{"id": "P1", "name": "other-pipeline"}
		]
	}`)
	response := &http.Response{
		Status:     http.StatusText(200),
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}
	client := newTestClient(response)

	pipeline, err := client.FindDataPipelineByName("prod-log-pipeline")
	if err != nil {
		t.Fatalf("Expected FindDataPipelineByName to succeed, received: %s", err)
	}
	if pipeline != nil {
		t.Errorf("Expected FindDataPipelineByName to return nil, instead got %v", *pipeline)
	}
}

func TestFindDataPipelineByNameMultipleMatches(t *testing.T) {
	body := []byte(`{
		"pipelines": [
			{"id": "P1", "name": "duplicate-pipeline"},
			{"id": "P2", "name": "duplicate-pipeline"}
		]
	}`)
	response := &http.Response{
		Status:     http.StatusText(200),
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}
	client := newTestClient(response)

	_, err := client.FindDataPipelineByName("duplicate-pipeline")
	if err == nil {
		t.Fatal("Expected FindDataPipelineByName to fail, but it succeeded")
	}
	want := `found 2 data pipelines named "duplicate-pipeline", expected at most 1`
	if got := err.Error(); got != want {
		t.Errorf("error = %q, want %q", got, want)
	}
}
