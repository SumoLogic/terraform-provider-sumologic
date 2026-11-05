package sumologic

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// The Terraform schema key for panel queries is "queries"; reading "query" used to panic.
func TestGetTracesListPanel(t *testing.T) {
	tfPanel := map[string]interface{}{
		"key":             "traces-key",
		"title":           "Traces",
		"visual_settings": `{"table":{}}`,
		"keep_visual_settings_consistent_with_parent": true,
		"queries": []interface{}{
			map[string]interface{}{
				"query_string": "",
				"query_type":   "Traces",
				"query_key":    "A",
			},
		},
	}

	panel := getTracesListPanel(tfPanel).(TracesListPanel)

	if panel.PanelType != "TracesListPanel" || panel.Key != "traces-key" {
		t.Errorf("unexpected panel: %+v", panel)
	}
	if len(panel.Queries) != 1 || panel.Queries[0].QueryType != "Traces" || panel.Queries[0].QueryString != "" {
		t.Errorf("unexpected queries: %+v", panel.Queries)
	}
}

// The API returns camelCase keys; the export must map them to the snake_case Terraform attributes.
func TestGetTerraformTracesListPanel(t *testing.T) {
	apiPanel := map[string]interface{}{
		"key":                                    "traces-key",
		"title":                                  "Traces",
		"visualSettings":                         `{"table":{}}`,
		"keepVisualSettingsConsistentWithParent": true,
		"queries": []interface{}{
			map[string]interface{}{
				"queryString": "",
				"queryType":   "Traces",
				"queryKey":    "A",
			},
		},
	}

	tfPanel := getTerraformTracesListPanel(apiPanel)[0]

	if got := tfPanel["visual_settings"]; got != `{"table":{}}` {
		t.Errorf("visual_settings: got %v", got)
	}
	queries := tfPanel["queries"].([]map[string]interface{})
	if len(queries) != 1 || queries[0]["query_type"] != "Traces" {
		t.Errorf("unexpected queries: %+v", queries)
	}
}

// Shape taken from the exported appservicehealth.json TracesListPanel.
func TestTracesQueryDataRoundTrip(t *testing.T) {
	apiQuery := map[string]interface{}{
		"queryString": "",
		"queryType":   "Traces",
		"queryKey":    "A",
		"tracesQueryData": map[string]interface{}{
			"filters": []interface{}{
				map[string]interface{}{
					"type": "FieldDescriptor", "operator": "=", "field": "deployment.environment",
					"value": "{{deployment.environment}}", "valueType": "StringTracingValue", "target": "any span",
				},
			},
		},
	}

	tfQueries := getTerraformSearchPanelQuery([]interface{}{apiQuery})
	tfFilter := tfQueries[0]["traces_query_data"].(TerraformObject)[0]["filter"].([]map[string]interface{})[0]
	if tfFilter["field"] != "deployment.environment" || tfFilter["value_type"] != "StringTracingValue" ||
		tfFilter["target"] != "any span" || tfFilter["operator"] != "=" || tfFilter["type"] != "FieldDescriptor" {
		t.Fatalf("export lost filter data: %+v", tfFilter)
	}

	// feed the exported form back through deploy, as terraform would after reading state
	tfQuery := map[string]interface{}{
		"query_string": "", "query_type": "Traces", "query_key": "A",
		"traces_query_data": []interface{}{map[string]interface{}{
			"filter": []interface{}{map[string]interface{}{
				"type": tfFilter["type"], "operator": tfFilter["operator"], "field": tfFilter["field"],
				"value": tfFilter["value"], "value_type": tfFilter["value_type"], "target": tfFilter["target"],
			}},
		}},
	}
	query := getSearchPanelQuery(tfQuery)
	want := TracesFilter{Type: "FieldDescriptor", Operator: "=", Field: "deployment.environment",
		Value: "{{deployment.environment}}", ValueType: "StringTracingValue", Target: "any span"}
	if query.TracesQueryData == nil || len(query.TracesQueryData.Filters) != 1 || query.TracesQueryData.Filters[0] != want {
		t.Fatalf("deploy lost filter data: %+v", query.TracesQueryData)
	}
}

func tracesDashboardConfig(panel map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"title": "test",
		"time_range": []interface{}{map[string]interface{}{
			"complete_literal_time_range": []interface{}{map[string]interface{}{"range_name": "today"}},
		}},
		"panel": []interface{}{panel},
	}
}

func tracesListPanelConfig(queries ...interface{}) map[string]interface{} {
	return map[string]interface{}{
		"traces_list_panel": []interface{}{map[string]interface{}{
			"key":     "traces-key",
			"queries": queries,
		}},
	}
}

func tracesQueryConfig(extra map[string]interface{}) map[string]interface{} {
	query := map[string]interface{}{"query_type": "Traces", "query_key": "A"}
	for k, v := range extra {
		query[k] = v
	}
	return query
}

func tracesFilterConfig(filter map[string]interface{}) []interface{} {
	return []interface{}{map[string]interface{}{"filter": []interface{}{filter}}}
}

// Runs the full schema validation + CustomizeDiff, so defaults, Required/Optional and cross-field checks are covered.
func planDashboard(t *testing.T, config map[string]interface{}) error {
	t.Helper()
	res := resourceSumologicDashboard()
	if diags := res.Validate(terraform.NewResourceConfigRaw(config)); diags.HasError() {
		return fmt.Errorf("%v", diags)
	}
	_, err := res.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(config), nil)
	return err
}

// Each config must pass schema validation and plan; none of them is rejected by the provider.
func TestDashboardTracesPanelSchemaAccepts(t *testing.T) {
	validFilter := map[string]interface{}{"field": "deployment.environment", "value": "prod"}

	cases := []struct {
		name  string
		panel map[string]interface{}
	}{
		{
			name: "traces query with filter and defaults",
			panel: tracesListPanelConfig(tracesQueryConfig(map[string]interface{}{
				"traces_query_data": tracesFilterConfig(validFilter),
			})),
		},
		{
			// the API accepts an empty filters list, and Terraformer omits an empty one on export
			name:  "traces query with no filter",
			panel: tracesListPanelConfig(tracesQueryConfig(map[string]interface{}{"traces_query_data": []interface{}{map[string]interface{}{}}})),
		},
		{
			// the API accepts empty filter fields; only a value it does not recognise is rejected, at apply time
			name: "empty operator",
			panel: tracesListPanelConfig(tracesQueryConfig(map[string]interface{}{
				"traces_query_data": tracesFilterConfig(map[string]interface{}{"field": "f", "value": "v", "operator": ""}),
			})),
		},
		{
			// the API accepts an empty queryString for any query type, so this must stay valid
			name: "logs query without query_string",
			panel: tracesListPanelConfig(map[string]interface{}{
				"query_type": "Logs", "query_key": "A",
			}),
		},
		{
			name:  "traces list panel without queries",
			panel: map[string]interface{}{"traces_list_panel": []interface{}{map[string]interface{}{"key": "traces-key"}}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := planDashboard(t, tracesDashboardConfig(tc.panel)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// Exports through d.Set against the real schema, so key names and value types are checked, not just the converter output.
func TestDashboardTracesPanelExportSetsState(t *testing.T) {
	d := resourceSumologicDashboard().TestResourceData()
	apiPanels := []interface{}{map[string]interface{}{
		"panelType":      "TracesListPanel",
		"key":            "traces-key",
		"visualSettings": `{"table":{}}`,
		"queries": []interface{}{map[string]interface{}{
			"queryString": "",
			"queryType":   "Traces",
			"queryKey":    "A",
			"tracesQueryData": map[string]interface{}{
				"filters": []interface{}{
					map[string]interface{}{"field": "http.method", "value": "GET", "target": "any span"},
					"not-a-map",
				},
			},
		}},
	}}

	if err := d.Set("panel", getTerraformPanels(apiPanels)); err != nil {
		t.Fatalf("d.Set failed: %v", err)
	}
	prefix := "panel.0.traces_list_panel.0."
	if got := d.Get(prefix + "visual_settings"); got != `{"table":{}}` {
		t.Errorf("visual_settings: got %v", got)
	}
	if got := d.Get(prefix + "queries.0.traces_query_data.0.filter.#"); got != 1 {
		t.Errorf("expected 1 filter, got %v", got)
	}
	if got := d.Get(prefix + "queries.0.traces_query_data.0.filter.0.field"); got != "http.method" {
		t.Errorf("filter field: got %v", got)
	}
}

func TestGetTerraformTracesListPanelWithoutQueries(t *testing.T) {
	tfPanel := getTerraformTracesListPanel(map[string]interface{}{"key": "traces-key"})[0]
	if _, ok := tfPanel["queries"]; ok {
		t.Errorf("queries should be left unset, got %v", tfPanel["queries"])
	}
}

// An empty traces_query_data block reaches the converter as a list holding nil.
func TestGetSearchPanelQueryEmptyTracesQueryData(t *testing.T) {
	query := getSearchPanelQuery(map[string]interface{}{
		"query_string": "", "query_type": "Traces", "query_key": "A",
		"traces_query_data": []interface{}{nil},
	})
	if query.TracesQueryData == nil || len(query.TracesQueryData.Filters) != 0 {
		t.Fatalf("expected empty traces query data, got %+v", query.TracesQueryData)
	}
}
