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
			old:        "-24h",
			new:        "-1d",
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
			old:        "-48h",
			new:        "-2d",
			equivalent: true,
		},
		{
			old:        "0s",
			new:        "0d",
			equivalent: true,
		},
		{
			old:        "-60s",
			new:        "1m",
			equivalent: false,
		},
		{
			old:        "",
			new:        "-0s",
			equivalent: false,
		},
		{
			old:        "-0s",
			new:        "",
			equivalent: false,
		},
		{
			old:        "0",
			new:        "0s",
			equivalent: false,
		},
		{
			old:        "0s",
			new:        "-0s",
			equivalent: true,
		},
		{
			old:        "-1h30m",
			new:        "90m",
			equivalent: false,
		},
		{
			old:        "-1h30m",
			new:        "-90m",
			equivalent: true,
		},
		{
			old:        "1h30m2h4s2h",
			new:        "4s30m3h1h1h",
			equivalent: true,
		},
		{
			old:        "1h30m",
			new:        "30m1h",
			equivalent: true,
		},
		{
			old:        "30min",
			new:        "30m",
			equivalent: false,
		},
		{
			old:        "1hour",
			new:        "1h",
			equivalent: false,
		},
		{
			old:        "-1second",
			new:        "-1s",
			equivalent: false,
		},
	}

	for i, tc := range testCases {
		value := SuppressEquivalentTimeDiff("test_time_property", tc.old, tc.new, nil)

		if tc.equivalent && !value {
			t.Fatalf("Expected test case %d to be equivalent. Old value %s, new value %s", i, tc.old, tc.new)
		}

		if !tc.equivalent && value {
			t.Fatalf("Expected test case %d to not be equivalent. Old value %s, new value %s", i, tc.old, tc.new)
		}
	}
}
