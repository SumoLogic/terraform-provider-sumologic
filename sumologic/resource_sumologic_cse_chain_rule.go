package sumologic

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSumologicCSEChainRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSEChainRuleCreate,
		Read:   resourceSumologicCSEChainRuleRead,
		Delete: resourceSumologicCSEChainRuleDelete,
		Update: resourceSumologicCSEChainRuleUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"entity_selectors": getEntitySelectorsSchema(),
			"expressions_and_limits": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expression": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppressSpaceDiff,
						},
						"limit": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
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
			"ordered": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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

func resourceSumologicCSEChainRuleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSEChainRuleGet *CSEChainRule
	id := d.Id()

	CSEChainRuleGet, err := c.GetCSEChainRule(id)
	if err != nil {
		log.Printf("[WARN] CSE Chain Rule not found when looking by id: %s, err: %v", id, err)
	}

	if CSEChainRuleGet == nil {
		log.Printf("[WARN] CSE Chain Rule not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("description", CSEChainRuleGet.Description)
	d.Set("enabled", CSEChainRuleGet.Enabled)
	d.Set("entity_selectors", entitySelectorArrayToResource(CSEChainRuleGet.EntitySelectors))
	d.Set("expressions_and_limits", expressionsAndLimitsArrayToResource(CSEChainRuleGet.ExpressionsAndLimits))
	d.Set("group_by_fields", CSEChainRuleGet.GroupByFields)
	d.Set("is_prototype", CSEChainRuleGet.IsPrototype)
	d.Set("ordered", CSEChainRuleGet.Ordered)
	d.Set("name", CSEChainRuleGet.Name)
	d.Set("severity", CSEChainRuleGet.Severity)
	d.Set("summary_expression", CSEChainRuleGet.SummaryExpression)
	d.Set("tags", CSEChainRuleGet.Tags)
	d.Set("window_size", CSEChainRuleGet.WindowSizeName)
	if strings.EqualFold(CSEChainRuleGet.WindowSizeName, "CUSTOM") {
		d.Set("window_size_millis", CSEChainRuleGet.WindowSize)
	}
	if CSEChainRuleGet.SuppressionWindowSize != nil {
		d.Set("suppression_window_size", CSEChainRuleGet.SuppressionWindowSize)
	}

	return nil
}

func resourceSumologicCSEChainRuleDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteCSEChainRule(d.Id())

}

func resourceSumologicCSEChainRuleCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {

		var suppressionWindowSize *int = nil
		if suppression, ok := d.GetOk("suppression_window_size"); ok {
			suppressionInt := suppression.(int)
			suppressionWindowSize = &suppressionInt
		}

		id, err := c.CreateCSEChainRule(CSEChainRule{
			Description:            d.Get("description").(string),
			Enabled:                d.Get("enabled").(bool),
			EntitySelectors:        resourceToEntitySelectorArray(d.Get("entity_selectors").([]interface{})),
			ExpressionsAndLimits:   resourceToExpressionsAndLimitsArray(d.Get("expressions_and_limits").([]interface{})),
			GroupByFields:          resourceToStringArray(d.Get("group_by_fields").([]interface{})),
			IsPrototype:            d.Get("is_prototype").(bool),
			Ordered:                d.Get("ordered").(bool),
			Name:                   d.Get("name").(string),
			Severity:               d.Get("severity").(int),
			Stream:                 "record",
			SummaryExpression:      d.Get("summary_expression").(string),
			Tags:                   resourceToStringArray(d.Get("tags").([]interface{})),
			WindowSize:             windowSizeField(d.Get("window_size").(string)),
			WindowSizeMilliseconds: d.Get("window_size_millis").(string),
			SuppressionWindowSize:  suppressionWindowSize,
		})

		if err != nil {
			return err
		}
		d.SetId(id)
	}

	return resourceSumologicCSEChainRuleRead(d, meta)
}

func resourceSumologicCSEChainRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	CSEChainRule, err := resourceToCSEChainRule(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSEChainRule(CSEChainRule); err != nil {
		return err
	}

	return resourceSumologicCSEChainRuleRead(d, meta)
}

func resourceToExpressionsAndLimitsArray(resourceExpressionsAndLimits []interface{}) []ExpressionAndLimit {
	result := make([]ExpressionAndLimit, len(resourceExpressionsAndLimits))

	for i, resourceExpressionAndLimit := range resourceExpressionsAndLimits {
		result[i] = ExpressionAndLimit{
			Expression: resourceExpressionAndLimit.(map[string]interface{})["expression"].(string),
			Limit:      resourceExpressionAndLimit.(map[string]interface{})["limit"].(int),
		}
	}

	return result
}

func expressionsAndLimitsArrayToResource(expressionsAndLimits []ExpressionAndLimit) []map[string]interface{} {
	result := make([]map[string]interface{}, len(expressionsAndLimits))

	for i, expressionAndLimit := range expressionsAndLimits {
		result[i] = map[string]interface{}{
			"expression": expressionAndLimit.Expression,
			"limit":      expressionAndLimit.Limit,
		}
	}

	return result
}

func resourceToCSEChainRule(d *schema.ResourceData) (CSEChainRule, error) {
	id := d.Id()
	if id == "" {
		return CSEChainRule{}, nil
	}

	var suppressionWindowSize *int = nil
	if suppression, ok := d.GetOk("suppression_window_size"); ok {
		suppressionInt := suppression.(int)
		suppressionWindowSize = &suppressionInt
	}

	return CSEChainRule{
		ID:                     id,
		Description:            d.Get("description").(string),
		Enabled:                d.Get("enabled").(bool),
		EntitySelectors:        resourceToEntitySelectorArray(d.Get("entity_selectors").([]interface{})),
		ExpressionsAndLimits:   resourceToExpressionsAndLimitsArray(d.Get("expressions_and_limits").([]interface{})),
		GroupByFields:          resourceToStringArray(d.Get("group_by_fields").([]interface{})),
		IsPrototype:            d.Get("is_prototype").(bool),
		Ordered:                d.Get("ordered").(bool),
		Name:                   d.Get("name").(string),
		Severity:               d.Get("severity").(int),
		Stream:                 "record",
		SummaryExpression:      d.Get("summary_expression").(string),
		Tags:                   resourceToStringArray(d.Get("tags").([]interface{})),
		WindowSize:             windowSizeField(d.Get("window_size").(string)),
		WindowSizeMilliseconds: d.Get("window_size_millis").(string),
		SuppressionWindowSize:  suppressionWindowSize,
	}, nil
}
