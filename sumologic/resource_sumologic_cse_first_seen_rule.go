package sumologic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSumologicCSEFirstSeenRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSEFirstSeenRuleCreate,
		Read:   resourceSumologicCSEFirstSeenRuleRead,
		Delete: resourceSumologicCSEFirstSeenRuleDelete,
		Update: resourceSumologicCSEFirstSeenRuleUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"baseline_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"baseline_window_size": {
				Type:     schema.TypeString,
				Required: true,
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
			"filter_expression": {
				Type:     schema.TypeString,
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
			"value_fields": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
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

func resourceSumologicCSEFirstSeenRuleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSEFirstSeenRuleGet *CSEFirstSeenRule
	id := d.Id()

	CSEFirstSeenRuleGet, err := c.GetCSEFirstSeenRule(id)
	if err != nil {
		log.Printf("[WARN] CSE FirstSeen Rule not found when looking by id: %s, err: %v", id, err)
	}

	if CSEFirstSeenRuleGet == nil {
		log.Printf("[WARN] CSE FirstSeen Rule not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("baseline_type", CSEFirstSeenRuleGet.BaselineType)
	d.Set("baseline_window_size", CSEFirstSeenRuleGet.BaselineWindowSize)
	d.Set("description_expression", CSEFirstSeenRuleGet.DescriptionExpression)
	d.Set("enabled", CSEFirstSeenRuleGet.Enabled)
	d.Set("entity_selectors", entitySelectorArrayToResource(CSEFirstSeenRuleGet.EntitySelectors))
	d.Set("filter_expression", CSEFirstSeenRuleGet.FilterExpression)
	d.Set("group_by_fields", CSEFirstSeenRuleGet.GroupByFields)
	d.Set("is_prototype", CSEFirstSeenRuleGet.IsPrototype)
	d.Set("name", CSEFirstSeenRuleGet.Name)
	d.Set("name_expression", CSEFirstSeenRuleGet.NameExpression)
	d.Set("retention_window_size", CSEFirstSeenRuleGet.RetentionWindowSize)
	d.Set("severity", CSEFirstSeenRuleGet.Severity)
	d.Set("summary_expression", CSEFirstSeenRuleGet.SummaryExpression)
	d.Set("tags", CSEFirstSeenRuleGet.Tags)
	d.Set("value_fields", CSEFirstSeenRuleGet.ValueFields)
	if CSEFirstSeenRuleGet.SuppressionWindowSize != nil {
		d.Set("suppression_window_size", CSEFirstSeenRuleGet.SuppressionWindowSize)
	}
	d.Set("tuning_expression_ids", CSEFirstSeenRuleGet.TuningExpressionIDs)
	return nil
}

func resourceSumologicCSEFirstSeenRuleDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteCSEFirstSeenRule(d.Id())

}

func resourceSumologicCSEFirstSeenRuleCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		var suppressionWindowSize *int = nil
		if suppression, ok := d.GetOk("suppression_window_size"); ok {
			suppressionInt := suppression.(int)
			suppressionWindowSize = &suppressionInt
		}

		id, err := c.CreateCSEFirstSeenRule(CSEFirstSeenRule{
			BaselineType:          d.Get("baseline_type").(string),
			BaselineWindowSize:    d.Get("baseline_window_size").(string),
			DescriptionExpression: d.Get("description_expression").(string),
			Enabled:               d.Get("enabled").(bool),
			EntitySelectors:       resourceToEntitySelectorArray(d.Get("entity_selectors").([]interface{})),
			FilterExpression:      d.Get("filter_expression").(string),
			GroupByFields:         resourceToStringArray(d.Get("group_by_fields").([]interface{})),
			IsPrototype:           d.Get("is_prototype").(bool),
			Name:                  d.Get("name").(string),
			NameExpression:        d.Get("name_expression").(string),
			RetentionWindowSize:   d.Get("retention_window_size").(string),
			Severity:              d.Get("severity").(int),
			SummaryExpression:     d.Get("summary_expression").(string),
			Tags:                  resourceToStringArray(d.Get("tags").([]interface{})),
			ValueFields:           resourceToStringArray(d.Get("value_fields").([]interface{})),
			Version:               1,
			SuppressionWindowSize: suppressionWindowSize,
			TuningExpressionIDs:   resourceToStringArray(d.Get("tuning_expression_ids").([]interface{})),
		})

		if err != nil {
			return err
		}
		d.SetId(id)
	}

	return resourceSumologicCSEFirstSeenRuleRead(d, meta)
}

func resourceSumologicCSEFirstSeenRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	CSEFirstSeenRule, err := resourceToCSEFirstSeenRule(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSEFirstSeenRule(CSEFirstSeenRule); err != nil {
		return err
	}

	return resourceSumologicCSEFirstSeenRuleRead(d, meta)
}

func resourceToCSEFirstSeenRule(d *schema.ResourceData) (CSEFirstSeenRule, error) {
	id := d.Id()
	if id == "" {
		return CSEFirstSeenRule{}, nil
	}

	var suppressionWindowSize *int = nil
	if suppression, ok := d.GetOk("suppression_window_size"); ok {
		suppressionInt := suppression.(int)
		suppressionWindowSize = &suppressionInt
	}

	return CSEFirstSeenRule{
		ID:                    id,
		BaselineType:          d.Get("baseline_type").(string),
		BaselineWindowSize:    d.Get("baseline_window_size").(string),
		DescriptionExpression: d.Get("description_expression").(string),
		Enabled:               d.Get("enabled").(bool),
		EntitySelectors:       resourceToEntitySelectorArray(d.Get("entity_selectors").([]interface{})),
		FilterExpression:      d.Get("filter_expression").(string),
		GroupByFields:         resourceToStringArray(d.Get("group_by_fields").([]interface{})),
		IsPrototype:           d.Get("is_prototype").(bool),
		Name:                  d.Get("name").(string),
		NameExpression:        d.Get("name_expression").(string),
		RetentionWindowSize:   d.Get("retention_window_size").(string),
		Severity:              d.Get("severity").(int),
		SummaryExpression:     d.Get("summary_expression").(string),
		Tags:                  resourceToStringArray(d.Get("tags").([]interface{})),
		ValueFields:           resourceToStringArray(d.Get("value_fields").([]interface{})),
		Version:               1,
		SuppressionWindowSize: suppressionWindowSize,
		TuningExpressionIDs:   resourceToStringArray(d.Get("tuning_expression_ids").([]interface{})),
	}, nil
}
