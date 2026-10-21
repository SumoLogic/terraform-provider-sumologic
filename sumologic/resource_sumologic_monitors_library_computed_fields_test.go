package sumologic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Server-managed fields must be Computed-only; Optional+Computed lets tools that
// copy computed values into desired state (e.g. Crossplane) pin stale values.
func TestMonitorsLibraryServerManagedFieldsAreComputedOnly(t *testing.T) {
	fields := []string{"version", "created_at", "created_by", "modified_at", "modified_by", "is_mutable", "is_system", "is_locked"}
	resources := map[string]*schema.Resource{
		"sumologic_monitor":        resourceSumologicMonitorsLibraryMonitor(),
		"sumologic_monitor_folder": resourceSumologicMonitorsLibraryFolder(),
	}
	for name, res := range resources {
		for _, f := range fields {
			s, ok := res.Schema[f]
			if !ok {
				t.Errorf("%s: missing field %q", name, f)
				continue
			}
			if s.Optional || !s.Computed {
				t.Errorf("%s.%s: want Computed-only, got Optional=%v Computed=%v", name, f, s.Optional, s.Computed)
			}
		}
	}
}
