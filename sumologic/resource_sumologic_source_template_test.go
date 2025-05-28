package sumologic

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strconv"
	"strings"
	"testing"
)

func TestAccSumologicSourceTemplate_basic(t *testing.T) {
	var sourceTemplate SourceTemplate
	testSchemaRef := "type =     \"Mac\""

	testSelector :=
		"tags =  [\n[\n{\n key = \"tag\"\n values= [\"Value\"]\n}\n]\n] \n names = [\"TestCollector1\"]"

	testConfig := "apache.yaml.example"
	testInputJson := "jsonencode({\n\"name\": \"hostmetrics_test_source_template_acc\",\n\"description\": \"Host metric source\" ,\n\"receivers\": {\n\"hostmetrics\": {\n\"receiverType\": \"hostmetrics\",\n\"collection_interval\": \"5m\",\n\"cpu_scraper_enabled\": true,\n\"disk_scraper_enabled\": true,\n\"load_scraper_enabled\": true,\n\"filesystem_scraper_enabled\": true,\n\"memory_scraper_enabled\": true,\n\"network_scraper_enabled\": true,\n\"process_scraper_enabled\": true,\n\"paging_scraper_enabled\": true\n}\n},\n\"processors\": {\n\"resource\": {\n\"processorType\": \"resource\",\n\"user_attributes\": [\n{\n\"key\": \"_sourceCategory\",\n\"value\": \"otel/host\"\n}\n],\n\"default_attributes\": [\n{\n\"key\": \"sumo.datasource\",\n\"value\": \"apache\"\n},\n]\n}\n}\n})"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSourceTemplateDestroy(sourceTemplate),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSumologicSourceTemplateConfigImported(testSchemaRef, testSelector, testConfig, testInputJson),
			},
			{
				ResourceName:      "sumologic_source_template.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSumologicSourceTemplate_create(t *testing.T) {
	var sourceTemplate SourceTemplate
	testSchemaRef := "type =     \"Mac\""
	testSelector := "tags =  [\n[\n{\n key = \"tag\"\n values= [\"Value\"]\n}\n]\n]\n names = [\"TestCollector1\"]"

	testInputJson := "jsonencode({\n\"name\": \"hostmetrics_test_source_template_acc\",\n\"description\": \"Host metric source\" ,\n\"receivers\": {\n\"hostmetrics\": {\n\"receiverType\": \"hostmetrics\",\n\"collection_interval\": \"5m\",\n\"cpu_scraper_enabled\": true,\n\"disk_scraper_enabled\": true,\n\"load_scraper_enabled\": true,\n\"filesystem_scraper_enabled\": true,\n\"memory_scraper_enabled\": true,\n\"network_scraper_enabled\": true,\n\"process_scraper_enabled\": true,\n\"paging_scraper_enabled\": true\n}\n},\n\"processors\": {\n\"resource\": {\n\"processorType\": \"resource\",\n\"user_attributes\": [\n{\n\"key\": \"_sourceCategory\",\n\"value\": \"otel/host\"\n}\n],\n\"default_attributes\": [\n{\n\"key\": \"sumo.datasource\",\n\"value\": \"apache\"\n},\n]\n}\n}\n})"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSourceTemplateDestroy(sourceTemplate),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicSourceTemplate(testSchemaRef, testSelector, testInputJson),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourceTemplateExists("sumologic_source_template.test", &sourceTemplate, t),
					testAccCheckSourceTemplateAttributes("sumologic_source_template.test"),
					resource.TestCheckResourceAttr("sumologic_source_template.test", "schema_ref.0.type", "Mac"),
					resource.TestCheckResourceAttr("sumologic_source_template.test", "selector.0.names.0", "TestCollector1"),
				),
			},
		},
	})
}

func TestAccSumologicSourceTemplate_update(t *testing.T) {
	var sourceTemplate SourceTemplate

	testSchemaRef := "type =     \"Mac\""
	testSelector := "tags =  [\n[\n{\n key = \"tag\"\n values= [\"Value\"]\n}\n]\n]"
	testInputJson := "jsonencode({\n\"name\": \"hostmetrics_test_source_template_acc\",\n\"description\": \"Host metric source\" ,\n\"receivers\": {\n\"hostmetrics\": {\n\"receiverType\": \"hostmetrics\",\n\"collection_interval\": \"5m\",\n\"cpu_scraper_enabled\": true,\n\"disk_scraper_enabled\": true,\n\"load_scraper_enabled\": true,\n\"filesystem_scraper_enabled\": true,\n\"memory_scraper_enabled\": true,\n\"network_scraper_enabled\": true,\n\"process_scraper_enabled\": true,\n\"paging_scraper_enabled\": true\n}\n},\n\"processors\": {\n\"resource\": {\n\"processorType\": \"resource\",\n\"user_attributes\": [\n{\n\"key\": \"_sourceCategory\",\n\"value\": \"otel/host\"\n}\n],\n\"default_attributes\": [\n{\n\"key\": \"sumo.datasource\",\n\"value\": \"apache\"\n},\n]\n}\n}\n})"

	testUpdatedSchemaRef := "type =     \"Mac\""
	testUpdatedSelector := "tags =  [\n[\n{\n key = \"updatedTag\"\n values= [\"Value\"]\n}\n]\n] \n names = [\"TestCollector1\"]"
	testUpdatedInputJson := "jsonencode({\n\"name\": \"hostmetrics_test_source_template_acc\",\n\"description\": \"Host metric source\" ,\n\"receivers\": {\n\"hostmetrics\": {\n\"receiverType\": \"hostmetrics\",\n\"collection_interval\": \"5m\",\n\"cpu_scraper_enabled\": true,\n\"disk_scraper_enabled\": true,\n\"load_scraper_enabled\": true,\n\"filesystem_scraper_enabled\": true,\n\"memory_scraper_enabled\": true,\n\"network_scraper_enabled\": true,\n\"process_scraper_enabled\": true,\n\"paging_scraper_enabled\": true\n}\n},\n\"processors\": {\n\"resource\": {\n\"processorType\": \"resource\",\n\"user_attributes\": [\n{\n\"key\": \"_sourceCategory\",\n\"value\": \"otel/hostupdated\"\n}\n],\n\"default_attributes\": [\n{\n\"key\": \"sumo.datasource\",\n\"value\": \"apache\"\n},\n]\n}\n}\n})"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSourceTemplateDestroy(sourceTemplate),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicSourceTemplate(testSchemaRef, testSelector, testInputJson),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourceTemplateExists("sumologic_source_template.test", &sourceTemplate, t),
					testAccCheckSourceTemplateAttributes("sumologic_source_template.test"),
					resource.TestCheckResourceAttr("sumologic_source_template.test", "schema_ref.0.type", "Mac"),
					resource.TestCheckResourceAttr("sumologic_source_template.test", "selector.0.tags.0.0.key", "tag"),
				),
			},
			{
				Config: testAccSumologicSourceTemplateUpdate(testUpdatedSchemaRef, testUpdatedSelector, testUpdatedInputJson),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_source_template.test", "schema_ref.0.type", "Mac"),
					resource.TestCheckResourceAttr("sumologic_source_template.test", "selector.0.tags.0.0.key", "updatedTag"),
				),
			},
		},
	})
}

func testAccCheckSourceTemplateDestroy(sourceTemplate SourceTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			u, err := client.GetSourceTemplate(id)
			if err != nil {
				return fmt.Errorf("Encountered an error: " + err.Error())
			}
			if u != nil {
				return fmt.Errorf("SourceTemplate %s still exists", id)
			}
		}
		return nil
	}
}

func testAccCheckSourceTemplateExists(name string, sourceTemplate *SourceTemplate, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			//need this so that we don't get an unused import error for strconv in some cases
			return fmt.Errorf("Error = %s. SourceTemplate not found: %s", strconv.FormatBool(ok), name)
		}

		//need this so that we don't get an unused import error for strings in some cases
		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("SourceTemplate ID is not set")
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)
		newSourceTemplate, err := client.GetSourceTemplate(id)
		if err != nil {
			return fmt.Errorf("SourceTemplate %s not found", id)
		}
		sourceTemplate = newSourceTemplate
		return nil
	}
}

func testAccCheckSumologicSourceTemplateConfigImported(schemaRef string, selector string, config string, inputJson string) string {
	return fmt.Sprintf(`
resource "sumologic_source_template" "foo" {
      schema_ref {
      %s
      }
      selector {
      %s
      }

      input_json = %s
}
`, schemaRef, selector, inputJson)
}

func testAccSumologicSourceTemplate(schemaRef string, selector string, inputJson string) string {
	return fmt.Sprintf(`
resource "sumologic_source_template" "test" {
    schema_ref {
    %s
    }
    selector {
    %s
    }
    input_json = %s
}
`, schemaRef, selector, inputJson)
}

func testAccSumologicSourceTemplateUpdate(schemaRef string, selector string, inputJson string) string {
	return fmt.Sprintf(`
resource "sumologic_source_template" "test" {
      schema_ref {
      %s
      }
      selector {
      %s
      }
      input_json = %s
}
`, schemaRef, selector, inputJson)
}

func testAccCheckSourceTemplateAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(name, "total_collector_linked"),
			resource.TestCheckResourceAttrSet(name, "input_json"),
			resource.TestCheckResourceAttrSet(name, "config"),
			resource.TestCheckResourceAttrSet(name, "total_collector_linked"),
			resource.TestCheckResourceAttrSet(name, "modified_at"),
			resource.TestCheckResourceAttrSet(name, "modified_by"),
			resource.TestCheckResourceAttrSet(name, "created_at"),
			resource.TestCheckResourceAttrSet(name, "created_by"),
		)
		return f(s)
	}
}
