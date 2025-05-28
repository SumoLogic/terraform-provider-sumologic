package sumologic

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicCloudToCloudSource_create(t *testing.T) {
	var cloudToCloudSource CloudToCloudSource
	var collector Collector
	cName, cDescription, cCategory := getRandomizedParams()
	cloudToCloudResourceName := "sumologic_cloud_to_cloud_source.okta"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudToCloudSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCloudToCloudSourceConfig(cName, cDescription, cCategory, configJSON),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudToCloudSourceExists(cloudToCloudResourceName, &cloudToCloudSource),
					testAccCheckCollectorExists("sumologic_collector.test", &collector),
					testAccCheckCollectorValues(&collector, cName, cDescription, cCategory, "Etc/UTC", ""),
					resource.TestCheckResourceAttrSet(cloudToCloudResourceName, "id"),
				),
			},
		},
	})
}

func TestAccSumologicCloudToCloudSource_update(t *testing.T) {
	var cloudToCloudSource CloudToCloudSource
	cName, cDescription, cCategory := getRandomizedParams()
	cloudToCloudResourceName := "sumologic_cloud_to_cloud_source.okta"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHTTPSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicCloudToCloudSourceConfig(cName, cDescription, cCategory, configJSON),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudToCloudSourceExists(cloudToCloudResourceName, &cloudToCloudSource),
					resource.TestCheckResourceAttrSet(cloudToCloudResourceName, "id"),
					testAccWaitCloudToCloudSource(),
				),
			},
			{
				Config: testAccSumologicCloudToCloudSourceConfig(cName, cDescription, cCategory, updatedConfigJSON),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudToCloudSourceExists(cloudToCloudResourceName, &cloudToCloudSource),
					resource.TestCheckResourceAttrSet(cloudToCloudResourceName, "id"),
				),
			},
		},
	})
}

// wait function for waiting before perfroming an update to cloud-to-cloud source as updates are blocked when the source is in pending state
func testAccWaitCloudToCloudSource() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		//lintignore:R018
		time.Sleep(90 * time.Second)
		return nil
	}
}

func testAccCheckCloudToCloudSourceDestroy(s *terraform.State) error {

	client := testAccProvider.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cloud_to_cloud_source" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Cloud-to-Cloud Source destruction check: Cloud-to-Cloud Source ID is not set")
		}
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		s, err := client.GetCloudToCloudSource(collectorID, id)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("Cloud-to-Cloud Source still exists")
		}
	}
	return nil
}
func testAccCheckCloudToCloudSourceExists(n string, cloudToCloudSource *CloudToCloudSource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Cloud-to-Cloud Source ID is not set")
		}
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Cloud-to-Cloud Source id should be int; got %s", rs.Primary.ID)
		}
		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		c := testAccProvider.Meta().(*Client)
		cloudToCloudSourceResp, err := c.GetCloudToCloudSource(collectorID, id)
		if err != nil {
			return err
		}
		*cloudToCloudSource = *cloudToCloudSourceResp
		return nil
	}
}

func testAccSumologicCloudToCloudSourceConfig(cName, cDescription, cCategory, sConfig string) string {
	return fmt.Sprintf(`
resource "sumologic_collector" "test" {
	name = "%s"
	description = "%s"
	category = "%s"
}

resource "sumologic_cloud_to_cloud_source" "okta" {
	collector_id    = sumologic_collector.test.id
	schema_ref = {
	  type = "Okta"
	  }
	config = <<JSON
	%s
	JSON
   }
`, cName, cDescription, cCategory, sConfig)
}

var configJSON = `{
	"name":"Okta_source",
	"domain":"demo.okta.com",
	"collectAll":true,
	"apiKey":"secret",
	"fields":{
	"_siemForward":false
	},
	"pollingInterval": 30
}`

var updatedConfigJSON = `{
	"name":"Okta_source_new",
	"domain":"demo.okta.com",
	"collectAll":true,
	"apiKey":"secret",
	"fields":{
	"_siemForward":false
	},
	"pollingInterval": 300
}`
