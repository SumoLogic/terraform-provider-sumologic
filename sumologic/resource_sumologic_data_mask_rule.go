package sumologic

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSumologicDataMaskRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicDataMaskRuleCreate,
		Read:   resourceSumologicDataMaskRuleRead,
		Update: resourceSumologicDataMaskRuleUpdate,
		Delete: resourceSumologicDataMaskRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"regex_pattern": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDataMaskRuleRegex,
			},
			"mask_string": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "##redactedPII##",
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 512),
				Default:      "",
			},
		},
	}
}

func resourceSumologicDataMaskRuleCreate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	if d.Id() == "" {
		rule, err := resourceToDataMaskRule(d)
		if err != nil {
			return err
		}

		createdRule, err := c.CreateDataMaskRule(rule)
		if err != nil {
			return err
		}

		d.SetId(createdRule.ID)
	}

	return resourceSumologicDataMaskRuleRead(d, meta)
}

func resourceSumologicDataMaskRuleRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	id := d.Id()

	rule, err := c.GetDataMaskRule(id)
	if err != nil {
		return err
	}

	if rule == nil {
		log.Printf("[WARN] Data mask rule not found, removing from state: %v", id)
		d.SetId("")
		return nil
	}

	if err = d.Set("name", rule.Name); err != nil {
		return fmt.Errorf("error setting name for data mask rule %s: %s", d.Id(), err)
	}
	if err = d.Set("regex_pattern", rule.RegexPattern); err != nil {
		return fmt.Errorf("error setting regex_pattern for data mask rule %s: %s", d.Id(), err)
	}
	if err = d.Set("mask_string", rule.MaskString); err != nil {
		return fmt.Errorf("error setting mask_string for data mask rule %s: %s", d.Id(), err)
	}
	if err = d.Set("enabled", rule.Enabled); err != nil {
		return fmt.Errorf("error setting enabled for data mask rule %s: %s", d.Id(), err)
	}
	if err = d.Set("description", rule.Description); err != nil {
		return fmt.Errorf("error setting description for data mask rule %s: %s", d.Id(), err)
	}

	return nil
}

func resourceSumologicDataMaskRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	rule, err := resourceToDataMaskRule(d)
	if err != nil {
		return err
	}
	rule.ID = d.Id()

	if _, err = c.UpdateDataMaskRule(rule); err != nil {
		return err
	}

	return resourceSumologicDataMaskRuleRead(d, meta)
}

func resourceSumologicDataMaskRuleDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)
	err := c.DeleteDataMaskRule(d.Id())
	if err != nil && !isDataMaskRuleNotFoundErr(err) {
		return err
	}
	return nil
}

func resourceToDataMaskRule(d *schema.ResourceData) (DataMaskRule, error) {
	return DataMaskRule{
		ID:           d.Id(),
		Name:         d.Get("name").(string),
		RegexPattern: d.Get("regex_pattern").(string),
		MaskString:   d.Get("mask_string").(string),
		Enabled:      d.Get("enabled").(bool),
		Description:  d.Get("description").(string),
	}, nil
}

func validateDataMaskRuleRegex(v interface{}, _ string) (warnings []string, errors []error) {
	pattern := v.(string)
	if _, err := regexp.Compile(pattern); err != nil {
		errors = append(errors, fmt.Errorf("invalid regex pattern: %s", err))
	}
	return warnings, errors
}

func isDataMaskRuleNotFoundErr(err error) bool {
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "not found") || strings.Contains(errMsg, "does not exist")
}


