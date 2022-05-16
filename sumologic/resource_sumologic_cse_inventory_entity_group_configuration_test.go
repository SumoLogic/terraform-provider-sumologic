package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicCSEInventoryEntityGroupConfiguration_createAndUpdate(t *testing.T) {
	SkipCseTest(t)

	var InventoryEntityGroupConfiguration CSEEntityGroupConfiguration
	criticality := "HIGH"
	description := "Test description"
	group := "goo"
	inventoryType := "computer"
	name := "Entity Group Configuration Tf test"
	suppressed := false
	tag := "foo"
	nameUpdated := "Updated Entity Group Configuration"

	resourceName := "sumologic_cse_inventory_entity_group_configuration.inventory_entity_group_configuration"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEInventoryEntityGroupConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSEInventoryEntityGroupConfigurationConfig(criticality, description, group,
					inventoryType, name, suppressed, tag),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEInventoryEntityGroupConfigurationExists(resourceName, &InventoryEntityGroupConfiguration),
					testCheckInventoryEntityGroupConfigurationValues(&InventoryEntityGroupConfiguration, criticality,
						description, group, inventoryType, name, suppressed, tag),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testCreateCSEInventoryEntityGroupConfigurationConfig(criticality, description, group,
					inventoryType, nameUpdated, suppressed, tag),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEInventoryEntityGroupConfigurationExists(resourceName, &InventoryEntityGroupConfiguration),
					testCheckInventoryEntityGroupConfigurationValues(&InventoryEntityGroupConfiguration, criticality,
						description, group, inventoryType, nameUpdated, suppressed, tag),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccCSEInventoryEntityGroupConfigurationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_inventory_entity_group_configuration" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Inventory Entity Group Configuration destruction check: CSE Inventory Entity" +
				" Group Configuration ID is not set")
		}

		s, err := client.GetCSEntityGroupConfiguration(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("Inventory Entity Group Configuration still exists")
		}
	}
	return nil
}

func testCreateCSEInventoryEntityGroupConfigurationConfig(
	criticality string, description string, group string,
	inventoryType string, name string, suppressed bool, tag string) string {
	return fmt.Sprintf(`
resource "sumologic_cse_inventory_entity_group_configuration" "inventory_entity_group_configuration" {
	criticality = "%s"
    description = "%s"
	groups = ["%s"]
	inventory_type = "%s"
	name = "%s"
	suppressed = %t
 	tags = ["%s"]
}
`, criticality, description, group, inventoryType, name, suppressed, tag)
}

func testCheckCSEInventoryEntityGroupConfigurationExists(n string, InventoryEntityGroupConfiguration *CSEEntityGroupConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("inventory entity group configuration ID is not set")
		}

		c := testAccProvider.Meta().(*Client)
		InventoryEntityGroupConfigurationResp, err := c.GetCSEntityGroupConfiguration(rs.Primary.ID)
		if err != nil {
			return err
		}

		*InventoryEntityGroupConfiguration = *InventoryEntityGroupConfigurationResp

		return nil
	}
}

func testCheckInventoryEntityGroupConfigurationValues(InventoryEntityGroupConfiguration *CSEEntityGroupConfiguration,
	criticality string, description string, group string,
	inventoryType string, name string, suppressed bool, tag string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if InventoryEntityGroupConfiguration.Criticality != criticality {
			return fmt.Errorf("bad criticality, expected \"%s\", got %#v", criticality, InventoryEntityGroupConfiguration.Criticality)
		}
		if InventoryEntityGroupConfiguration.Description != description {
			return fmt.Errorf("bad description, expected \"%s\", got %#v", description, InventoryEntityGroupConfiguration.Description)
		}
		if InventoryEntityGroupConfiguration.Groups[0] != group {
			return fmt.Errorf("bad group, expected \"%s\", got %#v", tag, InventoryEntityGroupConfiguration.Groups[0])
		}
		if InventoryEntityGroupConfiguration.InventoryType != inventoryType {
			return fmt.Errorf("bad inventoryType, expected \"%s\", got %#v", inventoryType, InventoryEntityGroupConfiguration.InventoryType)
		}
		if InventoryEntityGroupConfiguration.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got %#v", name, InventoryEntityGroupConfiguration.Name)
		}
		if InventoryEntityGroupConfiguration.Suppressed != suppressed {
			return fmt.Errorf("bad suppressed, expected \"%t\", got %#v", suppressed, InventoryEntityGroupConfiguration.Suppressed)
		}
		if InventoryEntityGroupConfiguration.Tags[0] != tag {
			return fmt.Errorf("bad tag, expected \"%s\", got %#v", tag, InventoryEntityGroupConfiguration.Tags[0])
		}

		return nil
	}
}
