package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSumologicSCECustomEntityType_create(t *testing.T) {
	SkipCseTest(t)

	var customEntityType CSECustomEntityType
	nName := "New Custom Entity Type"
	nIdentifier := "identifier"
	nFields := []string{"field1"}
	resourceName := "sumologic_cse_custom_entity_type.custom_entity_type"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSECustomEntityTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSECustomEntityTypeConfig(nName, nIdentifier, nFields),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSECustomEntityTypeExists(resourceName, &customEntityType),
					testCheckCustomEntityTypeValues(&customEntityType, nName, nIdentifier, nFields),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func TestAccSumologicSCECustomEntityType_update(t *testing.T) {
	SkipCseTest(t)

	var customEntityType CSECustomEntityType
	nName := "New Custom Entity Type"
	nIdentifier := "identifier"
	nFields := []string{"field1"}
	uName := "Changed type"
	uFields := []string{"field2"}
	resourceName := "sumologic_cse_custom_entity_type.custom_entity_type"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSECustomEntityTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSECustomEntityTypeConfig(nName, nIdentifier, nFields),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSECustomEntityTypeExists(resourceName, &customEntityType),
					testCheckCustomEntityTypeValues(&customEntityType, nName, nIdentifier, nFields),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testCreateCSECustomEntityTypeConfig(uName, nIdentifier, uFields),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSECustomEntityTypeExists(resourceName, &customEntityType),
					testCheckCustomEntityTypeValues(&customEntityType, uName, nIdentifier, uFields),
				),
			},
		},
	})
}

func testAccCSECustomEntityTypeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_custom_entity_type" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Custom Entity Type Config destruction check: CSE Custom Entity Type Config ID is not set")
		}

		s, err := client.GetCSECustomEntityType(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("entity Custom Entity Type still exists")
		}
	}
	return nil
}

func testCreateCSECustomEntityTypeConfig(nName string, nIdentifier string, nFields []string) string {

	return fmt.Sprintf(`
resource "sumologic_cse_custom_entity_type" "custom_entity_type" {
	name = "%s"
	identifier = "%s"
	fields = ["%s"]
}
`, nName, nIdentifier, nFields[0])
}

func testCheckCSECustomEntityTypeExists(n string, customEntityType *CSECustomEntityType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("entity Custom Entity Type ID is not set")
		}

		c := testAccProvider.Meta().(*Client)
		customEntityTypeResp, err := c.GetCSECustomEntityType(rs.Primary.ID)
		if err != nil {
			return err
		}

		*customEntityType = *customEntityTypeResp

		return nil
	}
}

func testCheckCustomEntityTypeValues(customEntityType *CSECustomEntityType, nName string, nIdentifier string, nFields []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if customEntityType.Name != nName {
			return fmt.Errorf("bad name, expected \"%s\", got: %#v", nName, customEntityType.Name)
		}
		if customEntityType.Identifier != nIdentifier {
			return fmt.Errorf("bad identifier, expected \"%s\", got: %#v", nIdentifier, customEntityType.Identifier)
		}
		if customEntityType.Fields != nil {
			if len(customEntityType.Fields) != len(nFields) {
				return fmt.Errorf("bad fields list, expected \"%d\", got: %d", len(nFields), len(customEntityType.Fields))
			}
			if customEntityType.Fields[0] != nFields[0] {
				return fmt.Errorf("bad field in list, expected \"%s\", got: %s", customEntityType.Fields[0], nFields[0])
			}
		}

		return nil
	}
}
