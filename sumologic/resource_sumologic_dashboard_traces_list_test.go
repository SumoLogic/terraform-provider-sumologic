package sumologic

import "testing"

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
