package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicSCEInsightsStatus_create(t *testing.T) {
	SkipCseTest(t)

	var insightStatus CSEInsightsStatusGet
	nName := "New Test Status"
	nDescription := "New Test Status Description"
	resourceName := "sumologic_cse_insights_status.insights_status"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEInsightsStatusDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSEInsightsStatusConfig(nName, nDescription),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEInsightsStatusExists(resourceName, &insightStatus),
					testCheckInsightsStatusValues(&insightStatus, nName, nDescription),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccCSEInsightsStatusDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_insights_status" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Insights Status destruction check: CSE Insights Status ID is not set")
		}

		s, err := client.GetCSEInsightsStatus(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("insight Status still exists")
		}
	}
	return nil
}

func testCreateCSEInsightsStatusConfig(nName string, nDescription string) string {
	return fmt.Sprintf(`
resource "sumologic_cse_insights_status" "insights_status" {
	name = "%s"
	description = "%s"
}
`, nName, nDescription)
}

func testCheckCSEInsightsStatusExists(n string, insightStatus *CSEInsightsStatusGet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("insight Status ID is not set")
		}

		c := testAccProvider.Meta().(*Client)
		insightStatusResp, err := c.GetCSEInsightsStatus(rs.Primary.ID)
		if err != nil {
			return err
		}

		*insightStatus = *insightStatusResp

		return nil
	}
}

func testCheckInsightsStatusValues(insightStatus *CSEInsightsStatusGet, nName string, nDescription string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if insightStatus.Name != nName {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", nName, insightStatus.Name)
		}
		if insightStatus.Description != nDescription {
			return fmt.Errorf("bad description, expected \"%s\", got: %#v", nDescription, insightStatus.Description)
		}

		return nil
	}
}
