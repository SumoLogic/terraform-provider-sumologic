package sumologic

import "testing"

func TestValidateDataMaskRuleRegex(t *testing.T) {
	_, errs := validateDataMaskRuleRegex("[a-z]+", "regex_pattern")
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors for valid regex, got %v", errs)
	}

	_, errs = validateDataMaskRuleRegex("[a-z", "regex_pattern")
	if len(errs) == 0 {
		t.Fatal("expected validation error for invalid regex")
	}
}

func TestListDataMaskRuleResponseReset(t *testing.T) {
	resp := &ListDataMaskRuleResponse{
		Data: []DataMaskRule{{ID: "a"}},
		Next: "next-token",
	}

	resp.Reset()

	if len(resp.Data) != 0 {
		t.Fatalf("expected Data to be empty, got %d", len(resp.Data))
	}
	if resp.Next != "" {
		t.Fatalf("expected Next to be empty, got %s", resp.Next)
	}
}
