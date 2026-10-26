package sumologic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSumologicCSESigmaRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicCSESigmaRuleCreate,
		Read:   resourceSumologicCSESigmaRuleRead,
		Delete: resourceSumologicCSESigmaRuleDelete,
		Update: resourceSumologicCSESigmaRuleUpdate,

		Schema: map[string]*schema.Schema{
			"sigma_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name_expression": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description_expression": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"summary_expression": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"expression": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppressSpaceDiff,
			},
			"is_prototype": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"entity_selectors": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expression": {
							Type:     schema.TypeString,
							Required: true,
						},
						"entity_type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"severity_mapping": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"default": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"field": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"mapping": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"from": {
										Type:     schema.TypeString,
										Required: true,
									},
									"to": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
			"suppression_window_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 7*24*60*60*1000),
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"original_source": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSumologicCSESigmaRuleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	id := d.Id()
	rule, err := c.GetCSEMatchRule(id)
	if err != nil {
		return err
	}

	if rule == nil {
		log.Printf("[WARN] CSE Sigma Rule not found, removing from state: %v", id)
		d.SetId("")
		return nil
	}

	d.Set("enabled", rule.Enabled)
	d.Set("name", rule.Name)
	d.Set("name_expression", rule.NameExpression)
	d.Set("description_expression", rule.DescriptionExpression)
	d.Set("summary_expression", rule.SummaryExpression)
	d.Set("expression", rule.Expression)
	d.Set("is_prototype", rule.IsPrototype)
	d.Set("entity_selectors", entitySelectorArrayToResource(rule.EntitySelectors))
	d.Set("severity_mapping", severityMappingToResource(rule.SeverityMapping))
	d.Set("tags", rule.Tags)
	d.Set("category", rule.Category)
	d.Set("original_source", rule.OriginalSource)
	if rule.SuppressionWindowSize != nil {
		d.Set("suppression_window_size", rule.SuppressionWindowSize)
	}

	return nil
}

func resourceSumologicCSESigmaRuleDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	return c.DeleteCSEMatchRule(d.Id())
}

func resourceSumologicCSESigmaRuleCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	sigmaRule, err := c.GetCSESigmaRule(d.Get("sigma_id").(string))
	if err != nil {
		return err
	}

	rule := sigmaRule.toMatchRule()
	rule.Enabled = d.Get("enabled").(bool)

	// Values set in the configuration take precedence over the converted Sigma rule.
	if v, ok := d.GetOk("name"); ok {
		rule.Name = v.(string)
	}
	if v, ok := d.GetOk("name_expression"); ok {
		rule.NameExpression = v.(string)
	}
	if v, ok := d.GetOk("description_expression"); ok {
		rule.DescriptionExpression = v.(string)
	}
	if v, ok := d.GetOk("summary_expression"); ok {
		rule.SummaryExpression = v.(string)
	}
	if v, ok := d.GetOk("expression"); ok {
		rule.Expression = v.(string)
	}
	if !d.GetRawConfig().GetAttr("is_prototype").IsNull() {
		rule.IsPrototype = d.Get("is_prototype").(bool)
	}
	if v, ok := d.GetOk("entity_selectors"); ok {
		rule.EntitySelectors = resourceToEntitySelectorArray(v.([]interface{}))
	}
	if v, ok := d.GetOk("severity_mapping"); ok {
		rule.SeverityMapping = resourceToSeverityMapping(v.([]interface{})[0])
	}
	if v, ok := d.GetOk("tags"); ok {
		rule.Tags = resourceToStringArray(v.([]interface{}))
	}
	if v, ok := d.GetOk("suppression_window_size"); ok {
		size := v.(int)
		rule.SuppressionWindowSize = &size
	}

	id, err := c.CreateCSEMatchRule(rule)
	if err != nil {
		return err
	}
	d.SetId(id)

	return resourceSumologicCSESigmaRuleRead(d, meta)
}

func resourceSumologicCSESigmaRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	var suppressionWindowSize *int
	if v, ok := d.GetOk("suppression_window_size"); ok {
		size := v.(int)
		suppressionWindowSize = &size
	}

	// The original source is overwritten on every update, so it must always be sent.
	err := c.UpdateCSEMatchRule(CSEMatchRule{
		ID:                    d.Id(),
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
		SuppressionWindowSize: suppressionWindowSize,
		OriginalSource:        cseSigmaOriginalSource,
	})
	if err != nil {
		return err
	}

	return resourceSumologicCSESigmaRuleRead(d, meta)
}
