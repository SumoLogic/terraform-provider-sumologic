package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var DefaultPolicies = Policies{
	Audit:                              AuditPolicy{Enabled: false},
	DataAccessLevel:                    DataAccessLevelPolicy{Enabled: false},
	MaxUserSessionTimeout:              MaxUserSessionTimeoutPolicy{MaxUserSessionTimeout: "7d"},
	SearchAudit:                        SearchAuditPolicy{Enabled: false},
	ShareDashboardsOutsideOrganization: ShareDashboardsOutsideOrganizationPolicy{Enabled: false},
	UserConcurrentSessionsLimit:        UserConcurrentSessionsLimitPolicy{Enabled: false, MaxConcurrentSessions: 100},
}

func resourceSumologicPolicies() *schema.Resource {
	return &schema.Resource{
		Create: resourceSumologicPoliciesCreate,
		Read:   resourceSumologicPoliciesRead,
		Update: resourceSumologicPoliciesUpdate,
		Delete: resourceSumologicPoliciesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"audit": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  DefaultPolicies.Audit.Enabled,
			},
			"data_access_level": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  DefaultPolicies.DataAccessLevel.Enabled,
			},
			"max_user_session_timeout": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  DefaultPolicies.MaxUserSessionTimeout.MaxUserSessionTimeout,
			},
			"search_audit": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  DefaultPolicies.SearchAudit.Enabled,
			},
			"share_dashboards_outside_organization": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  DefaultPolicies.ShareDashboardsOutsideOrganization.Enabled,
			},
			"user_concurrent_sessions_limit": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  DefaultPolicies.UserConcurrentSessionsLimit.Enabled,
						},
						"max_concurrent_sessions": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  DefaultPolicies.UserConcurrentSessionsLimit.MaxConcurrentSessions,
						},
					},
				},
			},
		},
	}
}

func resourceSumologicPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*Client)

	policies, err := c.GetPolicies()
	if err != nil {
		return err
	}

	setPoliciesResource(d, policies)
	return nil
}

func resourceSumologicPoliciesCreate(d *schema.ResourceData, meta interface{}) error {
	// Since policies can only be set and not created, we just update the policies with the given fields.
	err := resourceSumologicPoliciesUpdate(d, meta)
	if err != nil {
		return err
	}

	d.SetId("org-policies")
	return nil
}

func resourceSumologicPoliciesDelete(d *schema.ResourceData, meta interface{}) error {
	// Since policies cannot be deleted, we just reset to the default policies configuration.
	c := meta.(*Client)

	_, err := c.UpdatePolicies(DefaultPolicies)
	return err
}

func resourceSumologicPoliciesUpdate(d *schema.ResourceData, meta interface{}) error {
	policies := resourceToPolicies(d)

	c := meta.(*Client)
	updatedPolicies, err := c.UpdatePolicies(policies)
	if err != nil {
		return err
	}

	setPoliciesResource(d, updatedPolicies)
	return nil
}

func setPoliciesResource(d *schema.ResourceData, policies *Policies) {
	d.Set("audit", policies.Audit.Enabled)
	d.Set("data_access_level", policies.DataAccessLevel.Enabled)
	d.Set("max_user_session_timeout", policies.MaxUserSessionTimeout.MaxUserSessionTimeout)
	d.Set("search_audit", policies.SearchAudit.Enabled)
	d.Set("share_dashboards_outside_organization", policies.ShareDashboardsOutsideOrganization.Enabled)
	setUserConcurrentSessionsLimitPolicy(d, &policies.UserConcurrentSessionsLimit)
}

func resourceToPolicies(d *schema.ResourceData) Policies {
	var policies Policies
	policies.Audit = AuditPolicy{d.Get("audit").(bool)}
	policies.DataAccessLevel = DataAccessLevelPolicy{d.Get("data_access_level").(bool)}
	policies.MaxUserSessionTimeout = MaxUserSessionTimeoutPolicy{d.Get("max_user_session_timeout").(string)}
	policies.SearchAudit = SearchAuditPolicy{d.Get("search_audit").(bool)}
	policies.ShareDashboardsOutsideOrganization = ShareDashboardsOutsideOrganizationPolicy{d.Get("share_dashboards_outside_organization").(bool)}
	policies.UserConcurrentSessionsLimit = getUserConcurrentSessionsLimitPolicy(d)
	return policies
}

func setUserConcurrentSessionsLimitPolicy(d *schema.ResourceData, policy *UserConcurrentSessionsLimitPolicy) {
	userConcurrentSessionsLimitPolicyMap := make(map[string]interface{})
	userConcurrentSessionsLimitPolicyMap["enabled"] = policy.Enabled
	userConcurrentSessionsLimitPolicyMap["max_concurrent_sessions"] = policy.MaxConcurrentSessions

	userConcurrentSessionsLimitPolicy := make([]map[string]interface{}, 1)
	userConcurrentSessionsLimitPolicy[0] = userConcurrentSessionsLimitPolicyMap

	d.Set("user_concurrent_sessions_limit", userConcurrentSessionsLimitPolicy)
}

func getUserConcurrentSessionsLimitPolicy(d *schema.ResourceData) UserConcurrentSessionsLimitPolicy {
	resourceAsMap := d.Get("user_concurrent_sessions_limit").([]interface{})[0].(map[string]interface{})
	var userConcurrentSessionsLimitPolicy UserConcurrentSessionsLimitPolicy
	userConcurrentSessionsLimitPolicy.Enabled = resourceAsMap["enabled"].(bool)
	userConcurrentSessionsLimitPolicy.MaxConcurrentSessions = resourceAsMap["max_concurrent_sessions"].(int)
	return userConcurrentSessionsLimitPolicy
}
