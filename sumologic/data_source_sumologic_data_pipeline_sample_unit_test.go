package sumologic

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const dataPipelineSampleTestResponseJSON = `{
  "messages": [
    {
      "messageIndex": 0,
      "input": "line one",
      "output": "{\"_raw\":\"line one\",\"w\":\"one\"}",
      "pipelineMatched": true,
      "truncated": false,
      "nodeResults": [
        {
          "nodeName": "Triage",
          "nodeType": "router",
          "matched": true,
          "input": "{\"_raw\":\"line one\"}",
          "output": "{\"_raw\":\"line one\"}",
          "selectedRouteCondition": "env=prod",
          "selectedRouteTarget": "Parse"
        },
        {
          "nodeName": "Parse",
          "nodeType": "processing_group",
          "matched": true,
          "input": "{\"_raw\":\"line one\"}",
          "output": "{\"_raw\":\"line one\",\"w\":\"one\"}",
          "processorResults": [
            {
              "processorName": "p",
              "processorType": "REGEX_PARSE",
              "input": "{\"_raw\":\"line one\"}",
              "output": "{\"_raw\":\"line one\",\"w\":\"one\"}",
              "messageDropped": false
            }
          ],
          "selectedRouteCondition": null,
          "selectedRouteTarget": null
        }
      ],
      "edgeResults": [
        {"sourceNode": "Routing Expression", "targetNode": "Triage", "conditionExpression": null, "edgeOrder": 1},
        {"sourceNode": "Triage", "targetNode": "Parse", "conditionExpression": "env=prod", "edgeOrder": 1}
      ]
    },
    {
      "messageIndex": 1,
      "input": "line two",
      "output": "",
      "pipelineMatched": true,
      "truncated": true,
      "nodeResults": [
        {
          "nodeName": "Parse",
          "nodeType": "processing_group",
          "matched": true,
          "input": "{\"_raw\":\"line two\"}",
          "output": "",
          "processorResults": [
            {
              "processorName": "p",
              "processorType": "REGEX_PARSE",
              "input": "{\"_raw\":\"line two\"}",
              "output": "",
              "messageDropped": true
            }
          ]
        }
      ],
      "edgeResults": []
    }
  ]
}`

type dataPipelineCapturingHttpClient struct {
	status   int
	response string
	method   string
	path     string
	body     []byte
}

func (c *dataPipelineCapturingHttpClient) Do(req *http.Request) (*http.Response, error) {
	c.method = req.Method
	c.path = strings.TrimPrefix(req.URL.Path, "/api/")
	if req.Body != nil {
		c.body, _ = io.ReadAll(req.Body)
	}

	return &http.Response{
		Status:     http.StatusText(c.status),
		StatusCode: c.status,
		Body:       io.NopCloser(strings.NewReader(c.response)),
	}, nil
}

func newDataPipelineCapturingTestClient(status int, response string) (*Client, *dataPipelineCapturingHttpClient) {
	client := newTestClient(nil)
	httpClient := &dataPipelineCapturingHttpClient{status: status, response: response}
	client.httpClient = httpClient
	return client, httpClient
}

func dataPipelineSampleTestRawConfig() map[string]interface{} {
	return map[string]interface{}{
		"name":             "sample-pipeline",
		"route_expression": "_sourceCategory=prod/*",
		"sample_data":      "line one\nline two",
		"node": []interface{}{
			map[string]interface{}{
				"name":      "Routing Expression",
				"node_type": "source",
				"output":    []interface{}{map[string]interface{}{"target": "Parse"}},
			},
			map[string]interface{}{
				"name":      "Parse",
				"node_type": "processing_group",
				"processor": []interface{}{
					map[string]interface{}{
						"name":           "p",
						"processor_type": "REGEX_PARSE",
						"config":         `{"fieldToParseFrom":"_raw","regexPattern":"(?<w>\\w+)$"}`,
					},
				},
				"output": []interface{}{map[string]interface{}{"target": "Sumo Logic"}},
			},
			map[string]interface{}{
				"name":      "Sumo Logic",
				"node_type": "destination",
			},
		},
	}
}

func TestDataSourceDataPipelineTestReadSendsRequestAndFlattensResults(t *testing.T) {
	client, httpClient := newDataPipelineCapturingTestClient(200, dataPipelineSampleTestResponseJSON)

	raw := dataPipelineSampleTestRawConfig()
	raw["message_separator"] = `\n`
	raw["metadata_fields"] = map[string]interface{}{"_sourceCategory": "prod/nginx"}
	d := schema.TestResourceDataRaw(t, dataSourceDataPipelineTestSchema(), raw)

	if err := dataSourceSumologicDataPipelineTestRead(d, client); err != nil {
		t.Fatalf("Expected read to succeed, received: %s", err)
	}

	if httpClient.method != http.MethodPost || httpClient.path != "v1/pipelines/test" {
		t.Errorf("request = %s %s, want POST v1/pipelines/test", httpClient.method, httpClient.path)
	}

	var sent map[string]interface{}
	if err := json.Unmarshal(httpClient.body, &sent); err != nil {
		t.Fatalf("request body is not JSON: %v", err)
	}
	if got := sent["sampleData"]; got != "line one\nline two" {
		t.Errorf("sampleData = %q, want %q", got, "line one\nline two")
	}
	if got := sent["messageSeparator"]; got != `\n` {
		t.Errorf("messageSeparator = %q, want %q", got, `\n`)
	}
	if got := sent["metadataFields"]; !reflect.DeepEqual(got, map[string]interface{}{"_sourceCategory": "prod/nginx"}) {
		t.Errorf("metadataFields = %v", got)
	}

	pipeline := sent["pipeline"].(map[string]interface{})
	if pipeline["name"] != "sample-pipeline" || pipeline["pipelineType"] != "route_based" ||
		pipeline["routeExpression"] != "_sourceCategory=prod/*" || pipeline["isEnabled"] != true {
		t.Errorf("pipeline = %v", pipeline)
	}
	if _, ok := pipeline["version"]; ok {
		t.Errorf("pipeline must not carry a version: %v", pipeline)
	}
	nodes := pipeline["nodes"].([]interface{})
	if len(nodes) != 3 {
		t.Fatalf("sent %d nodes, want 3", len(nodes))
	}
	group := nodes[1].(map[string]interface{})
	processor := group["processors"].([]interface{})[0].(map[string]interface{})
	if processor["processorType"] != "REGEX_PARSE" || processor["order"] != float64(1) {
		t.Errorf("processor = %v", processor)
	}
	if !reflect.DeepEqual(processor["config"], map[string]interface{}{"fieldToParseFrom": "_raw", "regexPattern": `(?<w>\w+)$`}) {
		t.Errorf("processor config = %v", processor["config"])
	}

	if d.Id() == "" {
		t.Error("expected a non-empty id")
	}
	if got := d.Get("message.#").(int); got != 2 {
		t.Fatalf("message.# = %d, want 2", got)
	}

	first := d.Get("message.0").(map[string]interface{})
	if first["message_index"] != 0 || first["input"] != "line one" ||
		first["output"] != `{"_raw":"line one","w":"one"}` || first["truncated"] != false {
		t.Errorf("message.0 = %v", first)
	}

	router := d.Get("message.0.node_result.0").(map[string]interface{})
	if router["node_name"] != "Triage" || router["node_type"] != "router" ||
		router["selected_route_condition"] != "env=prod" || router["selected_route_target"] != "Parse" {
		t.Errorf("message.0.node_result.0 = %v", router)
	}

	group0 := d.Get("message.0.node_result.1").(map[string]interface{})
	if group0["selected_route_condition"] != "" || group0["selected_route_target"] != "" {
		t.Errorf("null route fields should flatten to empty strings: %v", group0)
	}
	processorResult := d.Get("message.0.node_result.1.processor_result.0").(map[string]interface{})
	if processorResult["processor_name"] != "p" || processorResult["processor_type"] != "REGEX_PARSE" ||
		processorResult["output"] != `{"_raw":"line one","w":"one"}` || processorResult["message_dropped"] != false {
		t.Errorf("processor_result = %v", processorResult)
	}

	edge := d.Get("message.0.edge_result.1").(map[string]interface{})
	if edge["source_node"] != "Triage" || edge["target_node"] != "Parse" ||
		edge["condition_expression"] != "env=prod" || edge["edge_order"] != 1 {
		t.Errorf("edge_result = %v", edge)
	}
	if got := d.Get("message.0.edge_result.0.condition_expression"); got != "" {
		t.Errorf("null condition should flatten to empty string, got %v", got)
	}

	second := d.Get("message.1").(map[string]interface{})
	if second["message_index"] != 1 || second["input"] != "line two" {
		t.Errorf("message.1 = %v", second)
	}
	if second["truncated"] != true {
		t.Errorf("message.1 truncated = %v, want true", second["truncated"])
	}
	if got := d.Get("message.1.node_result.0.processor_result.0.message_dropped"); got != true {
		t.Errorf("message_dropped = %v, want true", got)
	}
}

func TestDataSourceDataPipelineTestReadOmitsUnsetOptionalFields(t *testing.T) {
	client, httpClient := newDataPipelineCapturingTestClient(200, `{"messages": []}`)
	d := schema.TestResourceDataRaw(t, dataSourceDataPipelineTestSchema(), dataPipelineSampleTestRawConfig())

	if err := dataSourceSumologicDataPipelineTestRead(d, client); err != nil {
		t.Fatalf("Expected read to succeed, received: %s", err)
	}

	var sent map[string]interface{}
	if err := json.Unmarshal(httpClient.body, &sent); err != nil {
		t.Fatalf("request body is not JSON: %v", err)
	}
	for _, key := range []string{"messageSeparator", "metadataFields"} {
		if _, ok := sent[key]; ok {
			t.Errorf("%s should be omitted when unset, body = %s", key, httpClient.body)
		}
	}
	if got := d.Get("message.#").(int); got != 0 {
		t.Errorf("message.# = %d, want 0", got)
	}
}

func TestDataSourceDataPipelineTestReadIDIsStableAndInputSensitive(t *testing.T) {
	read := func(sampleData string) string {
		client, _ := newDataPipelineCapturingTestClient(200, `{"messages": []}`)
		raw := dataPipelineSampleTestRawConfig()
		raw["sample_data"] = sampleData
		d := schema.TestResourceDataRaw(t, dataSourceDataPipelineTestSchema(), raw)
		if err := dataSourceSumologicDataPipelineTestRead(d, client); err != nil {
			t.Fatalf("Expected read to succeed, received: %s", err)
		}
		return d.Id()
	}

	a, b, c := read("one"), read("one"), read("two")
	if a != b {
		t.Errorf("same input produced different ids: %q vs %q", a, b)
	}
	if a == c {
		t.Errorf("different sample data produced the same id: %q", a)
	}
}

func TestDataSourceDataPipelineTestReadSurfacesAPIError(t *testing.T) {
	client, _ := newDataPipelineCapturingTestClient(400, `{
		"id": "QCI5F-0NC1D-ZZSW5",
		"errors": [{"code": "PIPELINE-ERR-6", "message": "Invalid message separator: ';'."}]
	}`)
	d := schema.TestResourceDataRaw(t, dataSourceDataPipelineTestSchema(), dataPipelineSampleTestRawConfig())

	err := dataSourceSumologicDataPipelineTestRead(d, client)
	if err == nil {
		t.Fatal("Expected read to fail on a 400 response")
	}
	for _, want := range []string{`error testing data pipeline "sample-pipeline"`, "Invalid message separator"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error = %q, want it to contain %q", err.Error(), want)
		}
	}
	if d.Id() != "" {
		t.Errorf("id = %q, want empty after a failed read", d.Id())
	}
}

func TestDataSourceDataPipelineTestMessageSeparatorValidation(t *testing.T) {
	validate := dataSourceDataPipelineTestSchema()["message_separator"].ValidateFunc

	for _, separator := range dataPipelineMessageSeparators {
		if _, errs := validate(separator, "message_separator"); len(errs) != 0 {
			t.Errorf("separator %q should be valid, got %v", separator, errs)
		}
	}

	for _, separator := range []string{"", "\n", "\r\n", ";", ",", "|", "empty"} {
		if _, errs := validate(separator, "message_separator"); len(errs) == 0 {
			t.Errorf("separator %q should be rejected", separator)
		}
	}
}
