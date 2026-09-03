package sumologic

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicLookupTable_basic(t *testing.T) {
	var lookupTable LookupTable
	testName := "SampleLookupTable"
	testFieldName := "FieldName1"
	testFieldType := "boolean"
	testTtl := 100
	testPrimaryKeys := "FieldName1"
	testSizeLimitAction := "StopIncomingMessages"
	testDescription := "This is a sample lookup table description."

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLookupTableDestroy(lookupTable),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSumologicLookupTableConfigImported(testName, testFieldName, testFieldType, testTtl, testPrimaryKeys, testSizeLimitAction, testDescription),
			},
			{
				ResourceName:      "sumologic_lookup_table.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSumologicLookupTable_create(t *testing.T) {
	var lookupTable LookupTable
	testName := "SampleLookupTable"
	testFieldName := "FieldName1"
	testFieldType := "boolean"
	testTtl := 100
	testPrimaryKeys := "FieldName1"
	testSizeLimitAction := "StopIncomingMessages"
	testDescription := "This is a sample lookup table description."
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLookupTableDestroy(lookupTable),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicLookupTable(testName, testFieldName, testFieldType, testTtl, testPrimaryKeys, testSizeLimitAction, testDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLookupTableExists("sumologic_lookup_table.test", &lookupTable, t),
					testAccCheckLookupTableAttributes("sumologic_lookup_table.test"),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "fields.#", "1"),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "fields.0.field_name", testFieldName),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "fields.0.field_type", testFieldType),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "ttl", strconv.Itoa(testTtl)),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "primary_keys.0", testPrimaryKeys),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "size_limit_action", testSizeLimitAction),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "description", testDescription),
				),
			},
		},
	})
}

func TestAccSumologicLookupTable_update(t *testing.T) {
	var lookupTable LookupTable
	testName := "SampleLookupTable"
	testFieldName := "FieldName1"
	testFieldType := "boolean"
	testTtl := 100
	testPrimaryKeys := "FieldName1"
	testSizeLimitAction := "StopIncomingMessages"
	testDescription := "This is a sample lookup table description."

	testUpdatedTtl := 101
	testUpdatedSizeLimitAction := "DeleteOldData"
	testUpdatedDescription := "This is a sample lookup table description Updated"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLookupTableDestroy(lookupTable),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicLookupTable(testName, testFieldName, testFieldType, testTtl, testPrimaryKeys, testSizeLimitAction, testDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLookupTableExists("sumologic_lookup_table.test", &lookupTable, t),
					testAccCheckLookupTableAttributes("sumologic_lookup_table.test"),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "fields.#", "1"),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "fields.0.field_name", testFieldName),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "fields.0.field_type", testFieldType),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "ttl", strconv.Itoa(testTtl)),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "primary_keys.0", testPrimaryKeys),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "size_limit_action", testSizeLimitAction),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "description", testDescription),
				),
			},
			{
				Config: testAccSumologicLookupTableUpdate(testName, testFieldName, testFieldType, testUpdatedTtl, testPrimaryKeys, testUpdatedSizeLimitAction, testUpdatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "fields.#", "1"),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "fields.0.field_name", testFieldName),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "fields.0.field_type", testFieldType),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "ttl", strconv.Itoa(testUpdatedTtl)),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "primary_keys.0", testPrimaryKeys),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "size_limit_action", testUpdatedSizeLimitAction),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "description", testUpdatedDescription),
				),
			},
		},
	})
}

func TestAccSumologicLookupTable_content(t *testing.T) {
	var lookupTable LookupTable
	testName := "SampleLookupTable"
	testFieldName := "FieldName1"
	testFieldType := "string"
	testTtl := 100
	testPrimaryKeys := "FieldName1"
	testSizeLimitAction := "StopIncomingMessages"
	testDescription := "This is a sample lookup table description."

	// Exercise the real-world path: content read from disk via the HCL
	// file() function, not an inline string, since that's how the feature
	// is actually meant to be used (see website/docs/r/lookup_table.html.markdown).
	contentPath, err := filepath.Abs("testdata/lookup_table_content.csv")
	if err != nil {
		t.Fatal(err)
	}
	updatedContentPath, err := filepath.Abs("testdata/lookup_table_content_updated.csv")
	if err != nil {
		t.Fatal(err)
	}

	contentBytes, err := os.ReadFile(contentPath)
	if err != nil {
		t.Fatal(err)
	}
	updatedContentBytes, err := os.ReadFile(updatedContentPath)
	if err != nil {
		t.Fatal(err)
	}

	// content cannot be read back from the API, so state stores only its
	// hash - assert against that rather than the raw CSV. Hash what's
	// actually on disk so the expectation can't drift from the fixture.
	testContentHash := hashLookupContent(string(contentBytes))
	testUpdatedContentHash := hashLookupContent(string(updatedContentBytes))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLookupTableDestroy(lookupTable),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicLookupTableContent(testName, testFieldName, testFieldType, testTtl, testPrimaryKeys, testSizeLimitAction, testDescription, fmt.Sprintf("file(%q)", contentPath)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLookupTableExists("sumologic_lookup_table.test", &lookupTable, t),
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "content", testContentHash),
				),
			},
			{
				Config: testAccSumologicLookupTableContent(testName, testFieldName, testFieldType, testTtl, testPrimaryKeys, testSizeLimitAction, testDescription, fmt.Sprintf("file(%q)", updatedContentPath)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "content", testUpdatedContentHash),
				),
			},
			{
				// Removing the argument entirely must truncate the table
				// rather than leaving it unmanaged.
				Config: testAccSumologicLookupTable(testName, testFieldName, testFieldType, testTtl, testPrimaryKeys, testSizeLimitAction, testDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_lookup_table.test", "content", ""),
				),
			},
		},
	})
}

// testAccSumologicLookupTableContent embeds contentExpr verbatim as the
// "content" argument's value - callers pass a full HCL expression (e.g.
// `file("/abs/path.csv")`), not a string to be quoted.
func testAccSumologicLookupTableContent(name string, testFieldName string, testFieldType string, ttl int, primaryKeys string, sizeLimitAction string, description string, contentExpr string) string {
	return fmt.Sprintf(`
data "sumologic_personal_folder" "personalFolder" {}
resource "sumologic_lookup_table" "test" {
    name = "%s"
    fields {
      field_name = "%s"
      field_type = "%s"
    }
    ttl = %d
    primary_keys = ["%s"]
    parent_folder_id = "${data.sumologic_personal_folder.personalFolder.id}"
    size_limit_action = "%s"
    description = "%s"
    content = %s
}
`, name, testFieldName, testFieldType, ttl, primaryKeys, sizeLimitAction, description, contentExpr)
}

func testAccCheckLookupTableDestroy(lookupTable LookupTable) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		_, err := client.GetLookupTable(lookupTable.ID)
		if err == nil {
			return fmt.Errorf("Lookup Table still exists")
		}
		return nil
	}
}

func testAccCheckLookupTableExists(name string, lookupTable *LookupTable, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			//need this so that we don't get an unused import error for strconv in some cases
			return fmt.Errorf("Error = %s. LookupTable not found: %s", strconv.FormatBool(ok), name)
		}

		//need this so that we don't get an unused import error for strings in some cases
		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("LookupTable ID is not set")
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)
		newLookupTable, err := client.GetLookupTable(id)
		if err != nil {
			return fmt.Errorf("LookupTable %s not found", id)
		}
		lookupTable = newLookupTable
		return nil
	}
}

func testAccCheckSumologicLookupTableConfigImported(name string, testFieldName string, testFieldType string, ttl int, primaryKeys string, sizeLimitAction string, description string) string {
	return fmt.Sprintf(`
data "sumologic_personal_folder" "personalFolder" {}
resource "sumologic_lookup_table" "foo" {
      name = "%s"
      fields {
        field_name = "%s"
        field_type = "%s"
      }
      ttl = %d
      primary_keys = ["%s"]
      parent_folder_id = "${data.sumologic_personal_folder.personalFolder.id}"
      size_limit_action = "%s"
      description = "%s"
}
`, name, testFieldName, testFieldType, ttl, primaryKeys, sizeLimitAction, description)
}

func testAccSumologicLookupTable(name string, testFieldName string, testFieldType string, ttl int, primaryKeys string, sizeLimitAction string, description string) string {
	return fmt.Sprintf(`
data "sumologic_personal_folder" "personalFolder" {}
resource "sumologic_lookup_table" "test" {
    name = "%s"
    fields {
      field_name = "%s"
      field_type = "%s"
    }
    ttl = %d
    primary_keys = ["%s"]
    parent_folder_id = "${data.sumologic_personal_folder.personalFolder.id}"
    size_limit_action = "%s"
    description = "%s"
}
`, name, testFieldName, testFieldType, ttl, primaryKeys, sizeLimitAction, description)
}

func testAccSumologicLookupTableUpdate(name string, testFieldName string, testFieldType string, ttl int, primaryKeys string, sizeLimitAction string, description string) string {
	return fmt.Sprintf(`
data "sumologic_personal_folder" "personalFolder" {}
resource "sumologic_lookup_table" "test" {
      name = "%s"
      fields {
        field_name = "%s"
        field_type = "%s"
      }
      ttl = %d
      primary_keys = ["%s"]
      parent_folder_id = "${data.sumologic_personal_folder.personalFolder.id}"
      size_limit_action = "%s"
      description = "%s"
}
`, name, testFieldName, testFieldType, ttl, primaryKeys, sizeLimitAction, description)
}

func testAccCheckLookupTableAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(name, "name"),
			resource.TestCheckResourceAttrSet(name, "ttl"),
			resource.TestCheckResourceAttrSet(name, "parent_folder_id"),
			resource.TestCheckResourceAttrSet(name, "size_limit_action"),
			resource.TestCheckResourceAttrSet(name, "description"),
		)
		return f(s)
	}
}
