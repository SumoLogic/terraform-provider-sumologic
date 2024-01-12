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
					testCheckCustomInsightValues(t, &CustomInsight, description, enabled,
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
					testCheckCustomInsightValues(t, &CustomInsight, description, enabled,
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
        minimum_signal_severity = %d
        insight_severity = "%s"
    }
    dynamic_severity {
        minimum_signal_severity = %d
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

func testCheckCustomInsightValues(t *testing.T, CustomInsight *CSECustomInsight, description string,
	enabled bool, ordered bool, name string, severity string, signalName1 string,
	signalName2 string, tag string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		assert.Equal(t, description, CustomInsight.Description)
		assert.Equal(t, enabled, CustomInsight.Enabled)
		assert.Equal(t, ordered, CustomInsight.Ordered)
		assert.Equal(t, name, CustomInsight.Name)
		assert.Equal(t, severity, CustomInsight.Severity)
		assert.Equal(t, signalName1, CustomInsight.SignalNames[0])
		assert.Equal(t, signalName2, CustomInsight.SignalNames[1])
		assert.Equal(t, tag, CustomInsight.Tags[0])
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
