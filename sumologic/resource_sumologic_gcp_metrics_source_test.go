package sumologic

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func shouldTestGcpMetricsSource() bool {
	return !strings.EqualFold(os.Getenv("SUMOLOGIC_ENABLE_GCP_METRICS_ACC_TESTS"), "false")
}

func TestAccSumologicGcpMetricsSource_create(t *testing.T) {
	if shouldTestGcpMetricsSource() {
		var GcpMetricsSource PollingSource
		var collector Collector
		cName, cDescription, cCategory := getRandomizedParams()
		sName, sDescription, sCategory := getRandomizedParams()
		customServicePrefix := acctest.RandomWithPrefix("compute.googleapis.com")
		GcpMetricsResourceName := "sumologic_gcp_metrics_source.gcp_metrics_source"
		resource.Test(t, resource.TestCase{
			PreCheck:     func() { getServiceAccountCreds(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckGcpMetricsSourceDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccSumologicGcpMetricsSourceConfig(t, cName, cDescription, cCategory, sName, sDescription, sCategory, customServicePrefix),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckGcpMetricsSourceExists(GcpMetricsResourceName, &GcpMetricsSource),
						testAccCheckGcpMetricsSourceValues(&GcpMetricsSource, sName, sDescription, sCategory),
						testAccCheckCollectorExists("sumologic_collector.test", &collector),
						testAccCheckCollectorValues(&collector, cName, cDescription, cCategory, "Etc/UTC", ""),
						resource.TestCheckResourceAttrSet(GcpMetricsResourceName, "id"),
						resource.TestCheckResourceAttr(GcpMetricsResourceName, "name", sName),
						resource.TestCheckResourceAttr(GcpMetricsResourceName, "description", sDescription),
						resource.TestCheckResourceAttr(GcpMetricsResourceName, "category", sCategory),
						resource.TestCheckResourceAttr(GcpMetricsResourceName, "content_type", "GcpMetrics"),
						resource.TestCheckResourceAttr(GcpMetricsResourceName, "path.0.type", "GcpMetricsPath"),
						resource.TestCheckResourceAttr(GcpMetricsResourceName, "path.0.custom_services.1.service_name", "compute_instance_and_guests"),
						resource.TestCheckResourceAttr(GcpMetricsResourceName, "path.0.custom_services.1.prefixes.0", customServicePrefix),
						resource.TestCheckResourceAttr(GcpMetricsResourceName, "path.0.custom_services.1.prefixes.1", "compute.googleapis.com/guest/"),
					),
				},
			},
		})
	}
}

func TestAccSumologicGcpMetricsSource_update(t *testing.T) {
	if shouldTestGcpMetricsSource() {
		var GcpMetricsSource PollingSource
		cName, cDescription, cCategory := getRandomizedParams()
		sName, sDescription, sCategory := getRandomizedParams()
		sNameUpdated, sDescriptionUpdated, sCategoryUpdated := getRandomizedParams()
		customServicePrefix := acctest.RandomWithPrefix("compute.googleapis.com")
		updatedCustomServicePrefix := acctest.RandomWithPrefix("compute.googleapis.com")
		GcpMetricsResourceName := "sumologic_gcp_metrics_source.gcp_metrics_source"
		resource.Test(t, resource.TestCase{
			PreCheck:     func() { getServiceAccountCreds(t) },
			Providers:    testAccProviders,
			CheckDestroy: testAccCheckHTTPSourceDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccSumologicGcpMetricsSourceConfig(t, cName, cDescription, cCategory, sName, sDescription, sCategory, customServicePrefix),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckGcpMetricsSourceExists(GcpMetricsResourceName, &GcpMetricsSource),
						testAccCheckGcpMetricsSourceValues(&GcpMetricsSource, sName, sDescription, sCategory),
						resource.TestCheckResourceAttrSet(GcpMetricsResourceName, "id"),
						resource.TestCheckResourceAttr(GcpMetricsResourceName, "name", sName),
						resource.TestCheckResourceAttr(GcpMetricsResourceName, "description", sDescription),
						resource.TestCheckResourceAttr(GcpMetricsResourceName, "category", sCategory),
						resource.TestCheckResourceAttr(GcpMetricsResourceName, "content_type", "GcpMetrics"),
						resource.TestCheckResourceAttr(GcpMetricsResourceName, "path.0.type", "GcpMetricsPath"),
						resource.TestCheckResourceAttr(GcpMetricsResourceName, "path.0.custom_services.1.prefixes.0", customServicePrefix),
					),
				},
				{
					Config: testAccSumologicGcpMetricsSourceConfig(t, cName, cDescription, cCategory, sNameUpdated, sDescriptionUpdated, sCategoryUpdated, updatedCustomServicePrefix),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckGcpMetricsSourceExists(GcpMetricsResourceName, &GcpMetricsSource),
						testAccCheckGcpMetricsSourceValues(&GcpMetricsSource, sNameUpdated, sDescriptionUpdated, sCategoryUpdated),
						resource.TestCheckResourceAttrSet(GcpMetricsResourceName, "id"),
						resource.TestCheckResourceAttr(GcpMetricsResourceName, "name", sNameUpdated),
						resource.TestCheckResourceAttr(GcpMetricsResourceName, "description", sDescriptionUpdated),
						resource.TestCheckResourceAttr(GcpMetricsResourceName, "category", sCategoryUpdated),
						resource.TestCheckResourceAttr(GcpMetricsResourceName, "content_type", "GcpMetrics"),
						resource.TestCheckResourceAttr(GcpMetricsResourceName, "path.0.type", "GcpMetricsPath"),
						resource.TestCheckResourceAttr(GcpMetricsResourceName, "path.0.custom_services.1.prefixes.0", updatedCustomServicePrefix),
					),
				},
			},
		})
	}
}

func testAccCheckGcpMetricsSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_gcp_metrics_source" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("HTTP Source destruction check: HTTP Source ID is not set")
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

func testAccCheckGcpMetricsSourceExists(n string, pollingSource *PollingSource) resource.TestCheckFunc {
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

func testAccCheckGcpMetricsSourceValues(pollingSource *PollingSource, name, description, category string) resource.TestCheckFunc {
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

type ServiceAccountCreds struct {
	Type                    string `json:"type"`
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
}

func getServiceAccountCreds(t *testing.T) ServiceAccountCreds {
	var err error
	contentBytes := []byte("")
	serviceAccountDetailsJsonEnvName := "SUMOLOGIC_TEST_GOOGLE_APPLICATION_CREDENTIALS"

	serviceAccountDetailsJson, isEnvVarDefined := os.LookupEnv(serviceAccountDetailsJsonEnvName)
	if !isEnvVarDefined {
		t.Fatal(fmt.Sprintf("Environment variable %#v has to be defined", serviceAccountDetailsJsonEnvName))
	} else if len(serviceAccountDetailsJson) == 0 {
		t.Fatal(fmt.Sprintf("Environment variable %#v can not be empty string", serviceAccountDetailsJsonEnvName))
	}

	contentBytes = []byte(serviceAccountDetailsJson)
	var serviceAccountCreds ServiceAccountCreds
	err = json.Unmarshal(contentBytes, &serviceAccountCreds)
	if err != nil {
		t.Fatal(fmt.Sprintf("Failed to parse content pointed by environment variable %#v", serviceAccountDetailsJsonEnvName))
	}
	return serviceAccountCreds
}

func testAccSumologicGcpMetricsSourceConfig(t *testing.T, cName, cDescription, cCategory, sName, sDescription, sCategory, customServicePrefix string) string {
	cred := getServiceAccountCreds(t)
	srcStr := fmt.Sprintf(`
	resource "sumologic_collector" "test" {
		name = "%s"
		description = "%s"
		category = "%s"
	}
	resource "sumologic_gcp_metrics_source" "gcp_metrics_source" {
		name = "%s"
		description = "%s"
		category = "%s"
		content_type = "GcpMetrics"
		scan_interval = 300000
		paused = false
		collector_id = "${sumologic_collector.test.id}"
		authentication {
			type = "%s"
			project_id = "%s"
			private_key_id = "%s"
			private_key = <<EOPK
%sEOPK
			client_email = "%s"
			client_id = "%s"
			auth_uri = "%s"
			token_uri = "%s"
			auth_provider_x509_cert_url = "%s"
			client_x509_cert_url = "%s"
		}
		path {
			type = "GcpMetricsPath"
			limit_to_regions = ["asia-south1"]
			limit_to_services = ["Compute Engine", "CloudSQL"]
			custom_services {
				service_name = "mysql"
				prefixes = ["cloudsql.googleapis.com/database/mysql/","cloudsql.googleapis.com/database/memory/","cloudsql.googleapis.com/database/cpu","cloudsql.googleapis.com/database/disk"]
			}
			custom_services {
				service_name = "compute_instance_and_guests"
				prefixes = ["%s" ,"compute.googleapis.com/guest/", "compute.googleapis.com/instance/"]
			}
		}
		lifecycle {
			ignore_changes = [authentication[0].private_key]
		}
}
	`, cName, cDescription, cCategory, sName, sDescription, sCategory,
		cred.Type, cred.ProjectId, cred.PrivateKeyId, cred.PrivateKey, cred.ClientEmail, cred.ClientId, cred.AuthUri, cred.TokenUri, cred.AuthProviderX509CertUrl, cred.ClientX509CertUrl, customServicePrefix)
	return srcStr
}
