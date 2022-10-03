package sumologic

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode"
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
						DiffSuppressFunc: SuppressEquivalentTimeDiff,
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

func SuppressEquivalentTimeDiff(k, oldValue, newValue string, d *schema.ResourceData) (shouldSuppress bool) {
	defer func() {
		if err := recover(); err != nil {
			shouldSuppress = false
			log.Printf("[WARN] Time value could not be converted to seconds. Details: %s", err)
		}
	}()

	return getTimeInSeconds(newValue) == getTimeInSeconds(oldValue)
}

func getTimeInSeconds(timeValue string) int64 {
	// For creating new resources
	if timeValue == "" {
		timeValue = "0s"
	}

	if !regexp.MustCompile(`^-?((\d)+[smhdw])+$`).Match([]byte(timeValue)) {
		panic(fmt.Sprintf("Value %s is not in correct time value format.", timeValue))
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
			seconds += getSingleTimeValueInSeconds(value.String())
			value.Reset()
		}
	}

	if negative {
		seconds *= -1
	}

	return seconds
}

func getSingleTimeValueInSeconds(timeValue string) int64 {
	time, err := strconv.Atoi(timeValue[:len(timeValue)-1])
	var seconds int64 = int64(time)
	if err != nil {
		panic(fmt.Sprintf("Error when converting to integer from %s", timeValue[:len(timeValue)-1]))
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
		panic(fmt.Sprintf("Only [smhdw] time units are supported, but got %s", timeValue[len(timeValue)-1:]))
	}

	return seconds
}
