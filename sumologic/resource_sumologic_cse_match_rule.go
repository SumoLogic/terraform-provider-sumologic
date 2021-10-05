package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceSumologicCSEMatchRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSEMatchRuleCreate,
		Read:   resourceSumologicCSEMatchRuleRead,
		Delete: resourceSumologicCSEMatchRuleDelete,
		Update: resourceSumologicCSEMatchRuleUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description_expression": {
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
			"severity_mapping": getSeverityMappingSchema(),
			"summary_expression": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSumologicCSEMatchRuleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSEMatchRuleGet *CSEMatchRule
	id := d.Id()

	CSEMatchRuleGet, err := c.GetCSEMatchRule(id)
	if err != nil {
		log.Printf("[WARN] CSE Insights Status not found when looking by id: %s, err: %v", id, err)

	}

	if CSEMatchRuleGet == nil {
		log.Printf("[WARN] CSE Insights Status not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("description_expression", CSEMatchRuleGet.DescriptionExpression)
	d.Set("enabled", CSEMatchRuleGet.Enabled)
	d.Set("entity_selectors", CSEMatchRuleGet.EntitySelectors)
	d.Set("expression", CSEMatchRuleGet.Expression)
	d.Set("is_prototype", CSEMatchRuleGet.IsPrototype)
	d.Set("name", CSEMatchRuleGet.Name)
	d.Set("name_expression", CSEMatchRuleGet.NameExpression)
	d.Set("severity_mapping", CSEMatchRuleGet.SeverityMapping)
	d.Set("summary_expression", CSEMatchRuleGet.SummaryExpression)
	d.Set("tags", CSEMatchRuleGet.Tags)

	return nil
}

func resourceSumologicCSEMatchRuleDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteCSEMatchRule(d.Id())

}

func resourceSumologicCSEMatchRuleCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		id, err := c.CreateCSEMatchRule(CSEMatchRule{
			DescriptionExpression: d.Get("description_expression").(string),
			Enabled:               d.Get("enabled").(bool),
			EntitySelectors:       resourceToEntitySelectorArray(d.Get("entity_selectors").([]interface{})),
			Expression:            d.Get("expression").(string),
			IsPrototype:           d.Get("is_prototype").(bool),
			Name:                  d.Get("name").(string),
			NameExpression:        d.Get("name_expression").(string),
			SeverityMapping:       resourceToSeverityMapping(d.Get("severity_mapping").([]interface{})[0]),
			Stream:                "record",
			SummaryExpression:     d.Get("summary_expression").(string),
			Tags:                  resourceToStringArray(d.Get("tags").([]interface{})),
		})

		if err != nil {
			return err
		}
		log.Printf("[INFO] got id: %s", id)
		d.SetId(id)
	}

	return resourceSumologicCSEMatchRuleRead(d, meta)
}

func resourceSumologicCSEMatchRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	CSEMatchRule, err := resourceToCSEMatchRule(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSEMatchRule(CSEMatchRule); err != nil {
		return err
	}

	return resourceSumologicCSEMatchRuleRead(d, meta)
}

func resourceToCSEMatchRule(d *schema.ResourceData) (CSEMatchRule, error) {
	id := d.Id()
	if id == "" {
		return CSEMatchRule{}, nil
	}

	return CSEMatchRule{
		ID:                    id,
		DescriptionExpression: d.Get("description_expression").(string),
		Enabled:               d.Get("enabled").(bool),
		EntitySelectors:       resourceToEntitySelectorArray(d.Get("entity_selectors").([]interface{})),
		Expression:            d.Get("expression").(string),
		IsPrototype:           d.Get("is_prototype").(bool),
		Name:                  d.Get("name").(string),
		NameExpression:        d.Get("name_expression").(string),
		SeverityMapping:       resourceToSeverityMapping(d.Get("severity_mapping").([]interface{})[0]),
		Stream:                "record",
		SummaryExpression:     d.Get("summary_expression").(string),
		Tags:                  resourceToStringArray(d.Get("tags").([]interface{})),
	}, nil
}
