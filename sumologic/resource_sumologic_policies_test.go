package sumologic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPolicies_create(t *testing.T) {
	policies := Policies{
		Audit:                              AuditPolicy{Enabled: true},
		DataAccessLevel:                    DataAccessLevelPolicy{Enabled: true},
		MaxUserSessionTimeout:              MaxUserSessionTimeoutPolicy{MaxUserSessionTimeout: "1h"},
		SearchAudit:                        SearchAuditPolicy{Enabled: true},
		ShareDashboardsOutsideOrganization: ShareDashboardsOutsideOrganizationPolicy{Enabled: true},
		UserConcurrentSessionsLimit:        UserConcurrentSessionsLimitPolicy{Enabled: true, MaxConcurrentSessions: 70},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPoliciesDestroy(),
		Steps: []resource.TestStep{
			{
				Config: newPoliciesConfig("test_create", &policies),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPoliciesExists("sumologic_policies.test_create"),
					testPoliciesCheckResourceAttr("sumologic_policies.test_create", &policies),
				),
			},
		},
	})
}

func TestAccPolicies_update(t *testing.T) {
	policies := Policies{
		Audit:                              AuditPolicy{Enabled: true},
		DataAccessLevel:                    DataAccessLevelPolicy{Enabled: true},
		MaxUserSessionTimeout:              MaxUserSessionTimeoutPolicy{MaxUserSessionTimeout: "15m"},
		SearchAudit:                        SearchAuditPolicy{Enabled: true},
		ShareDashboardsOutsideOrganization: ShareDashboardsOutsideOrganizationPolicy{Enabled: true},
		UserConcurrentSessionsLimit:        UserConcurrentSessionsLimitPolicy{Enabled: true, MaxConcurrentSessions: 80},
	}

	updatedPolicies := Policies{
		Audit:                              AuditPolicy{Enabled: false},
		DataAccessLevel:                    DataAccessLevelPolicy{Enabled: false},
		MaxUserSessionTimeout:              MaxUserSessionTimeoutPolicy{MaxUserSessionTimeout: "7d"},
		SearchAudit:                        SearchAuditPolicy{Enabled: false},
		ShareDashboardsOutsideOrganization: ShareDashboardsOutsideOrganizationPolicy{Enabled: false},
		UserConcurrentSessionsLimit:        UserConcurrentSessionsLimitPolicy{Enabled: false, MaxConcurrentSessions: 100},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPoliciesDestroy(),
		Steps: []resource.TestStep{
			{
				Config: newPoliciesConfig("test_update", &policies),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPoliciesExists("sumologic_policies.test_update"),
					testPoliciesCheckResourceAttr("sumologic_policies.test_update", &policies),
				),
			},
			{
				Config: newPoliciesConfig("test_update", &updatedPolicies),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPoliciesExists("sumologic_policies.test_update"),
					testPoliciesCheckResourceAttr("sumologic_policies.test_update", &updatedPolicies),
				),
			},
		},
	})
}

func testAccCheckPoliciesDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "sumologic_policies" {
				continue
			}

			policies, err := client.GetPolicies()
			if err != nil {
				return fmt.Errorf("Encountered an error: " + err.Error())
			}

			if (*policies) != DefaultPolicies {
				return fmt.Errorf("Policies weren't reset properly")
			}
		}
		return nil
	}
}

func testAccCheckPoliciesExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Policies not found: %s", name)
		}
		return nil
	}
}

func testPoliciesCheckResourceAttr(resourceName string, policies *Policies) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceName, "audit", strconv.FormatBool(policies.Audit.Enabled)),
			resource.TestCheckResourceAttr(resourceName, "data_access_level", strconv.FormatBool(policies.DataAccessLevel.Enabled)),
			resource.TestCheckResourceAttr(resourceName, "max_user_session_timeout", policies.MaxUserSessionTimeout.MaxUserSessionTimeout),
			resource.TestCheckResourceAttr(resourceName, "search_audit", strconv.FormatBool(policies.SearchAudit.Enabled)),
			resource.TestCheckResourceAttr(resourceName, "share_dashboards_outside_organization", strconv.FormatBool(policies.ShareDashboardsOutsideOrganization.Enabled)),
			resource.TestCheckResourceAttr(resourceName, "user_concurrent_sessions_limit.0.enabled", strconv.FormatBool(policies.UserConcurrentSessionsLimit.Enabled)),
			resource.TestCheckResourceAttr(resourceName, "user_concurrent_sessions_limit.0.max_concurrent_sessions", strconv.Itoa(policies.UserConcurrentSessionsLimit.MaxConcurrentSessions)),
		)
		return f(s)
	}
}

func newPoliciesConfig(label string, policies *Policies) string {
	return fmt.Sprintf(`
resource "sumologic_policies" "%s" {
  audit = %t
  data_access_level = %t
  max_user_session_timeout = "%s"
  search_audit = %t
  share_dashboards_outside_organization = %t
  user_concurrent_sessions_limit {
    enabled = %t
    max_concurrent_sessions = %d
  }
}`, label,
		policies.Audit.Enabled,
		policies.DataAccessLevel.Enabled,
		policies.MaxUserSessionTimeout.MaxUserSessionTimeout,
		policies.SearchAudit.Enabled,
		policies.ShareDashboardsOutsideOrganization.Enabled,
		policies.UserConcurrentSessionsLimit.Enabled,
		policies.UserConcurrentSessionsLimit.MaxConcurrentSessions)
}
