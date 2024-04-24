package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
)

func resourceSumologicCSECustomInsight() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSECustomInsightCreate,
		Read:   resourceSumologicCSECustomInsightRead,
		Delete: resourceSumologicCSECustomInsightDelete,
		Update: resourceSumologicCSECustomInsightUpdate,
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ordered": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"rule_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"severity": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"HIGH", "MEDIUM", "LOW", "CRITICAL"}, false),
			},
			"dynamic_severity": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"minimum_signal_severity": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"insight_severity": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"HIGH", "MEDIUM", "LOW", "CRITICAL"}, false),
						},
					},
				},
			},
			"signal_names": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSumologicCSECustomInsightRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSECustomInsightGet *CSECustomInsight
	id := d.Id()

	CSECustomInsightGet, err := c.GetCSECustomInsight(id)
	if err != nil {
		log.Printf("[WARN] CSE Custom Insight not found when looking by id: %s, err: %v", id, err)
	}

	if CSECustomInsightGet == nil {
		log.Printf("[WARN] CSE Custom Insight not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("description", CSECustomInsightGet.Description)
	d.Set("enabled", CSECustomInsightGet.Enabled)
	d.Set("name", CSECustomInsightGet.Name)
	d.Set("ordered", CSECustomInsightGet.Ordered)
	d.Set("rule_ids", CSECustomInsightGet.RuleIds)
	d.Set("severity", CSECustomInsightGet.Severity)
	d.Set("dynamic_severity", dynamicSeverityArrayToResource(CSECustomInsightGet.DynamicSeverity))
	d.Set("signal_names", CSECustomInsightGet.SignalNames)
	d.Set("tags", CSECustomInsightGet.Tags)

	return nil
}

func resourceSumologicCSECustomInsightDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteCSECustomInsight(d.Id())
}

func resourceSumologicCSECustomInsightCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		id, err := c.CreateCSECustomInsight(CSECustomInsight{
			Description:     d.Get("description").(string),
			Enabled:         d.Get("enabled").(bool),
			RuleIds:         resourceToStringArray(d.Get("rule_ids").([]interface{})),
			Name:            d.Get("name").(string),
			Ordered:         d.Get("ordered").(bool),
			Severity:        d.Get("severity").(string),
			DynamicSeverity: resourceToDynamicSeverityArray(d.Get("dynamic_severity").([]interface{})),
			SignalNames:     resourceToStringArray(d.Get("signal_names").([]interface{})),
			Tags:            resourceToStringArray(d.Get("tags").([]interface{})),
		})

		if err != nil {
			return err
		}
		d.SetId(id)
	}

	return resourceSumologicCSECustomInsightRead(d, meta)
}

func resourceSumologicCSECustomInsightUpdate(d *schema.ResourceData, meta interface{}) error {
	CSECustomInsight, err := resourceToCSECustomInsight(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSECustomInsight(CSECustomInsight); err != nil {
		return err
	}

	return resourceSumologicCSECustomInsightRead(d, meta)
}

func resourceToCSECustomInsight(d *schema.ResourceData) (CSECustomInsight, error) {
	id := d.Id()
	if id == "" {
		return CSECustomInsight{}, nil
	}

	return CSECustomInsight{
		ID:              id,
		Description:     d.Get("description").(string),
		Enabled:         d.Get("enabled").(bool),
		RuleIds:         resourceToStringArray(d.Get("rule_ids").([]interface{})),
		Name:            d.Get("name").(string),
		Ordered:         d.Get("ordered").(bool),
		Severity:        d.Get("severity").(string),
		DynamicSeverity: resourceToDynamicSeverityArray(d.Get("dynamic_severity").([]interface{})),
		SignalNames:     resourceToStringArray(d.Get("signal_names").([]interface{})),
		Tags:            resourceToStringArray(d.Get("tags").([]interface{})),
	}, nil
}

func resourceToDynamicSeverityArray(resourceDynamicSeverity []interface{}) []DynamicSeverity {
	result := make([]DynamicSeverity, len(resourceDynamicSeverity))

	for i, resourceDynamicSeverity := range resourceDynamicSeverity {
		result[i] = DynamicSeverity{
			MinimumSignalSeverity: resourceDynamicSeverity.(map[string]interface{})["minimum_signal_severity"].(int),
			InsightSeverity:       resourceDynamicSeverity.(map[string]interface{})["insight_severity"].(string),
		}
	}

	return result
}

func dynamicSeverityArrayToResource(dynamicSeverities []DynamicSeverity) []map[string]interface{} {
	result := make([]map[string]interface{}, len(dynamicSeverities))

	for i, dynamicSeverity := range dynamicSeverities {
		result[i] = map[string]interface{}{
			"minimum_signal_severity": dynamicSeverity.MinimumSignalSeverity,
			"insight_severity":        dynamicSeverity.InsightSeverity,
		}
	}

	return result
}
