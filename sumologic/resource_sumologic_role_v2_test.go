package sumologic

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicRoleV2_basic(t *testing.T) {
	var roleV2 RoleV2
	testName := acctest.RandomWithPrefix("tf-acc-test")
	testAuditDataFilter := "info"
	testSelectionType := "All"
	testCapabilities := []string{"\"manageContent\""}
	testDescription := "Manage data of the org."
	testSecurityDataFilter := "error"
	testLogAnalyticsFilter := "!_sourceCategory=collector"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRoleV2Destroy(roleV2),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSumologicRoleV2ConfigImported(testName, testAuditDataFilter, testSelectionType, testCapabilities, testDescription, testSecurityDataFilter, testLogAnalyticsFilter),
			},
			{
				ResourceName:      "sumologic_role_v2.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccSumologicRoleV2_create(t *testing.T) {
	var roleV2 RoleV2
	testName := acctest.RandomWithPrefix("tf-acc-test")
	testAuditDataFilter := "info"
	testSelectionType := "Allow"
	testCapabilities := []string{"\"manageContent\""}
	testDescription := "Manage data of the org."
	testSecurityDataFilter := "error"
	testLogAnalyticsFilter := "!_sourceCategory=collector"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRoleV2Destroy(roleV2),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicRoleV2(testName, testAuditDataFilter, testSelectionType, testCapabilities, testDescription, testSecurityDataFilter, testLogAnalyticsFilter),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoleV2Exists("sumologic_role_v2.test", &roleV2, t),
					testAccCheckRoleV2Attributes("sumologic_role_v2.test"),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "audit_data_filter", testAuditDataFilter),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "selection_type", testSelectionType),
					resource.TestCheckTypeSetElemAttr("sumologic_role_v2.test", "capabilities.*", strings.Replace(testCapabilities[0], "\"", "", 2)),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "description", testDescription),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "security_data_filter", testSecurityDataFilter),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "log_analytics_filter", testLogAnalyticsFilter),
				),
			},
		},
	})
}

func TestAccSumologicRoleV2_update(t *testing.T) {
	var roleV2 RoleV2
	testName := acctest.RandomWithPrefix("tf-acc-test")
	testAuditDataFilter := "info"
	testSelectionType := "Allow"
	testCapabilities := []string{"\"manageContent\""}
	testDescription := "Manage data of the org."
	testSecurityDataFilter := "error"
	testLogAnalyticsFilter := "!_sourceCategory=collector"

	testUpdatedName := acctest.RandomWithPrefix("tf-acc-test")
	testUpdatedAuditDataFilter := "infoUpdate"
	testUpdatedSelectionType := "All"
	testUpdatedCapabilities := []string{"\"manageContent\""}
	testUpdatedDescription := "Manage data of the org.Update"
	testUpdatedSecurityDataFilter := "errorUpdate"
	testUpdatedLogAnalyticsFilter := "!_sourceCategory=collectorUpdate"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRoleV2Destroy(roleV2),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicRoleV2(testName, testAuditDataFilter, testSelectionType, testCapabilities, testDescription, testSecurityDataFilter, testLogAnalyticsFilter),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoleV2Exists("sumologic_role_v2.test", &roleV2, t),
					testAccCheckRoleV2Attributes("sumologic_role_v2.test"),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "audit_data_filter", testAuditDataFilter),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "selection_type", testSelectionType),
					resource.TestCheckTypeSetElemAttr("sumologic_role_v2.test", "capabilities.*", strings.Replace(testCapabilities[0], "\"", "", 2)),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "description", testDescription),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "security_data_filter", testSecurityDataFilter),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "log_analytics_filter", testLogAnalyticsFilter),
				),
			},
			{
				Config: testAccSumologicRoleV2Update(testUpdatedName, testUpdatedAuditDataFilter, testUpdatedSelectionType, testUpdatedCapabilities, testUpdatedDescription, testUpdatedSecurityDataFilter, testUpdatedLogAnalyticsFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "name", testUpdatedName),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "audit_data_filter", testUpdatedAuditDataFilter),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "selection_type", testUpdatedSelectionType),
					resource.TestCheckTypeSetElemAttr("sumologic_role_v2.test", "capabilities.*", strings.Replace(testUpdatedCapabilities[0], "\"", "", 2)),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "description", testUpdatedDescription),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "security_data_filter", testUpdatedSecurityDataFilter),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "log_analytics_filter", testUpdatedLogAnalyticsFilter),
				),
			},
		},
	})
}

func testAccCheckRoleV2Destroy(roleV2 RoleV2) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			u, err := client.GetRoleV2(id)
			if err != nil {
				return fmt.Errorf("Encountered an error: %w", err)
			}
			if u != nil {
				return fmt.Errorf("RoleV2 %s still exists", id)
			}
		}
		return nil
	}
}

func testAccCheckRoleV2Exists(name string, roleV2 *RoleV2, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			//need this so that we don't get an unused import error for strconv in some cases
			return fmt.Errorf("Error = %s. RoleV2 not found: %s", strconv.FormatBool(ok), name)
		}

		//need this so that we don't get an unused import error for strings in some cases
		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("RoleV2 ID is not set")
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)
		newRoleV2, err := client.GetRoleV2(id)
		if err != nil {
			return fmt.Errorf("RoleV2 %s not found", id)
		}
		roleV2 = newRoleV2
		return nil
	}
}
func testAccCheckSumologicRoleV2ConfigImported(name string, auditDataFilter string, selectionType string, capabilities []string, description string, securityDataFilter string, logAnalyticsFilter string) string {
	return fmt.Sprintf(`
resource "sumologic_role_v2" "foo" {
      name = "%s"
      audit_data_filter = "%s"
      selection_type = "%s"
      capabilities = %v
      description = "%s"
      security_data_filter = "%s"
      log_analytics_filter = "%s"
}
`, name, auditDataFilter, selectionType, capabilities, description, securityDataFilter, logAnalyticsFilter)
}

func testAccSumologicRoleV2(name string, auditDataFilter string, selectionType string, capabilities []string, description string, securityDataFilter string, logAnalyticsFilter string) string {
	return fmt.Sprintf(`
resource "sumologic_role_v2" "test" {
    selected_views {
		view_name = "sumologic_default"
	}
    name = "%s"
    audit_data_filter = "%s"
    selection_type = "%s"
    capabilities = %v
    description = "%s"
    security_data_filter = "%s"
    log_analytics_filter = "%s"
}
`, name, auditDataFilter, selectionType, capabilities, description, securityDataFilter, logAnalyticsFilter)
}

func testAccSumologicRoleV2Update(name string, auditDataFilter string, selectionType string, capabilities []string, description string, securityDataFilter string, logAnalyticsFilter string) string {
	return fmt.Sprintf(`
resource "sumologic_role_v2" "test" {
      name = "%s"
      audit_data_filter = "%s"
      selection_type = "%s"
      capabilities = %v
      description = "%s"
      security_data_filter = "%s"
      log_analytics_filter = "%s"
}
`, name, auditDataFilter, selectionType, capabilities, description, securityDataFilter, logAnalyticsFilter)
}

func TestAccSumologicRoleV2_selectionType_DefaultsToAll(t *testing.T) {
	var roleV2 RoleV2
	testName := acctest.RandomWithPrefix("tf-acc-test")
	testCapabilities := []string{"\"manageContent\""}
	testDescription := "Role with default selection_type"
	testAuditDataFilter := "info"
	testSecurityDataFilter := "error"
	testLogAnalyticsFilter := "!_sourceCategory=collector"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRoleV2Destroy(roleV2),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "sumologic_role_v2" "test" {
    name                 = "%s"
    audit_data_filter    = "%s"
    capabilities         = %v
    description          = "%s"
    security_data_filter = "%s"
    log_analytics_filter = "%s"
}
`, testName, testAuditDataFilter, testCapabilities, testDescription, testSecurityDataFilter, testLogAnalyticsFilter),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoleV2Exists("sumologic_role_v2.test", &roleV2, t),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "selection_type", "All"),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "sumologic_role_v2" "test" {
    name                 = "%s"
    audit_data_filter    = "%s"
    selection_type       = ""
    capabilities         = %v
    description          = "%s"
    security_data_filter = "%s"
    log_analytics_filter = "%s"
}
`, testName, testAuditDataFilter, testCapabilities, testDescription, testSecurityDataFilter, testLogAnalyticsFilter),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoleV2Exists("sumologic_role_v2.test", &roleV2, t),
					// Even if user sets "", internally defaults to "All"
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "selection_type", "All"),
				),
			},
		},
	})
}

func testAccCheckRoleV2Attributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(name, "name"),
			resource.TestCheckResourceAttrSet(name, "audit_data_filter"),
			resource.TestCheckResourceAttrSet(name, "selection_type"),
			resource.TestCheckResourceAttrSet(name, "description"),
			resource.TestCheckResourceAttrSet(name, "security_data_filter"),
			resource.TestCheckResourceAttrSet(name, "log_analytics_filter"),
		)
		return f(s)
	}
}

func TestSumologicRoleV2_selectionType_Validation(t *testing.T) {
	validate := resourceSumologicRoleV2().Schema["selection_type"].ValidateFunc
	if validate == nil {
		t.Fatal("selection_type must define a ValidateFunc")
	}

	for _, valid := range []string{"", "All", "Allow", "Deny"} {
		if _, errs := validate(valid, "selection_type"); len(errs) != 0 {
			t.Errorf("expected %q to be accepted, got %v", valid, errs)
		}
	}

	for _, invalid := range []string{"AllViews", "all", "None"} {
		if _, errs := validate(invalid, "selection_type"); len(errs) == 0 {
			t.Errorf("expected %q to be rejected", invalid)
		}
	}
}

func TestRoleV2SelectedViewsIgnoresOrder(t *testing.T) {
	views := func(names ...string) []interface{} {
		out := []interface{}{}
		for _, n := range names {
			out = append(out, map[string]interface{}{"view_name": n, "view_filter": ""})
		}
		return out
	}
	r := resourceSumologicRoleV2()

	configured := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name":           "r",
		"selected_views": views("index_a", "index_b", "index_c"),
	})

	// the API may return the views in a different order than configured
	read := r.TestResourceData()
	if err := read.Set("selected_views", flattenSelectedViews([]ViewFilterDefinition{
		{ViewName: "index_c"}, {ViewName: "index_a"}, {ViewName: "index_b"},
	})); err != nil {
		t.Fatal(err)
	}

	want := configured.Get("selected_views").(*schema.Set)
	got := read.Get("selected_views").(*schema.Set)
	if !want.Equal(got) {
		t.Errorf("selected_views differ by order: want %v, got %v", want.List(), got.List())
	}
	if n := len(resourceToRoleV2(read).SelectedViews); n != 3 {
		t.Errorf("expected 3 selected views, got %d", n)
	}
}

func TestAccSumologicRoleV2_selectedViewsOrder(t *testing.T) {
	var roleV2 RoleV2
	testName := acctest.RandomWithPrefix("tf-acc-test")
	config := fmt.Sprintf(`
resource "sumologic_role_v2" "test" {
    name           = "%s"
    selection_type = "Allow"
    capabilities   = ["manageContent"]
    selected_views {
        view_name = "sumologic_default"
    }
    selected_views {
        view_name = "sumologic_audit_events"
    }
    selected_views {
        view_name = "sumologic_volume"
    }
}
`, testName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRoleV2Destroy(roleV2),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoleV2Exists("sumologic_role_v2.test", &roleV2, t),
					resource.TestCheckResourceAttr("sumologic_role_v2.test", "selected_views.#", "3"),
				),
			},
			{
				// the API does not guarantee view order; a re-plan must be empty
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}
