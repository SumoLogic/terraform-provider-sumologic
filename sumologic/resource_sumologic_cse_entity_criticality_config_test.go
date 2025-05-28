package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicSCEEntityCriticalityConfig_create(t *testing.T) {
	SkipCseTest(t)

	var entityCriticalityConfig CSEEntityCriticalityConfig
	nName := "New Entity Criticality"
	nSeverityExpression := "severity + 2"
	resourceName := "sumologic_cse_entity_criticality_config.entity_criticality_config"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEEntityCriticalityConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSEEntityCriticalityConfigConfig(nName, nSeverityExpression),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEEntityCriticalityConfigExists(resourceName, &entityCriticalityConfig),
					testCheckEntityCriticalityConfigValues(&entityCriticalityConfig, nName, nSeverityExpression),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func TestAccSumologicSCEEntityCriticalityConfig_update(t *testing.T) {
	SkipCseTest(t)

	var entityCriticalityConfig CSEEntityCriticalityConfig
	nName := "New Entity Criticality"
	nSeverityExpression := "severity + 2"
	resourceName := "sumologic_cse_entity_criticality_config.entity_criticality_config"
	uSeverity := "severity + 3"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEEntityCriticalityConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSEEntityCriticalityConfigConfig(nName, nSeverityExpression),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEEntityCriticalityConfigExists(resourceName, &entityCriticalityConfig),
					testCheckEntityCriticalityConfigValues(&entityCriticalityConfig, nName, nSeverityExpression),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testCreateCSEEntityCriticalityConfigConfig(nName, uSeverity),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEEntityCriticalityConfigExists(resourceName, &entityCriticalityConfig),
					testCheckEntityCriticalityConfigValues(&entityCriticalityConfig, nName, uSeverity),
				),
			},
		},
	})
}

func testAccCSEEntityCriticalityConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_entity_criticality_config" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Entity Criticality Config destruction check: CSE Entity Criticality Config ID is not set")
		}

		s, err := client.GetCSEEntityCriticalityConfig(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("entity Criticality Config still exists")
		}
	}
	return nil
}

func testCreateCSEEntityCriticalityConfigConfig(nName string, nSeverityExpression string) string {
	return fmt.Sprintf(`
resource "sumologic_cse_entity_criticality_config" "entity_criticality_config" {
	name = "%s"
	severity_expression = "%s"
}
`, nName, nSeverityExpression)
}

func testCheckCSEEntityCriticalityConfigExists(n string, entityCriticalityConfig *CSEEntityCriticalityConfig) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("entity Criticality Config ID is not set")
		}

		c := testAccProvider.Meta().(*Client)
		entityCriticalityConfigResp, err := c.GetCSEEntityCriticalityConfig(rs.Primary.ID)
		if err != nil {
			return err
		}

		*entityCriticalityConfig = *entityCriticalityConfigResp

		return nil
	}
}

func testCheckEntityCriticalityConfigValues(entityCriticalityConfig *CSEEntityCriticalityConfig, nName string, nSeverityExpression string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if entityCriticalityConfig.Name != nName {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", nName, entityCriticalityConfig.Name)
		}
		if entityCriticalityConfig.SeverityExpression != nSeverityExpression {
			return fmt.Errorf("bad severity_expression, expected \"%s\", got: %#v", nSeverityExpression, entityCriticalityConfig.SeverityExpression)
		}

		return nil
	}
}
