// ----------------------------------------------------------------------------
//
//	***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//	This file is automatically generated by Sumo Logic and manual
//	changes will be clobbered when the file is regenerated. Do not submit
//	changes to this file.
//
// ----------------------------------------------------------------------------
package sumologic

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicSubdomain_basic(t *testing.T) {
	var subdomain Subdomain
	testSubdomain := "my-company"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSubdomainDestroy(subdomain),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSumologicSubdomainConfigImported(testSubdomain),
			},
			{
				ResourceName:      "sumologic_subdomain.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSumologicSubdomain_create(t *testing.T) {
	var subdomain Subdomain
	testSubdomain := "my-company"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSubdomainDestroy(subdomain),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicSubdomain(testSubdomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubdomainExists("sumologic_subdomain.test", &subdomain, t),
					testAccCheckSubdomainAttributes("sumologic_subdomain.test"),
					resource.TestCheckResourceAttr("sumologic_subdomain.test", "subdomain", testSubdomain),
				),
			},
		},
	})
}

func TestAccSumologicSubdomain_update(t *testing.T) {
	var subdomain Subdomain
	testSubdomain := "my-company"

	testUpdatedSubdomain := "my-new-company"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSubdomainDestroy(subdomain),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicSubdomain(testSubdomain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSubdomainExists("sumologic_subdomain.test", &subdomain, t),
					testAccCheckSubdomainAttributes("sumologic_subdomain.test"),
					resource.TestCheckResourceAttr("sumologic_subdomain.test", "subdomain", testSubdomain),
				),
			},
			{
				Config: testAccSumologicSubdomainUpdate(testUpdatedSubdomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_subdomain.test", "subdomain", testUpdatedSubdomain),
				),
			},
		},
	})
}

func testAccCheckSubdomainDestroy(subdomain Subdomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, r := range s.RootModule().Resources {
			id := r.Primary.ID
			u, err := client.GetSubdomain()
			if err != nil {
				return fmt.Errorf("Encountered an error: " + err.Error())
			}
			if u != nil {
				return fmt.Errorf("Subdomain %s still exists", id)
			}
		}
		return nil
	}
}

func testAccCheckSubdomainExists(name string, subdomain *Subdomain, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			//need this so that we don't get an unused import error for strconv in some cases
			return fmt.Errorf("Error = %s. Subdomain not found: %s", strconv.FormatBool(ok), name)
		}

		//need this so that we don't get an unused import error for strings in some cases
		if strings.EqualFold(rs.Primary.ID, "") {
			return fmt.Errorf("Subdomain ID is not set")
		}

		id := rs.Primary.ID
		client := testAccProvider.Meta().(*Client)
		newSubdomain, err := client.GetSubdomain()
		if err != nil {
			return fmt.Errorf("Subdomain %s not found", id)
		}
		subdomain = newSubdomain
		return nil
	}
}

func testAccCheckSumologicSubdomainConfigImported(subdomain string) string {
	return fmt.Sprintf(`
resource "sumologic_subdomain" "foo" {
      subdomain = "%s"
}
`, subdomain)
}

func testAccSumologicSubdomain(subdomain string) string {
	return fmt.Sprintf(`
resource "sumologic_subdomain" "test" {
    subdomain = "%s"
}
`, subdomain)
}

func testAccSumologicSubdomainUpdate(subdomain string) string {
	return fmt.Sprintf(`
resource "sumologic_subdomain" "test" {
      subdomain = "%s"
}
`, subdomain)
}

func testAccCheckSubdomainAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(name, "subdomain"),
		)
		return f(s)
	}
}
