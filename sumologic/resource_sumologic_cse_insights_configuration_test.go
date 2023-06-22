package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicSCEInsightsConfiguration_create(t *testing.T) {
	SkipCseTest(t)

	var insightConfiguration CSEInsightsConfiguration
	nLookbackDays := 10.0
	nThreshold := 13.0
	nGlobalSignalSuppressionWindow := 48.0
	resourceName := "sumologic_cse_insights_configuration.insights_configuration"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEInsightsConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSEInsightsConfigurationConfig(nLookbackDays, nThreshold, nGlobalSignalSuppressionWindow),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEInsightsConfigurationExists(resourceName, &insightConfiguration),
					testCheckInsightsConfigurationValues(&insightConfiguration, nLookbackDays, nThreshold, nGlobalSignalSuppressionWindow),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccCSEInsightsConfigurationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_insights_configuration" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Insights Status destruction check: CSE Insights Status ID is not set")
		}

		s, err := client.GetCSEInsightsConfiguration()
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			if s.Threshold != nil && s.LookbackDays != nil && s.GlobalSignalSuppressionWindow != nil {
				return fmt.Errorf("insight Configuration still exists")
			}
		}
	}
	return nil
}

func testCreateCSEInsightsConfigurationConfig(nLookbackDays float64, nThreshold float64, nGlobalSignalSuppressionWindow float64) string {
	return fmt.Sprintf(`
resource "sumologic_cse_insights_configuration" "insights_configuration" {
	lookback_days = "%f"
	threshold = "%f"
	global_signal_suppression_window = "%f"
}
`, nLookbackDays, nThreshold, nGlobalSignalSuppressionWindow)
}

func testCheckCSEInsightsConfigurationExists(n string, insightConfiguration *CSEInsightsConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("insight Configuration ID is not set")
		}

		c := testAccProvider.Meta().(*Client)
		insightConfigurationResp, err := c.GetCSEInsightsConfiguration()
		if err != nil {
			return err
		}

		*insightConfiguration = *insightConfigurationResp

		return nil
	}
}

func testCheckInsightsConfigurationValues(insightConfiguration *CSEInsightsConfiguration, nLookbackDays float64, nThreshold float64, nGlobalSignalSuppressionWindow float64) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if *insightConfiguration.LookbackDays != nLookbackDays {
			return fmt.Errorf("bad lookback days, expected \"%f\", got: %#v", nLookbackDays, insightConfiguration.LookbackDays)
		}
		if *insightConfiguration.Threshold != nThreshold {
			return fmt.Errorf("bad threshold, expected \"%f\", got: %#v", nThreshold, insightConfiguration.Threshold)
		}
		if *insightConfiguration.GlobalSignalSuppressionWindow != nGlobalSignalSuppressionWindow {
			return fmt.Errorf("bad global signal suppression window, expected \"%f\", got: %#v", nGlobalSignalSuppressionWindow, insightConfiguration.GlobalSignalSuppressionWindow)
		}

		return nil
	}
}
