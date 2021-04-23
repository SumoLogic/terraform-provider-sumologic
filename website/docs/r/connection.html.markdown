---
layout: 'sumologic'
page_title: 'SumoLogic: sumologic_connection'
description: |-
  Provides the ability to create, read, delete, update connections.
---

# sumologic_connection

Provides the ability to create, read, delete, update connections.

## Example Usage

```hcl
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
  webhook_type    = "Webhook"
}
```

## Argument reference

The following arguments are supported:

- `type` - (Required) Type of connection. Only `WebhookDefinition` is implimented right now.
- `name` - (Required) Name of connection. Name should be a valid alphanumeric value.
- `description` - (Optional) Description of the connection.
- `url` - (Required) URL for the webhook connection.
- `headers` - (Optional) Map of access authorization headers.
- `custom_headers` - (Optional) Map of custom webhook headers
- `default_payload` - (Required) Default payload of the webhook.
- `webhook_type` - (Optional) Type of webhook. Valid values are `AWSLambda`, `Azure`, `Datadog`, `HipChat`, `PagerDuty`, `Slack`, `Webhook`, `NewRelic`, and `MicrosoftTeams`. Default: `Webhook`

Additional data provided in state

- `id` - (Computed) The Id for this connection.

## Import

Connections can be imported using the connection id, e.g.:

```hcl
terraform import sumologic_connection.test 1234567890
```
