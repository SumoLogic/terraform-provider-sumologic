package sumologic

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestCSESigmaRuleToMatchRule(t *testing.T) {
	var rule CSESigmaRule
	err := json.Unmarshal([]byte(`{
		"name": "Example",
		"enabled": true,
		"is_prototype": true,
		"name_expression": "Example signal",
		"description_expression": "Example description",
		"summary_expression": "",
		"expression": "foo = bar",
		"stream": "record",
		"entity_selectors": [{"entity_type": "_hostname", "expression": "device_hostname"}],
		"score_mapping": {"type": "constant", "default": 3, "mapping": []},
		"tags": ["_mitreAttackTactic:TA0004"],
		"category": "Privilege Escalation",
		"conversion_metadata": {"sigma_uid": "637f689e-b4a5-4a86-be0e-0100a0a33ba2"}
	}`), &rule)
	if err != nil {
		t.Fatal(err)
	}

	got := rule.toMatchRule()

	want := CSEMatchRule{
		DescriptionExpression: "Example description",
		Enabled:               true,
		EntitySelectors:       []EntitySelector{{EntityType: "_hostname", Expression: "device_hostname"}},
		Expression:            "foo = bar",
		IsPrototype:           true,
		Name:                  "Example",
		NameExpression:        "Example signal",
		SeverityMapping:       SeverityMapping{Type: "constant", Default: 3, Mapping: []SeverityMappingValueMapping{}},
		Stream:                "record",
		Tags:                  []string{"_mitreAttackTactic:TA0004"},
		OriginalSource:        "sigma",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected match rule\n got: %#v\nwant: %#v", got, want)
	}

	// Category is read back from the created rule, never sent.
	body, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	var sent map[string]interface{}
	if err := json.Unmarshal(body, &sent); err != nil {
		t.Fatal(err)
	}
	if _, ok := sent["category"]; ok {
		t.Errorf("category should not be sent, got %v", sent["category"])
	}
	if sent["originalSource"] != "sigma" {
		t.Errorf("expected originalSource sigma, got %v", sent["originalSource"])
	}
}

func TestCSESigmaRuleToMatchRuleNilTags(t *testing.T) {
	got := CSESigmaRule{}.toMatchRule()

	if got.Tags == nil {
		t.Error("expected tags to be an empty list, got nil")
	}
	if got.OriginalSource != "sigma" {
		t.Errorf("expected originalSource sigma, got %q", got.OriginalSource)
	}
}

func TestCSESigmaScoreMappingValueMappings(t *testing.T) {
	var mapping CSESigmaScoreMapping
	err := json.Unmarshal([]byte(`{
		"type": "fieldValueMapping",
		"default": 3,
		"field": "severity",
		"mapping": [{"type": "eq", "from": "high", "to": 8}, "not-an-object", {"type": "eq", "from": "", "to": 1}]
	}`), &mapping)
	if err != nil {
		t.Fatal(err)
	}

	got := mapping.valueMappings()
	want := []SeverityMappingValueMapping{{Type: "eq", From: "high", To: 8}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected value mappings\n got: %#v\nwant: %#v", got, want)
	}
}
