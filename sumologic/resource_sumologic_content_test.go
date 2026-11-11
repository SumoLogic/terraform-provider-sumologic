package sumologic

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

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

func TestAccContent_multipleNotifications(t *testing.T) {
	var content Content

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContentDestroy(content),
		Steps: []resource.TestStep{
			{
				Config: testAccSumologicContentWithWebhook(configJson),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContentExists("sumologic_content.test", &content, t),
					testAccCheckContentAttributes("sumologic_content.test"),
				),
			}, {
				Config: testAccSumologicContentWithWebhook(multiNotificationConfigJson),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckContentExists("sumologic_content.test", &content, t),
					testAccCheckContentAttributes("sumologic_content.test"),
					resource.TestCheckResourceAttrWith("sumologic_content.test", "config",
						testAccCheckContentNotificationTypes(
							"EmailSearchNotificationSyncDefinition",
							"WebhookSearchNotificationSyncDefinition",
						)),
				),
			},
		},
	})
}

func testAccCheckContentNotificationTypes(expected ...string) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		var config struct {
			SearchSchedule struct {
				Notification  interface{} `json:"notification"`
				Notifications []struct {
					TaskType string `json:"taskType"`
				} `json:"notifications"`
			} `json:"searchSchedule"`
		}
		if err := json.Unmarshal([]byte(value), &config); err != nil {
			return err
		}
		if config.SearchSchedule.Notification != nil {
			return fmt.Errorf("expected no singular notification, got %v", config.SearchSchedule.Notification)
		}
		if len(config.SearchSchedule.Notifications) != len(expected) {
			return fmt.Errorf("expected %d notifications, got %d", len(expected), len(config.SearchSchedule.Notifications))
		}
		for i, n := range config.SearchSchedule.Notifications {
			if n.TaskType != expected[i] {
				return fmt.Errorf("expected notifications[%d] to be %s, got %s", i, expected[i], n.TaskType)
			}
		}
		return nil
	}
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
			return fmt.Errorf("Content(id=%s) still exists", content.ID)
		}
		return nil
	}
}

var updateConfigJson = `{
	"type": "SavedSearchWithScheduleSyncDefinition",
	"name": "tf-content-scheduled-search-01",
	"search": {
		"queryText": "\"warn\"",
		"defaultTimeRange": "-15m",
		"byReceiptTime": false,
		"viewName": "",
		"viewStartTime": "1970-01-01T00:00:00Z",
		"queryParameters": [],
		"intervalTimeType": "messageTime",
		"parsingMode": "Manual"
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
			"subjectTemplate": "Search Results: {{Name}}",
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
	"name": "tf-content-scheduled-search-01",
	"search": {
		"queryText": "\"error\"",
		"defaultTimeRange": "-15m",
		"byReceiptTime": false,
		"viewName": "",
		"viewStartTime": "1970-01-01T00:00:00Z",
		"queryParameters": [],
		"intervalTimeType": "messageTime",
		"parsingMode": "Manual"
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
			"subjectTemplate": "Search Results: {{Name}}",
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

var multiNotificationConfigJson = `{
	"type": "SavedSearchWithScheduleSyncDefinition",
	"name": "tf-content-scheduled-search-01",
	"search": {
		"queryText": "\"error\"",
		"defaultTimeRange": "-15m",
		"byReceiptTime": false,
		"viewName": "",
		"viewStartTime": "1970-01-01T00:00:00Z",
		"queryParameters": [],
		"intervalTimeType": "messageTime",
		"parsingMode": "Manual"
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
		"notifications": [
			{
				"taskType": "EmailSearchNotificationSyncDefinition",
				"toList": ["ops@acme.org"],
				"subjectTemplate": "Search Results: {{SearchName}}",
				"includeQuery": true,
				"includeResultSet": true,
				"includeHistogram": true,
				"includeCsvAttachment": false
			},
			{
				"taskType": "WebhookSearchNotificationSyncDefinition",
				"webhookId": "${sumologic_connection.tf_test_webhook.id}",
				"payload": "{\"text\": \"Alert: {{SearchName}} triggered with {{NumRawResults}} results\"}",
				"itemizeAlerts": false,
				"maxItemizedAlerts": 0
			}
		],
		"scheduleType": "1Hour",
		"muteErrorEmails": false,
		"parameters": []
	},
	"description": "Runs every hour with timerange of 15m and sends an email and a webhook notification"
}`

func testAccSumologicContentWithWebhook(configJson string) string {
	return fmt.Sprintf(`
resource "sumologic_connection" "tf_test_webhook" {
  name            = "TF Test Webhook for Content Multi Notif"
  type            = "WebhookConnection"
  description     = "Webhook connection for content multi-notification test"
  url             = "https://example.com"
  webhook_type    = "Webhook"
  default_payload = <<JSON
{"eventType": "{{Name}}"}
JSON
}
%s`, testAccSumologicContent(configJson))
}

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
