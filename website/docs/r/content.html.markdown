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
            "viewStartTime": "2025-01-01T00:00:00Z",
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
            "viewStartTime": "2025-01-01T00:00:00Z",
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
        "description": "Runs every hour with timerange of 15m and sends a webhook notification"
    })
}
```

## Example Scheduled Search with Save To View/Index Notification
```hcl
data "sumologic_personal_folder" "personalFolder" {}
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
            "viewStartTime": "2025-01-01T00:00:00Z",
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
        "description": "Runs every hour with timerange of 15m and saves to a view"
    })
}
```

## Example Scheduled Search with Save To Lookup Notification
```hcl
data "sumologic_personal_folder" "personalFolder" {}
variable "email" {
  description = "Email to be used in the library path"
  type = string
  default = "user@company.com"
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
            "viewStartTime": "2025-01-01T00:00:00Z",
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
        "description": "Runs every hour with timerange of 15m and saves to a lookup table."
    })
}
```

## Example Scheduled Search with Multiple Notifications
Use the `notifications` array in `searchSchedule` to send more than one notification from a single scheduled search.
Each entry takes the same shape as a single `notification` object. Exactly one of `notification` or `notifications`
must be specified. By default, up to 10 notifications are allowed per scheduled search.
```hcl
data "sumologic_personal_folder" "personalFolder" {}
resource "sumologic_connection" "connection" {
  type            = "WebhookConnection"
  name            = "test-connection"
  description     = "My description"
  url             = "https://connection-endpoint.com"
  default_payload = <<JSON
{
  "client" : "Sumo Logic",
  "eventType" : "{{Name}}",
  "search_url" : "{{QueryUrl}}"
}
JSON
  webhook_type    = "Webhook"
}
resource "sumologic_content" "test" {
    parent_id = "${data.sumologic_personal_folder.personalFolder.id}"
    config = jsonencode({
        "type": "SavedSearchWithScheduleSyncDefinition",
        "name": "Scheduled Search with Multiple Notifications Test",
        "search": {
            "queryText": "\"warn\"",
            "defaultTimeRange": "-15m",
            "byReceiptTime": false,
            "viewName": "",
            "viewStartTime": "2025-01-01T00:00:00Z",
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
                    "webhookId": "${sumologic_connection.connection.id}",
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
    })
}
```

## Example Dashboard
Use a `DashboardV2SyncDefinition` to manage a dashboard through the `sumologic_content` resource. The UI has no
"Export as Terraform" option, so the usual way to build `config` is to open the dashboard in the Sumo Logic UI,
choose **More > Export**, and paste the exported JSON into `jsonencode(...)`. To create dashboards from native
Terraform blocks instead, use the [`sumologic_dashboard`](dashboard.html) resource.

The `config` follows the dashboard request model, with these differences:
- `type` and `name` are required, because `config` is a content definition.
- `title`, `theme` and `timeRange` are required.
- The folder is set with `parent_id`, so `folderId` is not part of `config`.
- `domain`, `hierarchies`, `isPublic`, `highlightViolations` and `organizations` are not supported in `config`.
- Panel `visualSettings` and layout `structure` are JSON encoded as strings.
```hcl
data "sumologic_personal_folder" "personalFolder" {}

resource "sumologic_content" "dashboard" {
    parent_id = "${data.sumologic_personal_folder.personalFolder.id}"
    config = jsonencode({
        "type": "DashboardV2SyncDefinition",
        "name": "API Health Dashboard",
        "title": "API Health Dashboard",
        "description": "Errors, CPU and latency for the API service",
        "theme": "Light",
        "topologyLabelMap": {
            "data": {}
        },
        "refreshInterval": 300,
        "timeRange": {
            "type": "BeginBoundedTimeRange",
            "from": {
                "type": "RelativeTimeRangeBoundary",
                "relativeTime": "-1h"
            },
            "to": null
        },
        "layout": {
            "layoutType": "Grid",
            "layoutStructures": [
                {
                    "key": "panelA1B2C3D4E5F60001",
                    "structure": "{\"height\":3,\"width\":24,\"x\":0,\"y\":0}"
                },
                {
                    "key": "panelA1B2C3D4E5F60002",
                    "structure": "{\"height\":8,\"width\":12,\"x\":0,\"y\":3}"
                },
                {
                    "key": "panelA1B2C3D4E5F60003",
                    "structure": "{\"height\":8,\"width\":12,\"x\":12,\"y\":3}"
                },
                {
                    "key": "panelA1B2C3D4E5F60004",
                    "structure": "{\"height\":8,\"width\":24,\"x\":0,\"y\":11}"
                }
            ]
        },
        "panels": [
            {
                "id": null,
                "key": "panelA1B2C3D4E5F60001",
                "title": "About",
                "visualSettings": "{\"text\":{\"format\":\"markdownV2\"}}",
                "keepVisualSettingsConsistentWithParent": true,
                "panelType": "TextPanel",
                "text": "Errors, CPU and latency for the API service."
            },
            {
                "id": null,
                "key": "panelA1B2C3D4E5F60002",
                "title": "API errors",
                "visualSettings": "{\"general\":{\"mode\":\"timeSeries\",\"type\":\"line\"}}",
                "keepVisualSettingsConsistentWithParent": true,
                "panelType": "SumoSearchPanel",
                "queries": [
                    {
                        "transient": false,
                        "queryString": "_sourceCategory=api/{{env}} error | timeslice 5m | count by _timeslice",
                        "queryType": "Logs",
                        "queryKey": "A",
                        "metricsQueryMode": null,
                        "metricsQueryData": null,
                        "tracesQueryData": null,
                        "spansQueryData": null,
                        "parseMode": "Auto",
                        "timeSource": "Message"
                    }
                ],
                "description": "Error count in 5 minute buckets",
                "timeRange": null,
                "coloringRules": null,
                "linkedDashboards": []
            },
            {
                "id": null,
                "key": "panelA1B2C3D4E5F60003",
                "title": "CPU user",
                "visualSettings": "{\"general\":{\"mode\":\"timeSeries\",\"type\":\"line\"}}",
                "keepVisualSettingsConsistentWithParent": true,
                "panelType": "SumoSearchPanel",
                "queries": [
                    {
                        "transient": false,
                        "queryString": "metric=CPU_user _sourceCategory=api/{{env}} | avg by _sourceHost",
                        "queryType": "Metrics",
                        "queryKey": "A",
                        "metricsQueryMode": "Basic",
                        "metricsQueryData": {
                            "metric": "CPU_user",
                            "aggregationType": "Average",
                            "groupBy": "_sourceHost",
                            "filters": [
                                {
                                    "key": "_sourceCategory",
                                    "value": "api/{{env}}",
                                    "negation": false
                                }
                            ],
                            "operators": []
                        },
                        "tracesQueryData": null,
                        "spansQueryData": null,
                        "outputCardinalityLimit": 1000
                    }
                ],
                "description": "",
                "timeRange": null,
                "coloringRules": null,
                "linkedDashboards": []
            },
            {
                "id": null,
                "key": "panelA1B2C3D4E5F60004",
                "title": "Recommendations latency (p95)",
                "visualSettings": "{\"general\":{\"mode\":\"timeSeries\",\"type\":\"line\"}}",
                "keepVisualSettingsConsistentWithParent": true,
                "panelType": "SumoSearchPanel",
                "queries": [
                    {
                        "transient": false,
                        "queryString": "* | (endtimestamp - starttimestamp) / 1000 as duration | where toLowerCase(operation) matches toLowerCase(\"*/api/recommendations\") | timeslice 10m | pct(duration, 95) as duration_pct_95 by _timeslice",
                        "queryType": "Spans",
                        "queryKey": "A",
                        "metricsQueryMode": null,
                        "metricsQueryData": null,
                        "tracesQueryData": null,
                        "spansQueryData": {
                            "filters": [
                                {
                                    "type": "KeyValuePair",
                                    "fieldName": "operation",
                                    "operator": "=",
                                    "fieldValue": "*/api/recommendations"
                                }
                            ],
                            "visualizations": [
                                {
                                    "type": "calculation",
                                    "name": "duration_pct_95",
                                    "field": "duration",
                                    "aggregator": {
                                        "key": "pct",
                                        "percentile": 95.0
                                    }
                                }
                            ],
                            "groupBy": [
                                {
                                    "type": "time",
                                    "fieldValue": "10m"
                                }
                            ],
                            "limit": []
                        },
                        "parseMode": "Auto",
                        "timeSource": "Message",
                        "outputCardinalityLimit": 1000
                    }
                ],
                "description": "",
                "timeRange": {
                    "type": "BeginBoundedTimeRange",
                    "from": {
                        "type": "RelativeTimeRangeBoundary",
                        "relativeTime": "-1d"
                    },
                    "to": null
                },
                "coloringRules": null,
                "linkedDashboards": []
            }
        ],
        "variables": [
            {
                "id": null,
                "name": "env",
                "displayName": "env",
                "defaultValue": "prod",
                "sourceDefinition": {
                    "variableSourceType": "CsvVariableSourceDefinition",
                    "values": "prod,staging,dev"
                },
                "allowMultiSelect": false,
                "includeAllOption": false,
                "hideFromUI": false,
                "valueType": "Any"
            }
        ],
        "coloringRules": []
    })
}
```

## Argument reference

The following arguments are supported:

- `parent_id` - (Required) The identifier of the folder to import into. Identifiers from the Library in the Sumo user interface are provided in decimal format which is incompatible with Terraform. The identifier needs to be in hexadecimal format.
- `config` - (Required) JSON block for the content to import. NOTE: Updating the name will create a new object and leave a untracked content item (delete the existing content item and create a new content item if you want to update the name).
  For a `SavedSearchWithScheduleSyncDefinition`, `searchSchedule` takes either a single `notification` object or a
  `notifications` array of notification objects (up to 10 by default), but not both. If both are specified,
  `terraform apply` fails with an error from the Content import API; the conflict is not caught at plan time.

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
