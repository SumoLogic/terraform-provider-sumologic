---
layout: "sumologic"
page_title: "SumoLogic: sumologic_content"
description: |-
  Provides a way to interact with Sumologic Content
---

# sumologic_content
Provides a way to interact with Sumologic Content.

~> **NOTE:** If working with many content items it is recommended to reduce [Terraform Parallelism](https://www.terraform.io/docs/cli/commands/apply.html#parallelism-n) to 2 in order to not be rate limited.

## Example Scheduled Search with Email Notification 
```hcl
data "sumologic_personal_folder" "personalFolder" {}

resource "sumologic_content" "test" {
    parent_id = "${data.sumologic_personal_folder.personalFolder.id}"
    config = jsonencode({
        "type": "SavedSearchWithScheduleSyncDefinition",
        "name": "Scheduled Search with Email Notificiation Test",
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
            "threshold": {
                "operator": "gt",
                "count": 0
            },
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
    })
}
```

## Example Scheduled Search with Webhook Notification
```hcl
data "sumologic_personal_folder" "personalFolder" {}

resource "sumologic_connection" "connection" { 
  type        = "WebhookConnection"
  name        = "test-connection"
  description = "My description"
  url         = "https://connection-endpoint.com"
  headers = {
    "X-Header" : "my-header"
  }
  custom_headers = {
    "X-custom" : "my-custom-header"
  }
  default_payload = <<JSON
{
  "client" : "Sumo Logic",
  "eventType" : "{{Name}}",
  "description" : "{{Description}}",
  "search_url" : "{{QueryUrl}}",
  "num_records" : "{{NumRawResults}}",
  "search_results" : "{{AggregateResultsJson}}"
}
JSON
  resolution_payload = <<JSON
{
  "client" : "Sumo Logic",
  "eventType" : "{{Name}}",
  "description" : "{{Description}}",
  "search_url" : "{{QueryUrl}}"
}
JSON
  webhook_type    = "Webhook"
}

resource "sumologic_content" "test" {
    parent_id = "${data.sumologic_personal_folder.personalFolder.id}"
    config = jsonencode({
        "type": "SavedSearchWithScheduleSyncDefinition",
        "name": "Scheduled Search with Webhook Notification Test",
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
            "threshold": {
                "operator": "gt",
                "count": 0
            },
            "notification": {
              "taskType": "WebhookSearchNotificationSyncDefinition",
              "webhookId": "${sumologic_connection.connection.id}",
              "payload": "{}",
              "itemizeAlerts": false,
              "maxItemizedAlerts": 50
            },
            "scheduleType": "1Hour",
            "muteErrorEmails": false,
            "parameters": []
        },
        "description": "Runs every hour with timerange of 15m and sends email notifications"
    })
}
```

## Example Scheduled Search with Save To View/Index Notification 
```hcl
data "sumologic_personal_folder" "personalFolder" {}

resource "sumologic_scheduled_view" "test_view" {
  index_name = "test_view"
  query = <<QUERY
_view=connections connectionStats
| parse "connectionStats.CS *" as body
| json field=body "exitCode", "isHttp2"
| lookup org_name from shared/partners on partner_id=partnerid
| timeslice 10m
QUERY
  start_time = "2019-09-01T00:00:00Z"
  retention_period = 365
  lifecycle {
    prevent_destroy = true
    ignore_changes = [index_id]
  }
}

resource "sumologic_content" "test" {
    parent_id = "${data.sumologic_personal_folder.personalFolder.id}"
    config = jsonencode({
        "type": "SavedSearchWithScheduleSyncDefinition",
        "name": "Scheduled Search with Save To View Notification Test",
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
            "threshold": {
                "operator": "gt",
                "count": 0
            },
            "notification": {
              "taskType": "SaveToViewNotificationSyncDefinition",
              "viewName": "test_view",
            },
            "scheduleType": "1Hour",
            "muteErrorEmails": false,
            "parameters": []
        },
        "description": "Runs every hour with timerange of 15m and sends email notifications"
    })
}
```

## Example Scheduled Search with Save To Lookup  Notification
```hcl
data "sumologic_personal_folder" "personalFolder" {}

variable "email" {
  description = "Email to be used in the library path"
  type = string
  default = "shahzaib.ali@sumologic.com"
}

resource "sumologic_lookup_table" "lookupTable" {
  name = "test_lookup_table"
  fields {
    field_name = "host"
    field_type = "string"
  }
  ttl               = 100
  primary_keys      = ["host"]
  parent_folder_id  = "${data.sumologic_personal_folder.personalFolder.id}"
  size_limit_action = "DeleteOldData"
  description       = "some description"
}


resource "sumologic_content" "test" {
    parent_id = "${data.sumologic_personal_folder.personalFolder.id}"
    config = jsonencode({
        "type": "SavedSearchWithScheduleSyncDefinition",
        "name": "Scheduled Search with Save To Lookup Notification Test",
        "search": {
            "queryText": "\"warn\" fields host",
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
            "threshold": {
                "operator": "gt",
                "count": 0
            },
            "notification": {
              "taskType": "SaveToLookupNotificationSyncDefinition",
              "lookupFilePath": "/Library/Users/${var.email}/test_lookup_table",
              "isLookupMergeOperation": false,
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
- `config` - (Required) JSON block for the content to import. NOTE: Updating the name will create a new object and leave a untracked content item (delete the existing content item and create a new content item if you want to update the name).

### Timeouts

`sumologic_content` provides the following [Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `read` - (Default `1 minute`) Used for waiting for the import job to be successful
- `create` - (Default `10 minutes`) Used for waiting for the import job to be successful
- `update` - (Default `10 minutes`) Used for waiting for the import job to be successful
- `delete` - (Default `1 minute`) Used for waiting for the deletion job to be successful

## Attributes reference

The following attributes are exported:

- `id` - Unique identifier for the content item.

[1]: https://help.sumologic.com/APIs/Content-Management-API
