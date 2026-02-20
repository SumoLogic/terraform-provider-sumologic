package sumologic

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccPreCheckWithO365(t *testing.T) {
	testAccPreCheck(t)
	if v := os.Getenv("SUMOLOGIC_TEST_O365_TENANT_ID"); v == "" {
		t.Fatal("SUMOLOGIC_TEST_O365_TENANT_ID must be set for O365 acceptance tests")
	}
	if v := os.Getenv("SUMOLOGIC_TEST_O365_CLIENT_ID"); v == "" {
		t.Fatal("SUMOLOGIC_TEST_O365_CLIENT_ID must be set for O365 acceptance tests")
	}
	if v := os.Getenv("SUMOLOGIC_TEST_O365_CLIENT_SECRET"); v == "" {
		t.Fatal("SUMOLOGIC_TEST_O365_CLIENT_SECRET must be set for O365 acceptance tests")
	}
}

func TestAccSumologicO365AuditSource_create(t *testing.T) {
	var o365Source HTTPSource
	var collector Collector
	cName := acctest.RandomWithPrefix("tf-acc-test")
	cDescription := acctest.RandomWithPrefix("tf-acc-test")
	cCategory := acctest.RandomWithPrefix("tf-acc-test")
	sName := acctest.RandomWithPrefix("tf-acc-test")
	sDescription := acctest.RandomWithPrefix("tf-acc-test")
	sCategory := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "sumologic_o365_audit_source.o365"
	testO365TenantId := os.Getenv("SUMOLOGIC_TEST_O365_TENANT_ID")
	testO365ClientId := os.Getenv("SUMOLOGIC_TEST_O365_CLIENT_ID")
	testO365ClientSecret := os.Getenv("SUMOLOGIC_TEST_O365_CLIENT_SECRET")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithO365(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckO365AuditSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicO365AuditSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testO365TenantId, testO365ClientId, testO365ClientSecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckO365AuditSourceExists(resourceName, &o365Source),
					testAccCheckO365AuditSourceValues(&o365Source, sName, sDescription, sCategory),
					testAccCheckCollectorExists("sumologic_collector.test", &collector),
					testAccCheckCollectorValues(&collector, cName, cDescription, cCategory, "Etc/UTC", ""),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
					resource.TestCheckResourceAttr(resourceName, "name", sName),
					resource.TestCheckResourceAttr(resourceName, "description", sDescription),
					resource.TestCheckResourceAttr(resourceName, "category", sCategory),
					resource.TestCheckResourceAttr(resourceName, "third_party_ref.0.resources.0.service_type", "O365AuditNotification"),
					resource.TestCheckResourceAttr(resourceName, "third_party_ref.0.resources.0.path.0.type", "O365NotificationPath"),
					resource.TestCheckResourceAttr(resourceName, "third_party_ref.0.resources.0.path.0.workload", "Audit.Exchange"),
					resource.TestCheckResourceAttr(resourceName, "third_party_ref.0.resources.0.path.0.region", "Commercial"),
					resource.TestCheckResourceAttr(resourceName, "third_party_ref.0.resources.0.authentication.0.type", "O365AppRegistrationAuthentication"),
					resource.TestCheckResourceAttr(resourceName, "third_party_ref.0.resources.0.authentication.0.tenant_id", testO365TenantId),
					resource.TestCheckResourceAttr(resourceName, "third_party_ref.0.resources.0.authentication.0.client_id", testO365ClientId),
					resource.TestCheckResourceAttrSet(resourceName, "third_party_ref.0.resources.0.authentication.0.client_secret"),
				),
			},
		},
	})
}

func TestAccSumologicO365AuditSource_update(t *testing.T) {
	var o365Source HTTPSource
	cName := acctest.RandomWithPrefix("tf-acc-test")
	cDescription := acctest.RandomWithPrefix("tf-acc-test")
	cCategory := acctest.RandomWithPrefix("tf-acc-test")
	sName := acctest.RandomWithPrefix("tf-acc-test")
	sDescription := acctest.RandomWithPrefix("tf-acc-test")
	sCategory := acctest.RandomWithPrefix("tf-acc-test")
	sNameUpdated := acctest.RandomWithPrefix("tf-acc-test")
	sDescriptionUpdated := acctest.RandomWithPrefix("tf-acc-test")
	sCategoryUpdated := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "sumologic_o365_audit_source.o365"
	testO365TenantId := os.Getenv("SUMOLOGIC_TEST_O365_TENANT_ID")
	testO365ClientId := os.Getenv("SUMOLOGIC_TEST_O365_CLIENT_ID")
	testO365ClientSecret := os.Getenv("SUMOLOGIC_TEST_O365_CLIENT_SECRET")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithO365(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckO365AuditSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicO365AuditSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testO365TenantId, testO365ClientId, testO365ClientSecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckO365AuditSourceExists(resourceName, &o365Source),
					testAccCheckO365AuditSourceValues(&o365Source, sName, sDescription, sCategory),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
					resource.TestCheckResourceAttr(resourceName, "name", sName),
					resource.TestCheckResourceAttr(resourceName, "description", sDescription),
					resource.TestCheckResourceAttr(resourceName, "category", sCategory),
				),
			},
			{
				Config: testAccSumologicO365AuditSourceConfig(cName, cDescription, cCategory, sNameUpdated, sDescriptionUpdated, sCategoryUpdated, testO365TenantId, testO365ClientId, testO365ClientSecret),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckO365AuditSourceExists(resourceName, &o365Source),
					testAccCheckO365AuditSourceValues(&o365Source, sNameUpdated, sDescriptionUpdated, sCategoryUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
					resource.TestCheckResourceAttr(resourceName, "name", sNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "description", sDescriptionUpdated),
					resource.TestCheckResourceAttr(resourceName, "category", sCategoryUpdated),
				),
			},
		},
	})
}

func testAccCheckO365AuditSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_o365_audit_source" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("O365 Audit Source destruction check: O365 Audit Source ID is not set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}

		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}

		s, err := client.GetHTTPSource(collectorID, id)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("O365 Audit Source still exists")
		}
	}
	return nil
}

func testAccCheckO365AuditSourceExists(n string, o365Source *HTTPSource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("O365 Audit Source ID is not set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("O365 Audit Source id should be int; got %s", rs.Primary.ID)
		}

		collectorID, err := strconv.Atoi(rs.Primary.Attributes["collector_id"])
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}

		c := testAccProvider.Meta().(*Client)
		o365SourceResp, err := c.GetHTTPSource(collectorID, id)
		if err != nil {
			return err
		}

		*o365Source = *o365SourceResp

		return nil
	}
}

func testAccCheckO365AuditSourceValues(o365Source *HTTPSource, name, description, category string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o365Source.Name != name {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", name, o365Source.Name)
		}
		if o365Source.Description != description {
			return fmt.Errorf("bad description, expected \"%s\", got: %#v", description, o365Source.Description)
		}
		if o365Source.Category != category {
			return fmt.Errorf("bad category, expected \"%s\", got: %#v", category, o365Source.Category)
		}
		return nil
	}
}

func testAccSumologicO365AuditSourceConfig(cName, cDescription, cCategory, sName, sDescription, sCategory, testO365TenantId, testO365ClientId, testO365ClientSecret string) string {
	return fmt.Sprintf(`
resource "sumologic_collector" "test" {
	name = "%s"
	description = "%s"
	category = "%s"
}

resource "sumologic_o365_audit_source" "o365" {
	name = "%s"
	description = "%s"
	category = "%s"
	collector_id = "${sumologic_collector.test.id}"

	third_party_ref {
		resources {
			service_type = "O365AuditNotification"

			path {
				type = "O365NotificationPath"
				workload = "Audit.Exchange"
				region = "Commercial"
			}

			authentication {
				type = "O365AppRegistrationAuthentication"
				tenant_id = "%s"
				client_id = "%s"
				client_secret = "%s"
			}
		}
	}
}
`, cName, cDescription, cCategory, sName, sDescription, sCategory, testO365TenantId, testO365ClientId, testO365ClientSecret)
}
