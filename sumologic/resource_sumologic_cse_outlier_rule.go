package sumologic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSumologicCSEOutlierRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSEOutlierRuleCreate,
		Read:   resourceSumologicCSEOutlierRuleRead,
		Delete: resourceSumologicCSEOutlierRuleDelete,
		Update: resourceSumologicCSEOutlierRuleUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"aggregation_functions": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
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

			"baseline_window_size": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description_expression": {
				Type:     schema.TypeString,
				Required: true,
			},
			"deviation_threshold": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"entity_selectors": getEntitySelectorsSchema(),
			"floor_value": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"group_by_fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
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
			"retention_window_size": {
				Type:     schema.TypeString,
				Required: true,
			},
			"severity": {
				Type:     schema.TypeInt,
				Required: true,
			},
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
			"window_size": {
				Type:     schema.TypeString,
				Required: true,
			},
			"suppression_window_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 7*24*60*60*1000),
				ForceNew:     false,
			},
			"tuning_expression_ids": getTuningExpressionIDsSchema(),
		},
	}
}

func resourceSumologicCSEOutlierRuleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSEOutlierRuleGet *CSEOutlierRule
	id := d.Id()

	CSEOutlierRuleGet, err := c.GetCSEOutlierRule(id)
	if err != nil {
		log.Printf("[WARN] CSE Outlier Rule not found when looking by id: %s, err: %v", id, err)
	}

	if CSEOutlierRuleGet == nil {
		log.Printf("[WARN] CSE Outlier Rule not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("aggregation_functions", aggregationFunctionsArrayToResource(CSEOutlierRuleGet.AggregationFunctions))
	d.Set("baseline_window_size", CSEOutlierRuleGet.BaselineWindowSize)
	d.Set("description_expression", CSEOutlierRuleGet.DescriptionExpression)
	d.Set("deviation_threshold", CSEOutlierRuleGet.DeviationThreshold)
	d.Set("enabled", CSEOutlierRuleGet.Enabled)
	d.Set("entity_selectors", entitySelectorArrayToResource(CSEOutlierRuleGet.EntitySelectors))
	d.Set("floor_value", CSEOutlierRuleGet.FloorValue)
	d.Set("group_by_fields", CSEOutlierRuleGet.GroupByFields)
	d.Set("is_prototype", CSEOutlierRuleGet.IsPrototype)
	d.Set("match_expression", CSEOutlierRuleGet.MatchExpression)
	d.Set("name", CSEOutlierRuleGet.Name)
	d.Set("name_expression", CSEOutlierRuleGet.NameExpression)
	d.Set("retention_window_size", CSEOutlierRuleGet.RetentionWindowSize)
	d.Set("severity", CSEOutlierRuleGet.Severity)
	d.Set("summary_expression", CSEOutlierRuleGet.SummaryExpression)
	d.Set("tags", CSEOutlierRuleGet.Tags)
	d.Set("window_size", CSEOutlierRuleGet.WindowSizeName)
	if CSEOutlierRuleGet.SuppressionWindowSize != nil {
		d.Set("suppression_window_size", CSEOutlierRuleGet.SuppressionWindowSize)
	}
	d.Set("tuning_expression_ids", CSEOutlierRuleGet.TuningExpressionIDs)
	return nil
}

func resourceSumologicCSEOutlierRuleDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteCSEOutlierRule(d.Id())

}

func resourceSumologicCSEOutlierRuleCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		var suppressionWindowSize *int = nil
		if suppression, ok := d.GetOk("suppression_window_size"); ok {
			suppressionInt := suppression.(int)
			suppressionWindowSize = &suppressionInt
		}

		id, err := c.CreateCSEOutlierRule(CSEOutlierRule{
			AggregationFunctions:  resourceToAggregationFunctionsArray(d.Get("aggregation_functions").([]interface{})),
			BaselineWindowSize:    d.Get("baseline_window_size").(string),
			DescriptionExpression: d.Get("description_expression").(string),
			DeviationThreshold:    d.Get("deviation_threshold").(int),
			Enabled:               d.Get("enabled").(bool),
			EntitySelectors:       resourceToEntitySelectorArray(d.Get("entity_selectors").([]interface{})),
			FloorValue:            d.Get("floor_value").(int),
			GroupByFields:         resourceToStringArray(d.Get("group_by_fields").([]interface{})),
			IsPrototype:           d.Get("is_prototype").(bool),
			MatchExpression:       d.Get("match_expression").(string),
			Name:                  d.Get("name").(string),
			NameExpression:        d.Get("name_expression").(string),
			RetentionWindowSize:   d.Get("retention_window_size").(string),
			Severity:              d.Get("severity").(int),
			SummaryExpression:     d.Get("summary_expression").(string),
			Tags:                  resourceToStringArray(d.Get("tags").([]interface{})),
			WindowSize:            windowSizeField(d.Get("window_size").(string)),
			SuppressionWindowSize: suppressionWindowSize,
			TuningExpressionIDs:   resourceToStringArray(d.Get("tuning_expression_ids").([]interface{})),
		})

		if err != nil {
			return err
		}
		d.SetId(id)
	}

	return resourceSumologicCSEOutlierRuleRead(d, meta)
}

func resourceSumologicCSEOutlierRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	CSEOutlierRule, err := resourceToCSEOutlierRule(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSEOutlierRule(CSEOutlierRule); err != nil {
		return err
	}

	return resourceSumologicCSEOutlierRuleRead(d, meta)
}

func resourceToCSEOutlierRule(d *schema.ResourceData) (CSEOutlierRule, error) {
	id := d.Id()
	if id == "" {
		return CSEOutlierRule{}, nil
	}

	var suppressionWindowSize *int = nil
	if suppression, ok := d.GetOk("suppression_window_size"); ok {
		suppressionInt := suppression.(int)
		suppressionWindowSize = &suppressionInt
	}

	return CSEOutlierRule{
		ID:                    id,
		AggregationFunctions:  resourceToAggregationFunctionsArray(d.Get("aggregation_functions").([]interface{})),
		BaselineWindowSize:    d.Get("baseline_window_size").(string),
		DescriptionExpression: d.Get("description_expression").(string),
		DeviationThreshold:    d.Get("deviation_threshold").(int),
		Enabled:               d.Get("enabled").(bool),
		EntitySelectors:       resourceToEntitySelectorArray(d.Get("entity_selectors").([]interface{})),
		FloorValue:            d.Get("floor_value").(int),
		GroupByFields:         resourceToStringArray(d.Get("group_by_fields").([]interface{})),
		IsPrototype:           d.Get("is_prototype").(bool),
		MatchExpression:       d.Get("match_expression").(string),
		Name:                  d.Get("name").(string),
		NameExpression:        d.Get("name_expression").(string),
		RetentionWindowSize:   d.Get("retention_window_size").(string),
		Severity:              d.Get("severity").(int),
		SummaryExpression:     d.Get("summary_expression").(string),
		Tags:                  resourceToStringArray(d.Get("tags").([]interface{})),
		WindowSize:            windowSizeField(d.Get("window_size").(string)),
		SuppressionWindowSize: suppressionWindowSize,
		TuningExpressionIDs:   resourceToStringArray(d.Get("tuning_expression_ids").([]interface{})),
	}, nil
}
