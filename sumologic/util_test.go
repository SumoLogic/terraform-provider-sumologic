package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"reflect"
	"testing"
)

// Based on https://github.com/hashicorp/terraform-provider-aws/blob/232c3ed9c9d18aab6f6b70672d7c46d44c75a52e/internal/verify/diff_test.go#L62.
func TestSuppressTimeDiff(t *testing.T) {
	testCases := []struct {
		old        string
		new        string
		equivalent bool
		relative   bool
	}{
		{
			old:        "7d",
			new:        "1w",
			equivalent: true,
			relative:   false,
		},
		{
			old:        "-24h",
			new:        "-1d",
			equivalent: true,
			relative:   false,
		},
		{
			old:        "60m",
			new:        "1h",
			equivalent: true,
			relative:   false,
		},
		{
			old:        "60s",
			new:        "1m",
			equivalent: true,
			relative:   false,
		},
		{
			old:        "95s",
			new:        "2m",
			equivalent: false,
			relative:   false,
		},
		{
			old:        "3h",
			new:        "1d",
			equivalent: false,
			relative:   false,
		},
		{
			old:        "24d",
			new:        "1w",
			equivalent: false,
			relative:   false,
		},
		{
			old:        "25h",
			new:        "1d",
			equivalent: false,
			relative:   false,
		},
		{
			old:        "60h",
			new:        "1d",
			equivalent: false,
			relative:   false,
		},
		{
			old:        "-48h",
			new:        "-2d",
			equivalent: true,
			relative:   false,
		},
		{
			old:        "0s",
			new:        "0d",
			equivalent: true,
			relative:   false,
		},
		{
			old:        "-60s",
			new:        "1m",
			equivalent: true,
			relative:   false,
		},
		{
			old:        "",
			new:        "-0s",
			equivalent: false,
			relative:   false,
		},
		{
			old:        "-0s",
			new:        "",
			equivalent: false,
			relative:   false,
		},
		{
			old:        "0",
			new:        "0s",
			equivalent: false,
			relative:   false,
		},
		{
			old:        "0s",
			new:        "-0s",
			equivalent: true,
			relative:   false,
		},
		{
			old:        "-1h30m",
			new:        "90m",
			equivalent: true,
			relative:   false,
		},
		{
			old:        "-1h30m",
			new:        "90m",
			equivalent: false,
			relative:   true,
		},
		{
			old:        "-1h30m",
			new:        "-90m",
			equivalent: true,
			relative:   false,
		},
		{
			old:        "1h30m2h4s2h",
			new:        "4s30m3h1h1h",
			equivalent: true,
			relative:   false,
		},
		{
			old:        "1h30m",
			new:        "30m1h",
			equivalent: true,
			relative:   false,
		},
		{
			old:        "30min",
			new:        "30m",
			equivalent: false,
			relative:   false,
		},
		{
			old:        "1hour",
			new:        "1h",
			equivalent: false,
			relative:   false,
		},
		{
			old:        "-1second",
			new:        "-1s",
			equivalent: false,
			relative:   false,
		},
		{
			old:        "-1s",
			new:        "-1s",
			equivalent: true,
			relative:   true,
		},
		{
			old:        "-1s",
			new:        "1s",
			equivalent: false,
			relative:   true,
		},
		{
			old:        "-1s",
			new:        "1s",
			equivalent: true,
			relative:   false,
		},
		{
			old:        "-60m",
			new:        "1h",
			equivalent: false,
			relative:   true,
		},
		{
			old:        "-60m",
			new:        "1h",
			equivalent: true,
			relative:   false,
		},
	}

	for i, tc := range testCases {
		value := SuppressEquivalentTimeDiff(tc.relative)("test_time_property", tc.old, tc.new, nil)

		if tc.equivalent && !value {
			t.Fatalf("Expected test case %d to be equivalent. Old value %s, new value %s", i, tc.old, tc.new)
		}

		if !tc.equivalent && value {
			t.Fatalf("Expected test case %d to not be equivalent. Old value %s, new value %s", i, tc.old, tc.new)
		}
	}
}

func TestRemoveEmptyValues(t *testing.T) {
	inputJsonStr := `{
			"leaf": "b",
			"leafNull": null,
			"emptyList": [],
			"list": [
				1,
				2
			],
			"listOfObjects": [
				{
					"a": 1
				},
				{
					"b": null,
					"c": 3
				}
			],
			"map": {
				"e": 5,
				"a": null,
				"nestedMap": {
					"t": 6,
					"g": null,
					"nestedListEmpty": [],
					"nestedList": [
						1,
						2
					]
				}
			}
		}`

	cleanedJsonStr := `
		{
			"leaf": "b",
			"list": [
			  1,
			  2
			],
			"listOfObjects": [
			  {
				"a": 1
			  },
			  {
				"c": 3
			  }
			],
			"map": {
			  "e": 5,
			  "nestedMap": {
				"nestedList": [
				  1,
				  2
				],
				"t": 6
			  }
			}
		  }
		`

	inputMapObject, _ := structure.ExpandJsonFromString(inputJsonStr)
	removeEmptyValues(inputMapObject)
	cleanedMapObject, _ := structure.ExpandJsonFromString(cleanedJsonStr)
	isEqual := reflect.DeepEqual(inputMapObject, cleanedMapObject)

	if !isEqual {
		processedJsonStr, _ := structure.FlattenJsonToString(inputMapObject)
		t.Fatal("Expected json after removing empty values:", cleanedJsonStr, "but was:", processedJsonStr)
	}
}
