package sumologic

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var (
	ValidCompleteLiteralTimeRangeValues = []string{
		"today",
		"yesterday",
		"previous_week",
		"previous_month",
	}

	ValidLiteralTimeRangeValues = []string{
		"now",
		"second",
		"minute",
		"hour",
		"day",
		"today",
		"week",
		"month",
		"year",
	}
)

type TerraformObject [1]map[string]interface{}

// TimeRange related structs
type CompleteLiteralTimeRange struct {
	Type      string `json:"type"`
	RangeName string `json:"rangeName"`
}

type BeginBoundedTimeRange struct {
	Type string      `json:"type"`
	From interface{} `json:"from"`
	To   interface{} `json:"to"`
}

type RelativeTimeRangeBoundary struct {
	Type         string `json:"type"`
	RelativeTime string `json:"relativeTime"`
}

type EpochTimeRangeBoundary struct {
	Type        string `json:"type"`
	EpochMillis int64  `json:"epochMillis"`
}

type Iso8601TimeRangeBoundary struct {
	Type        string `json:"type"`
	Iso8601Time string `json:"iso8601Time"`
}

type LiteralTimeRangeBoundary struct {
	Type      string `json:"type"`
	RangeName string `json:"rangeName"`
}

func MakeTerraformObject() TerraformObject {
	terraformObject := [1]map[string]interface{}{}
	terraformObject[0] = make(map[string]interface{})
	return terraformObject
}

func GetTimeRangeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"complete_literal_time_range": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: GetCompleteLiteralTimeRangeSchema(),
			},
		},
		"begin_bounded_time_range": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: GetBeginBoundedTimeRangeSchema(),
			},
		},
	}
}

func GetCompleteLiteralTimeRangeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"range_name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(ValidCompleteLiteralTimeRangeValues, false),
		},
	}
}

func GetBeginBoundedTimeRangeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"from": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: GetTimeRangeBoundarySchema(),
			},
		},
		"to": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: GetTimeRangeBoundarySchema(),
			},
		},
	}
}

func GetTimeRangeBoundarySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"epoch_time_range": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"epoch_millis": {
						Type:     schema.TypeInt,
						Required: true,
					},
				},
			},
		},
		"iso8601_time_range": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"iso8601_time": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.IsRFC3339Time,
					},
				},
			},
		},
		"literal_time_range": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"range_name": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(ValidLiteralTimeRangeValues, false),
					},
				},
			},
		},
		"relative_time_range": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"relative_time": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateFunc:     validation.StringMatch(regexp.MustCompile(`^-?((\d)+[smhdw])+$`), "This value is not in correct format. Example: -2w5d3h"),
						DiffSuppressFunc: SuppressEquivalentTimeDiff(true),
					},
				},
			},
		},
	}
}

func GetTerraformTimeRange(timeRange map[string]interface{}) []map[string]interface{} {
	tfTimeRange := []map[string]interface{}{}
	tfTimeRange = append(tfTimeRange, make(map[string]interface{}))

	if timeRange["type"] == "BeginBoundedTimeRange" {
		boundedTimeRange := MakeTerraformObject()

		from := timeRange["from"].(map[string]interface{})
		rangeBoundary := GetTerraformTimeRangeBoundary(from)
		boundedTimeRange[0]["from"] = rangeBoundary

		if to := timeRange["to"]; to != nil {
			rangeBoundary := GetTerraformTimeRangeBoundary(to.(map[string]interface{}))
			boundedTimeRange[0]["to"] = rangeBoundary
		}

		tfTimeRange[0]["begin_bounded_time_range"] = boundedTimeRange
	} else if timeRange["type"] == "CompleteLiteralTimeRange" {
		rangeName := timeRange["rangeName"]

		completeLiteralTimeRange := MakeTerraformObject()
		completeLiteralTimeRange[0]["range_name"] = rangeName
		tfTimeRange[0]["complete_literal_time_range"] = completeLiteralTimeRange
	}

	return tfTimeRange
}

func GetTerraformTimeRangeBoundary(timeRangeBoundary map[string]interface{}) TerraformObject {
	tfTimeRangeBoundary := MakeTerraformObject()

	if timeRangeBoundary["type"] == "RelativeTimeRangeBoundary" {
		relativeRange := MakeTerraformObject()
		relativeRange[0]["relative_time"] = timeRangeBoundary["relativeTime"]
		tfTimeRangeBoundary[0]["relative_time_range"] = relativeRange
	} else if timeRangeBoundary["type"] == "EpochTimeRangeBoundary" {
		epochRange := MakeTerraformObject()
		epochRange[0]["epoch_millis"] = timeRangeBoundary["epochMillis"]
		tfTimeRangeBoundary[0]["epoch_time_range"] = epochRange
	} else if timeRangeBoundary["type"] == "Iso8601TimeRangeBoundary" {
		iso8601Range := MakeTerraformObject()
		iso8601Range[0]["iso8601_time"] = timeRangeBoundary["iso8601Time"]
		tfTimeRangeBoundary[0]["iso8601_time_range"] = iso8601Range
	} else if timeRangeBoundary["type"] == "LiteralTimeRangeBoundary" {
		literalRange := MakeTerraformObject()
		literalRange[0]["range_name"] = timeRangeBoundary["rangeName"]
		tfTimeRangeBoundary[0]["literal_time_range"] = literalRange
	}

	return tfTimeRangeBoundary
}

func GetTimeRange(tfTimeRange map[string]interface{}) interface{} {
	if val := tfTimeRange["complete_literal_time_range"].([]interface{}); len(val) == 1 {
		if literalRange, ok := val[0].(map[string]interface{}); ok {
			return CompleteLiteralTimeRange{
				Type:      "CompleteLiteralTimeRange",
				RangeName: literalRange["range_name"].(string),
			}
		}
	} else if val := tfTimeRange["begin_bounded_time_range"].([]interface{}); len(val) == 1 {
		if boundedRange, ok := val[0].(map[string]interface{}); ok {
			from := boundedRange["from"].([]interface{})
			boundaryStart := from[0].(map[string]interface{})
			var boundaryEnd map[string]interface{}
			if to := boundedRange["to"].([]interface{}); len(to) == 1 {
				boundaryEnd = to[0].(map[string]interface{})
			}

			return BeginBoundedTimeRange{
				Type: "BeginBoundedTimeRange",
				From: GetTimeRangeBoundary(boundaryStart),
				To:   GetTimeRangeBoundary(boundaryEnd),
			}
		}
	}

	return nil
}

func GetTimeRangeBoundary(tfRangeBoundary map[string]interface{}) interface{} {
	if len(tfRangeBoundary) == 0 {
		return nil
	}

	if val := tfRangeBoundary["epoch_time_range"].([]interface{}); len(val) == 1 {
		if epochBoundary, ok := val[0].(map[string](interface{})); ok {
			return EpochTimeRangeBoundary{
				Type:        "EpochTimeRangeBoundary",
				EpochMillis: int64(epochBoundary["epoch_millis"].(int)),
			}
		}
	} else if val := tfRangeBoundary["iso8601_time_range"].([]interface{}); len(val) == 1 {
		if iso8601Boundary, ok := val[0].(map[string](interface{})); ok {
			return Iso8601TimeRangeBoundary{
				Type:        "Iso8601TimeRangeBoundary",
				Iso8601Time: iso8601Boundary["iso8601_time"].(string),
			}
		}
	} else if val := tfRangeBoundary["literal_time_range"].([]interface{}); len(val) == 1 {
		if literalBoundary, ok := val[0].(map[string](interface{})); ok {
			return LiteralTimeRangeBoundary{
				Type:      "LiteralTimeRangeBoundary",
				RangeName: literalBoundary["range_name"].(string),
			}
		}
	} else if val := tfRangeBoundary["relative_time_range"].([]interface{}); len(val) == 1 {
		if relativeBoundary, ok := val[0].(map[string](interface{})); ok {
			return RelativeTimeRangeBoundary{
				Type:         "RelativeTimeRangeBoundary",
				RelativeTime: relativeBoundary["relative_time"].(string),
			}
		}
	}

	return nil
}

/*
This function returns a function (in 2 variants) determining whether two stringe representation
of a time value can be considered equivalent. One variant compares only absolute values (ignores '-' sign)
and the other compare relative values. For details see util_test.go::TestSuppressTimeDiff.

Some examples (we accept time units for seconds, minutes, hours, days and weeks):
-1h = -60m
1h20m = 80m
60m60m60m1h = 3h30m30m
-60m = 1h (only if we compare absolute values, so with isRelative = false)
1w = 604800s
1h != 61m
2m != 119s
-1h != 1h (only if we compare relative values, so with isRelative = true)
1d = 22h60m3600s
*/
func SuppressEquivalentTimeDiff(isRelative bool) func(k, oldValue, newValue string, d *schema.ResourceData) bool {
	return func(k, oldValue, newValue string, d *schema.ResourceData) bool {
		var handleError = func(err error) bool {
			log.Printf("[WARN] Time value could not be converted to seconds. Details: %s", err.Error())
			return false
		}

		if oldValue == "" || newValue == "" {
			return false
		}

		if oldValInSeconds, oldErr := getTimeInSeconds(oldValue); oldErr != nil {
			return handleError(oldErr)
		} else if newValInSeconds, newErr := getTimeInSeconds(newValue); newErr != nil {
			return handleError(newErr)
		} else {
			if !isRelative {
				oldValInSeconds = abs(oldValInSeconds)
				newValInSeconds = abs(newValInSeconds)
			}

			return oldValInSeconds == newValInSeconds
		}
	}
}

func getTimeInSeconds(timeValue string) (int64, error) {
	if !regexp.MustCompile(`^-?((\d)+[smhdw])+$`).Match([]byte(timeValue)) {
		return 0, fmt.Errorf("value %s is not in correct time value format", timeValue)
	}

	var negative bool = false
	var absTime string = timeValue
	if timeValue[0] == '-' {
		negative = true
		absTime = timeValue[1:]
	}

	seconds := int64(0)
	var value strings.Builder

	for _, ch := range absTime {
		value.WriteRune(ch)
		if !unicode.IsNumber(ch) {
			if valInSeconds, err := getSingleTimeValueInSeconds(value.String()); err != nil {
				return 0, err
			} else {
				seconds += valInSeconds
			}
			value.Reset()
		}
	}

	if negative {
		seconds *= -1
	}

	return seconds, nil
}

func getSingleTimeValueInSeconds(timeValue string) (int64, error) {
	time, err := strconv.Atoi(timeValue[:len(timeValue)-1])
	var seconds int64 = int64(time)
	if err != nil {
		return 0, fmt.Errorf("Error when converting to integer from %s", timeValue[:len(timeValue)-1])
	}

	switch timeValue[len(timeValue)-1:] {
	case "s":
		break
	case "m":
		seconds *= 60
	case "h":
		seconds *= 3600
	case "d":
		seconds *= 86400
	case "w":
		seconds *= 604800
	default:
		return 0, fmt.Errorf("only [smhdw] time units are supported, but got %s", timeValue[len(timeValue)-1:])
	}

	return seconds, nil
}

func abs(value int64) int64 {
	if value < 0 {
		return -value
	}

	return value
}

// traverse map and remove fields with null value and empty lists
func removeEmptyValues(object interface{}) {
	mapObject, isMap := object.(map[string]interface{})
	if isMap {
		for key, value := range mapObject {
			if value == nil {
				delete(mapObject, key)
			}
			listObject, isList := value.([]interface{})
			if isList {
				if len(listObject) == 0 {
					delete(mapObject, key)
				} else {
					for _, listItem := range listObject {
						removeEmptyValues(listItem)
					}
				}
			}
			removeEmptyValues(value)
		}
	}
}

func waitForJob(url string, timeout time.Duration, s *Client) (*Status, error) {
	conf := &resource.StateChangeConf{
		Pending: []string{
			"InProgress",
		},
		Target: []string{
			"Success",
		},
		Refresh: func() (interface{}, string, error) {
			var status Status
			b, err := s.Get(url)
			if err != nil {
				return nil, "", err
			}

			err = json.Unmarshal(b, &status)
			if err != nil {
				return nil, "", err
			}

			if status.Status == "Failed" {
				return status, status.Status, fmt.Errorf("async job failed - %s", status.Error)
			}

			return status, status.Status, nil
		},
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	result, err := conf.WaitForState()
	log.Printf("[DEBUG] Done waiting for job; err: %s, result: %v", err, result)
	if status, ok := result.(Status); ok {
		return &status, err
	} else {
		return nil, err
	}
}
