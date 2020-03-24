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
	url := acctest.RandomWithPrefix("https://")
	defaultPayload := `{"eventType" : "{{SearchName}}"}`
	webhookType := "Webhook"

	var connection Connection

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConnectionDestroy(connection),
		Steps: []resource.TestStep{
			{
				Config: createConnectionConfig(name, connectionType, description, url, webhookType, defaultPayload),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists("sumologic_connection.test", &connection, t),
					testAccCheckConnectionAttributes("sumologic_connection.test"),
					resource.TestCheckResourceAttr("sumologic_connection.test", "connectionType", connectionType),
					resource.TestCheckResourceAttr("sumologic_connection.test", "name", name),
					resource.TestCheckResourceAttr("sumologic_connection.test", "description", description),
					resource.TestCheckResourceAttr("sumologic_connection.test", "url", url),
					resource.TestCheckResourceAttr("sumologic_connection.test", "defaultPayload", defaultPayload),
					resource.TestCheckResourceAttr("sumologic_connection.test", "webhookType", webhookType),
				),
			},
		},
	})
}

func TestAccConnection_update(t *testing.T) {
	var connection Connection
	connectionType := "WebhookConnection"
	name := acctest.RandomWithPrefix("tf-connection-test-name")
	url := acctest.RandomWithPrefix("https://")
	defaultPayload := `{"eventType" : "{{SearchName}}"}`
	webhookType := "Webhook"
	fDescription := acctest.RandomWithPrefix("tf-connection-test-description")
	sDescription := acctest.RandomWithPrefix("tf-connection-test-description")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConnectionDestroy(connection),
		Steps: []resource.TestStep{
			{
				Config: createConnectionConfig(name, connectionType, fDescription, url, webhookType, defaultPayload),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists("sumologic_connection.test", &connection, t),
					testAccCheckConnectionAttributes("sumologic_connection.test"),
					resource.TestCheckResourceAttr("sumologic_connection.test", "connectionType", connectionType),
					resource.TestCheckResourceAttr("sumologic_connection.test", "name", name),
					resource.TestCheckResourceAttr("sumologic_connection.test", "description", fDescription),
					resource.TestCheckResourceAttr("sumologic_connection.test", "url", url),
					resource.TestCheckResourceAttr("sumologic_connection.test", "defaultPayload", defaultPayload),
					resource.TestCheckResourceAttr("sumologic_connection.test", "webhookType", webhookType),
				),
			}, {
				Config: createConnectionConfig(name, connectionType, sDescription, url, webhookType, defaultPayload),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists("sumologic_connection.test", &connection, t),
					testAccCheckConnectionAttributes("sumologic_connection.test"),
					resource.TestCheckResourceAttr("sumologic_connection.test", "connectionType", connectionType),
					resource.TestCheckResourceAttr("sumologic_connection.test", "name", name),
					resource.TestCheckResourceAttr("sumologic_connection.test", "description", sDescription),
					resource.TestCheckResourceAttr("sumologic_connection.test", "url", url),
					resource.TestCheckResourceAttr("sumologic_connection.test", "defaultPayload", defaultPayload),
					resource.TestCheckResourceAttr("sumologic_connection.test", "webhookType", webhookType),
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

func testAccCheckConnectionDestroy(connection Connection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		_, err := client.GetConnection(connection.ID)
		if err == nil {
			return fmt.Errorf("Connection still exists")
		}
		return nil
	}
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
