package sumologic

import "testing"

// The API returns camelCase keys; the export must map them to the snake_case Terraform attributes.
func TestGetTerraformServiceMapPanel(t *testing.T) {
	visualSettings := `{"serviceMap":{"mode":"serviceMap"},"latencyType":"{{latency_type}}"}`
	apiPanel := map[string]interface{}{
		"key":                                    "panel-key",
		"title":                                  "Service Map",
		"visualSettings":                         visualSettings,
		"keepVisualSettingsConsistentWithParent": true,
		"application":                            "app",
		"service":                                "svc",
		"showRemoteServices":                     true,
		"environment":                            "prod",
	}

	tfPanel := getTerraformServiceMapPanel(apiPanel)[0]

	expected := map[string]interface{}{
		"key":             "panel-key",
		"title":           "Service Map",
		"visual_settings": visualSettings,
		"keep_visual_settings_consistent_with_parent": true,
		"application":          "app",
		"service":              "svc",
		"show_remote_services": true,
		"environment":          "prod",
	}
	for k, want := range expected {
		if got := tfPanel[k]; got != want {
			t.Errorf("%s: got %v, want %v", k, got, want)
		}
	}
}
