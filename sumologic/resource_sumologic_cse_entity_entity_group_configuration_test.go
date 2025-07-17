package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicCSEEntityEntityGroupConfiguration_createAndUpdate(t *testing.T) {
	SkipCseTest(t)

	var EntityEntityGroupConfiguration CSEEntityGroupConfiguration
	criticality := "HIGH"
	description := "Test description"
	entityNamespace := "namespace"
	entityType := "_hostname"
	name := "Entity Group Configuration Tf test"
	suffix := "suffix"
	suppressed := false
	tag := "foo"
	nameUpdated := "Updated Entity Group Configuration"

	resourceName := "sumologic_cse_entity_entity_group_configuration.entity_entity_group_configuration"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEEntityEntityGroupConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSEEntityEntityGroupConfigurationConfig(criticality, description, entityNamespace,
					entityType, name, suffix, suppressed, tag),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEEntityEntityGroupConfigurationExists(resourceName, &EntityEntityGroupConfiguration),
					testCheckEntityEntityGroupConfigurationValues(&EntityEntityGroupConfiguration, criticality,
						description, entityNamespace, entityType, name, suffix, suppressed, tag),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testCreateCSEEntityEntityGroupConfigurationConfig(criticality, description, entityNamespace,
					entityType, nameUpdated, suffix, suppressed, tag),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEEntityEntityGroupConfigurationExists(resourceName, &EntityEntityGroupConfiguration),
					testCheckEntityEntityGroupConfigurationValues(&EntityEntityGroupConfiguration, criticality,
						description, entityNamespace, entityType, nameUpdated, suffix, suppressed, tag),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccCSEEntityEntityGroupConfigurationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_entity_entity_group_configuration" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Entity Entity Group Configuration destruction check: CSE Entity Entity" +
				" Group Configuration ID is not set")
		}

		s, err := client.GetCSEntityGroupConfiguration(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("Entity Entity Group Configuration still exists")
		}
	}
	return nil
}

func testCreateCSEEntityEntityGroupConfigurationConfig(
	criticality string, description string, entityNamespace string,
	entityType string, name string, suffix string, suppressed bool, tag string) string {
	return fmt.Sprintf(`
resource "sumologic_cse_entity_entity_group_configuration" "entity_entity_group_configuration" {
	criticality = "%s"
    description = "%s"
	entity_namespace = "%s"
	entity_type = "%s"
	name = "%s"
	suffix = "%s"
	suppressed = %t
 	tags = ["%s"]
}
`, criticality, description, entityNamespace, entityType, name, suffix, suppressed, tag)
}

func testCheckCSEEntityEntityGroupConfigurationExists(n string, EntityEntityGroupConfiguration *CSEEntityGroupConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("entity entity group configuration ID is not set")
		}

		c := testAccProvider.Meta().(*Client)
		EntityEntityGroupConfigurationResp, err := c.GetCSEntityGroupConfiguration(rs.Primary.ID)
		if err != nil {
			return err
		}

		*EntityEntityGroupConfiguration = *EntityEntityGroupConfigurationResp

		return nil
	}
}

func testCheckEntityEntityGroupConfigurationValues(EntityEntityGroupConfiguration *CSEEntityGroupConfiguration,
	criticality string, description string, entityNamespace string,
	entityType string, name string, suffix string, suppressed bool, tag string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if EntityEntityGroupConfiguration.Criticality != criticality {
			return fmt.Errorf("bad criticality, expected \"%s\", got %#v", criticality, EntityEntityGroupConfiguration.Criticality)
		}
		if EntityEntityGroupConfiguration.Description != description {
			return fmt.Errorf("bad description, expected \"%s\", got %#v", description, EntityEntityGroupConfiguration.Description)
		}
		if EntityEntityGroupConfiguration.EntityNamespace != entityNamespace {
			return fmt.Errorf("bad entityNamespace, expected \"%s\", got %#v", entityNamespace, EntityEntityGroupConfiguration.EntityNamespace)
		}
		if EntityEntityGroupConfiguration.EntityType != entityType {
			return fmt.Errorf("bad entityType, expected \"%s\", got %#v", entityType, EntityEntityGroupConfiguration.EntityType)
		}
		if EntityEntityGroupConfiguration.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got %#v", name, EntityEntityGroupConfiguration.Name)
		}
		if EntityEntityGroupConfiguration.Suffix != suffix {
			return fmt.Errorf("bad suffix, expected \"%s\", got %#v", suffix, EntityEntityGroupConfiguration.Suffix)
		}
		if EntityEntityGroupConfiguration.Suppressed != suppressed {
			return fmt.Errorf("bad suppressed, expected \"%t\", got %#v", suppressed, EntityEntityGroupConfiguration.Suppressed)
		}
		if EntityEntityGroupConfiguration.Tags[0] != tag {
			return fmt.Errorf("bad tag, expected \"%s\", got %#v", tag, EntityEntityGroupConfiguration.Tags[0])
		}

		return nil
	}
}
