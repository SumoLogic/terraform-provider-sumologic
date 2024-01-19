package sumologic

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicApp_basic(t *testing.T) {
	// uuid of Varnish - OpenTelemetry app
	uuid := "d2ef33c3-67f2-4438-9124-14a30ec2ecf3"
	version := "1.0.3"
	parameterKey := "key"
	parameterValue := "value"

	tfResourceName := "tf_app_test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicApp(tfResourceName, uuid, version, parameterKey, parameterValue),
			},
			{
				ResourceName:      fmt.Sprintf("sumologic_app.%s", tfResourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSumologicApp_create(t *testing.T) {
	// create config
	// uuid of Varnish - OpenTelemetry app
	uuid := "d2ef33c3-67f2-4438-9124-14a30ec2ecf3"
	version := "1.0.4"
	parameterKey := "key"
	parameterValue := "value"

	tfResourceName := "tf_app_test"
	tfAppResource := fmt.Sprintf("sumologic_app.%s", tfResourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicApp(tfResourceName, uuid, version, parameterKey, parameterValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppInstanceExists(tfAppResource, t),
					resource.TestCheckResourceAttr(tfAppResource, "uuid", uuid),
					resource.TestCheckResourceAttr(tfAppResource, "version", version),
					resource.TestCheckResourceAttr(tfAppResource, "parameters.%", "1"),
					resource.TestCheckResourceAttr(tfAppResource, "parameters.key", "value"),
				),
			},
		},
	})
}

func TestAccSumologicApp_update(t *testing.T) {

	// create config
	// uuid of Varnish - OpenTelemetry app
	uuid := "d2ef33c3-67f2-4438-9124-14a30ec2ecf3"
	version := "1.0.3"
	parameterKey := "key"
	parameterValue := "value"

	tfResourceName := "tf_app_test"
	tfAppResource := fmt.Sprintf("sumologic_app.%s", tfResourceName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicApp(tfResourceName, uuid, version, parameterKey, parameterValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppInstanceExists(tfAppResource, t),
					resource.TestCheckResourceAttr(tfAppResource, "uuid", uuid),
					resource.TestCheckResourceAttr(tfAppResource, "version", version),
					resource.TestCheckResourceAttr(tfAppResource, "parameters.%", "1"),
					resource.TestCheckResourceAttr(tfAppResource, "parameters.key", "value"),
				),
			},
			{
				Config: testAccSumologicApp(tfResourceName, uuid, "1.0.4", parameterKey, parameterValue),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppInstanceExists(tfAppResource, t),
					resource.TestCheckResourceAttr(tfAppResource, "uuid", uuid),
					resource.TestCheckResourceAttr(tfAppResource, "version", "1.0.4"),
					resource.TestCheckResourceAttr(tfAppResource, "parameters.%", "1"),
					resource.TestCheckResourceAttr(tfAppResource, "parameters.key", "value"),
				),
			},
		},
	})
}

func testAccCheckAppDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			appInstance, _ := client.GetAppInstance(id)
			if appInstance != nil {
				return fmt.Errorf("AppInstance %s still exists", id)
			}
		}
		return nil
	}
}

func testAccCheckAppInstanceExists(name string, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Error = %s. App Instance not found: %s", strconv.FormatBool(ok), name)
		}

		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("App Instance ID is not set")
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)
		_, err := client.GetAppInstance(id)
		if err != nil {
			return fmt.Errorf("App instance (id=%s) not found", id)
		}
		return nil
	}
}

func testAccSumologicApp(tfResourceName string, uuid string, version string, parameterKey string, parameterValue string) string {

	return fmt.Sprintf(`

	resource "sumologic_app" "%s" {
		uuid = "%s"
		version = "%s"
		parameters = {
			"%s": "%s"
		}
	}
	`, tfResourceName, uuid, version, parameterKey, parameterValue)
}
