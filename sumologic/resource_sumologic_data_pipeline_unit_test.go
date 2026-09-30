package sumologic

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"
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
	}

	raw := flattenDataPipelineResourceNodes(nodes)
	got := expandDataPipelineNodes(raw)

	if !reflect.DeepEqual(got, nodes) {
		t.Errorf("expand(flatten(nodes)) = %#v\nwant %#v", got, nodes)
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
