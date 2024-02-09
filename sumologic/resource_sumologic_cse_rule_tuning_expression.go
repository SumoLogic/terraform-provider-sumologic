package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceSumologicCSERuleTuningExpression() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSERuleTuningExpressionCreate,
		Read:   resourceSumologicCSERuleTuningExpressionRead,
		Delete: resourceSumologicCSERuleTuningExpressionDelete,
		Update: resourceSumologicCSERuleTuningExpressionUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"expression": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"exclude": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"is_global": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"rule_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSumologicCSERuleTuningExpressionRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var CSERuleTuningExpressionGet *CSERuleTuningExpression
	id := d.Id()

	CSERuleTuningExpressionGet, err := c.GetCSERuleTuningExpression(id)
	if err != nil {
		log.Printf("[WARN] CSE Insights Status not found when looking by id: %s, err: %v", id, err)

	}

	if CSERuleTuningExpressionGet == nil {
		log.Printf("[WARN] CSE Insights Status not found, removing from state: %v - %v", id, err)
		d.SetId("")
		return nil
	}

	d.Set("name", CSERuleTuningExpressionGet.Name)
	d.Set("description", CSERuleTuningExpressionGet.Description)
	d.Set("expression", CSERuleTuningExpressionGet.Expression)
	d.Set("enabled", CSERuleTuningExpressionGet.Enabled)
	d.Set("exclude", CSERuleTuningExpressionGet.Exclude)
	d.Set("is_global", CSERuleTuningExpressionGet.IsGlobal)
	d.Set("rule_ids", CSERuleTuningExpressionGet.RuleIds)

	return nil
}

func resourceSumologicCSERuleTuningExpressionDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteCSERuleTuningExpression(d.Id())

}

func resourceSumologicCSERuleTuningExpressionCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		id, err := c.CreateCSERuleTuningExpression(CSERuleTuningExpression{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Expression:  d.Get("expression").(string),
			Enabled:     d.Get("enabled").(bool),
			Exclude:     d.Get("exclude").(bool),
			IsGlobal:    d.Get("is_global").(bool),
			RuleIds:     resourceRuleIdsToStringArray(d.Get("rule_ids").(*schema.Set)),
		})

		if err != nil {
			return err
		}
		log.Printf("[INFO] got id: %s", id)
		d.SetId(id)
	}

	return resourceSumologicCSERuleTuningExpressionRead(d, meta)
}

func resourceSumologicCSERuleTuningExpressionUpdate(d *schema.ResourceData, meta interface{}) error {
	CSERuleTuningExpression, err := resourceToCSERuleTuningExpression(d)
	if err != nil {
		return err
	}

	c := meta.(*Client)
	if err = c.UpdateCSERuleTuningExpression(CSERuleTuningExpression); err != nil {
		return err
	}

	return resourceSumologicCSERuleTuningExpressionRead(d, meta)
}

func resourceRuleIdsToStringArray(resourceRuleIds *schema.Set) []string {
	rawRuleIds := resourceRuleIds.List()
	ruleIds := make([]string, len(rawRuleIds))
	for i, v := range rawRuleIds {
		ruleIds[i] = v.(string)
	}

	return ruleIds
}

func resourceToCSERuleTuningExpression(d *schema.ResourceData) (CSERuleTuningExpression, error) {
	id := d.Id()
	if id == "" {
		return CSERuleTuningExpression{}, nil
	}

	return CSERuleTuningExpression{
		ID:          id,
		Description: d.Get("description").(string),
		Name:        d.Get("name").(string),
		Expression:  d.Get("expression").(string),
		Enabled:     d.Get("enabled").(bool),
		Exclude:     d.Get("exclude").(bool),
		IsGlobal:    d.Get("is_global").(bool),
		RuleIds:     resourceRuleIdsToStringArray(d.Get("rule_ids").(*schema.Set)),
	}, nil
}
