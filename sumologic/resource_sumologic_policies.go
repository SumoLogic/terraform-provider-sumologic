package sumologic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Required: true,
			},
			"data_access_level": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"max_user_session_timeout": {
				Type:     schema.TypeString,
				Required: true,
			},
			"search_audit": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"share_dashboards_outside_organization": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"oauth_cimd_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enabled", "enabled-pre-registered-only"}, false),
			},
			"unused_access_keys_disable_after_in_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 365),
			},
			"access_keys_lifetime_in_days": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"0", "30", "45", "60", "90", "180", "365"}, false),
			},
			"user_concurrent_sessions_limit": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
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
	if policies.OAuthCimd != nil {
		d.Set("oauth_cimd_policy", policies.OAuthCimd.OAuthCimdPolicy)
	}
	if policies.DisableUnusedAccessKeys != nil {
		d.Set("unused_access_keys_disable_after_in_days", policies.DisableUnusedAccessKeys.UnusedAccessKeysDisableAfterInDays)
	}
	if policies.AccessKeysLifetime != nil {
		d.Set("access_keys_lifetime_in_days", policies.AccessKeysLifetime.AccessKeysLifetimeInDays)
	}
}

func resourceToPolicies(d *schema.ResourceData) Policies {
	var policies Policies
	policies.Audit = AuditPolicy{d.Get("audit").(bool)}
	policies.DataAccessLevel = DataAccessLevelPolicy{d.Get("data_access_level").(bool)}
	policies.MaxUserSessionTimeout = MaxUserSessionTimeoutPolicy{d.Get("max_user_session_timeout").(string)}
	policies.SearchAudit = SearchAuditPolicy{d.Get("search_audit").(bool)}
	policies.ShareDashboardsOutsideOrganization = ShareDashboardsOutsideOrganizationPolicy{d.Get("share_dashboards_outside_organization").(bool)}
	policies.UserConcurrentSessionsLimit = getUserConcurrentSessionsLimitPolicy(d)
	// Writing these policies can have side effects even when the value is unchanged. Computed attributes are filled
	// from state on update, so only send the ones that are set on create or that actually changed.
	// GetOkExists keeps an explicit 0 ("never") distinguishable from an unset attribute.
	isCreate := d.Id() == ""
	if v, ok := d.GetOk("oauth_cimd_policy"); ok && (isCreate || d.HasChange("oauth_cimd_policy")) {
		policies.OAuthCimd = &OAuthCimdPolicy{v.(string)}
	}
	if v, ok := d.GetOkExists("unused_access_keys_disable_after_in_days"); ok && (isCreate || d.HasChange("unused_access_keys_disable_after_in_days")) {
		policies.DisableUnusedAccessKeys = &DisableUnusedAccessKeysPolicy{v.(int)}
	}
	if v, ok := d.GetOk("access_keys_lifetime_in_days"); ok && (isCreate || d.HasChange("access_keys_lifetime_in_days")) {
		policies.AccessKeysLifetime = &AccessKeysLifetimePolicy{v.(string)}
	}
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
