package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicParser_crud(t *testing.T) {
	var parser ParsersLibraryParser
	testName := "terraform_test_parser_" + acctest.RandString(16)
	resourceName := "sumologic_parser.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckParserDestroy,
		Steps: []resource.TestStep{
			{
				// create
				Config: testAccSumologicParser(testName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParserExists(resourceName, &parser),
					resource.TestCheckResourceAttr(resourceName, "name", testName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "stanzas", "[sourcetype:"+testName+"]\nFORMAT = CSV"),
					resource.TestCheckResourceAttr(resourceName, "is_partial", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "parent_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttrSet(resourceName, "content_type"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
				),
			},
			{
				// update
				Config: testAccSumologicParserUpdate(testName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParserExists(resourceName, &parser),
					resource.TestCheckResourceAttr(resourceName, "name", testName),
					resource.TestCheckResourceAttr(resourceName, "description", "edited"),
					resource.TestCheckResourceAttr(resourceName, "stanzas", "[sourcetype:"+testName+"]\nFORMAT = JSON"),
					resource.TestCheckResourceAttr(resourceName, "local_stanzas", "[sourcetype:"+testName+"]\nWRAPPER = BSD_SYSLOG"),
					resource.TestCheckResourceAttr(resourceName, "families", "firewalls"),
					resource.TestCheckResourceAttr(resourceName, "is_partial", "true"),
				),
			},
			{
				// import
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSumologicParser_disappears(t *testing.T) {
	var parser ParsersLibraryParser
	testName := "terraform_test_parser_" + acctest.RandString(16)
	resourceName := "sumologic_parser.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckParserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicParser(testName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParserExists(resourceName, &parser),
					// Delete the parser out-of-band; the next plan must see it as gone and recreate it.
					func(s *terraform.State) error {
						client := testAccProvider.Meta().(*Client)
						return client.DeleteParsersLibraryParser(parser.ID)
					},
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckParserExists(name string, parser *ParsersLibraryParser) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("parser not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("parser ID is not set")
		}

		client := testAccProvider.Meta().(*Client)
		found, err := client.GetParsersLibraryParser(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error reading parser %s: %w", rs.Primary.ID, err)
		}
		if found == nil {
			return fmt.Errorf("parser %s does not exist", rs.Primary.ID)
		}
		*parser = *found
		return nil
	}
}

func testAccCheckParserDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "sumologic_parser" {
			continue
		}
		found, err := client.GetParsersLibraryParser(r.Primary.ID)
		if err != nil {
			return fmt.Errorf("encountered an error: %w", err)
		}
		if found != nil {
			return fmt.Errorf("parser %s still exists", r.Primary.ID)
		}
	}
	return nil
}

func testAccSumologicParser(name string) string {
	return fmt.Sprintf(`
resource "sumologic_parser" "test" {
  name    = "%s"
  stanzas = "[sourcetype:%s]\nFORMAT = CSV"
}
`, name, name)
}

func testAccSumologicParserUpdate(name string) string {
	return fmt.Sprintf(`
resource "sumologic_parser" "test" {
  name          = "%s"
  description   = "edited"
  stanzas       = "[sourcetype:%s]\nFORMAT = JSON"
  local_stanzas = "[sourcetype:%s]\nWRAPPER = BSD_SYSLOG"
  families      = "firewalls"
  is_partial    = true
}
`, name, name, name)
}
