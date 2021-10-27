package sumologic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourcSumologicUser_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAccSumologicUserConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceUserCheck("data.sumologic_user.by_email", "sumologic_user.test"),
					testAccDataSourceUserCheck("data.sumologic_user.by_id", "sumologic_user.test"),
				),
			},
		},
	})
}

func testAccDataSourceUserCheck(email, reference string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(email, "id"),
		resource.TestCheckResourceAttrPair(email, "id", reference, "id"),
		resource.TestCheckResourceAttrPair(email, "email", reference, "email"),
		resource.TestCheckResourceAttrPair(email, "first_name", reference, "first_name"),
		resource.TestCheckResourceAttrPair(email, "last_name", reference, "last_name"),
		resource.TestCheckResourceAttrPair(email, "is_active", reference, "is_active"),
	)
}

var testDataSourceAccSumologicUserConfig = `
resource "sumologic_user" "test1" {
  first_name = "Test1"
  last_name = "User1"
  email = "user1@example.com"
  is_active = "true"
  role_ids = ["${sumologic_role.test.id}"]
  transfer_to = "${sumologic_user.test2.id}"
}

resource "sumologic_user" "test2" {
	first_name = "Test2"
	last_name = "User2"
	email = "user2@example.com"
	is_active = "true"
	role_ids = ["${sumologic_role.test.id}"]
	transfer_to = "${sumologic_user.test1.id}"
  }

data "sumologic_user" "by_email" {
  email = "${sumologic_user.test.email}"
}

data "sumologic_user" "by_id" {
  id = "${sumologic_user.test.id}"
}
`
