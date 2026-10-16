package sumologic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestTimeRangeAllowedValuesValidation(t *testing.T) {
	for _, timeRange := range []string{"-15m", "15m", "-24h", "-1d", "1d", "4320m", "-4320m", "72h", "-72h", "3d", "-3d"} {
		_, errs := timeRangeAllowedValuesValidation(timeRange, "time_range")
		if len(errs) != 0 {
			t.Errorf("expected no validation errors for %q, got %v", timeRange, errs)
		}
	}

	for _, timeRange := range []string{"", "-45m", "-48h", "-2d", "-7d", "-72x", "72"} {
		_, errs := timeRangeAllowedValuesValidation(timeRange, "time_range")
		if len(errs) == 0 {
			t.Errorf("expected validation error for %q", timeRange)
		}
	}
}

func TestLogsStaticTriggerConditionTimeRangeThreeDays(t *testing.T) {
	critical := logsStaticTriggerConditionSchema["critical"].Elem.(*schema.Resource).Schema
	_, errs := critical["time_range"].ValidateFunc("-72h", "time_range")
	if len(errs) != 0 {
		t.Fatalf("expected no validation errors for -72h, got %v", errs)
	}
}
