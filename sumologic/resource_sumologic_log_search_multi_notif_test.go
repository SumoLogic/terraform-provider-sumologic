package sumologic

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogSearch_MultiNotificationEmailAndWebhook_CreateAndRead(t *testing.T) {
	var createdPayload map[string]interface{}
	var resourceID = "0000000000001234"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && r.URL.Path == "/api/v1/logSearches":
			body, _ := io.ReadAll(r.Body)
			err := json.Unmarshal(body, &createdPayload)
			if err != nil {
				t.Fatalf("Failed to unmarshal create request body: %v", err)
			}

			// Verify the request has notifications (plural), not notification (singular)
			schedule, ok := createdPayload["schedule"].(map[string]interface{})
			if !ok {
				t.Fatal("Request body missing 'schedule' field")
			}

			notifications, ok := schedule["notifications"].([]interface{})
			if !ok || len(notifications) == 0 {
				t.Fatal("Request body missing 'notifications' array in schedule")
			}

			if len(notifications) != 2 {
				t.Fatalf("Expected 2 notifications in request, got %d", len(notifications))
			}

			// Verify notification is absent (omitempty)
			if _, hasNotification := schedule["notification"]; hasNotification {
				t.Fatal("Request should NOT contain 'notification' field when 'notifications' is used")
			}

			// Verify first notification is Email
			n1 := notifications[0].(map[string]interface{})
			if n1["taskType"] != "EmailSearchNotificationSyncDefinition" {
				t.Fatalf("Expected first notification to be Email, got: %v", n1["taskType"])
			}

			// Verify second notification is Webhook
			n2 := notifications[1].(map[string]interface{})
			if n2["taskType"] != "WebhookSearchNotificationSyncDefinition" {
				t.Fatalf("Expected second notification to be Webhook, got: %v", n2["taskType"])
			}

			t.Logf("CREATE request validated: 2 notifications (Email + Webhook)")
			t.Logf("  Email notification toList: %v", n1["toList"])
			t.Logf("  Webhook notification webhookId: %v", n2["webhookId"])

			// Return the created resource with ID
			createdPayload["id"] = resourceID
			resp, _ := json.Marshal(createdPayload)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resp)

		case r.Method == "GET" && r.URL.Path == fmt.Sprintf("/api/v1/logSearches/%s", resourceID):
			// Return the same resource (simulating GET after create)
			resp, _ := json.Marshal(createdPayload)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resp)

		case r.Method == "DELETE" && r.URL.Path == fmt.Sprintf("/api/v1/logSearches/%s", resourceID):
			w.WriteHeader(http.StatusNoContent)

		default:
			t.Logf("Unhandled request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create a client pointing to our mock server
	client, err := NewClient("testAccessId", "testAccessKey", "", "", server.URL+"/api/", false)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Build the LogSearch with multi-notification (Email + Webhook)
	logSearch := LogSearch{
		Name:             "TF Multi Notif Email+Webhook Test",
		Description:      "TF Multi Notification with Email and Webhook",
		ParentId:         "0000000000000001",
		QueryString:      "error | count by _sourceCategory",
		RunByReceiptTime: false,
		ParsingMode:      "Manual",
		TimeRange: map[string]interface{}{
			"type": "BeginBoundedTimeRange",
			"from": map[string]interface{}{
				"type":      "LiteralTimeRange",
				"rangeName": "today",
			},
		},
		Schedule: &LogSearchSchedule{
			CronExpression:  "0 0 6 ? * 3 *",
			MuteErrorEmails: false,
			TimeZone:        "America/Los_Angeles",
			ScheduleType:    "Custom",
			ParseableTimeRange: map[string]interface{}{
				"type": "BeginBoundedTimeRange",
				"from": map[string]interface{}{
					"type":         "RelativeTimeRange",
					"relativeTime": "-15m",
				},
			},
			Threshold: &SearchNotificationThreshold{
				ThresholdType: "group",
				Operator:      "gt",
				Count:         10,
			},
			Notifications: []interface{}{
				map[string]interface{}{
					"taskType":             "EmailSearchNotificationSyncDefinition",
					"toList":               []string{"tf_multi_notif_email_webhook@sumologic.com"},
					"subjectTemplate":      "Alert: {{TriggerCondition}} for {{SearchName}}",
					"includeQuery":         true,
					"includeResultSet":     true,
					"includeHistogram":     true,
					"includeCsvAttachment": false,
				},
				map[string]interface{}{
					"taskType":  "WebhookSearchNotificationSyncDefinition",
					"webhookId": "0000000000005678",
					"payload":   "{\"alertName\": \"{{SearchName}}\"}",
				},
			},
		},
	}

	// Step 1: CREATE
	t.Log("=== Step 1: Creating LogSearch with multi-notification (Email + Webhook) ===")
	id, err := client.CreateLogSearch(logSearch)
	if err != nil {
		t.Fatalf("CreateLogSearch failed: %v", err)
	}
	if id != resourceID {
		t.Fatalf("Expected resource ID %s, got %s", resourceID, id)
	}
	t.Logf("CREATE successful. Resource ID: %s", id)

	// Step 2: GET (Read back)
	t.Log("=== Step 2: Reading LogSearch back via GET API ===")
	fetched, err := client.GetLogSearch(id)
	if err != nil {
		t.Fatalf("GetLogSearch failed: %v", err)
	}

	if fetched.Name != logSearch.Name {
		t.Fatalf("Name mismatch: expected %q, got %q", logSearch.Name, fetched.Name)
	}

	// Step 3: Verify notifications in API response
	t.Log("=== Step 3: Verifying notifications array in API GET response ===")
	if fetched.Schedule == nil {
		t.Fatal("Fetched LogSearch has no schedule")
	}

	notifications := fetched.Schedule.Notifications
	if notifications == nil || len(notifications) == 0 {
		t.Fatal("API GET response missing 'notifications' field")
	}

	if len(notifications) != 2 {
		t.Fatalf("Expected 2 notifications in GET response, got %d", len(notifications))
	}

	// Verify notification types in the response
	foundEmail := false
	foundWebhook := false
	for _, n := range notifications {
		notifMap, ok := n.(map[string]interface{})
		if !ok {
			t.Fatalf("Notification is not a map: %T", n)
		}
		taskType, _ := notifMap["taskType"].(string)
		switch taskType {
		case "EmailSearchNotificationSyncDefinition":
			foundEmail = true
			toList := notifMap["toList"]
			t.Logf("  Email notification found: toList=%v", toList)
		case "WebhookSearchNotificationSyncDefinition":
			foundWebhook = true
			webhookId := notifMap["webhookId"]
			t.Logf("  Webhook notification found: webhookId=%v", webhookId)
		default:
			t.Logf("  Unknown notification type: %s", taskType)
		}
	}

	if !foundEmail {
		t.Fatal("GET response missing Email notification")
	}
	if !foundWebhook {
		t.Fatal("GET response missing Webhook notification")
	}

	// Step 4: Verify the notification field is NOT present (omitempty)
	t.Log("=== Step 4: Verifying 'notification' (singular) is omitted ===")
	if fetched.Schedule.Notification != nil {
		t.Fatalf("Expected 'notification' to be nil when 'notifications' is used, got: %v", fetched.Schedule.Notification)
	}
	t.Log("  Confirmed: 'notification' field is nil (omitempty)")

	// Step 5: Delete
	t.Log("=== Step 5: Deleting LogSearch ===")
	err = client.DeleteLogSearch(id)
	if err != nil {
		t.Fatalf("DeleteLogSearch failed: %v", err)
	}
	t.Log("DELETE successful.")

	t.Log("")
	t.Log("=== ALL VERIFICATIONS PASSED ===")
	t.Log("Summary:")
	t.Log("  1. Created scheduled search with notifications array containing Email + Webhook")
	t.Log("  2. Verified CREATE request body has correct structure (notifications array, no singular notification)")
	t.Log("  3. GET response correctly returns both notifications in array")
	t.Log("  4. Email notification has taskType=EmailSearchNotificationSyncDefinition")
	t.Log("  5. Webhook notification has taskType=WebhookSearchNotificationSyncDefinition")
	t.Log("  6. Singular 'notification' field is properly omitted")
	t.Log("  7. Resource deleted successfully")
}

func TestLogSearch_MultiNotificationEmailAndWebhook_ResourceConversion(t *testing.T) {
	// This test verifies the Terraform resource → Go struct conversion for multi-notification
	t.Log("=== Testing resourceToLogSearchSchedule with multi-notification ===")

	// Simulate the Terraform resource data structure that represents a schedule
	// with two notifications blocks (email + webhook)
	scheduleData := map[string]interface{}{
		"cron_expression":  "0 0 6 ? * 3 *",
		"mute_error_emails": false,
		"time_zone":        "America/Los_Angeles",
		"schedule_type":    "Custom",
		"notification":     []interface{}{}, // Empty - not used
		"notifications": []interface{}{
			map[string]interface{}{
				"email_search_notification": []interface{}{
					map[string]interface{}{
						"to_list":                []interface{}{"test@sumologic.com"},
						"subject_template":       "Alert: {{SearchName}}",
						"include_query":          true,
						"include_result_set":     true,
						"include_histogram":      true,
						"include_csv_attachment": false,
					},
				},
				"webhook_search_notification":     []interface{}{},
				"save_to_view_notification":       []interface{}{},
				"save_to_lookup_notification":     []interface{}{},
				"service_now_search_notification": []interface{}{},
				"alert_search_notification":       []interface{}{},
				"cse_signal_notification":         []interface{}{},
			},
			map[string]interface{}{
				"email_search_notification": []interface{}{},
				"webhook_search_notification": []interface{}{
					map[string]interface{}{
						"webhook_id":         "0000000000005678",
						"payload":            "{\"alertName\": \"{{SearchName}}\"}",
						"itemize_alerts":     false,
						"max_itemized_alerts": 0,
					},
				},
				"save_to_view_notification":       []interface{}{},
				"save_to_lookup_notification":     []interface{}{},
				"service_now_search_notification": []interface{}{},
				"alert_search_notification":       []interface{}{},
				"cse_signal_notification":         []interface{}{},
			},
		},
		"parameter": []interface{}{},
		"parseable_time_range": []interface{}{
			map[string]interface{}{
				"begin_bounded_time_range": []interface{}{
					map[string]interface{}{
						"from": []interface{}{
							map[string]interface{}{
								"relative_time_range": []interface{}{
									map[string]interface{}{
										"relative_time": "-15m",
									},
								},
								"epoch_time_range":   []interface{}{},
								"literal_time_range": []interface{}{},
								"iso8601_time_range": []interface{}{},
							},
						},
						"to": []interface{}{},
					},
				},
				"complete_literal_time_range": []interface{}{},
			},
		},
		"threshold": []interface{}{
			map[string]interface{}{
				"threshold_type": "group",
				"operator":       "gt",
				"count":          10,
			},
		},
	}

	schedule := resourceToLogSearchSchedule([]interface{}{scheduleData})

	// Verify the schedule was converted correctly
	if schedule == nil {
		t.Fatal("resourceToLogSearchSchedule returned nil")
	}

	t.Logf("Schedule converted: cronExpression=%s, scheduleType=%s", schedule.CronExpression, schedule.ScheduleType)

	// Verify notifications
	if schedule.Notifications == nil || len(schedule.Notifications) == 0 {
		t.Fatal("Schedule has no notifications after conversion")
	}

	if len(schedule.Notifications) != 2 {
		t.Fatalf("Expected 2 notifications, got %d", len(schedule.Notifications))
	}

	t.Logf("Notifications count: %d", len(schedule.Notifications))

	// Verify first notification is an email type
	emailNotif, ok := schedule.Notifications[0].(EmailSearchNotification)
	if !ok {
		t.Fatalf("First notification is not EmailSearchNotification: %T", schedule.Notifications[0])
	}
	if emailNotif.TaskType != "EmailSearchNotificationSyncDefinition" {
		t.Fatalf("First notification should be Email, got taskType=%v", emailNotif.TaskType)
	}
	t.Logf("  Notification 1: Email (taskType=%s, toList=%v)", emailNotif.TaskType, emailNotif.ToList)

	// Verify second notification is a webhook type
	webhookNotif, ok := schedule.Notifications[1].(WebhookSearchNotification)
	if !ok {
		t.Fatalf("Second notification is not WebhookSearchNotification: %T", schedule.Notifications[1])
	}
	if webhookNotif.TaskType != "WebhookSearchNotificationSyncDefinition" {
		t.Fatalf("Second notification should be Webhook, got taskType=%v", webhookNotif.TaskType)
	}
	t.Logf("  Notification 2: Webhook (taskType=%s, webhookId=%s)", webhookNotif.TaskType, webhookNotif.WebhookId)

	// Verify notification (singular) is nil
	if schedule.Notification != nil {
		t.Fatalf("Singular 'notification' should be nil, got: %v", schedule.Notification)
	}
	t.Log("  Singular 'notification' is nil as expected")

	// Now test the JSON serialization
	t.Log("=== Verifying JSON serialization ===")
	jsonBytes, err := json.Marshal(schedule)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}

	var serialized map[string]interface{}
	json.Unmarshal(jsonBytes, &serialized)

	if val, hasNotification := serialized["notification"]; hasNotification && val != nil {
		t.Fatal("JSON should NOT have non-nil 'notification' key (omitempty)")
	}
	t.Log("  JSON omits/nulls 'notification' field (omitempty)")

	notifs, ok := serialized["notifications"].([]interface{})
	if !ok || len(notifs) != 2 {
		t.Fatalf("JSON should have 'notifications' array with 2 items, got: %v", serialized["notifications"])
	}
	t.Log("  JSON has 'notifications' array with 2 items")

	t.Log("")
	t.Log("=== RESOURCE CONVERSION TEST PASSED ===")
	t.Log("Summary:")
	t.Log("  1. Terraform HCL data → Go struct conversion works for multi-notification")
	t.Log("  2. Email notification correctly converted with taskType=EmailSearchNotificationSyncDefinition")
	t.Log("  3. Webhook notification correctly converted with taskType=WebhookSearchNotificationSyncDefinition")
	t.Log("  4. JSON serialization correctly omits singular 'notification' and includes 'notifications' array")
}
