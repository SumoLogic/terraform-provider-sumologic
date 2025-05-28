package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicSCEInsightsResolution_create(t *testing.T) {
	SkipCseTest(t)

	var insightResolution CSEInsightsResolutionGet
	nName := "New Insights Resolution"
	nDescription := "New Insights Resolution Description"
	nParentId := "Resolved"
	resourceName := "sumologic_cse_insights_resolution.insights_resolution"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEInsightsResolutionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSEInsightsResolutionConfig(nName, nDescription, nParentId),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEInsightsResolutionExists(resourceName, &insightResolution),
					testCheckInsightsResolutionValues(&insightResolution, nName, nDescription, nParentId),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccCSEInsightsResolutionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_insights_resolution" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Insights Resolution destruction check: CSE Insights Resolution ID is not set")
		}

		id := rs.Primary.ID

		s, err := client.GetCSEInsightsResolution(id)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("insight Resolution still exists")
		}
	}
	return nil
}

func testCreateCSEInsightsResolutionConfig(nName string, nDescription string, nParentId string) string {
	return fmt.Sprintf(`
resource "sumologic_cse_insights_resolution" "insights_resolution" {
	name = "%s"
	description = "%s"
	parent = "%s"
}
`, nName, nDescription, nParentId)
}

func testCheckCSEInsightsResolutionExists(n string, insightResolution *CSEInsightsResolutionGet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("insight Resolution ID is not set")
		}

		id := rs.Primary.ID

		c := testAccProvider.Meta().(*Client)
		insightResolutionResp, err := c.GetCSEInsightsResolution(id)
		if err != nil {
			return err
		}

		*insightResolution = *insightResolutionResp

		return nil
	}
}

func testCheckInsightsResolutionValues(insightResolution *CSEInsightsResolutionGet, nName string, nDescription string, nparentId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if insightResolution.Name != nName {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", nName, insightResolution.Name)
		}
		if insightResolution.Description != nDescription {
			return fmt.Errorf("bad description, expected \"%s\", got: %#v", nDescription, insightResolution.Description)
		}
		if insightResolution.Parent.Name != nparentId {
			return fmt.Errorf("bad parent id, expected \"%s\", got: %#v", nparentId, insightResolution.Parent.Name)
		}
		return nil
	}
}
