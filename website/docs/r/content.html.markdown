---
layout: "sumologic"
page_title: "SumoLogic: sumologic_content"
description: |-
  Provides a way to interact with Sumologic Content
---

# sumologic_content
Provides a way to interact with Sumologic Content.

## Example Usage
```hcl
data "sumologic_personal_folder" "personalFolder" {}

resource "sumologic_content" "test" {
parent_id = "${data.sumologic_personal_folder.personalFolder.id}"
config = 
    jsonencode({
        "type": "SavedSearchWithScheduleSyncDefinition",
        "name": "test-333",
        "search": {
            "queryText": "\"warn\"",
            "defaultTimeRange": "-15m",
            "byReceiptTime": false,
            "viewName": "",
            "viewStartTime": "1970-01-01T00:00:00Z",
            "queryParameters": [],
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
    })
}
```

## Argument reference

The following arguments are supported:

- `parent_id` - (Required) The identifier of the folder to import into. Identifiers from the Library in the Sumo user interface are provided in decimal format which is incompatible with Terraform. The identifier needs to be in hexadecimal format.
- `config` - (Required) JSON block for the content to import.

## Attributes reference

The following attributes are exported:

- `id` - Unique identifier for the contnet item.

[1]: https://help.sumologic.com/APIs/Content-Management-API
