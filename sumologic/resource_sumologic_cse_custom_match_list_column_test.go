package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicSCECustomMatchListColumn_create(t *testing.T) {
	SkipCseTest(t)

	var CustomMatchListColumn CSECustomMatchListColumn
	nName := "Custom Match List Column Test Terraform"
	nField := "device_ip"
	resourceName := "sumologic_cse_custom_match_list_column.custom_match_list_column"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSECustomMatchListColumnDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSECustomMatchListColumnConfig(nName, nField),
				Check: resource.ComposeTestCheckFunc(
					testCheckCustomMatchListColumnExists(resourceName, &CustomMatchListColumn),
					testCheckCustomMatchListColumnValues(&CustomMatchListColumn, nName, nField),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccCSECustomMatchListColumnDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_custom_match_list_column" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Custom Match List Column destruction check: CSE Custom Match List Column ID is not set")
		}

		CustomMatchListColumnID := rs.Primary.Attributes["id"]

		s, err := client.GetCSECustomMatchListColumn(CustomMatchListColumnID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("Custom Match List Column still exists")
		}
	}
	return nil
}

func testCreateCSECustomMatchListColumnConfig(nName string, nField string) string {
	return fmt.Sprintf(`
resource "sumologic_cse_custom_match_list_column" "custom_match_list_column" {
	name = "%s"
    fields = ["%s"]
}
`, nName, nField)
}

func testCheckCustomMatchListColumnExists(n string, CustomMatchListColumn *CSECustomMatchListColumn) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Custom Match List Column ID is not set")
		}

		CustomMatchListColumnID := rs.Primary.Attributes["id"]

		c := testAccProvider.Meta().(*Client)
		CustomMatchListColumnResp, err := c.GetCSECustomMatchListColumn(CustomMatchListColumnID)
		if err != nil {
			return err
		}

		*CustomMatchListColumn = *CustomMatchListColumnResp

		return nil
	}
}

func testCheckCustomMatchListColumnValues(CustomMatchListColumn *CSECustomMatchListColumn, nName string, nField string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if CustomMatchListColumn.Name != nName {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", nName, CustomMatchListColumn.Name)
		}
		if CustomMatchListColumn.Fields[0] != nField {
			return fmt.Errorf("bad field, expected \"%s\", got %#v", nField, CustomMatchListColumn.Fields[0])
		}
		return nil
	}
}
