package sumologic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const PUBLISHED_PLAYBOOK_NAME = "Test Playbook For Terraform"

func TestAccSumologicCsoarPlaybook_createShouldFail(t *testing.T) {
	rname := acctest.RandomWithPrefix("tf-acc-test-playbook")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccSumologicCsoarPlaybookConfigBasic(rname),
				ExpectError: regexp.MustCompile("playbooks cannot be created via Terraform"),
			},
		},
	})
}

func TestAccSumologicCsoarPlaybook_importAndUpdate(t *testing.T) {
	playbookName := PUBLISHED_PLAYBOOK_NAME
	resourceName := "sumologic_csoar_playbook.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCsoarPlaybookDestroy,
		Steps: []resource.TestStep{
			{
				Config:        testAccSumologicCsoarPlaybookConfigUpdate(playbookName),
				ImportState:   true,
				ResourceName:  resourceName,
				ImportStateId: playbookName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", playbookName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated via Terraform test"),
					resource.TestCheckResourceAttr(resourceName, "tags", "terraform,test"),
				),
			},
		},
	})
}

func TestAccSumologicCsoarPlaybook_booleanFieldUpdates(t *testing.T) {
	playbookName := PUBLISHED_PLAYBOOK_NAME
	resourceName := "sumologic_csoar_playbook.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCsoarPlaybookDestroy,
		Steps: []resource.TestStep{
			{
				Config:        testAccSumologicCsoarPlaybookConfigBooleanTest(playbookName, false),
				ImportState:   true,
				ResourceName:  resourceName,
				ImportStateId: playbookName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", playbookName),
					resource.TestCheckResourceAttr(resourceName, "is_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_published", "false"),
					resource.TestCheckResourceAttr(resourceName, "nested", "false"),
				),
			},
		},
	})
}

func TestAccSumologicCsoarPlaybook_updatedNameField(t *testing.T) {
	playbookName := PUBLISHED_PLAYBOOK_NAME
	resourceName := "sumologic_csoar_playbook.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCsoarPlaybookDestroy,
		Steps: []resource.TestStep{
			{
				Config:        testAccSumologicCsoarPlaybookConfigWithUpdatedName(playbookName),
				ImportState:   true,
				ResourceName:  resourceName,
				ImportStateId: playbookName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", playbookName),
					resource.TestCheckResourceAttr(resourceName, "updated_name", "New Playbook"),
					resource.TestCheckResourceAttr(resourceName, "description", "Test with updated_name"),
				),
			},
		},
	})
}

func testAccCheckCsoarPlaybookDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_csoar_playbook" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("CSOAR Playbook destruction check: playbook ID is not set")
		}
	}
	return nil
}

func testAccSumologicCsoarPlaybookConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "sumologic_csoar_playbook" "test" {
  name = "%s"
}`, name)
}

func testAccSumologicCsoarPlaybookConfigUpdate(name string) string {
	return fmt.Sprintf(`
resource "sumologic_csoar_playbook" "test" {
  name        = "%s"
  description = "Updated via Terraform test"
  tags        = "terraform,test"
}`, name)
}

func testAccSumologicCsoarPlaybookConfigBooleanTest(name string, enabledState bool) string {
	return fmt.Sprintf(`
resource "sumologic_csoar_playbook" "test" {
  name         = "%s"
  description  = "Boolean test"
  tags         = ""
  is_deleted   = false
  draft        = false
  is_published = %t
  nested       = %t
  type         = "General"
  is_enabled   = %t
  nodes        = jsonencode([])
  links        = jsonencode([])
}`, name, enabledState, enabledState, enabledState)
}

func testAccSumologicCsoarPlaybookConfigWithUpdatedName(name string) string {
	return fmt.Sprintf(`
resource "sumologic_csoar_playbook" "test" {
  name         = "%s"
  description  = "Test with updated_name"
  updated_name = "New Playbook Name"
  tags         = ""
  is_deleted   = false
  draft        = false
  is_published = true
  nested       = false
  type         = "General"
  is_enabled   = true
  nodes        = jsonencode([])
  links        = jsonencode([])
}`, name)
}
