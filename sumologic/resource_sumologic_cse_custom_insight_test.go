package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccSumologicCSECustomInsight_createAndUpdate(t *testing.T) {
	SkipCseTest(t)

	var CustomInsight CSECustomInsight
	description := "Test description"
	enabled := true
	ordered := true
	name := "Test Custom Insight"
	severity := "HIGH"
	signalName1 := "Some Signal Name *"
	signalName2 := "Some Other Signal Name *"
	tag := "foo"

	nameUpdated := "Updated Custom Insight"
	severityUpdated := "LOW"

	resourceName := "sumologic_cse_custom_insight.custom_insight"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSECustomInsightDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSECustomInsightConfig(description, enabled,
					ordered, name, severity, signalName1, signalName2, tag),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSECustomInsightExists(resourceName, &CustomInsight),
					testCheckCustomInsightValues(&CustomInsight, description, enabled,
						ordered, name, severity, signalName1, signalName2, tag),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testCreateCSECustomInsightConfig(description, enabled,
					ordered, nameUpdated, severityUpdated, signalName1,
					signalName2, tag),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSECustomInsightExists(resourceName, &CustomInsight),
					testCheckCustomInsightValues(&CustomInsight, description, enabled,
						ordered, nameUpdated, severityUpdated, signalName1,
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

func TestAccSumologicCSECustomInsightWithDynamicSeverity_createAndUpdate(t *testing.T) {
	SkipCseTest(t)

	var CustomInsight CSECustomInsight
	minimumSignalSeverity1 := 5
	dynamicSeverity1 := "MEDIUM"
	minimumSignalSeverity2 := 8
	dynamicSeverity2 := "HIGH"

	updatedMinimumSignalSeverity2 := 9
	updatedDynamicSeverity2 := "CRITICAL"

	resourceName := "sumologic_cse_custom_insight.custom_insight2"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSECustomInsightDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSECustomInsightConfigWithDynamicSeverity(minimumSignalSeverity1, dynamicSeverity1, minimumSignalSeverity2, dynamicSeverity2),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSECustomInsightExists(resourceName, &CustomInsight),
					testCheckCustomInsightDynamicSeverity(t, &CustomInsight,
						minimumSignalSeverity1, dynamicSeverity1, minimumSignalSeverity2, dynamicSeverity2),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testCreateCSECustomInsightConfigWithDynamicSeverity(minimumSignalSeverity1, dynamicSeverity1, updatedMinimumSignalSeverity2, updatedDynamicSeverity2),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSECustomInsightExists(resourceName, &CustomInsight),
					testCheckCustomInsightDynamicSeverity(t, &CustomInsight,
						minimumSignalSeverity1, dynamicSeverity1, updatedMinimumSignalSeverity2, updatedDynamicSeverity2),
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
	severity string, signalName1 string, signalName2 string, tag string) string {
	return fmt.Sprintf(`
resource "sumologic_cse_custom_insight" "custom_insight" {
	description = "%s"
    enabled = %t
    ordered = %t
    name = "%s"
    severity = "%s"
    signal_names = ["%s", "%s"]
    tags = ["%s"]
}
`, description, enabled, ordered, name, severity, signalName1,
		signalName2, tag)
}

func testCreateCSECustomInsightConfigWithDynamicSeverity(
	minimumSignalSeverity1 int, dynamicSeverity1 string, minimumSignalSeverity2 int, dynamicSeverity2 string) string {
	return fmt.Sprintf(`
resource "sumologic_cse_custom_insight" "custom_insight2" {
	description = "Dynamic severity insight"
    enabled = true
    ordered = true
    name = "Dynamic severity insight"
    severity = "LOW"
    dynamic_severity {
        minimum_signal_severity = "%d"
        insight_severity = "%s"
    }
    dynamic_severity {
        minimum_signal_severity = "%d"
        insight_severity = "%s"
    }
    tags = ["test tag"]
}
`, minimumSignalSeverity1, dynamicSeverity1, minimumSignalSeverity2, dynamicSeverity2)
}

func testCheckCSECustomInsightExists(n string, CustomInsight *CSECustomInsight) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CustomInsight ID is not set")
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
	enabled bool, ordered bool, name string, severity string, signalName1 string,
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
		if CustomInsight.SignalNames[0] != signalName1 {
			return fmt.Errorf("bad signalName1, expected \"%s\", got %#v", signalName1, CustomInsight.SignalNames[0])
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

func testCheckCustomInsightDynamicSeverity(t *testing.T, CustomInsight *CSECustomInsight,
	minimumSignalSeverity1 int, dynamicSeverity1 string, minimumSignalSeverity2 int, dynamicSeverity2 string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		assert.Equal(t, minimumSignalSeverity1, CustomInsight.DynamicSeverity[0].MinimumSignalSeverity)
		assert.Equal(t, dynamicSeverity1, CustomInsight.DynamicSeverity[0].InsightSeverity)
		assert.Equal(t, minimumSignalSeverity2, CustomInsight.DynamicSeverity[1].MinimumSignalSeverity)
		assert.Equal(t, dynamicSeverity2, CustomInsight.DynamicSeverity[1].InsightSeverity)
		return nil
	}
}
