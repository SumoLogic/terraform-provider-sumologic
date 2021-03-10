package sumologic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicGCPSource_create(t *testing.T) {
	var gcpSource GCPSource
	var collector Collector
	cName, cDescription, cCategory := getRandomizedParams()
	sName, sDescription, sCategory := getRandomizedParams()
	resourceName := "sumologic_gcp_source.gcp"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGCPSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicGCPSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGCPSourceExists(resourceName, &gcpSource),
					testAccCheckGCPSourceValues(&gcpSource, sName, sDescription, sCategory),
					testAccCheckCollectorExists("sumologic_collector.test", &collector),
					testAccCheckCollectorValues(&collector, cName, cDescription, cCategory, "Etc/UTC", ""),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
					resource.TestCheckResourceAttr(resourceName, "name", sName),
					resource.TestCheckResourceAttr(resourceName, "description", sDescription),
					resource.TestCheckResourceAttr(resourceName, "message_per_request", "false"),
					resource.TestCheckResourceAttr(resourceName, "category", sCategory),
					resource.TestCheckResourceAttr(resourceName, "content_type", "GoogleCloudLogs"),
				),
			},
		},
	})
}

func TestAccSumologicGCPSource_update(t *testing.T) {
	var gcpSource GCPSource
	cName, cDescription, cCategory := getRandomizedParams()
	sName, sDescription, sCategory := getRandomizedParams()
	sNameUpdated, sDescriptionUpdated, sCategoryUpdated := getRandomizedParams()
	resourceName := "sumologic_gcp_source.gcp"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGCPSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicGCPSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGCPSourceExists(resourceName, &gcpSource),
					testAccCheckGCPSourceValues(&gcpSource, sName, sDescription, sCategory),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
					resource.TestCheckResourceAttr(resourceName, "name", sName),
					resource.TestCheckResourceAttr(resourceName, "description", sDescription),
					resource.TestCheckResourceAttr(resourceName, "message_per_request", "false"),
					resource.TestCheckResourceAttr(resourceName, "category", sCategory),
					resource.TestCheckResourceAttr(resourceName, "content_type", "GoogleCloudLogs"),
				),
			},
			{
				Config: testAccSumologicGCPSourceConfig(cName, cDescription, cCategory, sNameUpdated, sDescriptionUpdated, sCategoryUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGCPSourceExists(resourceName, &gcpSource),
					testAccCheckGCPSourceValues(&gcpSource, sNameUpdated, sDescriptionUpdated, sCategoryUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
					resource.TestCheckResourceAttr(resourceName, "name", sNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", sDescriptionUpdated),
					resource.TestCheckResourceAttr(resourceName, "category", sCategoryUpdated),
					resource.TestCheckResourceAttr(resourceName, "content_type", "GoogleCloudLogs"),
				),
			},
		},
	})
}

func testAccCheckGCPSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_gcp_source" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("GCP Source destruction check: GCP Source ID is not set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}

		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}

		s, err := client.GetGCPSource(collectorID, id)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("GCP Source still exists")
		}
	}
	return nil
}

func testAccCheckGCPSourceExists(n string, gcpSource *GCPSource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("GCP Source ID is not set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("GCP Source id should be int; got %s", rs.Primary.ID)
		}

		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}

		c := testAccProvider.Meta().(*Client)
		gcpSourceResp, err := c.GetGCPSource(collectorID, id)
		if err != nil {
			return err
		}

		*gcpSource = *gcpSourceResp

		return nil
	}
}

func testAccCheckGCPSourceValues(gcpSource *GCPSource, name, description, category string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if gcpSource.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, gcpSource.Name)
		}
		if gcpSource.Description != description {
			return fmt.Errorf("bad description, expected \"%s\", got: %#v", description, gcpSource.Description)
		}
		if gcpSource.Category != category {
			return fmt.Errorf("bad category, expected \"%s\", got: %#v", category, gcpSource.Category)
		}
		return nil
	}
}

func testAccSumologicGCPSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory string) string {
	return fmt.Sprintf(`
resource "sumologic_collector" "test" {
	name = "%s"
	description = "%s"
	category = "%s"
}
	
resource "sumologic_gcp_source" "gcp" {
	name = "%s"
	description = "%s"
	message_per_request = false
	category = "%s"
	collector_id = "${sumologic_collector.test.id}"
	authentication {
	  type = "NoAuthentication"
	}
	
	path {
	  type = "NoPathExpression"
	}
}

`, cName, cDescription, cCategory, sName, sDescription, sCategory)
}
