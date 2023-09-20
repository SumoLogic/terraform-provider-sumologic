package sumologic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicSCEMatchList_createAndUpdate(t *testing.T) {
	SkipCseTest(t)

	var matchList CSEMatchListGet
	resourceName := "sumologic_cse_match_list.match_list"

	// Create values
	nName := fmt.Sprintf("Terraform Test Match List %s", uuid.New())
	nDefaultTtl := 10800
	nDescription := "Match List Description"
	nTargetColumn := "SrcIp"
	liDescription := "Match List Item Description"
	liExpiration := "2122-02-27T04:00:00"
	liValue := "value"
	liCount := 1

	// Update values
	uDefaultTtl := 3600
	uDescription := "Updated Match List Description"
	uliDescription := "Updated Match List item Description"
	uliExpiration := "2122-02-27T05:00:00"
	uliValue := "updated value"
	uliCount := 3

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSEMatchListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSEMatchListConfig(nDefaultTtl, nDescription, nName, nTargetColumn, liDescription, liExpiration, liValue, liCount),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEMatchListExists(resourceName, &matchList),
					testCheckMatchListValues(&matchList, nDefaultTtl, nDescription, nName, nTargetColumn),
					testCheckMatchListItemsValuesAndCount(resourceName, liDescription, liExpiration, liValue, liCount),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testCreateCSEMatchListConfig(uDefaultTtl, uDescription, nName, nTargetColumn, uliDescription, uliExpiration, uliValue, uliCount),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSEMatchListExists(resourceName, &matchList),
					testCheckMatchListValues(&matchList, uDefaultTtl, uDescription, nName, nTargetColumn),
					testCheckMatchListItemsValuesAndCount(resourceName, uliDescription, uliExpiration, uliValue, uliCount),
				),
			},
			{
				Config: testDeleteCSEMatchListItemConfig(uDefaultTtl, uDescription, nName, nTargetColumn),
				Check: resource.ComposeTestCheckFunc(
					testCheckMatchListItemsValuesAndCount(resourceName, "", "", "", 0),
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

func testCreateCSEMatchListConfig(nDefaultTtl int, nDescription string, nName string, nTargetColumn string, liDescription string, liExpiration string, liValue string, numItems int) string {
	var itemsStr = ""

	for i := 0; i < numItems; i++ {
		id := uuid.New()

		itemsStr += fmt.Sprintf(`
    items {
	description = "%s %d %s"
	expiration = "%s"
	value = "%s %d %s"
    }`, liDescription, i, id, liExpiration, liValue, i, id)
	}

	var str = fmt.Sprintf(`
resource "sumologic_cse_match_list" "match_list" {
    default_ttl = "%d"
    description = "%s"
    name = "%s"
    target_column = "%s" %s
}`, nDefaultTtl, nDescription, nName, nTargetColumn, itemsStr)

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

func testCheckMatchListItemsValuesAndCount(resourceName string, expectedDescription string, expectedExpiration string, expectedValue string, expectedCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("expected match list ID to be non-empty, but found empty string instead")
		}

		c := testAccProvider.Meta().(*Client)
		matchListResp, err := c.GetCSEMatchListItemsInMatchList(rs.Primary.ID)
		if err != nil {
			return err
		}

		actualCount := len(matchListResp.CSEMatchListItemsGetObjects)
		if actualCount != expectedCount {
			return fmt.Errorf("expected %d match list items, but found %d instead", expectedCount, actualCount)
		}

		if expectedCount == 0 {
			return nil
		}

		for _, item := range matchListResp.CSEMatchListItemsGetObjects {
			if item.ID == "" {
				return fmt.Errorf("expected match list item ID to be non-empty, but found empty string instead")
			}
			if !strings.Contains(item.Meta.Description, expectedDescription) {
				return fmt.Errorf("expected match list item description to contain \"%s\", but found \"%s\" instead", expectedDescription, item.Meta.Description)
			}
			if item.Expiration != expectedExpiration {
				return fmt.Errorf("expected expiration to be \"%s\", but found \"%s\" instead", expectedExpiration, item.Expiration)
			}
			if !strings.Contains(item.Value, expectedValue) {
				return fmt.Errorf("expected match list item value to contain \"%s\", but found \"%s\" instead", expectedValue, item.Value)
			}
		}
		return nil
	}
}
