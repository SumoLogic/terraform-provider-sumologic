package sumologic

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicUniversalSource_create(t *testing.T) {
	var universalSource UniversalSource
	var collector Collector
	cName, cDescription, cCategory := getRandomizedParams()
	universalResourceName := "sumologic_universal_source.universal"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckUniversalSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicUniversalSourceConfig(cName, cDescription, cCategory, configJSON),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUniversalSourceExists(universalResourceName, &universalSource),
					testAccCheckCollectorExists("sumologic_collector.test", &collector),
					testAccCheckCollectorValues(&collector, cName, cDescription, cCategory, "Etc/UTC", ""),
					resource.TestCheckResourceAttrSet(universalResourceName, "id"),
				),
			},
		},
	})
}

func TestAccSumologicUniversalSource_update(t *testing.T) {
	var universalSource UniversalSource
	cName, cDescription, cCategory := getRandomizedParams()
	universalResourceName := "sumologic_universal_source.universal"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHTTPSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicUniversalSourceConfig(cName, cDescription, cCategory, configJSON),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUniversalSourceExists(universalResourceName, &universalSource),
					resource.TestCheckResourceAttrSet(universalResourceName, "id"),
					testAccWaitUniversalSource(),
				),
			},
			{
				Config: testAccSumologicUniversalSourceConfig(cName, cDescription, cCategory, updatedConfigJSON),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUniversalSourceExists(universalResourceName, &universalSource),
					resource.TestCheckResourceAttrSet(universalResourceName, "id"),
				),
			},
		},
	})
}

// wait function for waiting before perfroming an update to universal source as updates are blocked when the source is in pending state
func testAccWaitUniversalSource() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		//lintignore:R018
		time.Sleep(30 * time.Second)
		return nil
	}
}

func testAccCheckUniversalSourceDestroy(s *terraform.State) error {

	client := testAccProvider.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_universal_source" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Universal Source destruction check: Universal Source ID is not set")
		}
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		s, err := client.GetUniversalSource(collectorID, id)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("Universal Source still exists")
		}
	}
	return nil
}
func testAccCheckUniversalSourceExists(n string, universalSource *UniversalSource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Universal Source ID is not set")
		}
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Universal Source id should be int; got %s", rs.Primary.ID)
		}
		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		c := testAccProvider.Meta().(*Client)
		universalSourceResp, err := c.GetUniversalSource(collectorID, id)
		if err != nil {
			return err
		}
		*universalSource = *universalSourceResp
		return nil
	}
}

func testAccSumologicUniversalSourceConfig(cName, cDescription, cCategory, sConfig string) string {
	return fmt.Sprintf(`
resource "sumologic_collector" "test" {
	name = "%s"
	description = "%s"
	category = "%s"
}

resource "sumologic_universal_source" "universal" {
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
