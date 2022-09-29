package sumologic

import (
	"testing"
)

// Based on https://github.com/hashicorp/terraform-provider-aws/blob/232c3ed9c9d18aab6f6b70672d7c46d44c75a52e/internal/verify/diff_test.go#L62.
func TestSuppressTimeDiff(t *testing.T) {
	testCases := []struct {
		old        string
		new        string
		equivalent bool
	}{
		{
			old:        "7d",
			new:        "1w",
			equivalent: true,
		},
		{
			old:        "24h",
			new:        "1d",
			equivalent: true,
		},
		{
			old:        "60m",
			new:        "1h",
			equivalent: true,
		},
		{
			old:        "60s",
			new:        "1m",
			equivalent: true,
		},
		{
			old:        "95s",
			new:        "2m",
			equivalent: false,
		},
		{
			old:        "3h",
			new:        "1d",
			equivalent: false,
		},
		{
			old:        "24d",
			new:        "1w",
			equivalent: false,
		},
		{
			old:        "25h",
			new:        "1d",
			equivalent: false,
		},
		{
			old:        "60h",
			new:        "1d",
			equivalent: false,
		},
		{
			old:        "48h",
			new:        "2d",
			equivalent: true,
		},
		{
			old:        "0s",
			new:        "0d",
			equivalent: true,
		},
	}

	for i, tc := range testCases {
		value := SuppressEquivalentTimeDiff("test_time_property", tc.old, tc.new, nil)

		if tc.equivalent && !value {
			t.Fatalf("Expected test case %d to be equivalent. Old value %s, new value %s", i, tc.old, tc.new)
		}

		if !tc.equivalent && value {
			t.Fatalf("Expected test case %d to not be equivalen. Old value %s, new value %s", i, tc.old, tc.new)
		}
	}
}
