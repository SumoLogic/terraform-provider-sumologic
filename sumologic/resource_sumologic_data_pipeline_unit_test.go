package sumologic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestCreateDataPipelineSuccess(t *testing.T) {
	response := &http.Response{
		Status:     http.StatusText(201),
		StatusCode: 201,
		Body:       io.NopCloser(bytes.NewReader([]byte(publishedPipelineJSON))),
	}
	client := newTestClient(response)

	request := DataPipelineRequest{
		Name:            "prod-log-pipeline",
		PipelineType:    "route_based",
		IsEnabled:       true,
		RouteExpression: "_sourceCategory=prod/*",
		Nodes:           []DataPipelineNode{{Name: "Sumo Logic", NodeType: "destination"}},
	}

	pipeline, err := client.CreateDataPipeline(request)
	if err != nil {
		t.Fatalf("Expected CreateDataPipeline to succeed, received: %s", err)
	}
	if pipeline.ID != "00000000HX546" {
		t.Errorf("ID = %q, want %q", pipeline.ID, "00000000HX546")
	}
}

func TestUpdateDataPipelineVersionConflict(t *testing.T) {
	body := []byte(`{
		"status": 409,
		"code": "pipelines.pipeline.versionConflict",
		"message": "The pipeline has been modified since it was last read."
	}`)
	response := &http.Response{
		Status:     http.StatusText(409),
		StatusCode: 409,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}
	client := newTestClient(response)

	_, err := client.UpdateDataPipeline("00000000HX546", DataPipelineUpdateRequest{Version: 3})
	if err == nil {
		t.Fatal("Expected UpdateDataPipeline to fail on a version conflict, but it succeeded")
	}
}

func TestDataPipelineUpdateRequestVersionZeroSerialized(t *testing.T) {
	body, err := json.Marshal(DataPipelineUpdateRequest{Version: 0})
	if err != nil {
		t.Fatalf("Marshal failed: %s", err)
	}
	if !bytes.Contains(body, []byte(`"version":0`)) {
		t.Errorf("marshaled request = %s, want it to contain \"version\":0", body)
	}
}

func TestPublishDataPipelineSuccess(t *testing.T) {
	response := &http.Response{
		Status:     http.StatusText(200),
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(publishedPipelineJSON))),
	}
	client := newTestClient(response)

	pipeline, err := client.PublishDataPipeline("00000000HX546")
	if err != nil {
		t.Fatalf("Expected PublishDataPipeline to succeed, received: %s", err)
	}
	if pipeline.State != "published" {
		t.Errorf("State = %q, want %q", pipeline.State, "published")
	}
}

func TestSetDataPipelineEnabledSuccess(t *testing.T) {
	response := &http.Response{
		Status:     http.StatusText(200),
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{}`))),
	}
	client := newTestClient(response)

	if err := client.SetDataPipelineEnabled("00000000HX546", false); err != nil {
		t.Fatalf("Expected SetDataPipelineEnabled to succeed, received: %s", err)
	}
}

func TestUpdateDataPipelineMetadataSuccess(t *testing.T) {
	response := &http.Response{
		Status:     http.StatusText(200),
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(publishedPipelineJSON))),
	}
	client := newTestClient(response)

	pipeline, err := client.UpdateDataPipelineMetadata("00000000HX546", DataPipelineMetadataRequest{
		Name:        "renamed-pipeline",
		Description: "",
	})
	if err != nil {
		t.Fatalf("Expected UpdateDataPipelineMetadata to succeed, received: %s", err)
	}
	if pipeline.ID != "00000000HX546" {
		t.Errorf("ID = %q, want %q", pipeline.ID, "00000000HX546")
	}
}

func TestDataPipelineMetadataRequestEmptyDescriptionSerialized(t *testing.T) {
	body, err := json.Marshal(DataPipelineMetadataRequest{Name: "renamed-pipeline", Description: ""})
	if err != nil {
		t.Fatalf("Marshal failed: %s", err)
	}
	if !bytes.Contains(body, []byte(`"description":""`)) {
		t.Errorf("marshaled request = %s, want it to contain \"description\":\"\"", body)
	}
}

func TestDeleteDataPipelineSuccess(t *testing.T) {
	response := &http.Response{
		Status:     http.StatusText(200),
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"pipelines": []}`))),
	}
	client := newTestClient(response)

	if err := client.DeleteDataPipeline("00000000HX546"); err != nil {
		t.Fatalf("Expected DeleteDataPipeline to succeed, received: %s", err)
	}
}

func TestDeleteDataPipelineNotFound(t *testing.T) {
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

	if err := client.DeleteDataPipeline("00000000HX999"); err == nil {
		t.Fatal("Expected DeleteDataPipeline to fail for a missing pipeline, but it succeeded")
	}
}

func TestReorderDataPipelineNodesMatchesReferenceOrder(t *testing.T) {
	apiOrder := []DataPipelineNode{
		{ID: "N1", Name: "Routing Expression", NodeType: "source"},
		{ID: "N3", Name: "Sumo Logic", NodeType: "destination"},
		{ID: "N2", Name: "Triage", NodeType: "router"},
	}

	got := reorderDataPipelineNodes(apiOrder, []string{"Routing Expression", "Triage", "Sumo Logic"})

	want := []string{"Routing Expression", "Triage", "Sumo Logic"}
	if len(got) != len(want) {
		t.Fatalf("len(got) = %d, want %d", len(got), len(want))
	}
	for i, name := range want {
		if got[i].Name != name {
			t.Errorf("got[%d].Name = %q, want %q", i, got[i].Name, name)
		}
	}
}

func TestReorderDataPipelineNodesAppendsUnreferencedNodes(t *testing.T) {
	apiOrder := []DataPipelineNode{
		{Name: "Routing Expression"},
		{Name: "Triage"},
	}

	got := reorderDataPipelineNodes(apiOrder, nil)

	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if got[0].Name != "Routing Expression" || got[1].Name != "Triage" {
		t.Errorf("got = %v, want original API order preserved", got)
	}
}

func TestSortDataPipelineProcessorsByOrderFollowsServerOrder(t *testing.T) {
	apiOrder := []DataPipelineProcessor{
		{ID: "P1", Name: "Add Fields", Order: 1},
		{ID: "P3", Name: "Filter by Condition", Order: 3},
		{ID: "P2", Name: "JSON Parse", Order: 2},
	}

	got := sortDataPipelineProcessorsByOrder(apiOrder)

	want := []string{"Add Fields", "JSON Parse", "Filter by Condition"}
	if len(got) != len(want) {
		t.Fatalf("len(got) = %d, want %d", len(got), len(want))
	}
	for i, name := range want {
		if got[i].Name != name {
			t.Errorf("got[%d].Name = %q, want %q", i, got[i].Name, name)
		}
	}
}

func TestResourceDataPipelineReadDetectsOutOfBandProcessorReorder(t *testing.T) {
	pipelineJSON := func(regexOrder, jsonOrder int) string {
		return fmt.Sprintf(`{
  "id": "00000000HX546",
  "name": "prod-log-pipeline",
  "pipelineType": "route_based",
  "isEnabled": true,
  "routeExpression": "_sourceCategory=prod/*",
  "version": 3,
  "state": "published",
  "nodes": [
    {"id": "NODE1", "name": "Routing Expression", "nodeType": "source", "outputs": [{"target": "Normalize"}]},
    {"id": "NODE2", "name": "Normalize", "nodeType": "processing_group", "isEnabled": true,
     "outputs": [{"target": "Sumo Logic"}],
     "processors": [
       {"id": "P1", "name": "Extract IP", "processorType": "REGEX_PARSE", "order": %d, "isEnabled": true, "config": {"regexPattern": "ip"}},
       {"id": "P2", "name": "Parse Body", "processorType": "JSON_PARSE", "order": %d, "isEnabled": true, "config": {"fieldContainingJson": "_raw"}}
     ]},
    {"id": "NODE3", "name": "Sumo Logic", "nodeType": "destination"}
  ]
}`, regexOrder, jsonOrder)
	}
	config := map[string]interface{}{
		"name":             "prod-log-pipeline",
		"route_expression": "_sourceCategory=prod/*",
		"node": []interface{}{
			map[string]interface{}{
				"name":      "Routing Expression",
				"node_type": "source",
				"output":    []interface{}{map[string]interface{}{"target": "Normalize"}},
			},
			map[string]interface{}{
				"name":      "Normalize",
				"node_type": "processing_group",
				"output":    []interface{}{map[string]interface{}{"target": "Sumo Logic"}},
				"processor": []interface{}{
					map[string]interface{}{"name": "Extract IP", "processor_type": "REGEX_PARSE", "config": `{"regexPattern":"ip"}`},
					map[string]interface{}{"name": "Parse Body", "processor_type": "JSON_PARSE", "config": `{"fieldContainingJson":"_raw"}`},
				},
			},
			map[string]interface{}{"name": "Sumo Logic", "node_type": "destination"},
		},
	}
	r := resourceSumologicDataPipeline()

	inSync, _ := newDataPipelineRoutingTestClient(map[string]string{"v1/pipelines/00000000HX546": pipelineJSON(1, 2)})
	d := r.Data(&terraform.InstanceState{ID: "00000000HX546"})
	if err := resourceSumologicDataPipelineRead(d, inSync); err != nil {
		t.Fatalf("Expected read to succeed, received: %s", err)
	}
	diff, err := r.Diff(context.Background(), d.State(), terraform.NewResourceConfigRaw(config), inSync)
	if err != nil {
		t.Fatalf("Expected diff to succeed, received: %s", err)
	}
	if diff != nil && !diff.Empty() {
		t.Fatalf("Expected no diff when server order matches config, got %v", diff.Attributes)
	}

	reordered, _ := newDataPipelineRoutingTestClient(map[string]string{"v1/pipelines/00000000HX546": pipelineJSON(2, 1)})
	d = r.Data(d.State())
	if err := resourceSumologicDataPipelineRead(d, reordered); err != nil {
		t.Fatalf("Expected read to succeed, received: %s", err)
	}
	processors := d.Get("node.1.processor").([]interface{})
	if got := processors[0].(map[string]interface{})["name"]; got != "Parse Body" {
		t.Errorf("Expected state to follow server order, first processor = %v", got)
	}
	diff, err = r.Diff(context.Background(), d.State(), terraform.NewResourceConfigRaw(config), reordered)
	if err != nil {
		t.Fatalf("Expected diff to succeed, received: %s", err)
	}
	if diff == nil || diff.Empty() {
		t.Fatal("Expected a diff after the processors were reordered outside Terraform")
	}
}

func TestTopologicalDataPipelineNodeOrderLinear(t *testing.T) {
	apiOrder := []DataPipelineNode{
		{Name: "Routing Expression", NodeType: "source", Outputs: []DataPipelineOutput{{Target: "Triage"}}},
		{Name: "Sumo Logic", NodeType: "destination"},
		{Name: "Triage", NodeType: "router", Outputs: []DataPipelineOutput{{Target: "Sumo Logic"}}},
	}

	got := topologicalDataPipelineNodeOrder(apiOrder)

	want := []string{"Routing Expression", "Triage", "Sumo Logic"}
	if len(got) != len(want) {
		t.Fatalf("got = %v, want %v", got, want)
	}
	for i, name := range want {
		if got[i] != name {
			t.Errorf("got[%d] = %q, want %q (got = %v)", i, got[i], name, got)
		}
	}
}

func TestTopologicalDataPipelineNodeOrderBranching(t *testing.T) {
	one := 1
	apiOrder := []DataPipelineNode{
		{Name: "Sumo Logic", NodeType: "destination"},
		{
			Name: "Triage", NodeType: "router",
			Outputs: []DataPipelineOutput{
				{Target: "Sumo Logic"},
				{Target: "Archive Router", Order: &one},
			},
		},
		{Name: "Archive Router", NodeType: "router", Outputs: []DataPipelineOutput{{Target: "Sumo Logic"}}},
		{Name: "Routing Expression", NodeType: "source", Outputs: []DataPipelineOutput{{Target: "Triage"}}},
	}

	got := topologicalDataPipelineNodeOrder(apiOrder)

	want := []string{"Routing Expression", "Triage", "Archive Router", "Sumo Logic"}
	if len(got) != len(want) {
		t.Fatalf("got = %v, want %v", got, want)
	}
	for i, name := range want {
		if got[i] != name {
			t.Errorf("got[%d] = %q, want %q (got = %v)", i, got[i], name, got)
		}
	}
}

func TestTopologicalDataPipelineNodeOrderUnreachableNodeAppended(t *testing.T) {
	apiOrder := []DataPipelineNode{
		{Name: "Routing Expression", NodeType: "source", Outputs: []DataPipelineOutput{{Target: "Sumo Logic"}}},
		{Name: "Orphan", NodeType: "router"},
		{Name: "Sumo Logic", NodeType: "destination"},
	}

	got := topologicalDataPipelineNodeOrder(apiOrder)

	if len(got) != 3 {
		t.Fatalf("len(got) = %d, want 3 (got = %v)", len(got), got)
	}
	if got[len(got)-1] != "Orphan" {
		t.Errorf("got = %v, want unreachable node appended last", got)
	}
}

func TestExpandDataPipelineNodesRoundTrip(t *testing.T) {
	order := 1
	isEnabled := true
	nodes := []DataPipelineNode{
		{
			ID:       "NODE2",
			Name:     "Triage",
			NodeType: "router",
			Outputs: []DataPipelineOutput{
				{Target: "Security Processing", Condition: "event_category=security", Order: &order},
				{Target: "Sumo Logic"},
			},
		},
		{
			ID:               "NODE3",
			Name:             "Log Normalization",
			NodeType:         "processing_group",
			FilterExpression: "_sourceCategory=prod/app/*",
			IsEnabled:        &isEnabled,
			Outputs:          []DataPipelineOutput{{Target: "Sumo Logic"}},
			Processors: []DataPipelineProcessor{
				{
					ID:            "PROC1",
					Name:          "json-parser-1",
					ProcessorType: "parsejson",
					Order:         1,
					IsEnabled:     true,
					Config:        json.RawMessage(`{"field":"_raw"}`),
				},
			},
		},
	}

	raw := flattenDataPipelineResourceNodes(nodes)
	got := expandDataPipelineNodes(raw)

	if !reflect.DeepEqual(got, nodes) {
		t.Errorf("expand(flatten(nodes)) = %#v\nwant %#v", got, nodes)
	}
}

func TestExpandDataPipelineNodesOmitsProcessingGroupFieldsOnOtherNodes(t *testing.T) {
	raw := []interface{}{
		map[string]interface{}{
			"id":                "",
			"name":              "Triage",
			"node_type":         "router",
			"filter_expression": "",
			"is_enabled":        true,
			"output":            []interface{}{map[string]interface{}{"target": "Sumo Logic", "condition": "", "order": 0}},
			"processor":         []interface{}{},
		},
	}

	body, err := json.Marshal(expandDataPipelineNodes(raw))
	if err != nil {
		t.Fatalf("Expected marshal to succeed, received: %s", err)
	}

	var nodes []map[string]interface{}
	if err := json.Unmarshal(body, &nodes); err != nil {
		t.Fatalf("Expected unmarshal to succeed, received: %s", err)
	}
	for _, field := range []string{"filterExpression", "isEnabled", "processors"} {
		if _, ok := nodes[0][field]; ok {
			t.Errorf("router node should not send %q, got %s", field, body)
		}
	}
}

func TestValidateDataPipelineNodeFields(t *testing.T) {
	node := func(nodeType string, overrides map[string]interface{}) map[string]interface{} {
		m := map[string]interface{}{
			"name":              "n",
			"node_type":         nodeType,
			"filter_expression": "",
			"is_enabled":        true,
			"processor":         []interface{}{},
		}
		for k, v := range overrides {
			m[k] = v
		}
		return m
	}
	processor := []interface{}{map[string]interface{}{"name": "p"}}

	cases := []struct {
		name    string
		node    map[string]interface{}
		wantErr string
	}{
		{"router defaults", node("router", nil), ""},
		{"processing group with all fields", node("processing_group", map[string]interface{}{
			"filter_expression": "env=prod", "is_enabled": false, "processor": processor,
		}), ""},
		{"router filter_expression", node("router", map[string]interface{}{"filter_expression": "env=prod"}), "filter_expression"},
		{"source is_enabled false", node("source", map[string]interface{}{"is_enabled": false}), "is_enabled"},
		{"destination processor", node("destination", map[string]interface{}{"processor": processor}), "processor"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateDataPipelineNodeFields([]interface{}{tc.node})
			if tc.wantErr == "" {
				if err != nil {
					t.Errorf("Expected no error, received: %s", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("Expected error mentioning %q, received: %v", tc.wantErr, err)
			}
		})
	}
}

func TestResourceDataPipelineDiffRejectsProcessingGroupFieldsOnOtherNodes(t *testing.T) {
	client, _ := newDataPipelineRoutingTestClient(map[string]string{
		"v1/pipelines/00000000HX546": simplePipelineJSON,
	})
	r := resourceSumologicDataPipeline()
	d := r.Data(&terraform.InstanceState{ID: "00000000HX546"})
	if err := resourceSumologicDataPipelineRead(d, client); err != nil {
		t.Fatalf("Expected read to succeed, received: %s", err)
	}

	config := simplePipelineConfig("_sourceCategory=prod/*", true)
	config["node"].([]interface{})[0].(map[string]interface{})["filter_expression"] = "env=prod"

	_, err := r.Diff(context.Background(), d.State(), terraform.NewResourceConfigRaw(config), client)
	if err == nil || !strings.Contains(err.Error(), "filter_expression can only be set on processing_group nodes") {
		t.Errorf("Expected diff to reject filter_expression on a source node, received: %v", err)
	}
}

func TestExpandDataPipelineOutputsZeroOrderOmitted(t *testing.T) {
	raw := []interface{}{
		map[string]interface{}{"target": "Sumo Logic", "condition": "", "order": 0},
	}

	outputs := expandDataPipelineOutputs(raw)
	if got := len(outputs); got != 1 {
		t.Fatalf("len(outputs) = %d, want 1", got)
	}
	if outputs[0].Order != nil {
		t.Errorf("Order = %v, want nil", *outputs[0].Order)
	}
}

const simplePipelineJSON = `{
  "id": "00000000HX546",
  "name": "prod-log-pipeline",
  "description": "",
  "pipelineType": "route_based",
  "isEnabled": true,
  "routeExpression": "_sourceCategory=prod/*",
  "version": 3,
  "state": "published",
  "nodes": [
    {
      "id": "NODE1",
      "name": "Routing Expression",
      "nodeType": "source",
      "outputs": [{"target": "Sumo Logic"}]
    },
    {
      "id": "NODE2",
      "name": "Sumo Logic",
      "nodeType": "destination"
    }
  ]
}`

func simplePipelineConfig(routeExpression string, isEnabled bool) map[string]interface{} {
	return map[string]interface{}{
		"name":             "prod-log-pipeline",
		"route_expression": routeExpression,
		"is_enabled":       isEnabled,
		"node": []interface{}{
			map[string]interface{}{
				"name":      "Routing Expression",
				"node_type": "source",
				"output": []interface{}{
					map[string]interface{}{"target": "Sumo Logic"},
				},
			},
			map[string]interface{}{
				"name":      "Sumo Logic",
				"node_type": "destination",
			},
		},
	}
}

func applyDataPipelineUpdate(t *testing.T, config map[string]interface{}) []string {
	t.Helper()

	client, httpClient := newDataPipelineRoutingTestClient(map[string]string{
		"v1/pipelines/00000000HX546":          simplePipelineJSON,
		"v1/pipelines/00000000HX546/publish":  simplePipelineJSON,
		"v1/pipelines/00000000HX546/enable":   `{}`,
		"v1/pipelines/00000000HX546/disable":  `{}`,
		"v1/pipelines/00000000HX546/metadata": simplePipelineJSON,
	})
	r := resourceSumologicDataPipeline()

	d := r.Data(&terraform.InstanceState{ID: "00000000HX546"})
	if err := resourceSumologicDataPipelineRead(d, client); err != nil {
		t.Fatalf("Expected read to succeed, received: %s", err)
	}
	state := d.State()

	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), client)
	if err != nil {
		t.Fatalf("Expected diff to succeed, received: %s", err)
	}
	if diff == nil || diff.Empty() {
		t.Fatal("Expected a non-empty diff")
	}

	httpClient.requested = nil
	if _, diags := r.Apply(context.Background(), state, diff, client); diags.HasError() {
		t.Fatalf("Expected apply to succeed, received: %v", diags)
	}

	return httpClient.requested
}

func TestResourceDataPipelineUpdateEnabledOnlySkipsUpdateAndPublish(t *testing.T) {
	requested := applyDataPipelineUpdate(t, simplePipelineConfig("_sourceCategory=prod/*", false))

	want := []string{
		"PATCH v1/pipelines/00000000HX546/disable",
		"GET v1/pipelines/00000000HX546",
	}
	if !reflect.DeepEqual(requested, want) {
		t.Errorf("requested = %v, want %v", requested, want)
	}
}

func TestResourceDataPipelineUpdateConfigChangeUpdatesAndPublishes(t *testing.T) {
	requested := applyDataPipelineUpdate(t, simplePipelineConfig("_sourceCategory=prod/updated", true))

	want := []string{
		"PUT v1/pipelines/00000000HX546",
		"POST v1/pipelines/00000000HX546/publish",
		"PATCH v1/pipelines/00000000HX546/enable",
		"GET v1/pipelines/00000000HX546",
	}
	if !reflect.DeepEqual(requested, want) {
		t.Errorf("requested = %v, want %v", requested, want)
	}
}

func TestResourceDataPipelineUpdateNameOnlySendsOnlyMetadata(t *testing.T) {
	config := simplePipelineConfig("_sourceCategory=prod/*", true)
	config["name"] = "prod-log-pipeline-renamed"
	requested := applyDataPipelineUpdate(t, config)

	want := []string{
		"PATCH v1/pipelines/00000000HX546/metadata",
		"GET v1/pipelines/00000000HX546",
	}
	if !reflect.DeepEqual(requested, want) {
		t.Errorf("requested = %v, want %v", requested, want)
	}
}

func TestResourceDataPipelineUpdateDescriptionOnlySendsOnlyMetadata(t *testing.T) {
	config := simplePipelineConfig("_sourceCategory=prod/*", true)
	config["description"] = "Handles prod log parsing."
	requested := applyDataPipelineUpdate(t, config)

	want := []string{
		"PATCH v1/pipelines/00000000HX546/metadata",
		"GET v1/pipelines/00000000HX546",
	}
	if !reflect.DeepEqual(requested, want) {
		t.Errorf("requested = %v, want %v", requested, want)
	}
}

func TestResourceDataPipelineUpdateConfigAndMetadataChangeRunsBoth(t *testing.T) {
	config := simplePipelineConfig("_sourceCategory=prod/updated", true)
	config["name"] = "prod-log-pipeline-renamed"
	requested := applyDataPipelineUpdate(t, config)

	want := []string{
		"PUT v1/pipelines/00000000HX546",
		"POST v1/pipelines/00000000HX546/publish",
		"PATCH v1/pipelines/00000000HX546/enable",
		"PATCH v1/pipelines/00000000HX546/metadata",
		"GET v1/pipelines/00000000HX546",
	}
	if !reflect.DeepEqual(requested, want) {
		t.Errorf("requested = %v, want %v", requested, want)
	}
}
