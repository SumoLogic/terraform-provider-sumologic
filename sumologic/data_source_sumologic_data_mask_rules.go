package sumologic

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSumologicDataMaskRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSumologicDataMaskRulesRead,

		Schema: map[string]*schema.Schema{
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: dataSourceDataMaskRuleComputedSchema(),
				},
			},
		},
	}
}

func dataSourceSumologicDataMaskRulesRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	rules, err := c.ListDataMaskRules()
	if err != nil {
		return fmt.Errorf("error retrieving data mask rules: %v", err)
	}

	terraformRules := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		terraformRules = append(terraformRules, map[string]interface{}{
			"id":                   rule.ID,
			"name":                 rule.Name,
			"pattern":              rule.Pattern,
			"pii_type":             rule.PiiType,
			"replacement":          rule.Replacement,
			"scope":                normalizeDataMaskRuleScope(rule.Scope),
			"scope_target_org_ids": rule.ScopeTargetOrgIds,
			"enabled":              rule.Enabled,
			"description":          rule.Description,
			"is_active":            rule.IsActive,
		})
	}

	d.Set("rules", terraformRules)
	d.SetId(generateDataMaskRulesID(rules))

	return nil
}

func generateDataMaskRulesID(rules []DataMaskRule) string {
	ids := []string{"data_mask_rule_ids"}
	for _, rule := range rules {
		ids = append(ids, rule.ID)
	}
	sort.Strings(ids)

	idString := strings.Join(ids, "|")
	hash := sha256.Sum256([]byte(idString))
	return hex.EncodeToString(hash[:])
}


