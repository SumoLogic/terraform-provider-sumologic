package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
)

func resourceSumologicCSEAggregationRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSEAggregationRuleCreate,
		Read:   resourceSumologicCSEAggregationRuleRead,
		Delete: resourceSumologicCSEAggregationRuleDelete,
		Update: resourceSumologicCSEAggregationRuleUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"aggregation_functions": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"function": {
							Type:     schema.TypeString,
							Required: true,
						},
						"arguments": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"description_expression": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"entity_selectors": getEntitySelectorsSchema(),
			"group_by_entity": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"group_by_fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_prototype": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"match_expression": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name_expression": {
				Type:     schema.TypeString,
				Required: true,
			},
			"severity_mapping": getSeverityMappingSchema(),
			"summary_expression": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
			"trigger_expression": {
				Type:     schema.TypeString,
				Required: true,
			},
			"window_size": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSumologicCSEAggregationRuleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSEAggregationRuleGet *CSEAggregationRule
	id := d.Id()

	CSEAggregationRuleGet, err := c.GetCSEAggregationRule(id)
	if err != nil {
		log.Printf("[WARN] CSE Aggregation Rule not found when looking by id: %s, err: %v", id, err)
	}

	if CSEAggregationRuleGet == nil {
		log.Printf("[WARN] CSE Aggregation Rule not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("aggregation_functions", aggregationFunctionsArrayToResource(CSEAggregationRuleGet.AggregationFunctions))
	d.Set("description_expression", CSEAggregationRuleGet.DescriptionExpression)
	d.Set("enabled", CSEAggregationRuleGet.Enabled)
	d.Set("entity_selectors", entitySelectorArrayToResource(CSEAggregationRuleGet.EntitySelectors))
	d.Set("match_expression", CSEAggregationRuleGet.MatchExpression)
	d.Set("group_by_entity", CSEAggregationRuleGet.GroupByEntity)
	d.Set("group_by_fields", CSEAggregationRuleGet.GroupByFields)
	d.Set("is_prototype", CSEAggregationRuleGet.IsPrototype)
	d.Set("match_expression", CSEAggregationRuleGet.MatchExpression)
	d.Set("name", CSEAggregationRuleGet.Name)
	d.Set("name_expression", CSEAggregationRuleGet.NameExpression)
	d.Set("severity_mapping", severityMappingToResource(CSEAggregationRuleGet.SeverityMapping))
	d.Set("summary_expression", CSEAggregationRuleGet.SummaryExpression)
	d.Set("tags", CSEAggregationRuleGet.Tags)
	d.Set("trigger_expression", CSEAggregationRuleGet.TriggerExpression)
	d.Set("window_size", CSEAggregationRuleGet.WindowSizeName)

	return nil
}

func resourceSumologicCSEAggregationRuleDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteCSEAggregationRule(d.Id())

}

func resourceSumologicCSEAggregationRuleCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		id, err := c.CreateCSEAggregationRule(CSEAggregationRule{
			AggregationFunctions:  resourceToAggregationFunctionsArray(d.Get("aggregation_functions").([]interface{})),
			DescriptionExpression: d.Get("description_expression").(string),
			Enabled:               d.Get("enabled").(bool),
			EntitySelectors:       resourceToEntitySelectorArray(d.Get("entity_selectors").([]interface{})),
			GroupByEntity:         d.Get("group_by_entity").(bool),
			GroupByFields:         resourceToStringArray(d.Get("group_by_fields").([]interface{})),
			IsPrototype:           d.Get("is_prototype").(bool),
			MatchExpression:       d.Get("match_expression").(string),
			Name:                  d.Get("name").(string),
			NameExpression:        d.Get("name_expression").(string),
			SeverityMapping:       resourceToSeverityMapping(d.Get("severity_mapping").([]interface{})[0]),
			Stream:                "record",
			SummaryExpression:     d.Get("summary_expression").(string),
			Tags:                  resourceToStringArray(d.Get("tags").([]interface{})),
			TriggerExpression:     d.Get("trigger_expression").(string),
			WindowSize:            windowSizeField(d.Get("window_size").(string)),
		})

		if err != nil {
			return err
		}
		d.SetId(id)
	}

	return resourceSumologicCSEAggregationRuleRead(d, meta)
}

func resourceSumologicCSEAggregationRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	CSEAggregationRule, err := resourceToCSEAggregationRule(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSEAggregationRule(CSEAggregationRule); err != nil {
		return err
	}

	return resourceSumologicCSEAggregationRuleRead(d, meta)
}

func resourceToAggregationFunctionsArray(resourceAggregationFunctions []interface{}) []AggregationFunction {
	result := make([]AggregationFunction, len(resourceAggregationFunctions))

	for i, resourceAggregationFunction := range resourceAggregationFunctions {
		aggregationFunctionMap := resourceAggregationFunction.(map[string]interface{})
		result[i] = AggregationFunction{
			Name:      aggregationFunctionMap["name"].(string),
			Function:  aggregationFunctionMap["function"].(string),
			Arguments: resourceToStringArray(aggregationFunctionMap["arguments"].([]interface{})),
		}
	}

	return result
}

func aggregationFunctionsArrayToResource(aggregationFunctions []AggregationFunction) []map[string]interface{} {
	result := make([]map[string]interface{}, len(aggregationFunctions))

	for i, aggregationFunction := range aggregationFunctions {
		result[i] = map[string]interface{}{
			"name":      aggregationFunction.Name,
			"function":  aggregationFunction.Function,
			"arguments": aggregationFunction.Arguments,
		}
	}

	return result
}

func resourceToCSEAggregationRule(d *schema.ResourceData) (CSEAggregationRule, error) {
	id := d.Id()
	if id == "" {
		return CSEAggregationRule{}, nil
	}

	return CSEAggregationRule{
		ID:                    id,
		AggregationFunctions:  resourceToAggregationFunctionsArray(d.Get("aggregation_functions").([]interface{})),
		DescriptionExpression: d.Get("description_expression").(string),
		Enabled:               d.Get("enabled").(bool),
		EntitySelectors:       resourceToEntitySelectorArray(d.Get("entity_selectors").([]interface{})),
		GroupByEntity:         d.Get("group_by_entity").(bool),
		GroupByFields:         resourceToStringArray(d.Get("group_by_fields").([]interface{})),
		IsPrototype:           d.Get("is_prototype").(bool),
		MatchExpression:       d.Get("match_expression").(string),
		Name:                  d.Get("name").(string),
		NameExpression:        d.Get("name_expression").(string),
		SeverityMapping:       resourceToSeverityMapping(d.Get("severity_mapping").([]interface{})[0]),
		Stream:                "record",
		SummaryExpression:     d.Get("summary_expression").(string),
		Tags:                  resourceToStringArray(d.Get("tags").([]interface{})),
		TriggerExpression:     d.Get("trigger_expression").(string),
		WindowSize:            windowSizeField(d.Get("window_size").(string)),
	}, nil
}
