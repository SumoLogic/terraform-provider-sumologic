package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				ValidateFunc: validation.StringInSlice([]string{"HIGH", "MEDIUM", "LOW"}, false),
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
			Description: d.Get("description").(string),
			Enabled:     d.Get("enabled").(bool),
			RuleIds:     resourceToStringArray(d.Get("rule_ids").([]interface{})),
			Name:        d.Get("name").(string),
			Ordered:     d.Get("ordered").(bool),
			Severity:    d.Get("severity").(string),
			SignalNames: resourceToStringArray(d.Get("signal_names").([]interface{})),
			Tags:        resourceToStringArray(d.Get("tags").([]interface{})),
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
		ID:          id,
		Description: d.Get("description").(string),
		Enabled:     d.Get("enabled").(bool),
		RuleIds:     resourceToStringArray(d.Get("rule_ids").([]interface{})),
		Name:        d.Get("name").(string),
		Ordered:     d.Get("ordered").(bool),
		Severity:    d.Get("severity").(string),
		SignalNames: resourceToStringArray(d.Get("signal_names").([]interface{})),
		Tags:        resourceToStringArray(d.Get("tags").([]interface{})),
	}, nil
}
