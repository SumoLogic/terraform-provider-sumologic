package sumologic

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var supportedDataMaskRuleScopes = []string{"org", "child_org", "all_orgs", "all_child_orgs"}
var supportedDataMaskRulePIITypes = []string{"phone", "email", "ip", "ssn", "credit_card", "custom"}

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
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"pattern": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDataMaskRuleRegex,
			},
			"pii_type": {
				Type:         schema.TypeString,
				Required:     true,
				StateFunc:    normalizeLowerState,
				ValidateFunc: validation.StringInSlice(supportedDataMaskRulePIITypes, false),
			},
			"replacement": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 512),
			},
			"scope": {
				Type:         schema.TypeString,
				Required:     true,
				StateFunc:    normalizeDataMaskRuleScopeState,
				ValidateFunc: validation.StringInSlice(supportedDataMaskRuleScopes, false),
			},
			"scope_target_org_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1024),
				Default:      "",
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
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
	if err = d.Set("pattern", rule.Pattern); err != nil {
		return fmt.Errorf("error setting pattern for data mask rule %s: %s", d.Id(), err)
	}
	if err = d.Set("pii_type", strings.ToLower(rule.PiiType)); err != nil {
		return fmt.Errorf("error setting pii_type for data mask rule %s: %s", d.Id(), err)
	}
	if err = d.Set("replacement", rule.Replacement); err != nil {
		return fmt.Errorf("error setting replacement for data mask rule %s: %s", d.Id(), err)
	}
	if err = d.Set("scope", normalizeDataMaskRuleScope(rule.Scope)); err != nil {
		return fmt.Errorf("error setting scope for data mask rule %s: %s", d.Id(), err)
	}
	if err = d.Set("scope_target_org_ids", rule.ScopeTargetOrgIds); err != nil {
		return fmt.Errorf("error setting scope_target_org_ids for data mask rule %s: %s", d.Id(), err)
	}
	if err = d.Set("enabled", rule.Enabled); err != nil {
		return fmt.Errorf("error setting enabled for data mask rule %s: %s", d.Id(), err)
	}
	if err = d.Set("description", rule.Description); err != nil {
		return fmt.Errorf("error setting description for data mask rule %s: %s", d.Id(), err)
	}
	if err = d.Set("is_active", rule.IsActive); err != nil {
		return fmt.Errorf("error setting is_active for data mask rule %s: %s", d.Id(), err)
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
	scope := normalizeDataMaskRuleScopeState(d.Get("scope").(string))
	scopeTargetOrgIDs := expandStringList(d.Get("scope_target_org_ids").([]interface{}))
	if scope == "child_org" && len(scopeTargetOrgIDs) == 0 {
		return DataMaskRule{}, fmt.Errorf("scope_target_org_ids must be set when scope is child_org")
	}

	return DataMaskRule{
		ID:                d.Id(),
		Name:              d.Get("name").(string),
		Pattern:           d.Get("pattern").(string),
		PiiType:           strings.ToLower(d.Get("pii_type").(string)),
		Replacement:       d.Get("replacement").(string),
		Scope:             normalizeDataMaskRuleScopeForAPI(scope),
		ScopeTargetOrgIds: scopeTargetOrgIDs,
		Enabled:           d.Get("enabled").(bool),
		Description:       d.Get("description").(string),
	}, nil
}

func normalizeLowerState(v interface{}) string {
	return strings.ToLower(v.(string))
}

func normalizeDataMaskRuleScopeState(v interface{}) string {
	return normalizeDataMaskRuleScope(v.(string))
}

func normalizeDataMaskRuleScope(scope string) string {
	scope = strings.ToLower(scope)
	if scope == "all_child_orgs" {
		return "all_orgs"
	}
	return scope
}

func normalizeDataMaskRuleScopeForAPI(scope string) string {
	if scope == "all_orgs" {
		return "all_child_orgs"
	}
	return scope
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

func expandStringList(items []interface{}) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, item.(string))
	}
	return result
}

