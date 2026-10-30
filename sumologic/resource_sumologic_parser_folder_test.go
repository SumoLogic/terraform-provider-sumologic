package sumologic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicParserFolder_crud(t *testing.T) {
	var folder ParsersLibraryFolder
	testName := "terraform_test_parser_folder_" + acctest.RandString(16)
	resourceName := "sumologic_parser_folder.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckParserFolderDestroy,
		Steps: []resource.TestStep{
			{
				// create
				Config: testAccSumologicParserFolder(testName, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParserFolderExists(resourceName, &folder),
					resource.TestCheckResourceAttr(resourceName, "name", testName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "parent_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttrSet(resourceName, "content_type"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
				),
			},
			{
				// update
				Config: testAccSumologicParserFolder(testName+"_renamed", "edited"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParserFolderExists(resourceName, &folder),
					resource.TestCheckResourceAttr(resourceName, "name", testName+"_renamed"),
					resource.TestCheckResourceAttr(resourceName, "description", "edited"),
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

func TestAccSumologicParserFolder_move(t *testing.T) {
	var folder ParsersLibraryFolder
	testName := "terraform_test_parser_folder_" + acctest.RandString(16)
	childName := "sumologic_parser_folder.child"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckParserFolderDestroy,
		Steps: []resource.TestStep{
			{
				// child under folder a
				Config: testAccSumologicParserFolderNested(testName, "a"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParserFolderExists(childName, &folder),
					resource.TestCheckResourceAttrPair(childName, "parent_id", "sumologic_parser_folder.a", "id"),
				),
			},
			{
				// move child under folder b
				Config: testAccSumologicParserFolderNested(testName, "b"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParserFolderExists(childName, &folder),
					resource.TestCheckResourceAttrPair(childName, "parent_id", "sumologic_parser_folder.b", "id"),
				),
			},
		},
	})
}

func TestAccSumologicParserFolder_withParser(t *testing.T) {
	testName := "terraform_test_parser_folder_" + acctest.RandString(16)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckParserDestroy,
			testAccCheckParserFolderDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicParserInFolder(testName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"sumologic_parser.test", "parent_id",
						"sumologic_parser_folder.test", "id"),
				),
			},
		},
	})
}

func TestAccSumologicParserFolder_disappears(t *testing.T) {
	var folder ParsersLibraryFolder
	testName := "terraform_test_parser_folder_" + acctest.RandString(16)
	resourceName := "sumologic_parser_folder.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckParserFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicParserFolder(testName, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParserFolderExists(resourceName, &folder),
					// Delete the folder out-of-band; the next plan must see it as gone and recreate it.
					func(s *terraform.State) error {
						client := testAccProvider.Meta().(*Client)
						return client.DeleteParsersLibraryFolder(folder.ID)
					},
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Deleting a folder deletes everything under it, so the resource must refuse to delete a folder
// that has children it does not manage.
func TestAccSumologicParserFolder_deleteWithChildren(t *testing.T) {
	var folder ParsersLibraryFolder
	var childID string
	testName := "terraform_test_parser_folder_" + acctest.RandString(16)
	resourceName := "sumologic_parser_folder.test"

	deleteChild := func() {
		if childID == "" {
			return
		}
		client := testAccProvider.Meta().(*Client)
		_ = client.DeleteParsersLibraryParser(childID)
		childID = ""
	}
	t.Cleanup(deleteChild)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckParserFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicParserFolder(testName, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParserFolderExists(resourceName, &folder),
					// Add a parser to the folder outside of Terraform.
					func(s *terraform.State) error {
						client := testAccProvider.Meta().(*Client)
						id, err := client.CreateParsersLibraryParser(ParsersLibraryParser{
							Name:     testName + "_child",
							Stanzas:  "[sourcetype:" + testName + "_child]\nFORMAT = CSV",
							ParentID: folder.ID,
						})
						if err != nil {
							return err
						}
						childID = id
						return nil
					},
				),
			},
			{
				// Destroy must be refused while the folder has a child.
				Config:      testAccSumologicParserFolder(testName, ""),
				Destroy:     true,
				ExpectError: regexp.MustCompile("is not empty"),
			},
			{
				// Once the child is gone, the folder can be destroyed.
				PreConfig: deleteChild,
				Config:    testAccSumologicParserFolder(testName, ""),
				Destroy:   true,
			},
		},
	})
}

func testAccCheckParserFolderExists(name string, folder *ParsersLibraryFolder) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("parser folder not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("parser folder ID is not set")
		}

		client := testAccProvider.Meta().(*Client)
		found, err := client.GetParsersLibraryFolder(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error reading parser folder %s: %w", rs.Primary.ID, err)
		}
		if found == nil {
			return fmt.Errorf("parser folder %s does not exist", rs.Primary.ID)
		}
		*folder = *found
		return nil
	}
}

func testAccCheckParserFolderDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "sumologic_parser_folder" {
			continue
		}
		found, err := client.GetParsersLibraryFolder(r.Primary.ID)
		if err != nil {
			return fmt.Errorf("encountered an error: %w", err)
		}
		if found != nil {
			return fmt.Errorf("parser folder %s still exists", r.Primary.ID)
		}
	}
	return nil
}

func testAccSumologicParserFolder(name, description string) string {
	return fmt.Sprintf(`
resource "sumologic_parser_folder" "test" {
  name        = "%s"
  description = "%s"
}
`, name, description)
}

// parent is "a" or "b": which folder the child is placed under.
func testAccSumologicParserFolderNested(name, parent string) string {
	return fmt.Sprintf(`
resource "sumologic_parser_folder" "a" {
  name = "%s_a"
}

resource "sumologic_parser_folder" "b" {
  name = "%s_b"
}

resource "sumologic_parser_folder" "child" {
  name      = "%s_child"
  parent_id = sumologic_parser_folder.%s.id
}
`, name, name, name, parent)
}

func testAccSumologicParserInFolder(name string) string {
	return fmt.Sprintf(`
resource "sumologic_parser_folder" "test" {
  name = "%s"
}

resource "sumologic_parser" "test" {
  name      = "%s_parser"
  parent_id = sumologic_parser_folder.test.id
  stanzas   = "[sourcetype:%s]\nFORMAT = CSV"
}
`, name, name, name)
}
