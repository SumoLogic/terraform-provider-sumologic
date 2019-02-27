package sumologic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccSumologicRoleMinimal(t *testing.T) {
	var role *Role
	resourceName := "sumologic_role.test"
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicRoleConfigMinimal,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoleExists(resourceName, &role, t),
					testAccCheckRoleAttributes(resourceName, &role),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "MyRole"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "filterPredicate", ""),
					resource.TestCheckResourceAttr(resourceName, "users.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "capabilities.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lookup_by_name", "destroy"},
			},
		},
	})
}

func TestAccSumologicRoleSimple(t *testing.T) {
	var role *Role
	resourceName := "sumologic_role.test"
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicRoleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoleExists(resourceName, &role, t),
					testAccCheckRoleAttributes(resourceName, &role),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "MyRole"),
					resource.TestCheckResourceAttr(resourceName, "description", "MyRoleDesc"),
					resource.TestCheckResourceAttr(resourceName, "filterPredicate", "Cat"),
					resource.TestCheckResourceAttr(resourceName, "users.0", "AAABBBCCCDDDEEEF"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "capabilities.0", "viewCollectors"),
					resource.TestCheckResourceAttr(resourceName, "capabilities.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lookup_by_name", "destroy"},
			},
		},
	})
}

func TestAccSumologicRoleLookupByName(t *testing.T) {
	var role *Role
	resourceName := "sumologic_role.test"
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRoleDestroy,
		// TODO: if we keep lookup_by_name, we need to beef up the tests and have 2 steps
		// TODO: first step creates the resource
		// TODO: second step looks it up by name
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicRoleConfigLookupByName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoleExists(resourceName, &role, t),
					testAccCheckRoleAttributes(resourceName, &role),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lookup_by_name", "destroy"},
			},
		},
	})
}

func TestAccSumologicRoleAllConfig(t *testing.T) {
	var role *Role
	resourceName := "sumologic_role.test"
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicRoleConfigAll,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoleExists(resourceName, &role, t),
					testAccCheckRoleAttributes(resourceName, &role),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "RoleName"),
					resource.TestCheckResourceAttr(resourceName, "description", "RoleDesc"),
					resource.TestCheckResourceAttr(resourceName, "filterPredicate", "Cat"),
					resource.TestCheckResourceAttr(resourceName, "users.0", "AAABBBCCCDDDEEEF"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "capabilities.0", "viewCollectors"),
					resource.TestCheckResourceAttr(resourceName, "capabilities.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lookup_by_name", "destroy"},
			},
		},
	})
}

func TestAccSumologicRoleChangeConfig(t *testing.T) {
	var role *Role
	resourceName := "sumologic_role.test"
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicRoleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoleExists(resourceName, &role, t),
					testAccCheckRoleAttributes(resourceName, &role),
					resource.TestCheckResourceAttr(resourceName, "name", "MyRole"),
					resource.TestCheckResourceAttr(resourceName, "description", "MyRoleDesc"),
					resource.TestCheckResourceAttr(resourceName, "filterPredicate", "Cat"),
					resource.TestCheckResourceAttr(resourceName, "users.0", "AAABBBCCCDDDEEEF"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "capabilities.0", "viewCollectors"),
					resource.TestCheckResourceAttr(resourceName, "capabilities.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"lookup_by_name", "destroy"},
			},
			{
				Config: testAccSumologicRoleConfigAll,
				Check: resource.ComposeTestCheckFunc(
					// check the id of this resource is the same as the one in the previous step
					testAccCheckRoleId(resourceName, &role),
					testAccCheckRoleExists(resourceName, &role, t),
					testAccCheckRoleAttributes(resourceName, &role),
					resource.TestCheckResourceAttr(resourceName, "name", "RoleName"),
					resource.TestCheckResourceAttr(resourceName, "description", "RoleDesc"),
					resource.TestCheckResourceAttr(resourceName, "filterPredicate", "Cat"),
					resource.TestCheckResourceAttr(resourceName, "users.0", "AAABBBCCCDDDEEEF"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "capabilities.0", "viewCollectors"),
					resource.TestCheckResourceAttr(resourceName, "capabilities.#", "1"),
				),
			},
		},
	})
}

func TestAccSumologicRoleManualDeletion(t *testing.T) {
	var role *Role

	deleteRole := func() {
		c := testAccProvider.Meta().(*Client)
		_, err := c.GetRole(role.ID)
		if err != nil {
			t.Fatal(fmt.Sprintf("attempted to delete role %s but it does not exist (%s)", role.ID, err))
		}
		err = c.DeleteRole(role.ID)
		if err != nil {
			t.Fatal(fmt.Sprintf("failed to delete role %s (%s)", role.ID, err))
		}
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicRoleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoleExists("sumologic_role.test", &role, t),
					resource.TestCheckResourceAttr("sumologic_role.test", "name", "MyRole"),
					resource.TestCheckResourceAttr("sumologic_role.test", "description", "MyRoleDesc"),
					resource.TestCheckResourceAttr("sumologic_role.test", "filterPredicate", "Cat"),
					resource.TestCheckResourceAttr("sumologic_role.test", "users.0", "AAABBBCCCDDDEEEF"),
					resource.TestCheckResourceAttr("sumologic_role.test", "users.#", "1"),
					resource.TestCheckResourceAttr("sumologic_role.test", "capabilities.0", "viewCollectors"),
					resource.TestCheckResourceAttr("sumologic_role.test", "capabilities.#", "1"),
				),
			},
			{
				PreConfig: deleteRole, // simulate a manual deletion by deleting the role between the 2 applies
				Config:    testAccSumologicRoleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRoleExists("sumologic_role.test", &role, t),
					resource.TestCheckResourceAttr("sumologic_role.test", "name", "MyRole"),
					resource.TestCheckResourceAttr("sumologic_role.test", "description", "MyRoleDesc"),
					resource.TestCheckResourceAttr("sumologic_role.test", "filterPredicate", "Cat"),
					resource.TestCheckResourceAttr("sumologic_role.test", "users.0", "AAABBBCCCDDDEEEF"),
					resource.TestCheckResourceAttr("sumologic_role.test", "users.#", "1"),
					resource.TestCheckResourceAttr("sumologic_role.test", "capabilities.0", "viewCollectors"),
					resource.TestCheckResourceAttr("sumologic_role.test", "capabilities.#", "1"),
				),
			},
		},
	})
}

// TODO: if we keep the role's destroy attribute we need to include a test checking if destroy=false works as expected

// Returns a function checking that the role with the id from the state file has an expected id.
// The expected id is specified in the role passed as parameter
func testAccCheckRoleId(name string, role **Role) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("role ID is not set")
		}

		id := rs.Primary.ID

		expectedId := (**role).ID
		if id != expectedId {
			return fmt.Errorf("incorrect role id: got %s; expected %s", id, expectedId)
		}
		return nil
	}
}

// Returns a function checking that the role with the id from the state exists.
// If the collecor exists, its attributes are updated in *role
func testAccCheckRoleExists(name string, role **Role, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("role ID is not set")
		}

		id := rs.Primary.ID
		c := testAccProvider.Meta().(*Client)
		_, err := c.GetRole(id)
		if err != nil {
			return fmt.Errorf("role %s not found", id)
		}
		return nil
	}
}

// Returns a function checking that the attributes in the state match that attributes of the actual resource created
func testAccCheckRoleAttributes(name string, expected **Role) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(name, "name", (**expected).Name),
			resource.TestCheckResourceAttr(name, "description", (**expected).Description),
			resource.TestCheckResourceAttr(name, "filterPredicate", (**expected).FilterPredicate),
			resource.TestCheckResourceAttr(name, "users.0", (**expected).Users[0]),
			resource.TestCheckResourceAttr(name, "users.#", string(len((**expected).Users))),
			resource.TestCheckResourceAttr(name, "capabilities.0", (**expected).Capabilities[0]),
			resource.TestCheckResourceAttr(name, "capabilities.#", string(len((**expected).Capabilities))),
		)
		return f(s)
	}
}

func testAccCheckRoleDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_role" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("role destruction check: role ID is not set")
		}

		id := rs.Primary.ID
		_, err := c.GetRole(id)
		if err == nil {
			return fmt.Errorf("role destruction check: role %s is still present", id)
		}
		// check that the error is what we expect
		if !strings.Contains(err.Error(), "404") {
			return fmt.Errorf("role destruction check: unexpected error %s", err)
		}
	}
	return nil
}

var testAccSumologicRoleConfigMinimal = `

resource "sumologic_role" "test" {
  name = "MyRole"
}
`

var testAccSumologicRoleConfig = `

resource "sumologic_role" "test" {
  name = "MyRole"
  description = "MyRoleDesc"
  filter_predicate = "Cat"
  users = [
    "AAABBBCCCDDDEEEF"
  ]
  capabilities = [
    "viewCollectors"
  ]
}
`

var testAccSumologicRoleConfigLookupByName = `

resource "sumologic_role" "test" {
  name = "MyRole"
  description = "MyRoleDesc"
  filter_predicate = "Cat"
  users = [
    "AAABBBCCCDDDEEEF"
  ]
  capabilities = [
    "viewCollectors"
  ]
  lookup_by_name=true
}
`

var testAccSumologicRoleConfigAll = `
resource "sumologic_role" "test" {
  name = "MyRole"
  description = "MyRoleDesc"
  filter_predicate = "Cat"
  users = [
    "AAABBBCCCDDDEEEF"
  ]
  capabilities = [
    "viewCollectors"
  ]
  lookup_by_name=true
  destroy=true
}
`
