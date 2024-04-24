package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicSCEEntityNormalizationConfiguration_create(t *testing.T) {
	SkipCseTest(t)

	var entityNormalizationConfiguration CSEEntityNormalizationConfiguration
	nWindowsNormalizationEnabled := true
	nFqdnNormalizationEnabled := true
	nAwsNormalizationEnabled := true
	nDefaultNormalizedDomain := "test.domain"
	nNormalizeHostnames := true
	nNormalizeUsernames := true
	nNormalizedDomain := "normalized.domain"
	nRawDomain := "raw.domain"

	resourceName := "sumologic_cse_entity_normalization_configuration.entity_normalization_configuration"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEEntityNormalizationConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSEEntityNormalizationConfigurationConfig(nWindowsNormalizationEnabled, nFqdnNormalizationEnabled, nAwsNormalizationEnabled, nDefaultNormalizedDomain, nNormalizeHostnames, nNormalizeUsernames, nNormalizedDomain, nRawDomain),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEEntityNormalizationConfigurationExists(resourceName, &entityNormalizationConfiguration),
					testCheckEntityNormalizationConfigurationValues(&entityNormalizationConfiguration, nWindowsNormalizationEnabled, nFqdnNormalizationEnabled, nAwsNormalizationEnabled, nDefaultNormalizedDomain,
						nNormalizeHostnames, nNormalizeUsernames, nNormalizedDomain, nRawDomain),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccCSEEntityNormalizationConfigurationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_entity_normalization_configuration" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Entity Normalization destruction check: CSE Entity Normalization ID is not set")
		}

		s, err := client.GetCSEEntityNormalizationConfiguration()
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			if s.NormalizeHostnames != false && s.NormalizeUsernames != false {
				return fmt.Errorf("entity normalization Configuration still exists")
			}
		}
	}
	return nil
}

func testCreateCSEEntityNormalizationConfigurationConfig(nWindowsNormalizationEnabled bool, nFqdnNormalizationEnabled bool, nAwsNormalizationEnabled bool, nDefaultNormalizedDomain string,
	nNormalizeHostnames bool, nNormalizeUsernames bool, nNormalizedDomain string, nRawDomain string) string {
	return fmt.Sprintf(`
resource "sumologic_cse_entity_normalization_configuration" "entity_normalization_configuration" {
	windows_normalization_enabled = "%t"
	fqdn_normalization_enabled = "%t"
	aws_normalization_enabled = "%t"
	default_normalized_domain = "%s"
	normalize_hostnames = "%t"
	normalize_usernames = "%t"
	domain_mappings{
		normalized_domain = "%s"
		raw_domain = "%s"
	}
}
`, nWindowsNormalizationEnabled, nFqdnNormalizationEnabled, nAwsNormalizationEnabled, nDefaultNormalizedDomain,
		nNormalizeHostnames, nNormalizeUsernames, nNormalizedDomain, nRawDomain)
}

func testCheckCSEEntityNormalizationConfigurationExists(n string, insightConfiguration *CSEEntityNormalizationConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("insight Configuration ID is not set")
		}

		c := testAccProvider.Meta().(*Client)
		insightConfigurationResp, err := c.GetCSEEntityNormalizationConfiguration()
		if err != nil {
			return err
		}

		*insightConfiguration = *insightConfigurationResp

		return nil
	}
}

func testCheckEntityNormalizationConfigurationValues(entityNormalizationConfiguration *CSEEntityNormalizationConfiguration, nWindowsNormalizationEnabled bool, nFqdnNormalizationEnabled bool, nAwsNormalizationEnabled bool, nDefaultNormalizedDomain string,
	nNormalizeHostnames bool, nNormalizeUsernames bool, nNormalizedDomain string, nRawDomain string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if entityNormalizationConfiguration.WindowsNormalizationEnabled != nWindowsNormalizationEnabled {
			return fmt.Errorf("bad windows_normalization_enabled, expected \"%t\", got: %#v", nWindowsNormalizationEnabled, entityNormalizationConfiguration.WindowsNormalizationEnabled)
		}
		if entityNormalizationConfiguration.FqdnNormalizationEnabled != nFqdnNormalizationEnabled {
			return fmt.Errorf("bad fqdn_normalization_enabled, expected \"%t\", got: %#v", nFqdnNormalizationEnabled, entityNormalizationConfiguration.FqdnNormalizationEnabled)
		}
		if entityNormalizationConfiguration.AwsNormalizationEnabled != nAwsNormalizationEnabled {
			return fmt.Errorf("bad aws_normalization_enabled, expected \"%t\", got: %#v", nAwsNormalizationEnabled, entityNormalizationConfiguration.AwsNormalizationEnabled)
		}
		if entityNormalizationConfiguration.DefaultNormalizedDomain != nDefaultNormalizedDomain {
			return fmt.Errorf("bad default_normalized_domain, expected \"%s\", got: %#v", nDefaultNormalizedDomain, entityNormalizationConfiguration.DefaultNormalizedDomain)
		}
		if entityNormalizationConfiguration.NormalizeHostnames != nNormalizeHostnames {
			return fmt.Errorf("bad normalize_hostnames, expected \"%t\", got: %#v", nNormalizeHostnames, entityNormalizationConfiguration.NormalizeHostnames)
		}
		if entityNormalizationConfiguration.NormalizeUsernames != nNormalizeUsernames {
			return fmt.Errorf("bad normalize_usernames, expected \"%t\", got: %#v", nNormalizeUsernames, entityNormalizationConfiguration.NormalizeUsernames)
		}
		if entityNormalizationConfiguration.DomainMappings[0].NormalizedDomain != nNormalizedDomain {
			return fmt.Errorf("bad normalized_domain, expected \"%s\", got: %#v", nNormalizedDomain, entityNormalizationConfiguration.DomainMappings[0].NormalizedDomain)
		}
		if entityNormalizationConfiguration.DomainMappings[0].RawDomain != nRawDomain {
			return fmt.Errorf("bad raw_domain, expected \"%s\", got: %#v", nNormalizedDomain, entityNormalizationConfiguration.DomainMappings[0].RawDomain)
		}

		return nil
	}
}
