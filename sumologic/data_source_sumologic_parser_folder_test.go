package sumologic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// A parser created without a parent_id is placed in the Parsers library root folder,
// so its parent_id gives us a real folder to look up without needing a pre-existing one.
func TestAccDataSourceParserFolder_byID(t *testing.T) {
	testName := "terraform_test_parser_" + acctest.RandString(16)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckParserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceParserFolderByID(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.sumologic_parser_folder.test", "id"),
					resource.TestCheckResourceAttrSet("data.sumologic_parser_folder.test", "name"),
					resource.TestCheckResourceAttrPair(
						"data.sumologic_parser_folder.test", "id",
						"sumologic_parser.test", "parent_id"),
				),
			},
		},
	})
}

func TestAccDataSourceParserFolder_idIsNotAFolder(t *testing.T) {
	testName := "terraform_test_parser_" + acctest.RandString(16)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckParserDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceParserFolderIDOfParser(testName),
				ExpectError: regexp.MustCompile(`is not a folder`),
			},
		},
	})
}

func TestAccDataSourceParserFolder_idDoesNotExist(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
data "sumologic_parser_folder" "test" {
  id = "0000000000000000"
}
`,
				ExpectError: regexp.MustCompile(`parser folder with id '0000000000000000' does not exist`),
			},
		},
	})
}

func TestAccDataSourceParserFolder_nameDoesNotExist(t *testing.T) {
	name := "terraform_test_missing_folder_" + acctest.RandString(16)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "sumologic_parser_folder" "test" {
  name = "%s"
}
`, name),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`parser folder with name '%s' does not exist`, name)),
			},
		},
	})
}

func TestAccDataSourceParserFolder_requiresIDOrName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
data "sumologic_parser_folder" "test" {
}
`,
				ExpectError: regexp.MustCompile(`one of .*id,name.* must be specified`),
			},
		},
	})
}

func testAccDataSourceParserFolderByID(name string) string {
	return fmt.Sprintf(`
resource "sumologic_parser" "test" {
  name    = "%s"
  stanzas = "[sourcetype:%s]\nFORMAT = CSV"
}

data "sumologic_parser_folder" "test" {
  id = sumologic_parser.test.parent_id
}
`, name, name)
}

func testAccDataSourceParserFolderIDOfParser(name string) string {
	return fmt.Sprintf(`
resource "sumologic_parser" "test" {
  name    = "%s"
  stanzas = "[sourcetype:%s]\nFORMAT = CSV"
}

data "sumologic_parser_folder" "test" {
  id = sumologic_parser.test.id
}
`, name, name)
}
