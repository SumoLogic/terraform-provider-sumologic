package sumologic

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceSumologicCSEThresholdRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSEThresholdRuleCreate,
		Read:   resourceSumologicCSEThresholdRuleRead,
		Delete: resourceSumologicCSEThresholdRuleDelete,
		Update: resourceSumologicCSEThresholdRuleUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"count_distinct": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"count_field": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"entity_selectors": getEntitySelectorsSchema(),
			"expression": {
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
			"limit": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
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
			"window_size_millis": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"suppression_window_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 7*24*60*60*1000),
				ForceNew:     false,
			},
		},
	}
}

func resourceSumologicCSEThresholdRuleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSEThresholdRuleGet *CSEThresholdRule
	id := d.Id()

	CSEThresholdRuleGet, err := c.GetCSEThresholdRule(id)
	if err != nil {
		log.Printf("[WARN] CSE Threshold Rule not found when looking by id: %s, err: %v", id, err)
	}

	if CSEThresholdRuleGet == nil {
		log.Printf("[WARN] CSE Threshold Rule not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("count_distinct", CSEThresholdRuleGet.CountDistinct)
	d.Set("count_field", CSEThresholdRuleGet.CountField)
	d.Set("description", CSEThresholdRuleGet.Description)
	d.Set("enabled", CSEThresholdRuleGet.Enabled)
	d.Set("entity_selectors", entitySelectorArrayToResource(CSEThresholdRuleGet.EntitySelectors))
	d.Set("expression", CSEThresholdRuleGet.Expression)
	d.Set("group_by_fields", CSEThresholdRuleGet.GroupByFields)
	d.Set("is_prototype", CSEThresholdRuleGet.IsPrototype)
	d.Set("limit", CSEThresholdRuleGet.Limit)
	d.Set("name", CSEThresholdRuleGet.Name)
	d.Set("severity", CSEThresholdRuleGet.Severity)
	d.Set("summary_expression", CSEThresholdRuleGet.SummaryExpression)
	d.Set("tags", CSEThresholdRuleGet.Tags)
	d.Set("window_size", CSEThresholdRuleGet.WindowSizeName)
	if strings.EqualFold(CSEThresholdRuleGet.WindowSizeName, "CUSTOM") {
		d.Set("window_size_millis", CSEThresholdRuleGet.WindowSize)
	}
	if CSEThresholdRuleGet.SuppressionWindowSize != nil {
		d.Set("suppression_window_size", CSEThresholdRuleGet.SuppressionWindowSize)
	}
	return nil
}

func resourceSumologicCSEThresholdRuleDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteCSEThresholdRule(d.Id())

}

func resourceSumologicCSEThresholdRuleCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {

		var suppressionWindowSize *int = nil
		if suppression, ok := d.GetOk("suppression_window_size"); ok {
			suppressionInt := suppression.(int)
			suppressionWindowSize = &suppressionInt
		}

		id, err := c.CreateCSEThresholdRule(CSEThresholdRule{
			CountDistinct:          d.Get("count_distinct").(bool),
			CountField:             d.Get("count_field").(string),
			Description:            d.Get("description").(string),
			Enabled:                d.Get("enabled").(bool),
			EntitySelectors:        resourceToEntitySelectorArray(d.Get("entity_selectors").([]interface{})),
			Expression:             d.Get("expression").(string),
			GroupByFields:          resourceToStringArray(d.Get("group_by_fields").([]interface{})),
			IsPrototype:            d.Get("is_prototype").(bool),
			Limit:                  d.Get("limit").(int),
			Name:                   d.Get("name").(string),
			Severity:               d.Get("severity").(int),
			Stream:                 "record",
			SummaryExpression:      d.Get("summary_expression").(string),
			Tags:                   resourceToStringArray(d.Get("tags").([]interface{})),
			Version:                1,
			WindowSize:             windowSizeField(d.Get("window_size").(string)),
			WindowSizeMilliseconds: d.Get("window_size_millis").(string),
			SuppressionWindowSize:  suppressionWindowSize,
		})

		if err != nil {
			return err
		}
		d.SetId(id)
	}

	return resourceSumologicCSEThresholdRuleRead(d, meta)
}

func resourceSumologicCSEThresholdRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	CSEThresholdRule, err := resourceToCSEThresholdRule(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSEThresholdRule(CSEThresholdRule); err != nil {
		return err
	}

	return resourceSumologicCSEThresholdRuleRead(d, meta)
}

func resourceToCSEThresholdRule(d *schema.ResourceData) (CSEThresholdRule, error) {
	id := d.Id()
	if id == "" {
		return CSEThresholdRule{}, nil
	}

	var suppressionWindowSize *int = nil
	if suppression, ok := d.GetOk("suppression_window_size"); ok {
		suppressionInt := suppression.(int)
		suppressionWindowSize = &suppressionInt
	}

	return CSEThresholdRule{
		ID:                     id,
		CountDistinct:          d.Get("count_distinct").(bool),
		CountField:             d.Get("count_field").(string),
		Description:            d.Get("description").(string),
		Enabled:                d.Get("enabled").(bool),
		EntitySelectors:        resourceToEntitySelectorArray(d.Get("entity_selectors").([]interface{})),
		Expression:             d.Get("expression").(string),
		GroupByFields:          resourceToStringArray(d.Get("group_by_fields").([]interface{})),
		IsPrototype:            d.Get("is_prototype").(bool),
		Limit:                  d.Get("limit").(int),
		Name:                   d.Get("name").(string),
		Severity:               d.Get("severity").(int),
		Stream:                 "record",
		SummaryExpression:      d.Get("summary_expression").(string),
		Tags:                   resourceToStringArray(d.Get("tags").([]interface{})),
		Version:                1,
		WindowSize:             windowSizeField(d.Get("window_size").(string)),
		WindowSizeMilliseconds: d.Get("window_size_millis").(string),
		SuppressionWindowSize:  suppressionWindowSize,
	}, nil
}
