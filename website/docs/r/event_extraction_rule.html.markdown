---
layout: "sumologic"
page_title: "SumoLogic: sumologic_event_extraction_rule [Beta feature]"
description: |-
  Provides a Sumo Logic Event Extraction Rule
---

# sumologic_event_extraction_rule
Provides a **Sumo Logic Event Extraction Rule**, which allows you to extract structured **Events** from log data and ingest them into Sumo Logic’s Events index.

Event Extraction Rules are commonly used to capture **deployment events, configuration changes, feature flag changes, and infrastructure changes** from logs.

---

## Example Usage

### Basic Event Extraction Rule

```hcl
resource "sumologic_event_extraction_rule" "deployment_event" {
  name  = "deployment-event"
  query = "_sourceCategory=deployments"

  configuration = {
    eventType = {
      value_source = "Deployment"
    }
    eventPriority = {
      value_source = "High"
    }
    eventSource = {
      value_source = "Jenkins"
    }
    eventName = {
      value_source = "monitor-manager deployed"
    }
  }
}
```

### Event Extraction Rule with Correlation Expression

```hcl
resource "sumologic_event_extraction_rule" "deployment_event" {
  name        = "deployment-event"
  description = "Captures deployment events from Jenkins logs"
  query       = "_sourceCategory=deployments | json \"version\""
  enabled     = true

  correlation_expression {
    query_field_name          = "version"
    event_field_name          = "version"
    string_matching_algorithm = "ExactMatch"
  }

  configuration = {
    eventType = {
      value_source = "Deployment"
    }
    eventPriority = {
      value_source = "High"
    }
    eventSource = {
      value_source = "Jenkins"
    }
    eventName = {
      value_source = "monitor-manager deployed"
    }
    eventDescription = {
      value_source = "2 containers were upgraded"
    }
  }
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) Name of the Event Extraction Rule.
- `description` - (Optional) Description of the rule.
- `query` - (Required) Log query used to extract events.
- `enabled` - (Optional) Whether the macro will be enabled. Default true.
- `correlation_expression` - (Block List, Optional) Specifies how to determine related events for a log search query. See [correlationExpression schema](#schema-for-correlationExpression) for details.
- `configuration` - (Block List, Optional) Defines how event fields are mapped to their corresponding values. See [configuration schema](#schema-for-configuration)
for details.

## Attributes reference
In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Event extraction rule.

### Schema for `correlationExpression`
- `query_field_name` - (Required) Name of the field returned by the log query.
- `event_field_name` - (Required) Name of the corresponding event field.
- `string_matching_algorithm` - (Required) Algorithm used to match values.

### Schema for `configuration`

- `configuration` - (Required, Map) This object defines how event fields are mapped to their corresponding values. Each field specifies a valueSource, which provides the actual value, and an optional mappingType, indicating the value is hardcoded.

#### Configuration field mapping

- `value_source` - (Required) The value or extracted field used for the event field.
- `mapping_type` - (Optional) Specifies how the value is mapped.
  Supported value: `HardCoded` (default)

#### Required event fields

The following keys **must** be defined inside the `configuration` map:

- `eventType` - (Required) Type of the event.
  Accepted values: `Deployment`, `Feature Flag Change`, `Configuration Change`, `Infrastructure Change`
- `eventPriority` - (Required) Priority of the event.
  Accepted values: `High`, `Medium`, `Low`
- `eventSource` - (Required) Source system where the event originated (for example, `Jenkins`).
- `eventName` - (Required) Human-readable name of the event.

#### Optional event fields

- `eventDescription` - (Optional) Additional context or details about the event.
- Custom fields - (Optional) Additional keys can be added to capture domain-specific event metadata.
