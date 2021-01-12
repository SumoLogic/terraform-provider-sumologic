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
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeList,
								Elem: &schema.Schema{
									Type: schema.TypeString,
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
				ValidateFunc: validation.StringInSlice([]string{"Light", "Dark"}, false),
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
	panelSchema := getPanelBaseSchema()

	panelSchema["container_panel"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: getContainerPanelSchema(),
		},
	}

	return panelSchema
}

func getPanelBaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// TODO what is use of "id" field in panel?
		"id": {
			Type:     schema.TypeString,
			Optional: true,
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
		// Need to be set for EventsOfInterestScatterPanel only
		"panel_type": {
			Type:     schema.TypeString,
			Optional: true,
		},
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
				Schema: map[string]*schema.Schema{
					"text": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
	}
}

func getContainerPanelSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"layout": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: getLayoutSchema(),
			},
		},
		"panel": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: getPanelBaseSchema(),
			},
		},
		"variable": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: getVariablesSchema(),
			},
		},
		"coloring_rule": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: getColoringRulesSchema(),
			},
		},
	}
}

func getSumoSearchPanelSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
						Type: schema.TypeString,
						// TODO why is key not a required param but value is?
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
		"layout_structures": {
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
			Optional: true,
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
		"color_thresholds": {
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
						ValidateFunc: validation.ValidateRFC3339TimeString,
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
	// container_panel is not part of base panel schema to avoid infinite recursion.
	if val, ok := tfPanel["container_panel"].([]interface{}); ok && len(val) == 1 {
		return getContainerPanel(tfPanel)
	} else if len(tfPanel["text_panel"].([]interface{})) == 1 {
		return getTextPanel(tfPanel)
	} else if len(tfPanel["sumo_search_panel"].([]interface{})) == 1 {
		return getSumoSearchPanel(tfPanel)
	}
	return nil
}

func getContainerPanel(tfPanel map[string]interface{}) interface{} {
	var containerPanel ContainerPanel
	containerPanel.PanelType = "ContainerPanel"

	if key, ok := tfPanel["key"].(string); ok {
		containerPanel.Key = key
	}
	if title, ok := tfPanel["title"].(string); ok {
		containerPanel.Title = title
	}
	if visualSettings, ok := tfPanel["visual_settings"].(string); ok {
		containerPanel.VisualSettings = visualSettings
	}
	if consistentVisualSettings, ok := tfPanel["keep_visual_settings_consistent_with_parent"].(bool); ok {
		containerPanel.KeepVisualSettingsConsistentWithParent = consistentVisualSettings
	}

	// container panel specific properties
	tfContainerPanel := tfPanel["container_panel"].([]interface{})[0].(map[string]interface{})
	var panels []interface{}
	if val, ok := tfContainerPanel["panel"]; ok {
		tfPanels := val.([]interface{})
		for _, val := range tfPanels {
			panel := getPanel(val.(map[string]interface{}))
			panels = append(panels, panel)
		}
	}
	containerPanel.Panels = panels

	var layout interface{}
	if val, ok := tfContainerPanel["layout"]; ok {
		tfLayout := val.([]interface{})[0]
		layout = getLayout(tfLayout.(map[string]interface{}))
	}
	containerPanel.Layout = layout

	var variables []Variable
	if val, ok := tfContainerPanel["variable"]; ok {
		tfVariables := val.([]interface{})
		for _, tfVariable := range tfVariables {
			variable := getVariable(tfVariable.(map[string]interface{}))
			variables = append(variables, variable)
		}
	}
	containerPanel.Variables = variables

	var coloringRules []ColoringRule
	if val, ok := tfContainerPanel["coloring_rule"]; ok {
		tfColoringRules := val.([]interface{})
		for _, tfColoringRule := range tfColoringRules {
			coloringRule := getColoringRule(tfColoringRule.(map[string]interface{}))
			coloringRules = append(coloringRules, coloringRule)
		}
	}
	containerPanel.ColoringRules = coloringRules

	return containerPanel
}

func getTextPanel(tfPanel map[string]interface{}) interface{} {
	var textPanel TextPanel
	textPanel.PanelType = "TextPanel"

	if key, ok := tfPanel["key"].(string); ok {
		textPanel.Key = key
	}
	if title, ok := tfPanel["title"].(string); ok {
		textPanel.Title = title
	}
	if visualSettings, ok := tfPanel["visual_settings"].(string); ok {
		textPanel.VisualSettings = visualSettings
	}
	if consistentVisualSettings, ok := tfPanel["keep_visual_settings_consistent_with_parent"].(bool); ok {
		textPanel.KeepVisualSettingsConsistentWithParent = consistentVisualSettings
	}

	// text panel specific properties
	if val, ok := tfPanel["text_panel"].([]interface{}); ok {
		if tfTextPanel, ok := val[0].(map[string]interface{}); ok {
			textPanel.Text = tfTextPanel["text"].(string)
		}
	}
	return textPanel
}

func getSumoSearchPanel(tfPanel map[string]interface{}) interface{} {
	var searchPanel SumoSearchPanel
	searchPanel.PanelType = "SumoSearchPanel"

	if key, ok := tfPanel["key"].(string); ok {
		searchPanel.Key = key
	}
	if title, ok := tfPanel["title"].(string); ok {
		searchPanel.Title = title
	}
	if visualSettings, ok := tfPanel["visual_settings"].(string); ok {
		searchPanel.VisualSettings = visualSettings
	}
	if consistentVisualSettings, ok := tfPanel["keep_visual_settings_consistent_with_parent"].(bool); ok {
		searchPanel.KeepVisualSettingsConsistentWithParent = consistentVisualSettings
	}

	// search panel specific properties
	if val, ok := tfPanel["sumo_search_panel"].([]interface{}); ok {
		tfSearchPanel := val[0].(map[string]interface{})
		if description, ok := tfSearchPanel["description"].(string); ok {
			searchPanel.Description = description
		}
		if val, ok := tfSearchPanel["time_range"]; ok {
			tfTimeRange := val.([]interface{})[0]
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
	}
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
	if val, ok := tfTimeRange["complete_literal_time_range"].([]interface{}); ok {
		if literalRange, ok := val[0].(map[string]interface{}); ok {
			return CompleteLiteralTimeRange{
				Type:      "CompleteLiteralTimeRange",
				RangeName: literalRange["range_name"].(string),
			}
		}
	} else if val, ok := tfTimeRange["bounded_time_range"].([]interface{}); ok {
		if boundedRange, ok := val[0].(map[string]interface{}); ok {
			boundaryStart := boundedRange["from"].([]interface{})[0]
			var boundaryEnd interface{}
			if val, ok := boundedRange["to"].([]interface{}); ok {
				boundaryEnd = val[0]
			}

			return BeginBoundedTimeRange{
				Type: "BeginBoundedTimeRange",
				From: getTimeRangeBoundary(boundaryStart.(map[string]interface{})),
				To:   getTimeRangeBoundary(boundaryEnd.(map[string]interface{})),
			}
		}
	}

	return nil
}

func getTimeRangeBoundary(tfRangeBoundary map[string]interface{}) interface{} {
	if val, ok := tfRangeBoundary["epoch_time_range"].([]interface{}); ok {
		if epochBoundary, ok := val[0].(map[string](interface{})); ok {
			return EpochTimeRangeBoundary{
				Type:        "EpochTimeRangeBoundary",
				EpochMillis: epochBoundary["epoch_millis"].(int64),
			}
		}
	} else if val, ok := tfRangeBoundary["iso8601_time_range"].([]interface{}); ok {
		if iso8601Boundary, ok := val[0].(map[string](interface{})); ok {
			return Iso8601TimeRangeBoundary{
				Type:        "Iso8601TimeRangeBoundary",
				Iso8601Time: iso8601Boundary["iso8601_time"].(string),
			}
		}
	} else if val, ok := tfRangeBoundary["literal_time_range"].([]interface{}); ok {
		if literalBoundary, ok := val[0].(map[string](interface{})); ok {
			return LiteralTimeRangeBoundary{
				Type:      "LiteralTimeRangeBoundary",
				RangeName: literalBoundary["range_name"].(string),
			}
		}
	} else if val, ok := tfRangeBoundary["relative_time_range"].([]interface{}); ok {
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
	// data is a map from string to list of strings
	if val, ok := tfTopologyLabel["data"].(map[string]interface{}); ok {
		labelMap := make(map[string][]string)
		for k, v := range val {
			labelMap[k] = v.([]string)
		}
		return &TopologyLabel{
			Data: labelMap,
		}
	}
	return nil
}

func getLayout(tfLayout map[string]interface{}) interface{} {
	if val, ok := tfLayout["grid"].([]interface{}); ok {
		if gridLayout, ok := val[0].(map[string]interface{}); ok {
			if tfStructures, ok := gridLayout["layout_structures"].([]interface{}); ok {
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
	if len(tfSourceDef["log_query_variable_source_definition"].([]interface{})) == 1 {
		val := tfSourceDef["log_query_variable_source_definition"].([]interface{})
		logQuerySourceDef := val[0].(map[string]interface{})
		return LogQueryVariableSourceDefinition{
			VariableSourceType: "LogQueryVariableSourceDefinition",
			Query:              logQuerySourceDef["query"].(string),
			Field:              logQuerySourceDef["field"].(string),
		}
	} else if len(tfSourceDef["metadata_variable_source_definition"].([]interface{})) == 1 {
		val := tfSourceDef["metadata_variable_source_definition"].([]interface{})
		metadataSourceDef := val[0].(map[string]interface{})
		return MetadataVariableSourceDefinition{
			VariableSourceType: "MetadataVariableSourceDefinition",
			Filter:             metadataSourceDef["filter"].(string),
			Key:                metadataSourceDef["key"].(string),
		}
	} else if len(tfSourceDef["csv_variable_source_definition"].([]interface{})) == 1 {
		val := tfSourceDef["csv_variable_source_definition"].([]interface{})
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

	tfColorThresholds := tfColoringRule["color_thresholds"].([]interface{})
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

	// TODO: Set rest of the fields
	log.Printf("=====================================================================n")
	log.Printf("title: %+v\n", d.Get("title"))
	log.Printf("description: %+v\n", d.Get("description"))
	log.Printf("folder_id: %+v\n", d.Get("folder_id"))
	log.Printf("=====================================================================n")
	return nil
}

func resourceSumologicDashboardCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	if d.Id() == "" {
		dashboard := resourceToDashboard(d)
		log.Printf("=====================================================================n")
		log.Printf("Creating dashboard: %+v\n", dashboard)
		log.Printf("=====================================================================n")

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
	log.Printf("=====================================================================n")
	log.Printf("Read dashboard: %+v\n", dashboard)
	log.Printf("=====================================================================n")
	if err != nil {
		return err
	}

	if dashboard == nil {
		log.Printf("[WARN] Dashboard not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	return setDashboard(d, dashboard)
}

func resourceSumologicDashboardDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	return c.DeleteDashboard(d.Id())
}

func resourceSumologicDashboardUpdate(d *schema.ResourceData, meta interface{}) error {
	dashboard := resourceToDashboard(d)
	log.Printf("=====================================================================n")
	log.Printf("Updating dashboard: %+v\n", dashboard)
	log.Printf("=====================================================================n")

	c := meta.(*Client)
	err := c.UpdateDashboard(dashboard)

	if err != nil {
		return err
	}

	return resourceSumologicDashboardRead(d, meta)
}
