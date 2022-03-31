package sumologic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicLocalFileSource_create(t *testing.T) {
	var localFileSource LocalFileSource
	var localFileTraceSource LocalFileSource
	var kinesisLogSource LocalFileSource
	var collector Collector
	cName := acctest.RandomWithPrefix("tf-acc-test")
	cDescription := acctest.RandomWithPrefix("tf-acc-test")
	cCategory := acctest.RandomWithPrefix("tf-acc-test")
	sName := acctest.RandomWithPrefix("tf-acc-test")
	sDescription := acctest.RandomWithPrefix("tf-acc-test")
	sCategory := acctest.RandomWithPrefix("tf-acc-test")
	tName := acctest.RandomWithPrefix("tf-acc-test")
	tDescription := acctest.RandomWithPrefix("tf-acc-test")
	tCategory := acctest.RandomWithPrefix("tf-acc-test")
	kName := acctest.RandomWithPrefix("tf-acc-test")
	kDescription := acctest.RandomWithPrefix("tf-acc-test")
	kCategory := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "sumologic_localFile_source.localFile"
	tracingResourceName := "sumologic_localFile_source.traces"
	kinesisResourceName := "sumologic_localFile_source.kinesisLog"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLocalFileSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicLocalFileSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, tName, tDescription, tCategory, kName, kDescription, kCategory),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalFileSourceExists(resourceName, &localFileSource),
					testAccCheckLocalFileSourceValues(&localFileSource, sName, sDescription, sCategory),
					testAccCheckCollectorExists("sumologic_collector.test", &collector),
					testAccCheckCollectorValues(&collector, cName, cDescription, cCategory, "Etc/UTC", ""),
					testAccCheckLocalFileSourceExists(tracingResourceName, &localFileTraceSource),
					testAccCheckLocalFileSourceValues(&localFileTraceSource, tName, tDescription, tCategory),
					testAccCheckLocalFileSourceExists(kinesisResourceName, &kinesisLogSource),
					testAccCheckLocalFileSourceValues(&kinesisLogSource, kName, kDescription, kCategory),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
					resource.TestCheckResourceAttr(resourceName, "name", sName),
					resource.TestCheckResourceAttr(resourceName, "description", sDescription),
					resource.TestCheckResourceAttr(resourceName, "message_per_request", "false"),
					resource.TestCheckResourceAttr(resourceName, "category", sCategory),
					resource.TestCheckResourceAttrSet(tracingResourceName, "id"),
					resource.TestCheckResourceAttrSet(tracingResourceName, "url"),
					resource.TestCheckResourceAttr(tracingResourceName, "name", tName),
					resource.TestCheckResourceAttr(tracingResourceName, "description", tDescription),
					resource.TestCheckResourceAttr(tracingResourceName, "category", tCategory),
					resource.TestCheckResourceAttr(tracingResourceName, "content_type", "Zipkin"),
					resource.TestCheckResourceAttrSet(kinesisResourceName, "id"),
					resource.TestCheckResourceAttrSet(kinesisResourceName, "url"),
					resource.TestCheckResourceAttr(kinesisResourceName, "name", kName),
					resource.TestCheckResourceAttr(kinesisResourceName, "description", kDescription),
					resource.TestCheckResourceAttr(kinesisResourceName, "category", kCategory),
					resource.TestCheckResourceAttr(kinesisResourceName, "content_type", "KinesisLog"),
				),
			},
		},
	})
}

func TestAccSumologicLocalFileSource_update(t *testing.T) {
	var localFileSource LocalFileSource
	var localFileTraceSource LocalFileSource
	var kinesisLogSource LocalFileSource
	cName := acctest.RandomWithPrefix("tf-acc-test")
	cDescription := acctest.RandomWithPrefix("tf-acc-test")
	cCategory := acctest.RandomWithPrefix("tf-acc-test")
	sName := acctest.RandomWithPrefix("tf-acc-test")
	sDescription := acctest.RandomWithPrefix("tf-acc-test")
	sCategory := acctest.RandomWithPrefix("tf-acc-test")
	tName := acctest.RandomWithPrefix("tf-acc-test")
	tDescription := acctest.RandomWithPrefix("tf-acc-test")
	tCategory := acctest.RandomWithPrefix("tf-acc-test")
	kName := acctest.RandomWithPrefix("tf-acc-test")
	kDescription := acctest.RandomWithPrefix("tf-acc-test")
	kCategory := acctest.RandomWithPrefix("tf-acc-test")
	sNameUpdated := acctest.RandomWithPrefix("tf-acc-test")
	sDescriptionUpdated := acctest.RandomWithPrefix("tf-acc-test")
	sCategoryUpdated := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "sumologic_localFile_source.localFile"
	tracingResourceName := "sumologic_localFile_source.traces"
	kinesisResourceName := "sumologic_localFile_source.kinesisLog"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLocalFileSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicLocalFileSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, tName, tDescription, tCategory, kName, kDescription, kCategory),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalFileSourceExists(resourceName, &localFileSource),
					testAccCheckLocalFileSourceValues(&localFileSource, sName, sDescription, sCategory),
					testAccCheckLocalFileSourceExists(tracingResourceName, &localFileTraceSource),
					testAccCheckLocalFileSourceValues(&localFileTraceSource, tName, tDescription, tCategory),
					testAccCheckLocalFileSourceExists(kinesisResourceName, &kinesisLogSource),
					testAccCheckLocalFileSourceValues(&kinesisLogSource, kName, kDescription, kCategory),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
					resource.TestCheckResourceAttr(resourceName, "name", sName),
					resource.TestCheckResourceAttr(resourceName, "description", sDescription),
					resource.TestCheckResourceAttr(resourceName, "message_per_request", "false"),
					resource.TestCheckResourceAttr(resourceName, "category", sCategory),
					resource.TestCheckResourceAttrSet(tracingResourceName, "id"),
					resource.TestCheckResourceAttrSet(tracingResourceName, "url"),
					resource.TestCheckResourceAttr(tracingResourceName, "name", tName),
					resource.TestCheckResourceAttr(tracingResourceName, "description", tDescription),
					resource.TestCheckResourceAttr(tracingResourceName, "category", tCategory),
					resource.TestCheckResourceAttr(tracingResourceName, "content_type", "Zipkin"),
					resource.TestCheckResourceAttrSet(kinesisResourceName, "id"),
					resource.TestCheckResourceAttrSet(kinesisResourceName, "url"),
					resource.TestCheckResourceAttr(kinesisResourceName, "name", kName),
					resource.TestCheckResourceAttr(kinesisResourceName, "description", kDescription),
					resource.TestCheckResourceAttr(kinesisResourceName, "category", kCategory),
					resource.TestCheckResourceAttr(kinesisResourceName, "content_type", "KinesisLog"),
				),
			},
			{
				Config: testAccSumologicLocalFileSourceConfig(cName, cDescription, cCategory, sNameUpdated, sDescriptionUpdated, sCategoryUpdated, tName, tDescription, tCategory, kName, kDescription, kCategory),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocalFileSourceExists(resourceName, &localFileSource),
					testAccCheckLocalFileSourceValues(&localFileSource, sNameUpdated, sDescriptionUpdated, sCategoryUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "pathExpression"),
					resource.TestCheckResourceAttrSet(resourceName, "denylist"),
					resource.TestCheckResourceAttr(resourceName, "name", sNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", sDescriptionUpdated),
					resource.TestCheckResourceAttr(resourceName, "category", sCategoryUpdated),
					resource.TestCheckResourceAttr(tracingResourceName, "content_type", "Zipkin"),
					resource.TestCheckResourceAttr(kinesisResourceName, "content_type", "KinesisLog"),
				),
			},
		},
	})
}

func testAccCheckLocalFileSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_localFile_source" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("LocalFile Source destruction check: LocalFile Source ID is not set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}

		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}

		s, err := client.GetLocalFileSource(collectorID, id)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("LocalFile Source still exists")
		}
	}
	return nil
}

func testAccCheckLocalFileSourceExists(n string, localFileSource *LocalFileSource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("LocalFile Source ID is not set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("LocalFile Source id should be int; got %s", rs.Primary.ID)
		}

		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}

		c := testAccProvider.Meta().(*Client)
		localFileSourceResp, err := c.GetLocalFileSource(collectorID, id)
		if err != nil {
			return err
		}

		*localFileSource = *localFileSourceResp

		return nil
	}
}

func testAccCheckLocalFileSourceValues(localFileSource *LocalFileSource, name, description, category string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if localFileSource.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, localFileSource.Name)
		}
		if localFileSource.Description != description {
			return fmt.Errorf("bad description, expected \"%s\", got: %#v", description, localFileSource.Description)
		}
		if localFileSource.Category != category {
			return fmt.Errorf("bad category, expected \"%s\", got: %#v", category, localFileSource.Category)
		}
		return nil
	}
}

func testAccSumologicLocalFileSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, tName, tDescription, tCategory, kName, kDescription, kCategory string) string {
	return fmt.Sprintf(`
resource "sumologic_collector" "test" {
	name = "%s"
	description = "%s"
	category = "%s"
}
	
resource "sumologic_local_file_source" "localFile" {
	name = "%s"
	description = "%s"
	category = "%s"
	collector_id = "${sumologic_collector.test.id}"
}

resource "sumologic_local_file_source" "traces" {
	name = "%s"
	description = "%s"
	category = "%s"
	content_type = "Zipkin"
	collector_id = "${sumologic_collector.test.id}"
}

resource "sumologic_local_file_source" "kinesisLog" {
	name = "%s"
	description = "%s"
	category = "%s"
	content_type = "KinesisLog"
	collector_id = "${sumologic_collector.test.id}"
}
`, cName, cDescription, cCategory, sName, sDescription, sCategory, tName, tDescription, tCategory, kName, kDescription, kCategory)
}
