package sumologic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccSumologicSCETagSchema_create_update(t *testing.T) {
	SkipCseTest(t)

	var TagSchema CSETagSchema
	nKey := "location"
	nLabel := "Label"
	nContentTypes := []string{"entity"}
	nFreeForm := true
	nVOValue := "value"
	nVOLabel := "label"
	nVOLink := "http://foo.bar.com"
	uLabel := "uLabel"
	resourceName := "sumologic_cse_tag_schema.tag_schema"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSETagSchemaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCreateCSETagSchemaConfig(nKey, nLabel, nContentTypes, nFreeForm, nVOValue, nVOLabel, nVOLink),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSETagSchemaExists(resourceName, &TagSchema),
					testCheckTagSchemaValues(&TagSchema, nKey, nLabel, nContentTypes, nFreeForm, nVOValue, nVOLabel, nVOLink),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testCreateCSETagSchemaConfig(nKey, uLabel, nContentTypes, nFreeForm, nVOValue, nVOLabel, nVOLink),
				Check: resource.ComposeTestCheckFunc(
					testCheckCSETagSchemaExists(resourceName, &TagSchema),
					testCheckTagSchemaValues(&TagSchema, nKey, uLabel, nContentTypes, nFreeForm, nVOValue, nVOLabel, nVOLink),
				),
			},
		},
	})
}

func testAccCSETagSchemaDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sumologic_cse_tag_schema" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("CSE Tag Schema destruction check: CSE Tag Schema key is not set")
		}

		s, err := client.GetCSETagSchema(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Encountered an error: " + err.Error())
		}
		if s != nil {
			return fmt.Errorf("Tag Schema still exists")
		}
	}
	return nil
}

func testCreateCSETagSchemaConfig(nKey string, nLabel string, nContentTypes []string, nFreeForm bool, nVOValue string, nVOLabel string, nVOLink string) string {

	return fmt.Sprintf(`
resource "sumologic_cse_tag_schema" "tag_schema" {
	key = "%s"
	label = "%s"
	content_types = ["%s"]
	free_form = "%t"	    
	value_options {
    	value = "%s"
    	label = "%s"
		link = "%s"
    }
}
`, nKey, nLabel, nContentTypes[0], nFreeForm, nVOValue, nVOLabel, nVOLink)
}

func testCheckCSETagSchemaExists(n string, TagSchema *CSETagSchema) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Tag Schema key is not set")
		}

		c := testAccProvider.Meta().(*Client)
		TagSchemaResp, err := c.GetCSETagSchema(rs.Primary.ID)
		if err != nil {
			return err
		}

		*TagSchema = *TagSchemaResp

		return nil
	}
}

func testCheckTagSchemaValues(TagSchema *CSETagSchema, nKey string, nLabel string, nContentTypes []string, nFreeForm bool, nVOValue string, nVOLabel string, nVOLink string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if TagSchema.Key != nKey {
			return fmt.Errorf("bad key, expected \"%s\", got: %#v", nKey, TagSchema.Key)
		}
		if TagSchema.Label != nLabel {
			return fmt.Errorf("bad label, expected \"%s\", got: %#v", nLabel, TagSchema.Label)
		}
		if TagSchema.ContentTypes != nil {
			if len(TagSchema.ContentTypes) != len(nContentTypes) {
				return fmt.Errorf("bad content_types list lenght, expected \"%d\", got: %d", len(nContentTypes), len(TagSchema.ContentTypes))
			}
			if TagSchema.ContentTypes[0] != nContentTypes[0] {
				return fmt.Errorf("bad content_types in list, expected \"%s\", got: %s", nContentTypes[0], TagSchema.ContentTypes[0])
			}
		}
		if TagSchema.FreeForm != nFreeForm {
			return fmt.Errorf("bad free_form field, expected \"%t\", got: %#v", nFreeForm, TagSchema.FreeForm)
		}
		if TagSchema.ValueOptionObjects[0].Value != nVOValue {
			return fmt.Errorf("bad value_option.value field, expected \"%s\", got: %#v", nVOValue, TagSchema.ValueOptionObjects[0].Value)
		}
		if TagSchema.ValueOptionObjects[0].Label != nVOLabel {
			return fmt.Errorf("bad value_option.label field, expected \"%s\", got: %#v", nVOLabel, TagSchema.ValueOptionObjects[0].Label)
		}
		if TagSchema.ValueOptionObjects[0].Link != nVOLink {
			return fmt.Errorf("bad value_option.link field, expected \"%s\", got: %#v", nVOLink, TagSchema.ValueOptionObjects[0].Link)
		}

		return nil
	}
}
