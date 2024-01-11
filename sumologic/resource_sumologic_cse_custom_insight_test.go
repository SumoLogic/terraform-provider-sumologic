package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicCSECustomInsight_createAndUpdate(t *testing.T) {
	SkipCseTest(t)

	var CustomInsight CSECustomInsight
	description := "Test description"
	enabled := true
	ordered := true
	name := "Test Custom Insight"
	severity := "HIGH"
	minimumSignalSeverity := 5
	insightSeverity := "CRITICAL"
	signalName1 := "Some Signal Name *"
	signalName2 := "Some Other Signal Name *"
	tag := "foo"

	nameUpdated := "Updated Custom Insight"
	severityUpdated := "LOW"
	minimumSignalSeverityUpdated := 8

	resourceName := "sumologic_cse_custom_insight.custom_insight"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSECustomInsightDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSECustomInsightConfig(description, enabled,
					ordered, name, severity, minimumSignalSeverity, insightSeverity, signalName1, signalName2, tag),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSECustomInsightExists(resourceName, &CustomInsight),
					testCheckCustomInsightValues(&CustomInsight, description, enabled,
						ordered, name, severity, minimumSignalSeverity, insightSeverity, signalName1, signalName2, tag),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testCreateCSECustomInsightConfig(description, enabled,
					ordered, nameUpdated, severityUpdated, minimumSignalSeverityUpdated, insightSeverity, signalName1,
					signalName2, tag),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSECustomInsightExists(resourceName, &CustomInsight),
					testCheckCustomInsightValues(&CustomInsight, description, enabled,
						ordered, nameUpdated, severityUpdated, minimumSignalSeverityUpdated, insightSeverity, signalName1,
						signalName2, tag),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCSECustomInsightDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_custom_insight" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Custom Insight destruction check: CSE Custom Insight ID is not set")
		}

		s, err := client.GetCSECustomInsight(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("Custom Insight still exists")
		}
	}
	return nil
}

func testCreateCSECustomInsightConfig(
	description string, enabled bool, ordered bool, name string,
	severity string, minimumSignalSeverity int, insightSeverity string, signalName1 string, signalName2 string, tag string) string {
	return fmt.Sprintf(`
resource "sumologic_cse_custom_insight" "custom_insight" {
	description = "%s"
    enabled = %t
    ordered = %t
    name = "%s"
    severity = "%s"
    dynamic_severity {
		minimum_signal_severity = "%d"
		insight_severity = "%s"
	}
    signal_names = ["%s", "%s"]
    tags = ["%s"]
}
`, description, enabled, ordered, name, severity, minimumSignalSeverity, insightSeverity, signalName1, signalName2, tag)
}

func testCheckCSECustomInsightExists(n string, CustomInsight *CSECustomInsight) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("chain rule ID is not set")
		}

		c := testAccProvider.Meta().(*Client)
		CustomInsightResp, err := c.GetCSECustomInsight(rs.Primary.ID)
		if err != nil {
			return err
		}

		*CustomInsight = *CustomInsightResp

		return nil
	}
}

func testCheckCustomInsightValues(CustomInsight *CSECustomInsight, description string,
	enabled bool, ordered bool, name string, severity string, minimumSignalSeverity int, insightSeverity string, signalName1 string,
	signalName2 string, tag string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if CustomInsight.Description != description {
			return fmt.Errorf("bad description, expected \"%s\", got %#v", description, CustomInsight.Description)
		}
		if CustomInsight.Enabled != enabled {
			return fmt.Errorf("bad enabled, expected \"%t\", got %#v", enabled, CustomInsight.Enabled)
		}
		if CustomInsight.Ordered != ordered {
			return fmt.Errorf("bad ordered, expected \"%t\", got %#v", ordered, CustomInsight.Ordered)
		}
		if CustomInsight.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got %#v", name, CustomInsight.Name)
		}
		if CustomInsight.Severity != severity {
			return fmt.Errorf("bad severity, expected \"%s\", got %#v", severity, CustomInsight.Severity)
		}
		if CustomInsight.Severity != severity {
			return fmt.Errorf("bad severity, expected \"%s\", got %#v", severity, CustomInsight.Severity)
		}
		if CustomInsight.DynamicSeverity[0].MinimumSignalSeverity != minimumSignalSeverity {
			return fmt.Errorf("bad minimumSignalSeverity, expected \"%d\", got %#v", minimumSignalSeverity, CustomInsight.DynamicSeverity[0].MinimumSignalSeverity)
		}
		if CustomInsight.DynamicSeverity[0].InsightSeverity != insightSeverity {
			return fmt.Errorf("bad insightSeverity, expected \"%s\", got %#v", insightSeverity, CustomInsight.DynamicSeverity[0].InsightSeverity)
		}
		if CustomInsight.SignalNames[1] != signalName2 {
			return fmt.Errorf("bad signalName2, expected \"%s\", got %#v", signalName2, CustomInsight.SignalNames[1])
		}
		if CustomInsight.Tags[0] != tag {
			return fmt.Errorf("bad tag, expected \"%s\", got %#v", tag, CustomInsight.Tags[0])
		}

		return nil
	}
}
