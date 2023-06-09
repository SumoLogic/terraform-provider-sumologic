package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
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
			"aggregate_function": {
				Type:     schema.TypeMap,
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
				Type:     schema.TypeString,
				Required: true,
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

	d.Set("aggregate_function", CSEOutlierRuleGet.AggregateFunction)
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
	d.Set("window_size", CSEOutlierRuleGet.WindowSize)

	return nil
}

func resourceSumologicCSEOutlierRuleDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteCSEOutlierRule(d.Id())

}

func resourceSumologicCSEOutlierRuleCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		id, err := c.CreateCSEOutlierRule(CSEOutlierRule{
			AggregateFunction:     d.Get("aggregate_function").(map[string]interface{}),
			BaselineWindowSize:    d.Get("baseline_window_size").(string),
			DescriptionExpression: d.Get("description_expression").(string),
			DeviationThreshold:    d.Get("deviation_threshold").(int),
			Enabled:               d.Get("enabled").(bool),
			EntitySelectors:       resourceToEntitySelectorArray(d.Get("entity_selectors").([]interface{})),
			FloorValue:            d.Get("floor_value").(string),
			GroupByFields:         resourceToStringArray(d.Get("group_by_fields").([]interface{})),
			IsPrototype:           d.Get("is_prototype").(bool),
			MatchExpression:       d.Get("match_expression").(string),
			Name:                  d.Get("name").(string),
			NameExpression:        d.Get("name_expression").(string),
			RetentionWindowSize:   d.Get("retention_window_size").(string),
			Severity:              d.Get("severity").(int),
			SummaryExpression:     d.Get("summary_expression").(string),
			Tags:                  resourceToStringArray(d.Get("tags").([]interface{})),
			WindowSize:            d.Get("window_size").(string),
			Version:               1,
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

	return CSEOutlierRule{
		ID:                    id,
		AggregateFunction:     d.Get("aggregate_function").(map[string]interface{}),
		BaselineWindowSize:    d.Get("baseline_window_size").(string),
		DescriptionExpression: d.Get("description_expression").(string),
		DeviationThreshold:    d.Get("deviation_threshold").(int),
		Enabled:               d.Get("enabled").(bool),
		EntitySelectors:       resourceToEntitySelectorArray(d.Get("entity_selectors").([]interface{})),
		FloorValue:            d.Get("floor_value").(string),
		GroupByFields:         resourceToStringArray(d.Get("group_by_fields").([]interface{})),
		IsPrototype:           d.Get("is_prototype").(bool),
		MatchExpression:       d.Get("match_expression").(string),
		Name:                  d.Get("name").(string),
		NameExpression:        d.Get("name_expression").(string),
		RetentionWindowSize:   d.Get("retention_window_size").(string),
		Severity:              d.Get("severity").(int),
		SummaryExpression:     d.Get("summary_expression").(string),
		Tags:                  resourceToStringArray(d.Get("tags").([]interface{})),
		WindowSize:            d.Get("window_size").(string),
		Version:               1,
	}, nil
}
