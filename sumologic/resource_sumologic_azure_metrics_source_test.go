package sumologic

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicAzureMetricsSource_create(t *testing.T) {
	var azureMetricsSource PollingSource
	var collector Collector
	cName, cDescription, cCategory := getRandomizedParams()
	sName, sDescription, sCategory := getRandomizedParams()
	azureMetricsResourceName := "sumologic_azure_metrics_source.azure"
	testTenantId := os.Getenv("SUMOLOGIC_TEST_AZURE_TENANT_ID")
	testClientId := os.Getenv("SUMOLOGIC_TEST_AZURE_CLIENT_ID")
	testClientSecret := os.Getenv("SUMOLOGIC_TEST_AZURE_CLIENT_SECRET")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAzureMetricsSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicAzureMetricsSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testTenantId, testClientId, testClientSecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAzureMetricsSourceExists(azureMetricsResourceName, &azureMetricsSource),
					testAccCheckAzureMetricsSourceValues(&azureMetricsSource, sName, sDescription, sCategory),
					testAccCheckCollectorExists("sumologic_collector.test", &collector),
					testAccCheckCollectorValues(&collector, cName, cDescription, cCategory, "Etc/UTC", ""),
					resource.TestCheckResourceAttrSet(azureMetricsResourceName, "id"),
					resource.TestCheckResourceAttr(azureMetricsResourceName, "name", sName),
					resource.TestCheckResourceAttr(azureMetricsResourceName, "description", sDescription),
					resource.TestCheckResourceAttr(azureMetricsResourceName, "category", sCategory),
					resource.TestCheckResourceAttr(azureMetricsResourceName, "content_type", "AzureMetrics"),
					resource.TestCheckResourceAttr(azureMetricsResourceName, "path.0.type", "AzureMetricsPath"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSumologicAzureMetricsSource_update(t *testing.T) {
	var azureMetricsSource PollingSource
	cName, cDescription, cCategory := getRandomizedParams()
	sName, sDescription, sCategory := getRandomizedParams()
	sNameUpdated, sDescriptionUpdated, sCategoryUpdated := getRandomizedParams()
	azureMetricsResourceName := "sumologic_azure_metrics_source.azure"
	testTenantId := os.Getenv("SUMOLOGIC_TEST_AZURE_TENANT_ID")
	testClientId := os.Getenv("SUMOLOGIC_TEST_AZURE_CLIENT_ID")
	testClientSecret := os.Getenv("SUMOLOGIC_TEST_AZURE_CLIENT_SECRET")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAzureMetricsSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicAzureMetricsSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testTenantId, testClientId, testClientSecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAzureMetricsSourceExists(azureMetricsResourceName, &azureMetricsSource),
					testAccCheckAzureMetricsSourceValues(&azureMetricsSource, sName, sDescription, sCategory),
					resource.TestCheckResourceAttrSet(azureMetricsResourceName, "id"),
					resource.TestCheckResourceAttr(azureMetricsResourceName, "name", sName),
					resource.TestCheckResourceAttr(azureMetricsResourceName, "description", sDescription),
					resource.TestCheckResourceAttr(azureMetricsResourceName, "category", sCategory),
					resource.TestCheckResourceAttr(azureMetricsResourceName, "content_type", "AzureMetrics"),
					resource.TestCheckResourceAttr(azureMetricsResourceName, "path.0.type", "AzureMetricsPath"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccSumologicAzureMetricsSourceConfig(cName, cDescription, cCategory, sNameUpdated, sDescriptionUpdated, sCategoryUpdated, testTenantId, testClientId, testClientSecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAzureMetricsSourceExists(azureMetricsResourceName, &azureMetricsSource),
					testAccCheckAzureMetricsSourceValues(&azureMetricsSource, sNameUpdated, sDescriptionUpdated, sCategoryUpdated),
					resource.TestCheckResourceAttrSet(azureMetricsResourceName, "id"),
					resource.TestCheckResourceAttr(azureMetricsResourceName, "name", sNameUpdated),
					resource.TestCheckResourceAttr(azureMetricsResourceName, "description", sDescriptionUpdated),
					resource.TestCheckResourceAttr(azureMetricsResourceName, "category", sCategoryUpdated),
					resource.TestCheckResourceAttr(azureMetricsResourceName, "content_type", "AzureMetrics"),
					resource.TestCheckResourceAttr(azureMetricsResourceName, "path.0.type", "AzureMetricsPath"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})

}

func testAccCheckAzureMetricsSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_azure_metrics_source" {
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

func testAccCheckAzureMetricsSourceExists(n string, pollingSource *PollingSource) resource.TestCheckFunc {
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

func testAccCheckAzureMetricsSourceValues(pollingSource *PollingSource, name, description, category string) resource.TestCheckFunc {
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

func testAccSumologicAzureMetricsSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testTenantId, testClientId, testClientSecret string) string {
	return fmt.Sprintf(`
resource "sumologic_collector" "test" {
	name = "%s"
	description = "%s"
	category = "%s"
}
resource "sumologic_azure_metrics_source" "azure" {
	name = "%s"
	description = "%s"
	category = "%s"
	content_type = "AzureMetrics"
	collector_id = "${sumologic_collector.test.id}"

	authentication {
		type = "AzureClientSecretAuthentication"
		tenant_id = "%s"
		client_id = "%s"
		client_secret = "%s"
	}

	path {
		type = "AzureMetricsPath"
		environment = "Azure"
		limit_to_namespaces = ["Microsoft.ClassicStorage/storageAccounts"]
		azure_tag_filters {
			type = "AzureTagFilters"
			namespace = "Microsoft.ClassicStorage/storageAccounts"
			tags {
				name = "test-name-1"
				values = ["value1"]
			}
			tags {
				name = "test-name-2"
				values = ["value2"]
			}
		}
	}

	lifecycle {
		ignore_changes = [authentication.0.client_secret]
	}
}`, cName, cDescription, cCategory, sName, sDescription, sCategory, testTenantId, testClientId, testClientSecret)
}
