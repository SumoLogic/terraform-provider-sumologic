package sumologic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSumologicMetricTransformationRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicMetricTransformationRuleCreate,
		Read:   resourceSumologicMetricTransformationRuleRead,
		Update: resourceSumologicMetricTransformationRuleUpdate,
		Delete: resourceSumologicMetricTransformationRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: resourceSumologicMetricTransformationRuleCustomizeDiff,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"selector": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"retention": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      400,
				ValidateFunc: validation.IntInSlice([]int{0, 8, 400}),
			},
			"transformed_metrics_retention": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntInSlice([]int{0, 8, 400}),
			},
			"dimension_transformation": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transformation_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								metricTransformationTypeAggregateOn,
								metricTransformationTypeAddOrReplace,
							}, false),
						},
						"aggregate_on": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"dimension_to_replace": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			// Computed, server-assigned attributes.
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// resourceSumologicMetricTransformationRuleCustomizeDiff catches invalid combinations at plan time.
func resourceSumologicMetricTransformationRuleCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	transformations := d.Get("dimension_transformation").([]interface{})

	for i, t := range transformations {
		tm, ok := t.(map[string]interface{})
		if !ok {
			continue
		}
		switch tm["transformation_type"].(string) {
		case metricTransformationTypeAggregateOn:
			if aggregateOn, _ := tm["aggregate_on"].([]interface{}); len(aggregateOn) == 0 {
				return fmt.Errorf("dimension_transformation.%d: aggregate_on is required when transformation_type is %s", i, metricTransformationTypeAggregateOn)
			}
		case metricTransformationTypeAddOrReplace:
			if tm["dimension_to_replace"].(string) == "" || tm["value"].(string) == "" {
				return fmt.Errorf("dimension_transformation.%d: dimension_to_replace and value are required when transformation_type is %s", i, metricTransformationTypeAddOrReplace)
			}
		}
	}

	// Per the API docs: with no dimension transformations the transformed metrics retention is
	// forced to 0, and a retention of 0 (do not store) is only valid with transformations.
	if len(transformations) == 0 {
		if d.Get("transformed_metrics_retention").(int) != 0 {
			return fmt.Errorf("transformed_metrics_retention must be 0 when no dimension_transformation is defined")
		}
		if d.Get("retention").(int) == 0 {
			return fmt.Errorf("retention can only be 0 when at least one dimension_transformation is defined")
		}
	}
	return nil
}

func resourceSumologicMetricTransformationRuleCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	rule := resourceToMetricTransformationRule(d)

	id, err := c.CreateMetricTransformationRule(rule)
	if err != nil {
		return err
	}
	d.SetId(id)

	return resourceSumologicMetricTransformationRuleRead(d, meta)
}

func resourceSumologicMetricTransformationRuleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	rule, err := c.GetMetricTransformationRule(d.Id())
	if err != nil {
		return err
	}
	if rule == nil {
		log.Printf("[WARN] Metric transformation rule not found, removing from state: %v", d.Id())
		d.SetId("")
		return nil
	}

	def := rule.RuleDefinition

	// d.Set keys must be string literals (tfproviderlint R001).
	if err := d.Set("name", def.Name); err != nil {
		return err
	}
	if err := d.Set("selector", def.Selector); err != nil {
		return err
	}
	if err := d.Set("enabled", rule.Enabled); err != nil {
		return err
	}
	if err := d.Set("retention", int(def.Retention)); err != nil {
		return err
	}
	if err := d.Set("transformed_metrics_retention", int(def.TransformedMetricsRetention)); err != nil {
		return err
	}
	if err := d.Set("dimension_transformation", flattenMetricDimensionTransformations(def.DimensionTransformations)); err != nil {
		return err
	}
	if err := d.Set("created_at", rule.CreatedAt); err != nil {
		return err
	}
	if err := d.Set("created_by", rule.CreatedBy); err != nil {
		return err
	}
	if err := d.Set("modified_at", rule.ModifiedAt); err != nil {
		return err
	}
	return d.Set("modified_by", rule.ModifiedBy)
}

func resourceSumologicMetricTransformationRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	// The API has no version field; PUT replaces the properties sent in the request.
	rule := resourceToMetricTransformationRule(d)

	if err := c.UpdateMetricTransformationRule(rule); err != nil {
		return err
	}

	return resourceSumologicMetricTransformationRuleRead(d, meta)
}

func resourceSumologicMetricTransformationRuleDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	return c.DeleteMetricTransformationRule(d.Id())
}

func resourceToMetricTransformationRule(d *schema.ResourceData) MetricTransformationRule {
	return MetricTransformationRule{
		ID:      d.Id(),
		Enabled: d.Get("enabled").(bool),
		RuleDefinition: MetricTransformationRuleDefinition{
			Name:                        d.Get("name").(string),
			Selector:                    d.Get("selector").(string),
			Retention:                   int64(d.Get("retention").(int)),
			TransformedMetricsRetention: int64(d.Get("transformed_metrics_retention").(int)),
			DimensionTransformations:    expandMetricDimensionTransformations(d.Get("dimension_transformation").([]interface{})),
		},
	}
}

func expandMetricDimensionTransformations(raw []interface{}) []MetricDimensionTransformation {
	result := make([]MetricDimensionTransformation, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		t := MetricDimensionTransformation{
			TransformationType: m["transformation_type"].(string),
		}
		switch t.TransformationType {
		case metricTransformationTypeAggregateOn:
			aggregateOn := []string{}
			for _, v := range m["aggregate_on"].([]interface{}) {
				if s, ok := v.(string); ok {
					aggregateOn = append(aggregateOn, s)
				}
			}
			t.AggregateOn = aggregateOn
		case metricTransformationTypeAddOrReplace:
			t.DimensionToReplace = m["dimension_to_replace"].(string)
			t.Value = m["value"].(string)
		}
		result = append(result, t)
	}
	return result
}

func flattenMetricDimensionTransformations(transformations []MetricDimensionTransformation) []interface{} {
	result := make([]interface{}, 0, len(transformations))
	for _, t := range transformations {
		aggregateOn := make([]interface{}, 0, len(t.AggregateOn))
		for _, v := range t.AggregateOn {
			aggregateOn = append(aggregateOn, v)
		}
		result = append(result, map[string]interface{}{
			"transformation_type":  t.TransformationType,
			"aggregate_on":         aggregateOn,
			"dimension_to_replace": t.DimensionToReplace,
			"value":                t.Value,
		})
	}
	return result
}
