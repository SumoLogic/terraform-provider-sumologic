package sumologic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var (
	validCompleteLiteralTimeRangeValues = []string{
		"today",
		"yesterday",
		"previous_week",
		"previous_month",
	}

	validLiteralTimeRangeValues = []string{
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

	validMetricsAggregationDataValues = []string{
		"Count",
		"Minimum",
		"Maximum",
		"Sum",
		"Average",
		"None",
	}
)

func resourceSumologicDashboard() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicDashboardCreate,
		Read:   resourceSumologicDashboardRead,
		Delete: resourceSumologicDashboardDelete,
		Update: resourceSumologicDashboardUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"title": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"folder_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"topology_label_map": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"label": {
										Type:     schema.TypeString,
										Required: true,
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"refresh_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"time_range": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: getTimeRangeSchema(),
				},
			},
			"panel": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: getPanelSchema(),
				},
			},
			"layout": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: getLayoutSchema(),
				},
			},
			"variable": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: getVariablesSchema(),
				},
			},
			"theme": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Light", "Dark"}, true),
				Default:      "Light",
			},
			// TODO Do we need this field in terraform?
			"coloring_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: getColoringRulesSchema(),
				},
			},
		},
	}
}

func getPanelSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"sumo_search_panel": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: getSumoSearchPanelSchema(),
			},
		},
		"text_panel": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: getTextPanelSchema(),
			},
		},
	}
}

func getPanelBaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"title": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"visual_settings": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"keep_visual_settings_consistent_with_parent": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
	}
}

func getTextPanelSchema() map[string]*schema.Schema {
	panelSchema := getPanelBaseSchema()

	textPanelSchema := map[string]*schema.Schema{
		"text": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	for k, v := range textPanelSchema {
		panelSchema[k] = v
	}

	return panelSchema
}

func getSumoSearchPanelSchema() map[string]*schema.Schema {
	panelSchema := getPanelBaseSchema()

	searchPanelSchema := map[string]*schema.Schema{
		"query": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: getSumoSearchPanelQuerySchema(),
			},
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"time_range": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: getTimeRangeSchema(),
			},
		},
		"coloring_rule": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: getColoringRulesSchema(),
			},
		},
		"linked_dashboard": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeString,
						Required: true,
					},
					"relative_path": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"include_time_range": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},
					"include_variables": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},
	}

	for k, v := range searchPanelSchema {
		panelSchema[k] = v
	}

	return panelSchema
}

func getSumoSearchPanelQuerySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"query_string": {
			Type:     schema.TypeString,
			Required: true,
		},
		"query_type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"Logs", "Metrics"}, false),
		},
		"query_key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"metrics_query_mode": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"Basic", "Advanced"}, false),
		},
		"metrics_query_data": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: getMetricsQueryDataSchema(),
			},
		},
	}
}

func getMetricsQueryDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metric": {
			Type:     schema.TypeString,
			Required: true,
		},
		"aggregation_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(validMetricsAggregationDataValues, false),
		},
		"group_by": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"filter": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Required: true,
					},
					"value": {
						Type:     schema.TypeString,
						Required: true,
					},
					"negation": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"operator": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"operator_name": {
						Type:     schema.TypeString,
						Required: true,
					},
					"parameter": {
						Type:     schema.TypeList,
						Required: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Type:     schema.TypeString,
									Required: true,
								},
								"value": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func getLayoutSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"grid": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: getGridLayoutSchema(),
			},
		},
	}
}

func getGridLayoutSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"layout_structure": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Required: true,
					},
					"structure": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
	}
}

func getVariablesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 256),
		},
		"display_name": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(0, 256),
		},
		"default_value": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"source_definition": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: getSourceDefinitionSchema(),
			},
		},
		"allow_multi_select": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"include_all_option": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"hide_from_ui": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}
}

func getSourceDefinitionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"log_query_variable_source_definition": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"query": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 65536),
					},
					"field": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 65536),
					},
				},
			},
		},
		"metadata_variable_source_definition": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"filter": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 65536),
					},
					"key": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
		"csv_variable_source_definition": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"values": {
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 65536),
					},
				},
			},
		},
	}
}

func getColoringRulesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"scope": {
			Type:     schema.TypeString,
			Required: true,
		},
		"single_series_aggregate_function": {
			Type:     schema.TypeString,
			Required: true,
		},
		"multiple_series_aggregate_function": {
			Type:     schema.TypeString,
			Required: true,
		},
		"color_threshold": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"color": {
						Type:     schema.TypeString,
						Required: true,
					},
					"min": {
						Type:     schema.TypeFloat,
						Optional: true,
					},
					"max": {
						Type:     schema.TypeFloat,
						Optional: true,
					},
				},
			},
		},
	}
}

func getTimeRangeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"complete_literal_time_range": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: getCompleteLiteralTimeRangeSchema(),
			},
		},
		"begin_bounded_time_range": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: getBeginBoundedTimeRangeSchema(),
			},
		},
	}
}

func getCompleteLiteralTimeRangeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"range_name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(validCompleteLiteralTimeRangeValues, false),
		},
	}
}

func getBeginBoundedTimeRangeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"from": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: getTimeRangeBoundarySchema(),
			},
		},
		"to": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: getTimeRangeBoundarySchema(),
			},
		},
	}
}

func getTimeRangeBoundarySchema() map[string]*schema.Schema {
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
						ValidateFunc: validation.StringInSlice(validLiteralTimeRangeValues, false),
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
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
	}
}

func resourceToDashboard(d *schema.ResourceData) Dashboard {
	var timeRange interface{}
	if val, ok := d.GetOk("time_range"); ok {
		tfTimeRange := val.([]interface{})[0]
		timeRange = getTimeRange(tfTimeRange.(map[string]interface{}))
	}

	var topologyLabel *TopologyLabel
	if val, ok := d.GetOk("topology_label_map"); ok {
		tfTopologyLabel := val.([]interface{})[0]
		topologyLabel = getTopologyLabel(tfTopologyLabel.(map[string]interface{}))
	}

	var panels []interface{}
	if val, ok := d.GetOk("panel"); ok {
		tfPanels := val.([]interface{})
		for _, tfPanel := range tfPanels {
			panel := getPanel(tfPanel.(map[string]interface{}))
			panels = append(panels, panel)
		}
	}

	var layout interface{}
	if val, ok := d.GetOk("layout"); ok {
		tfLayout := val.([]interface{})[0]
		layout = getLayout(tfLayout.(map[string]interface{}))
	}

	var variables []Variable
	if val, ok := d.GetOk("variable"); ok {
		tfVariables := val.([]interface{})
		for _, tfVariable := range tfVariables {
			variable := getVariable(tfVariable.(map[string]interface{}))
			variables = append(variables, variable)
		}
	}

	var coloringRules []ColoringRule
	if val, ok := d.GetOk("coloring_rule"); ok {
		tfColoringRules := val.([]interface{})
		for _, tfColoringRule := range tfColoringRules {
			coloringRule := getColoringRule(tfColoringRule.(map[string]interface{}))
			coloringRules = append(coloringRules, coloringRule)
		}
	}

	return Dashboard{
		ID:               d.Id(),
		Title:            d.Get("title").(string),
		Description:      d.Get("description").(string),
		FolderId:         d.Get("folder_id").(string),
		TopologyLabelMap: topologyLabel,
		RefreshInterval:  d.Get("refresh_interval").(int),
		TimeRange:        timeRange,
		Panels:           panels,
		Layout:           layout,
		Variables:        variables,
		Theme:            d.Get("theme").(string),
		ColoringRules:    coloringRules,
	}
}

func getPanel(tfPanel map[string]interface{}) interface{} {
	if val := tfPanel["text_panel"].([]interface{}); len(val) == 1 {
		if tfTextPanel, ok := val[0].(map[string]interface{}); ok {
			return getTextPanel(tfTextPanel)
		}
	} else if val := tfPanel["sumo_search_panel"].([]interface{}); len(val) == 1 {
		if tfSearchPanel, ok := val[0].(map[string]interface{}); ok {
			return getSumoSearchPanel(tfSearchPanel)
		}
	}
	return nil
}

func getTextPanel(tfTextPanel map[string]interface{}) interface{} {
	var textPanel TextPanel
	textPanel.PanelType = "TextPanel"

	textPanel.Key = tfTextPanel["key"].(string)
	if title, ok := tfTextPanel["title"].(string); ok {
		textPanel.Title = title
	}
	if visualSettings, ok := tfTextPanel["visual_settings"].(string); ok {
		textPanel.VisualSettings = visualSettings
	}
	if consistentVisualSettings, ok := tfTextPanel["keep_visual_settings_consistent_with_parent"].(bool); ok {
		textPanel.KeepVisualSettingsConsistentWithParent = consistentVisualSettings
	}

	// text panel specific properties
	textPanel.Text = tfTextPanel["text"].(string)
	return textPanel
}

func getSumoSearchPanel(tfSearchPanel map[string]interface{}) interface{} {
	var searchPanel SumoSearchPanel
	searchPanel.PanelType = "SumoSearchPanel"

	searchPanel.Key = tfSearchPanel["key"].(string)
	if title, ok := tfSearchPanel["title"].(string); ok {
		searchPanel.Title = title
	}
	if visualSettings, ok := tfSearchPanel["visual_settings"].(string); ok {
		searchPanel.VisualSettings = visualSettings
	}
	if consistentVisualSettings, ok := tfSearchPanel["keep_visual_settings_consistent_with_parent"].(bool); ok {
		searchPanel.KeepVisualSettingsConsistentWithParent = consistentVisualSettings
	}

	// search panel specific properties
	if description, ok := tfSearchPanel["description"].(string); ok {
		searchPanel.Description = description
	}
	if val := tfSearchPanel["time_range"].([]interface{}); len(val) == 1 {
		tfTimeRange := val[0]
		searchPanel.TimeRange = getTimeRange(tfTimeRange.(map[string]interface{}))
	}

	tfQueries := tfSearchPanel["query"].([]interface{})
	var queries []SearchPanelQuery
	for _, tfQuery := range tfQueries {
		query := getSearchPanelQuery(tfQuery.(map[string]interface{}))
		queries = append(queries, query)
	}
	searchPanel.Queries = queries

	tfColoringRules := tfSearchPanel["coloring_rule"].([]interface{})
	var rules []ColoringRule
	for _, tfQuery := range tfColoringRules {
		rule := getColoringRule(tfQuery.(map[string]interface{}))
		rules = append(rules, rule)
	}
	searchPanel.ColoringRules = rules

	tfLinkedDashboards := tfSearchPanel["linked_dashboard"].([]interface{})
	var linkedDashboards []LinkedDashboard
	for _, tfLinkedDashboard := range tfLinkedDashboards {
		linkedDashboard := getLinkedDashboard(tfLinkedDashboard.(map[string]interface{}))
		linkedDashboards = append(linkedDashboards, linkedDashboard)
	}
	searchPanel.LinkedDashboards = linkedDashboards

	return searchPanel
}

func getSearchPanelQuery(tfQuery map[string]interface{}) SearchPanelQuery {
	var query SearchPanelQuery

	query.QueryString = tfQuery["query_string"].(string)
	query.QueryType = tfQuery["query_type"].(string)
	query.QueryKey = tfQuery["query_key"].(string)

	if val, ok := tfQuery["metrics_query_mode"]; ok {
		query.MetricsQueryMode = val.(string)
	}
	if val, ok := tfQuery["metrics_query_data"]; ok {
		if tfQueryData := val.([]interface{}); len(tfQueryData) == 1 {
			query.MetricsQueryData = getMetricsQueryData(tfQueryData[0].(map[string]interface{}))
		}
	}

	return query
}

func getMetricsQueryData(tfQueryData map[string]interface{}) *MetricsQueryData {
	var queryData MetricsQueryData

	queryData.Metric = tfQueryData["metric"].(string)

	if val, ok := tfQueryData["aggregation_type"]; ok {
		queryData.AggregationType = val.(string)
	}
	if val, ok := tfQueryData["group_by"]; ok {
		queryData.GroupBy = val.(string)
	}

	tfQueryFilters := tfQueryData["filter"].([]interface{})
	var filters []MetricsQueryFilter
	for _, tfQueryFilter := range tfQueryFilters {
		filter := getMetricsQueryFilter(tfQueryFilter.(map[string]interface{}))
		filters = append(filters, filter)
	}
	queryData.Filters = filters

	tfQueryOperators := tfQueryData["operator"].([]interface{})
	var operators []MetricsQueryOperator
	for _, tfQueryOperator := range tfQueryOperators {
		operator := getMetricsQueryOperator(tfQueryOperator.(map[string]interface{}))
		operators = append(operators, operator)
	}
	queryData.Operators = operators

	return &queryData
}

func getMetricsQueryFilter(tfQueryFilter map[string]interface{}) MetricsQueryFilter {
	var filter MetricsQueryFilter

	filter.Key = tfQueryFilter["key"].(string)
	filter.Value = tfQueryFilter["value"].(string)
	if val, ok := tfQueryFilter["negation"]; ok {
		filter.Negation = val.(bool)
	}

	return filter
}

func getMetricsQueryOperator(tfQueryOperator map[string]interface{}) MetricsQueryOperator {
	var operator MetricsQueryOperator

	operator.Name = tfQueryOperator["operator_name"].(string)

	tfQueryParameters := tfQueryOperator["parameter"].([]interface{})
	var parameters []MetricsQueryOperatorParameter
	for _, val := range tfQueryParameters {
		tfQueryParameter := val.(map[string]interface{})
		parameter := MetricsQueryOperatorParameter{
			Key:   tfQueryParameter["key"].(string),
			Value: tfQueryParameter["value"].(string),
		}
		parameters = append(parameters, parameter)
	}
	operator.Parameters = parameters

	return operator
}

func getTimeRange(tfTimeRange map[string]interface{}) interface{} {
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
				From: getTimeRangeBoundary(boundaryStart),
				To:   getTimeRangeBoundary(boundaryEnd),
			}
		}
	}

	return nil
}

func getTimeRangeBoundary(tfRangeBoundary map[string]interface{}) interface{} {
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

func getTopologyLabel(tfTopologyLabel map[string]interface{}) *TopologyLabel {
	if items := tfTopologyLabel["data"].([]interface{}); len(items) == 1 {
		labelMap := make(map[string][]string)
		for _, item := range items {
			dataItem := item.(map[string]interface{})
			key := dataItem["label"].(string)
			itemValues := dataItem["values"].([]interface{})
			values := make([]string, len(itemValues))
			for i := range itemValues {
				values[i] = itemValues[i].(string)
			}
			labelMap[key] = values
		}
		return &TopologyLabel{
			Data: labelMap,
		}
	}
	return nil
}

func getLayout(tfLayout map[string]interface{}) interface{} {
	if val := tfLayout["grid"].([]interface{}); len(val) == 1 {
		if gridLayout, ok := val[0].(map[string]interface{}); ok {
			if tfStructures, ok := gridLayout["layout_structure"].([]interface{}); ok {
				var structures []LayoutStructure
				for _, v := range tfStructures {
					if tfStructure, ok := v.(map[string]interface{}); ok {
						structure := LayoutStructure{
							Key:       tfStructure["key"].(string),
							Structure: tfStructure["structure"].(string),
						}
						structures = append(structures, structure)
					}
				}

				return GridLayout{
					LayoutType:       "Grid",
					LayoutStructures: structures,
				}
			}
		}
	}
	return nil
}

func getVariable(tfVariable map[string]interface{}) Variable {
	var variable Variable

	if val, ok := tfVariable["id"]; ok {
		variable.Id = val.(string)
	}
	variable.Name = tfVariable["name"].(string)
	if val, ok := tfVariable["display_name"]; ok {
		variable.DisplayName = val.(string)
	}
	if val, ok := tfVariable["default_value"]; ok {
		variable.DefaultValue = val.(string)
	}
	if val, ok := tfVariable["source_definition"]; ok {
		tfSourceDef := val.([]interface{})[0]
		variable.SourceDefinition = getSourceDefinition(tfSourceDef.(map[string]interface{}))
	}
	if val, ok := tfVariable["allow_multi_select"]; ok {
		variable.AllowMultiSelect = val.(bool)
	}
	if val, ok := tfVariable["include_all_option"]; ok {
		variable.IncludeAllOption = val.(bool)
	}
	if val, ok := tfVariable["hide_from_ui"]; ok {
		variable.HideFromUI = val.(bool)
	}

	return variable
}

func getSourceDefinition(tfSourceDef map[string]interface{}) interface{} {
	if val := tfSourceDef["log_query_variable_source_definition"].([]interface{}); len(val) == 1 {
		logQuerySourceDef := val[0].(map[string]interface{})
		return LogQueryVariableSourceDefinition{
			VariableSourceType: "LogQueryVariableSourceDefinition",
			Query:              logQuerySourceDef["query"].(string),
			Field:              logQuerySourceDef["field"].(string),
		}
	} else if val := tfSourceDef["metadata_variable_source_definition"].([]interface{}); len(val) == 1 {
		metadataSourceDef := val[0].(map[string]interface{})
		return MetadataVariableSourceDefinition{
			VariableSourceType: "MetadataVariableSourceDefinition",
			Filter:             metadataSourceDef["filter"].(string),
			Key:                metadataSourceDef["key"].(string),
		}
	} else if val := tfSourceDef["csv_variable_source_definition"].([]interface{}); len(val) == 1 {
		csvSourceDef := val[0].(map[string]interface{})
		return CsvVariableSourceDefinition{
			VariableSourceType: "CsvVariableSourceDefinition",
			Values:             csvSourceDef["values"].(string),
		}
	}
	return nil
}

func getColoringRule(tfColoringRule map[string]interface{}) ColoringRule {
	var coloringRule ColoringRule

	coloringRule.Scope = tfColoringRule["scope"].(string)
	coloringRule.SingleSeriesAggregateFunction = tfColoringRule["single_series_aggregate_function"].(string)
	coloringRule.MultipleSeriesAggregateFunction = tfColoringRule["multiple_series_aggregate_function"].(string)

	tfColorThresholds := tfColoringRule["color_threshold"].([]interface{})
	var colorThresholds []ColorThreshold
	for _, val := range tfColorThresholds {
		tfColorThreshold := val.(map[string]interface{})
		colorThreshold := ColorThreshold{
			Color: tfColorThreshold["color"].(string),
		}
		if val, ok := tfColorThreshold["min"]; ok {
			colorThreshold.Min = val.(float64)
		}
		if val, ok := tfColorThreshold["max"]; ok {
			colorThreshold.Max = val.(float64)
		}
		colorThresholds = append(colorThresholds, colorThreshold)
	}
	coloringRule.ColorThresholds = colorThresholds
	return coloringRule
}

func getLinkedDashboard(tfLinkedDashboard map[string]interface{}) LinkedDashboard {
	var linkedDashboard LinkedDashboard

	linkedDashboard.Id = tfLinkedDashboard["id"].(string)

	if val, ok := tfLinkedDashboard["relative_path"]; ok {
		linkedDashboard.RelativePath = val.(string)
	}
	if val, ok := tfLinkedDashboard["include_time_range"]; ok {
		linkedDashboard.IncludeTimeRange = val.(bool)
	}
	if val, ok := tfLinkedDashboard["include_variables"]; ok {
		linkedDashboard.IncludeVariables = val.(bool)
	}

	return linkedDashboard
}

func setDashboard(d *schema.ResourceData, dashboard *Dashboard) error {
	if err := d.Set("title", dashboard.Title); err != nil {
		return err
	}
	if err := d.Set("description", dashboard.Description); err != nil {
		return err
	}
	if err := d.Set("folder_id", dashboard.FolderId); err != nil {
		return err
	}
	if err := d.Set("refresh_interval", dashboard.RefreshInterval); err != nil {
		return err
	}
	if err := d.Set("theme", dashboard.Theme); err != nil {
		return err
	}

	topologyLabel := getTerraformTopologyLabel(dashboard.TopologyLabelMap)
	if err := d.Set("topology_label_map", topologyLabel); err != nil {
		return err
	}

	timeRange := getTerraformTimeRange(dashboard.TimeRange.(map[string]interface{}))
	if err := d.Set("time_range", timeRange); err != nil {
		return err
	}

	panels := getTerraformPanels(dashboard.Panels)
	if err := d.Set("panel", panels); err != nil {
		return err
	}

	layout := getTerraformLayout(dashboard.Layout.(map[string]interface{}))
	if err := d.Set("layout", layout); err != nil {
		return err
	}

	variables := getTerraformVariables(dashboard.Variables)
	if err := d.Set("variable", variables); err != nil {
		return err
	}

	coloringRules := getTerraformColoringRules(dashboard.ColoringRules)
	if err := d.Set("coloring_rule", coloringRules); err != nil {
		return err
	}

	log.Println("=====================================================================")
	log.Printf("title: %+v\n", d.Get("title"))
	log.Printf("description: %+v\n", d.Get("description"))
	log.Printf("folder_id: %+v\n", d.Get("folder_id"))
	log.Printf("topology_label_map: %+v\n", d.Get("topology_label_map"))
	log.Printf("time_range: %+v\n", d.Get("time_range"))
	log.Printf("panel: %+v\n", d.Get("panel"))
	log.Printf("layout: %+v\n", d.Get("layout"))
	log.Printf("variable: %+v\n", d.Get("variable"))
	log.Printf("coloring_rule: %+v\n", d.Get("coloring_rule"))
	log.Println("=====================================================================")
	return nil
}

type TerraformObject [1]map[string]interface{}

func makeTerraformObject() TerraformObject {
	terraformObject := [1]map[string]interface{}{}
	terraformObject[0] = make(map[string]interface{})
	return terraformObject
}

func getTerraformTopologyLabel(topologyLabel *TopologyLabel) []map[string]interface{} {
	// API returns an empty data map if we don't set topologyLabelMap.
	if len(topologyLabel.Data) == 0 {
		return nil
	}

	tfTopologyLabel := make([]map[string]interface{}, 0)
	tfTopologyLabel = append(tfTopologyLabel, make(map[string]interface{}))

	data := topologyLabel.Data
	tfDataItems := make([]map[string]interface{}, 0)
	for label, values := range data {
		tfDataItem := make(map[string]interface{})
		tfDataItem["label"] = label
		tfDataItem["values"] = values
		tfDataItems = append(tfDataItems, tfDataItem)
	}
	tfTopologyLabel[0]["data"] = tfDataItems
	return tfTopologyLabel
}

func getTerraformTimeRange(timeRange map[string]interface{}) []map[string]interface{} {
	tfTimeRange := []map[string]interface{}{}
	tfTimeRange = append(tfTimeRange, make(map[string]interface{}))

	if timeRange["type"] == "BeginBoundedTimeRange" {
		boundedTimeRange := makeTerraformObject()

		from := timeRange["from"].(map[string]interface{})
		rangeBoundary := getTerraformTimeRangeBoundary(from)
		boundedTimeRange[0]["from"] = rangeBoundary

		if to := timeRange["to"]; to != nil {
			rangeBoundary := getTerraformTimeRangeBoundary(to.(map[string]interface{}))
			boundedTimeRange[0]["to"] = rangeBoundary
		}

		tfTimeRange[0]["begin_bounded_time_range"] = boundedTimeRange
	} else if timeRange["type"] == "CompleteLiteralTimeRange" {
		rangeName := timeRange["rangeName"]

		completeLiteralTimeRange := makeTerraformObject()
		completeLiteralTimeRange[0]["range_name"] = rangeName
		tfTimeRange[0]["complete_literal_time_range"] = completeLiteralTimeRange
	}

	return tfTimeRange
}

func getTerraformTimeRangeBoundary(timeRangeBoundary map[string]interface{}) TerraformObject {
	tfTimeRangeBoundary := makeTerraformObject()

	if timeRangeBoundary["type"] == "RelativeTimeRangeBoundary" {
		relativeRange := makeTerraformObject()
		relativeRange[0]["relative_time"] = timeRangeBoundary["relativeTime"]
		tfTimeRangeBoundary[0]["relative_time_range"] = relativeRange
	} else if timeRangeBoundary["type"] == "EpochTimeRangeBoundary" {
		epochRange := makeTerraformObject()
		epochRange[0]["epoch_millis"] = timeRangeBoundary["epochMillis"]
		tfTimeRangeBoundary[0]["epoch_time_range"] = epochRange
	} else if timeRangeBoundary["type"] == "Iso8601TimeRangeBoundary" {
		iso8601Range := makeTerraformObject()
		iso8601Range[0]["iso8601_time"] = timeRangeBoundary["iso8601Time"]
		tfTimeRangeBoundary[0]["iso8601_time_range"] = iso8601Range
	} else if timeRangeBoundary["type"] == "LiteralTimeRangeBoundary" {
		literalRange := makeTerraformObject()
		literalRange[0]["range_name"] = timeRangeBoundary["rangeName"]
		tfTimeRangeBoundary[0]["literal_time_range"] = literalRange
	}

	return tfTimeRangeBoundary
}

func getTerraformPanels(panels []interface{}) []map[string]interface{} {
	tfPanels := make([]map[string]interface{}, len(panels))

	for i, val := range panels {
		panel := val.(map[string]interface{})

		tfPanel := map[string]interface{}{}
		if panel["panelType"] == "TextPanel" {
			tfPanel["text_panel"] = getTerraformTextPanel(panel)
		} else if panel["panelType"] == "SumoSearchPanel" {
			tfPanel["sumo_search_panel"] = getTerraformSearchPanel(panel)
		}

		tfPanels[i] = tfPanel
	}
	return tfPanels
}

func getTerraformTextPanel(textPanel map[string]interface{}) TerraformObject {
	tfTextPanel := makeTerraformObject()

	tfTextPanel[0]["key"] = textPanel["key"]
	if title, ok := textPanel["title"]; ok {
		tfTextPanel[0]["title"] = title
	}
	if visualSettings, ok := textPanel["visualSettings"]; ok {
		tfTextPanel[0]["visual_settings"] = visualSettings
	}
	if keepVisualSettingsConsistentWithParent, ok := textPanel["keepVisualSettingsConsistentWithParent"]; ok {
		tfTextPanel[0]["keep_visual_settings_consistent_with_parent"] = keepVisualSettingsConsistentWithParent
	}
	tfTextPanel[0]["text"] = textPanel["text"]

	return tfTextPanel
}

func getTerraformSearchPanel(searchPanel map[string]interface{}) TerraformObject {
	tfSearchPanel := makeTerraformObject()

	tfSearchPanel[0]["key"] = searchPanel["key"]
	if title, ok := searchPanel["title"]; ok {
		tfSearchPanel[0]["title"] = title
	}
	if visualSettings, ok := searchPanel["visualSettings"]; ok {
		tfSearchPanel[0]["visual_settings"] = visualSettings
	}
	if keepVisualSettingsConsistentWithParent, ok := searchPanel["keepVisualSettingsConsistentWithParent"]; ok {
		tfSearchPanel[0]["keep_visual_settings_consistent_with_parent"] = keepVisualSettingsConsistentWithParent
	}

	tfSearchPanel[0]["query"] = getTerraformSearchPanelQuery(searchPanel["queries"].([]interface{}))
	if description, ok := searchPanel["description"]; ok {
		tfSearchPanel[0]["description"] = description
	}
	if timeRange := searchPanel["timeRange"]; timeRange != nil {
		tfSearchPanel[0]["time_range"] = getTerraformTimeRange(timeRange.(map[string]interface{}))
	}
	if coloringRules := searchPanel["coloringRules"]; coloringRules != nil {
		tfSearchPanel[0]["coloring_rule"] = getTerraformTimeRange(coloringRules.(map[string]interface{}))
	}
	if linkedDashboards := searchPanel["linkedDashboards"]; linkedDashboards != nil {
		tfSearchPanel[0]["linked_dashboard"] = getTerraformLinkedDashboards(linkedDashboards.([]interface{}))
	}

	return tfSearchPanel
}

func getTerraformSearchPanelQuery(queries []interface{}) []map[string]interface{} {
	tfPanelQueries := make([]map[string]interface{}, len(queries))

	for i, val := range queries {
		query := val.(map[string]interface{})
		tfPanelQueries[i] = make(map[string]interface{})
		tfPanelQueries[i]["query_string"] = query["queryString"]
		tfPanelQueries[i]["query_type"] = query["queryType"]
		tfPanelQueries[i]["query_key"] = query["queryKey"]
		if metricsQueryMode, ok := query["metricsQueryMode"]; ok {
			tfPanelQueries[i]["metrics_query_mode"] = metricsQueryMode
		}
		if metricsQueryData, ok := query["metricsQueryData"]; ok && metricsQueryData != nil {
			tfPanelQueries[i]["metrics_query_data"] =
				getTerraformMetricsQueryDataScheme(metricsQueryData.(map[string]interface{}))
		}
	}
	return tfPanelQueries
}

func getTerraformMetricsQueryDataScheme(queryData map[string]interface{}) TerraformObject {
	tfMetricsQueryData := makeTerraformObject()

	tfMetricsQueryData[0]["metric"] = queryData["metric"]
	if aggregationType, ok := queryData["aggregationType"]; ok {
		tfMetricsQueryData[0]["aggregation_type"] = aggregationType
	}
	if groupBy, ok := queryData["groupBy"]; ok {
		tfMetricsQueryData[0]["group_by"] = groupBy
	}

	filters := queryData["filters"].([]interface{})
	tfFilters := make([]map[string]interface{}, len(filters))
	for i, val := range filters {
		filter := val.(map[string]interface{})
		tfFilters[i] = make(map[string]interface{})
		tfFilters[i]["key"] = filter["key"]
		tfFilters[i]["value"] = filter["value"]
		tfFilters[i]["negation"] = filter["negation"]
	}
	tfMetricsQueryData[0]["filter"] = tfFilters

	if val, ok := queryData["operators"]; ok && val != nil {
		operators := val.([]interface{})
		tfOperators := make([]map[string]interface{}, len(operators))
		for i, val := range operators {
			operator := val.(map[string]interface{})
			tfOperators[i] = getTerraformMetricsQueryOperator(operator)
		}
		tfMetricsQueryData[0]["operator"] = tfOperators
	}

	return tfMetricsQueryData
}

func getTerraformMetricsQueryOperator(operator map[string]interface{}) map[string]interface{} {
	tfOperator := make(map[string]interface{})
	tfOperator["operator_name"] = operator["operatorName"]

	parameters := operator["parameters"].([]interface{})
	tfParameters := make([]map[string]interface{}, len(parameters))
	for i, val := range parameters {
		parameter := val.(map[string]interface{})
		tfParameters[i] = make(map[string]interface{})
		tfParameters[i]["key"] = parameter["key"]
		tfParameters[i]["value"] = parameter["value"]
	}
	tfOperator["parameter"] = tfParameters

	return tfOperator
}

func getTerraformLinkedDashboards(dashboards []interface{}) []map[string]interface{} {
	tfLinkedDashboards := make([]map[string]interface{}, len(dashboards))

	for i, val := range dashboards {
		dashboard := val.(map[string]interface{})
		tfLinkedDashboards[i] = make(map[string]interface{})
		tfLinkedDashboards[i]["id"] = dashboard["id"]
		tfLinkedDashboards[i]["relative_path"] = dashboard["relativePath"]
		tfLinkedDashboards[i]["include_time_range"] = dashboard["includeTimeRange"]
		tfLinkedDashboards[i]["include_variables"] = dashboard["includeVariables"]
	}

	return tfLinkedDashboards
}

func getTerraformLayout(layout map[string]interface{}) []map[string]interface{} {
	tfLayout := []map[string]interface{}{}
	tfLayout = append(tfLayout, make(map[string]interface{}))

	if layout["layoutType"] == "Grid" {
		gridLayout := makeTerraformObject()

		layoutStructures := layout["layoutStructures"].([]interface{})
		tfLayoutStructures := make([]map[string]interface{}, len(layoutStructures))
		for i, structure := range layoutStructures {
			tfLayoutStructures[i] = structure.(map[string]interface{})
		}
		gridLayout[0]["layout_structure"] = tfLayoutStructures
		tfLayout[0]["grid"] = gridLayout
	}

	return tfLayout
}

func getTerraformVariables(variables []Variable) []map[string]interface{} {
	tfVariables := make([]map[string]interface{}, len(variables))

	for i, variable := range variables {
		tfVariables[i] = make(map[string]interface{})
		tfVariables[i]["name"] = variable.Name
		tfVariables[i]["display_name"] = variable.DisplayName
		tfVariables[i]["default_value"] = variable.DefaultValue
		tfVariables[i]["allow_multi_select"] = variable.AllowMultiSelect
		tfVariables[i]["include_all_option"] = variable.IncludeAllOption
		tfVariables[i]["hide_from_ui"] = variable.HideFromUI
		tfVariables[i]["source_definition"] =
			getTerraformVariableSourceDefinition(variable.SourceDefinition.(map[string]interface{}))
	}

	return tfVariables
}

func getTerraformVariableSourceDefinition(sourceDefinition map[string]interface{}) TerraformObject {
	tfSourceDefinition := makeTerraformObject()

	if sourceDefinition["variableSourceType"] == "MetadataVariableSourceDefinition" {
		metadataDefinition := makeTerraformObject()
		metadataDefinition[0]["filter"] = sourceDefinition["filter"]
		metadataDefinition[0]["key"] = sourceDefinition["key"]
		tfSourceDefinition[0]["metadata_variable_source_definition"] = metadataDefinition
	} else if sourceDefinition["variableSourceType"] == "CsvVariableSourceDefinition" {
		csvDefinition := makeTerraformObject()
		csvDefinition[0]["values"] = sourceDefinition["values"]
		tfSourceDefinition[0]["csv_variable_source_definition"] = csvDefinition
	} else if sourceDefinition["variableSourceType"] == "LogQueryVariableSourceDefinition" {
		logQueryDefinition := makeTerraformObject()
		logQueryDefinition[0]["query"] = sourceDefinition["query"]
		logQueryDefinition[0]["field"] = sourceDefinition["field"]
		tfSourceDefinition[0]["log_query_variable_source_definition"] = logQueryDefinition
	}

	return tfSourceDefinition
}

func getTerraformColoringRules(coloringRules []ColoringRule) []map[string]interface{} {
	tfColoringRules := make([]map[string]interface{}, len(coloringRules))

	for i, rule := range coloringRules {
		tfColoringRules[i] = make(map[string]interface{})
		tfColoringRules[i]["scope"] = rule.Scope
		tfColoringRules[i]["single_series_aggregate_function"] = rule.SingleSeriesAggregateFunction
		tfColoringRules[i]["multiple_series_aggregate_function"] = rule.MultipleSeriesAggregateFunction

		tfColorThresholds := make([]map[string]interface{}, len(rule.ColorThresholds))
		for j, threshold := range rule.ColorThresholds {
			tfColorThresholds[j] = make(map[string]interface{})
			tfColorThresholds[j]["color"] = threshold.Color
			tfColorThresholds[j]["min"] = threshold.Min
			tfColorThresholds[j]["max"] = threshold.Max
		}
		tfColoringRules[i]["color_threshold"] = tfColorThresholds
	}

	return tfColoringRules
}

func resourceSumologicDashboardCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	if d.Id() == "" {
		dashboard := resourceToDashboard(d)
		log.Println("=====================================================================")
		log.Printf("Creating dashboard: %+v\n", dashboard)
		log.Println("=====================================================================")

		createdDashboard, err := c.CreateDashboard(dashboard)
		if err != nil {
			return err
		}
		d.SetId(createdDashboard.ID)
	}

	return resourceSumologicDashboardRead(d, meta)
}

func resourceSumologicDashboardRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	dashboard, err := c.GetDashboard(id)
	log.Println("=====================================================================")
	log.Printf("Read dashboard: %+v\n", dashboard)
	log.Println("=====================================================================")
	if err != nil {
		return err
	}

	if dashboard == nil {
		log.Printf("[WARN] Dashboard not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	err = setDashboard(d, dashboard)
	return err
}

func resourceSumologicDashboardDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	log.Printf("Deleting dashboard: %+v\n", d.Id())
	return c.DeleteDashboard(d.Id())
}

func resourceSumologicDashboardUpdate(d *schema.ResourceData, meta interface{}) error {
	dashboard := resourceToDashboard(d)
	log.Println("=====================================================================")
	log.Printf("Updating dashboard: %+v\n", dashboard)
	log.Println("=====================================================================")

	c := meta.(*Client)
	err := c.UpdateDashboard(dashboard)

	if err != nil {
		return err
	}

	return resourceSumologicDashboardRead(d, meta)
}
