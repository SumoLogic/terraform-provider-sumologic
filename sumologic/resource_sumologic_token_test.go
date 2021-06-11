package sumologic

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicToken_basic(t *testing.T) {
	var token Token

	testName := "Test Terraform Token"
	testDescription := "Description tf test token"
	testStatus := "Active"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTokenDestroy(token),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSumologicTokenConfigImported(testName, testDescription, testStatus),
			},
			{
				ResourceName:      "sumologic_token.foo",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccSumologicToken_create(t *testing.T) {
	var token Token
	testName := "Test Terraform Token"
	testDescription := "Description tf test token"
	testStatus := "Active"
	testType := "CollectorRegistration"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTokenDestroy(token),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicToken(testName, testDescription, testStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTokenExists("sumologic_token.test", &token, t),
					testAccCheckTokenAttributes("sumologic_token.test"),
					resource.TestCheckResourceAttr("sumologic_token.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_token.test", "description", testDescription),
					resource.TestCheckResourceAttr("sumologic_token.test", "status", testStatus),
					resource.TestCheckResourceAttr("sumologic_token.test", "type", testType),
					resource.TestCheckResourceAttr("sumologic_token.test", "version", "0"),
				),
			},
		},
	})
}

func testAccCheckTokenDestroy(token Token) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			u, err := client.GetToken(id)
			if err != nil {
				return fmt.Errorf("Encountered an error: " + err.Error())
			}
			if u != nil {
				return fmt.Errorf("Token %s still exists", id)
			}
		}
		return nil
	}
}

func testAccCheckTokenExists(name string, token *Token, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			//need this so that we don't get an unused import error for strconv in some cases
			return fmt.Errorf("Error = %s. Token not found: %s", strconv.FormatBool(ok), name)
		}

		//need this so that we don't get an unused import error for strings in some cases
		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("Token ID is not set")
		}

		id := rs.Primary.ID
		c := testAccProvider.Meta().(*Client)
		newToken, err := c.GetToken(id)
		if err != nil {
			return fmt.Errorf("Token %s not found", id)
		}
		token = newToken
		return nil
	}
}

func TestAccSumologicToken_update(t *testing.T) {
	var token Token
	testName := "Test Terraform Token"
	testDescription := "Description tf test token"
	testStatus := "Active"
	testType := "CollectorRegistration"

	testUpdatedName := "Test Terraform Token Update"
	testUpdatedDescription := "Description tf test token with update"
	testUpdatedStatus := "Inactive"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTokenDestroy(token),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicToken(testName, testDescription, testStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTokenExists("sumologic_token.test", &token, t),
					testAccCheckTokenAttributes("sumologic_token.test"),
					resource.TestCheckResourceAttr("sumologic_token.test", "name", testName),
					resource.TestCheckResourceAttr("sumologic_token.test", "description", testDescription),
					resource.TestCheckResourceAttr("sumologic_token.test", "status", testStatus),
					resource.TestCheckResourceAttr("sumologic_token.test", "type", testType),
					resource.TestCheckResourceAttr("sumologic_token.test", "version", "0"),
				),
			},
			{
				Config: testAccSumologicTokenUpdate(testUpdatedName, testUpdatedDescription, testUpdatedStatus),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTokenExists("sumologic_token.test", &token, t),
					testAccCheckTokenAttributes("sumologic_token.test"),
					resource.TestCheckResourceAttr("sumologic_token.test", "name", testUpdatedName),
					resource.TestCheckResourceAttr("sumologic_token.test", "description", testUpdatedDescription),
					resource.TestCheckResourceAttr("sumologic_token.test", "status", testUpdatedStatus),
					resource.TestCheckResourceAttr("sumologic_token.test", "type", testType),
					resource.TestCheckResourceAttr("sumologic_token.test", "version", "1"),
				),
			},
		},
	})
}

func testAccCheckSumologicTokenConfigImported(name string, description string, status string) string {
	return fmt.Sprintf(`
resource "sumologic_token" "foo" {
      name = "%s"
      description = "%s"
      status = "%s"
      type = "CollectorRegistration"
}
`, name, description, status)
}

func testAccSumologicToken(name string, description string, status string) string {
	return fmt.Sprintf(`
resource "sumologic_token" "test" {
    name = "%s"
    description = "%s"
    status = "%s"
	type = "CollectorRegistration"
}
`, name, description, status)
}

func testAccSumologicTokenUpdate(name string, description string, status string) string {
	return fmt.Sprintf(`
resource "sumologic_token" "test" {
      name = "%s"
      description = "%s"
      status = "%s"
	  type = "CollectorRegistration"
}
`, name, description, status)
}

func testAccCheckTokenAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(name, "name"),
			resource.TestCheckResourceAttrSet(name, "description"),
			resource.TestCheckResourceAttrSet(name, "status"),
			resource.TestCheckResourceAttrSet(name, "type"),
			resource.TestCheckResourceAttrSet(name, "version"),
		)
		return f(s)
	}
}
