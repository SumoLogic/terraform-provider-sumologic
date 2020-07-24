package sumologic

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

//Testing create functionality for Content resources
func TestAccContent_create(t *testing.T) {
	var content Content

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContentDestroy(content),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicContent(configJson),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContentExists("sumologic_content.test", &content, t),
					testAccCheckContentAttributes("sumologic_content.test"),
				),
			},
		},
	})
}

func TestAccContent_update(t *testing.T) {
	var content Content

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContentDestroy(content),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicContent(configJson),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContentExists("sumologic_content.test", &content, t),
					testAccCheckContentAttributes("sumologic_content.test"),
				),
			}, {
				Config: testAccSumologicContent(updateConfigJson),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContentExists("sumologic_content.test", &content, t),
					testAccCheckContentAttributes("sumologic_content.test"),
				),
			},
		},
	})
}

func testAccCheckContentExists(name string, content *Content, t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Content not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Content ID is not set")
		}

		id := rs.Primary.ID
		c := testAccProvider.Meta().(*Client)
		newContent, err := c.GetContent(id, time.Minute)
		if err != nil {
			return fmt.Errorf("Content %s not found", id)
		}
		content = newContent
		return nil
	}
}

func testAccCheckContentAttributes(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		f := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttrSet(name, "parent_id"),
		)
		return f(s)
	}
}

func testAccCheckContentDestroy(content Content) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		_, err := client.GetContent(content.ID, time.Minute)
		if err == nil {
			return fmt.Errorf("Content still exists")
		}
		return nil
	}
}

var updateConfigJson = `{
	"type": "SavedSearchWithScheduleSyncDefinition",
	"name": "test-121",
	"search": {
		"queryText": "\"warn\"",
		"defaultTimeRange": "-15m",
		"byReceiptTime": false,
		"viewName": "",
		"viewStartTime": "1970-01-01T00:00:00Z",
		"queryParameters": [],
		"parsingMode": "AutoParse"
	},
	"searchSchedule": {
		"cronExpression": "0 0 * * * ? *",
		"displayableTimeRange": "-10m",
		"parseableTimeRange": {
			"type": "BeginBoundedTimeRange",
			"from": {
				"type": "RelativeTimeRangeBoundary",
				"relativeTime": "-50m"
			},
			"to": null
		},
		"timeZone": "America/Los_Angeles",
		"threshold": null,
		"notification": {
			"taskType": "EmailSearchNotificationSyncDefinition",
			"toList": ["ops@acme.org"],
			"subjectTemplate": "Search Results: {{SearchName}}",
			"includeQuery": true,
			"includeResultSet": true,
			"includeHistogram": false,
			"includeCsvAttachment": false
		},
		"scheduleType": "1Hour",
		"muteErrorEmails": false,
		"parameters": []
	},
	"description": "Runs every hour with timerange of 15m and sends email notifications updated"
}`

var configJson = `{
	"type": "SavedSearchWithScheduleSyncDefinition",
	"name": "test-121",
	"search": {
		"queryText": "\"error\"",
		"defaultTimeRange": "-15m",
		"byReceiptTime": false,
		"viewName": "",
		"viewStartTime": "1970-01-01T00:00:00Z",
		"queryParameters": [],
		"parsingMode": "AutoParse"
	},
	"searchSchedule": {
		"cronExpression": "0 0 * * * ? *",
		"displayableTimeRange": "-10m",
		"parseableTimeRange": {
			"type": "BeginBoundedTimeRange",
			"from": {
				"type": "RelativeTimeRangeBoundary",
				"relativeTime": "-50m"
			},
			"to": null
		},
		"timeZone": "America/Los_Angeles",
		"threshold": null,
		"notification": {
			"taskType": "EmailSearchNotificationSyncDefinition",
			"toList": ["ops@acme.org"],
			"subjectTemplate": "Search Results: {{SearchName}}",
			"includeQuery": true,
			"includeResultSet": true,
			"includeHistogram": false,
			"includeCsvAttachment": false
		},
		"scheduleType": "1Hour",
		"muteErrorEmails": false,
		"parameters": []
	},
	"description": "Runs every hour with timerange of 15m and sends email notifications"
}`

func testAccSumologicContent(configJson string) string {
	return fmt.Sprintf(`
data "sumologic_personal_folder" "personalFolder" {}
resource "sumologic_content" "test" {
  parent_id = "${data.sumologic_personal_folder.personalFolder.id}"
  config = <<JSON
%s
JSON
}
`, configJson)
}
