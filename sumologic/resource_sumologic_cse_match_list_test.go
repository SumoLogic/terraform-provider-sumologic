package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicSCEMatchList_createAndUpdate(t *testing.T) {
	SkipCseTest(t)

	var matchList CSEMatchListGet
	nDefaultTtl := 10800
	nDescription := "New Match List Description"
	nName := "Match List Name"
	nTargetColumn := "SrcIp"
	liDescription := "Match List Item Description"
	liValue := "value"
	liExpiration := "2122-02-27T04:00:00"
	uDefaultTtl := 3600
	uDescription := "Updated Match List Description"
	uliDescription := "Updated Match List item Description"
	resourceName := "sumologic_cse_match_list.match_list"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEMatchListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSEMatchListConfig(nDefaultTtl, nDescription, nName, nTargetColumn, liDescription, liExpiration, liValue),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEMatchListExists(resourceName, &matchList),
					testCheckMatchListValues(&matchList, nDefaultTtl, nDescription, nName, nTargetColumn),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testCreateCSEMatchListConfig(uDefaultTtl, uDescription, nName, nTargetColumn, uliDescription, liExpiration, liValue),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEMatchListExists(resourceName, &matchList),
					testCheckMatchListValues(&matchList, uDefaultTtl, uDescription, nName, nTargetColumn),
				),
			},
			{
				Config: testDeleteCSEMatchListItemConfig(uDefaultTtl, uDescription, nName, nTargetColumn),
				Check: resource.ComposeTestCheckFunc(
					testCheckMatchListItemsEmpty(resourceName),
				),
			},
		},
	})
}

func testAccCSEMatchListDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_match_list" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Match List destruction check: CSE Match List ID is not set")
		}

		s, err := client.GetCSEMatchList(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("match List still exists")
		}
	}
	return nil
}

func testCreateCSEMatchListConfig(nDefaultTtl int, nDescription string, nName string, nTargetColumn string, liDescription string, liExpiration string, liValue string) string {
	return fmt.Sprintf(`
resource "sumologic_cse_match_list" "match_list" {
	default_ttl = "%d"
	description = "%s"
	name = "%s"
	target_column = "%s"
	items {
		description = "%s"
		expiration = "%s"
		value = "%s"
	}
}
`, nDefaultTtl, nDescription, nName, nTargetColumn, liDescription, liExpiration, liValue)
}

func testCreateCSEMatchListConfigWithOver1000Items(nDefaultTtl int, nDescription string, nName string, nTargetColumn string, liDescription string, liExpiration string, liValue string) string {
	var str = fmt.Sprintf(`
resource "sumologic_cse_match_list" "match_list" {
    default_ttl = "%d"
    description = "%s"
    name = "%s"
    target_column = "%s"
    items {
`, nDefaultTtl, nDescription, nName, nTargetColumn)

	for i := 0; i < 1001; i++ {
		str += fmt.Sprintf(
			`[
            description = "%s"
            expiration = "%s"
            value = "%s"
        ],`, liDescription, liExpiration, liValue)
	}
	str += "}"
	return str
}

func testDeleteCSEMatchListItemConfig(nDefaultTtl int, nDescription string, nName string, nTargetColumn string) string {
	return fmt.Sprintf(`
resource "sumologic_cse_match_list" "match_list" {
	default_ttl = "%d"
	description = "%s"
	name = "%s"
	target_column = "%s"
}
`, nDefaultTtl, nDescription, nName, nTargetColumn)
}

func testCheckCSEMatchListExists(n string, matchList *CSEMatchListGet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("match List ID is not set")
		}
		c := testAccProvider.Meta().(*Client)
		matchListResp, err := c.GetCSEMatchList(rs.Primary.ID)
		if err != nil {
			return err
		}

		*matchList = *matchListResp

		return nil
	}
}

func testCheckMatchListItemsEmpty(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("match List ID is not set")
		}

		c := testAccProvider.Meta().(*Client)
		matchListResp, err := c.GetCSEMatchListItemsInMatchList(rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(matchListResp.CSEMatchListItemsGetObjects) != 0 {
			return fmt.Errorf("match list items not empty")
		}

		return nil
	}
}

func testCheckMatchListValues(matchList *CSEMatchListGet, nDefaultTtl int, nDescription string, nName string, nTargetColumn string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if matchList.DefaultTtl != nDefaultTtl {
			return fmt.Errorf("bad default ttl, expected \"%s\", got: %#v", nName, matchList.Name)
		}
		if matchList.Description != nDescription {
			return fmt.Errorf("bad description, expected \"%s\", got: %#v", nDescription, matchList.Description)
		}
		if matchList.Name != nName {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", nName, matchList.Name)
		}
		if matchList.TargetColumn != nTargetColumn {
			return fmt.Errorf("bad target column, expected \"%s\", got: %#v", nName, matchList.Name)
		}

		return nil
	}
}
