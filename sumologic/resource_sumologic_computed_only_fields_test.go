package sumologic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Server-managed fields must be Computed-only; Optional+Computed lets tools that
// copy computed values into desired state (e.g. Crossplane) pin stale values.
func TestServerManagedFieldsAreComputedOnly(t *testing.T) {
	library := []string{"version", "created_at", "created_by", "modified_at", "modified_by", "is_mutable", "is_system"}
	withLocked := append(append([]string{}, library...), "is_locked")

	cases := []struct {
		name     string
		resource *schema.Resource
		fields   []string
	}{
		{"sumologic_monitor", resourceSumologicMonitorsLibraryMonitor(), withLocked},
		{"sumologic_monitor_folder", resourceSumologicMonitorsLibraryFolder(), withLocked},
		{"sumologic_slo", resourceSumologicSLO(), library},
		{"sumologic_slo_folder", resourceSumologicSLOLibraryFolder(), library},
		{"sumologic_muting_schedule", resourceSumologicMutingSchedulesLibraryMutingSchedule(), library},
		{"sumologic_token", resourceSumologicToken(), []string{"version"}},
	}
	for _, c := range cases {
		for _, f := range c.fields {
			s, ok := c.resource.Schema[f]
			if !ok {
				t.Errorf("%s: missing field %q", c.name, f)
				continue
			}
			if s.Optional || !s.Computed {
				t.Errorf("%s.%s: want Computed-only, got Optional=%v Computed=%v", c.name, f, s.Optional, s.Computed)
			}
		}
	}
}
