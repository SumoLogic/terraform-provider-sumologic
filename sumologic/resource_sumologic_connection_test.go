package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

//Testing create functionality for Connection resources
func TestAccConnection_create(t *testing.T) {
	connectionType := "WebhookConnection"
	name := acctest.RandomWithPrefix("tf-connection-test-name")
	description := acctest.RandomWithPrefix("tf-connection-test-description")
	url := "https://example.com"
	defaultPayload := "{\"eventType\" : \"{{Name}}\"}"
	webhookType := "Webhook"

	var connection Connection

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConnectionConfig(name, connectionType, description, url, webhookType, defaultPayload),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists("sumologic_connection.test", &connection, t),
					testAccCheckConnectionAttributes("sumologic_connection.test"),
					resource.TestCheckResourceAttr("sumologic_connection.test", "type", connectionType),
					resource.TestCheckResourceAttr("sumologic_connection.test", "name", name),
					resource.TestCheckResourceAttr("sumologic_connection.test", "description", description),
					resource.TestCheckResourceAttr("sumologic_connection.test", "url", url),
					resource.TestCheckResourceAttr("sumologic_connection.test", "default_payload", defaultPayload+"\n"),
					resource.TestCheckResourceAttr("sumologic_connection.test", "webhook_type", webhookType),
				),
			},
		},
	})
}

func TestAccConnection_createServiceNowWebhook(t *testing.T) {
	connectionType := "WebhookConnection"
	name := acctest.RandomWithPrefix("tf-servicenow-webhook-connection-test-name")
	description := acctest.RandomWithPrefix("tf-servicenow-webhook-connection-test-description")
	url := "https://example.com"
	defaultPayload := "{\"eventType\" : \"{{Name}}\"}"
	webhookType := "ServiceNow"
	connectionSubtype := "Incident"
	headers := "{\"Authorization\": \"Basic ABC123\"}"

	var connection Connection

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: createServiceNowWebhookConnectionConfig(name, connectionType, description, url, connectionSubtype, defaultPayload, headers),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists("sumologic_connection.serviceNowTest", &connection, t),
					testAccCheckConnectionAttributes("sumologic_connection.serviceNowTest"),
					resource.TestCheckResourceAttr("sumologic_connection.serviceNowTest", "type", connectionType),
					resource.TestCheckResourceAttr("sumologic_connection.serviceNowTest", "name", name),
					resource.TestCheckResourceAttr("sumologic_connection.serviceNowTest", "description", description),
					resource.TestCheckResourceAttr("sumologic_connection.serviceNowTest", "url", url),
					resource.TestCheckResourceAttr("sumologic_connection.serviceNowTest", "headers", headers+"\n"),
					resource.TestCheckResourceAttr("sumologic_connection.serviceNowTest", "default_payload", defaultPayload+"\n"),
					resource.TestCheckResourceAttr("sumologic_connection.serviceNowTest", "webhook_type", webhookType),
					resource.TestCheckResourceAttr("sumologic_connection.serviceNowTest", "connection_subtype", connectionSubtype),
				),
			},
		},
	})
}

func TestAccConnection_update(t *testing.T) {
	var connection Connection
	connectionType := "WebhookConnection"
	name := acctest.RandomWithPrefix("tf-connection-test-name")
	url := "https://example.com"
	defaultPayload := `{"eventType" : "{{Name}}"}`
	webhookType := "Webhook"
	fDescription := acctest.RandomWithPrefix("tf-connection-test-description")
	sDescription := acctest.RandomWithPrefix("tf-connection-test-description")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: createConnectionConfig(name, connectionType, fDescription, url, webhookType, defaultPayload),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists("sumologic_connection.test", &connection, t),
					testAccCheckConnectionAttributes("sumologic_connection.test"),
					resource.TestCheckResourceAttr("sumologic_connection.test", "type", connectionType),
					resource.TestCheckResourceAttr("sumologic_connection.test", "name", name),
					resource.TestCheckResourceAttr("sumologic_connection.test", "description", fDescription),
					resource.TestCheckResourceAttr("sumologic_connection.test", "url", url),
					resource.TestCheckResourceAttr("sumologic_connection.test", "default_payload", defaultPayload+"\n"),
					resource.TestCheckResourceAttr("sumologic_connection.test", "webhook_type", webhookType),
				),
			}, {
				Config: createConnectionConfig(name, connectionType, sDescription, url, webhookType, defaultPayload),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists("sumologic_connection.test", &connection, t),
					testAccCheckConnectionAttributes("sumologic_connection.test"),
					resource.TestCheckResourceAttr("sumologic_connection.test", "type", connectionType),
					resource.TestCheckResourceAttr("sumologic_connection.test", "name", name),
					resource.TestCheckResourceAttr("sumologic_connection.test", "description", sDescription),
					resource.TestCheckResourceAttr("sumologic_connection.test", "url", url),
					resource.TestCheckResourceAttr("sumologic_connection.test", "default_payload", defaultPayload+"\n"),
					resource.TestCheckResourceAttr("sumologic_connection.test", "webhook_type", webhookType),
				),
			},
		},
	})
}

func testAccCheckConnectionExists(name string, connection *Connection, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Connection not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Connection ID is not set")
		}

		id := rs.Primary.ID
		c := testAccProvider.Meta().(*Client)
		newConnection, err := c.GetConnection(id)
		if err != nil {
			return fmt.Errorf("Connection %s not found", id)
		}
		connection = newConnection
		return nil
	}
}

func testAccCheckConnectionAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(name, "id"),
		)
		return f(s)
	}
}

func testAccCheckConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_connection" {
			continue
		}

		id := rs.Primary.ID
		c, err := client.GetConnection(id)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if c != nil {
			return fmt.Errorf("Connection still exists")
		}
	}
	return nil
}

func createConnectionConfig(name, connectionType, desc, url, webhookType, defaultPayload string) string {
	return fmt.Sprintf(`
resource "sumologic_connection" "test" {
	name = "%s"
	type = "%s"
	description = "%s"
	url = "%s"
	webhook_type = "%s"
	default_payload = <<JSON
%s
JSON
}
`, name, connectionType, desc, url, webhookType, defaultPayload)
}

func createServiceNowWebhookConnectionConfig(name, connectionType, desc, url, connectionSubtype, defaultPayload, headers string) string {
	return fmt.Sprintf(`
resource "sumologic_connection" "serviceNowTest" {
	name = "%s"
	type = "%s"
	description = "%s"
	url = "%s"
	headers = "%s"
	webhook_type = "ServiceNow"
	connection_subtype = "%s"
	default_payload = <<JSON
%s
JSON
}
`, name, connectionType, desc, url, connectionSubtype, defaultPayload, headers)
}
