package sumologic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicServiceAccount_basic(t *testing.T) {
	suffix := uuid.New().String()[:8]
	email := fmt.Sprintf("tf-acc-%s@example.com", suffix)
	name := "tf-acc-" + suffix

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountConfig(name, email, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_service_account.foo", "name", name),
					resource.TestCheckResourceAttr("sumologic_service_account.foo", "email", email),
					resource.TestCheckResourceAttr("sumologic_service_account.foo", "is_active", "true"),
					resource.TestCheckResourceAttr("data.sumologic_service_account.by_email", "name", name),
				),
			},
			{
				Config: testAccServiceAccountConfig(name+"-renamed", email, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_service_account.foo", "name", name+"-renamed"),
					resource.TestCheckResourceAttr("sumologic_service_account.foo", "is_active", "false"),
				),
			},
			{
				ResourceName:            "sumologic_service_account.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"transfer_to", "delete_content"},
			},
		},
	})
}

func testAccServiceAccountConfig(name string, email string, isActive bool) string {
	return fmt.Sprintf(`
resource "sumologic_role" "foo" {
  name         = "%[1]s"
  description  = "Role for a service account acceptance test"
  capabilities = ["viewCollectors"]
}

resource "sumologic_service_account" "foo" {
  name      = "%[1]s"
  email     = "%[2]s"
  role_ids  = [sumologic_role.foo.id]
  is_active = %[3]t
}

data "sumologic_service_account" "by_email" {
  email = sumologic_service_account.foo.email
}
`, name, email, isActive)
}

func testAccCheckServiceAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_service_account" {
			continue
		}
		account, err := client.GetServiceAccount(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error checking service account %s: %v", rs.Primary.ID, err)
		}
		if account != nil {
			return fmt.Errorf("service account %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func TestAccSumologicServiceAccount_createInactive(t *testing.T) {
	suffix := uuid.New().String()[:8]
	name := "tf-acc-" + suffix
	email := fmt.Sprintf("tf-acc-%s@example.com", suffix)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountConfig(name, email, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_service_account.foo", "is_active", "false"),
				),
			},
		},
	})
}

func TestAccSumologicServiceAccount_updateEmailInPlace(t *testing.T) {
	suffix := uuid.New().String()[:8]
	name := "tf-acc-" + suffix
	var firstID string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountConfig(name, fmt.Sprintf("tf-acc-%s@example.com", suffix), true),
				Check:  testAccCaptureServiceAccountID("sumologic_service_account.foo", &firstID),
			},
			{
				Config: testAccServiceAccountConfig(name, fmt.Sprintf("tf-acc-%s-new@example.com", suffix), true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_service_account.foo", "email", fmt.Sprintf("tf-acc-%s-new@example.com", suffix)),
					testAccCheckServiceAccountIDUnchanged("sumologic_service_account.foo", &firstID),
				),
			},
		},
	})
}

func TestAccSumologicServiceAccount_deleteWithTransferTo(t *testing.T) {
	suffix := uuid.New().String()[:8]

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "sumologic_role" "foo" {
  name         = "tf-acc-%[1]s"
  description  = "Role for a service account acceptance test"
  capabilities = ["viewCollectors"]
}

resource "sumologic_service_account" "receiver" {
  name     = "tf-acc-%[1]s-receiver"
  email    = "tf-acc-%[1]s-receiver@example.com"
  role_ids = [sumologic_role.foo.id]
}

resource "sumologic_service_account" "deleted" {
  name        = "tf-acc-%[1]s-deleted"
  email       = "tf-acc-%[1]s-deleted@example.com"
  role_ids    = [sumologic_role.foo.id]
  transfer_to = sumologic_service_account.receiver.id
}
`, suffix),
			},
		},
	})
}

func TestAccSumologicServiceAccount_deleteContent(t *testing.T) {
	suffix := uuid.New().String()[:8]

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "sumologic_role" "foo" {
  name         = "tf-acc-%[1]s"
  description  = "Role for a service account acceptance test"
  capabilities = ["viewCollectors"]
}

resource "sumologic_service_account" "foo" {
  name           = "tf-acc-%[1]s"
  email          = "tf-acc-%[1]s@example.com"
  role_ids       = [sumologic_role.foo.id]
  delete_content = true
}
`, suffix),
			},
		},
	})
}

func TestAccSumologicServiceAccount_conflictingDeleteOptions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
resource "sumologic_service_account" "foo" {
  name           = "tf-acc-conflict"
  email          = "tf-acc-conflict@example.com"
  role_ids       = ["0000000000000001"]
  transfer_to    = "0000000000000002"
  delete_content = true
}
`,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`conflicts with`),
			},
		},
	})
}

// A service account deleted in the UI is removed from state on the next refresh and created again on apply.
func TestAccSumologicServiceAccount_deletedOutsideTerraform(t *testing.T) {
	suffix := uuid.New().String()[:8]
	name := "tf-acc-" + suffix
	email := fmt.Sprintf("tf-acc-%s@example.com", suffix)
	var firstID string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountConfig(name, email, true),
				Check:  testAccCaptureServiceAccountID("sumologic_service_account.foo", &firstID),
			},
			{
				PreConfig: func() {
					client := testAccProvider.Meta().(*Client)
					if err := client.DeleteServiceAccount(firstID, "", false); err != nil {
						t.Fatalf("could not delete the service account out of band: %v", err)
					}
				},
				Config: testAccServiceAccountConfig(name, email, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sumologic_service_account.foo", "email", email),
					func(s *terraform.State) error {
						if s.RootModule().Resources["sumologic_service_account.foo"].Primary.ID == firstID {
							return fmt.Errorf("expected the service account to be created again with a new id")
						}
						return nil
					},
				),
			},
		},
	})
}

func testAccCaptureServiceAccountID(resourceName string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}
		*id = rs.Primary.ID
		return nil
	}
}

func testAccCheckServiceAccountIDUnchanged(resourceName string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}
		if rs.Primary.ID != *id {
			return fmt.Errorf("expected the service account to be updated in place, id changed from %s to %s", *id, rs.Primary.ID)
		}
		return nil
	}
}

func TestAccSumologicServiceAccount_dataSourceById(t *testing.T) {
	suffix := uuid.New().String()[:8]
	name := "tf-acc-" + suffix

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountConfig(name, fmt.Sprintf("tf-acc-%s@example.com", suffix), true) + `
data "sumologic_service_account" "by_id" {
  id = sumologic_service_account.foo.id
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.sumologic_service_account.by_id", "name", name),
					resource.TestCheckResourceAttrPair("data.sumologic_service_account.by_id", "role_ids.#", "sumologic_service_account.foo", "role_ids.#"),
				),
			},
		},
	})
}

func TestAccSumologicServiceAccount_dataSourceUnknownEmail(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
data "sumologic_service_account" "missing" {
  email = "tf-acc-does-not-exist@example.com"
}
`,
				ExpectError: regexp.MustCompile(`not found`),
			},
		},
	})
}

func TestAccSumologicServiceAccount_updateRoleIds(t *testing.T) {
	suffix := uuid.New().String()[:8]
	config := func(roles string) string {
		return fmt.Sprintf(`
resource "sumologic_role" "first" {
  name         = "tf-acc-%[1]s-first"
  description  = "First role for a service account acceptance test"
  capabilities = ["viewCollectors"]
}

resource "sumologic_role" "second" {
  name         = "tf-acc-%[1]s-second"
  description  = "Second role for a service account acceptance test"
  capabilities = ["viewFields"]
}

resource "sumologic_service_account" "foo" {
  name     = "tf-acc-%[1]s"
  email    = "tf-acc-%[1]s@example.com"
  role_ids = %[2]s
}
`, suffix, roles)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config("[sumologic_role.first.id]"),
				Check:  resource.TestCheckResourceAttr("sumologic_service_account.foo", "role_ids.#", "1"),
			},
			{
				Config: config("[sumologic_role.first.id, sumologic_role.second.id]"),
				Check:  resource.TestCheckResourceAttr("sumologic_service_account.foo", "role_ids.#", "2"),
			},
			{
				Config: config("[sumologic_role.second.id]"),
				Check:  resource.TestCheckResourceAttr("sumologic_service_account.foo", "role_ids.#", "1"),
			},
		},
	})
}
