package sumologic

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicAzureEventHubLogSource_create(t *testing.T) {
	var azureEventHubLogSource PollingSource
	var collector Collector
	cName, cDescription, cCategory := getRandomizedParams()
	sName, sDescription, sCategory := getRandomizedParams()
	azureEventHubLogResourceName := "sumologic_azure_event_hub_log_source.azure"
	testNamespace := os.Getenv("SUMOLOGIC_TEST_NAMESPACE")
	testEventHub := os.Getenv("SUMOLOGIC_TEST_EVENT_HUB")
	testConsumerGroup := os.Getenv("SUMOLOGIC_TEST_CONSUMER_GROUP")
	testRegion := os.Getenv("SUMOLOGIC_TEST_REGION")
	testSASKeyName := os.Getenv("SUMOLOGIC_TEST_SAS_KEY_NAME")
	testSASKey := os.Getenv("SUMOLOGIC_TEST_SAS_KEY")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAzureEventHubLogSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicAzureEventHubLogSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testSASKeyName, testSASKey, testNamespace, testEventHub, testConsumerGroup, testRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAzureEventHubLogSourceExists(azureEventHubLogResourceName, &azureEventHubLogSource),
					testAccCheckAzureEventHubLogSourceValues(&azureEventHubLogSource, sName, sDescription, sCategory),
					testAccCheckCollectorExists("sumologic_collector.test", &collector),
					testAccCheckCollectorValues(&collector, cName, cDescription, cCategory, "Etc/UTC", ""),
					resource.TestCheckResourceAttrSet(azureEventHubLogResourceName, "id"),
					resource.TestCheckResourceAttr(azureEventHubLogResourceName, "name", sName),
					resource.TestCheckResourceAttr(azureEventHubLogResourceName, "description", sDescription),
					resource.TestCheckResourceAttr(azureEventHubLogResourceName, "category", sCategory),
					resource.TestCheckResourceAttr(azureEventHubLogResourceName, "content_type", "AzureEventHubLog"),
					resource.TestCheckResourceAttr(azureEventHubLogResourceName, "path.0.type", "AzureEventHubPath"),
				),
			},
		},
	})
}
func TestAccSumologicAzureEventHubLogSource_update(t *testing.T) {
	var azureEventHubLogSource PollingSource
	cName, cDescription, cCategory := getRandomizedParams()
	sName, sDescription, sCategory := getRandomizedParams()
	sNameUpdated, sDescriptionUpdated, sCategoryUpdated := getRandomizedParams()
	azureEventHubLogResourceName := "sumologic_azure_event_hub_log_source.azure"
	testNamespace := os.Getenv("SUMOLOGIC_TEST_NAMESPACE")
	testEventHub := os.Getenv("SUMOLOGIC_TEST_EVENT_HUB")
	testConsumerGroup := os.Getenv("SUMOLOGIC_TEST_CONSUMER_GROUP")
	testRegion := os.Getenv("SUMOLOGIC_TEST_REGION")
	testSASKeyName := os.Getenv("SUMOLOGIC_TEST_SAS_KEY_NAME")
	testSASKey := os.Getenv("SUMOLOGIC_TEST_SAS_KEY")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAzureEventHubLogSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicAzureEventHubLogSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testSASKeyName, testSASKey, testNamespace, testEventHub, testConsumerGroup, testRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAzureEventHubLogSourceExists(azureEventHubLogResourceName, &azureEventHubLogSource),
					testAccCheckAzureEventHubLogSourceValues(&azureEventHubLogSource, sName, sDescription, sCategory),
					resource.TestCheckResourceAttrSet(azureEventHubLogResourceName, "id"),
					resource.TestCheckResourceAttr(azureEventHubLogResourceName, "name", sName),
					resource.TestCheckResourceAttr(azureEventHubLogResourceName, "description", sDescription),
					resource.TestCheckResourceAttr(azureEventHubLogResourceName, "category", sCategory),
					resource.TestCheckResourceAttr(azureEventHubLogResourceName, "content_type", "AzureEventHubLog"),
					resource.TestCheckResourceAttr(azureEventHubLogResourceName, "path.0.type", "AzureEventHubPath"),
				),
			},
			{
				Config: testAccSumologicAzureEventHubLogSourceConfig(cName, cDescription, cCategory, sNameUpdated, sDescriptionUpdated, sCategoryUpdated, testSASKeyName, testSASKey, testNamespace, testEventHub, testConsumerGroup, testRegion),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAzureEventHubLogSourceExists(azureEventHubLogResourceName, &azureEventHubLogSource),
					testAccCheckAzureEventHubLogSourceValues(&azureEventHubLogSource, sNameUpdated, sDescriptionUpdated, sCategoryUpdated),
					resource.TestCheckResourceAttrSet(azureEventHubLogResourceName, "id"),
					resource.TestCheckResourceAttr(azureEventHubLogResourceName, "name", sNameUpdated),
					resource.TestCheckResourceAttr(azureEventHubLogResourceName, "description", sDescriptionUpdated),
					resource.TestCheckResourceAttr(azureEventHubLogResourceName, "category", sCategoryUpdated),
					resource.TestCheckResourceAttr(azureEventHubLogResourceName, "content_type", "AzureEventHubLog"),
					resource.TestCheckResourceAttr(azureEventHubLogResourceName, "path.0.type", "AzureEventHubPath"),
				),
			},
		},
	})
}
func testAccCheckAzureEventHubLogSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_azure_event_hub_log_source" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Azure Event Hub Log Source destruction check: Azure Event Hub Log Source ID is not set")
		}
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		s, err := client.GetPollingSource(collectorID, id)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("Polling Source still exists")
		}
	}
	return nil
}
func testAccCheckAzureEventHubLogSourceExists(n string, pollingSource *PollingSource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Polling Source ID is not set")
		}
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Polling Source id should be int; got %s", rs.Primary.ID)
		}
		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		c := testAccProvider.Meta().(*Client)
		pollingSourceResp, err := c.GetPollingSource(collectorID, id)
		if err != nil {
			return err
		}
		*pollingSource = *pollingSourceResp
		return nil
	}
}
func testAccCheckAzureEventHubLogSourceValues(pollingSource *PollingSource, name, description, category string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if pollingSource.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, pollingSource.Name)
		}
		if pollingSource.Description != description {
			return fmt.Errorf("bad description, expected \"%s\", got: %#v", description, pollingSource.Description)
		}
		if pollingSource.Category != category {
			return fmt.Errorf("bad category, expected \"%s\", got: %#v", category, pollingSource.Category)
		}
		return nil
	}
}
func testAccSumologicAzureEventHubLogSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testSASKeyName, testSASKey, testNamespace, testEventHub, testConsumerGroup, testRegion string) string {
	return fmt.Sprintf(`
resource "sumologic_collector" "test" {
	name = "%s"
	description = "%s"
	category = "%s"
}
resource "sumologic_azure_event_hub_log_source" "azure" {
	name          = "%s"
	description   = "%s"
	category      = "%s"
	content_type  = "AzureEventHubLog"
	collector_id  = "${sumologic_collector.test.id}"
  
	authentication {
	  type = "AzureEventHubAuthentication"
	  shared_access_policy_name = "%s"
	  shared_access_policy_key = "%s"
	}
  
	path {
	  type = "AzureEventHubPath"
	  namespace = "%s"
	  event_hub_name = "%s"
	  consumer_group = "%s"
	  region = "%s"
	}

	lifecycle {
		ignore_changes = [authentication.0.shared_access_policy_key]
	}
}`, cName, cDescription, cCategory, sName, sDescription, sCategory, testSASKeyName, testSASKey, testNamespace, testEventHub, testConsumerGroup, testRegion)
}
