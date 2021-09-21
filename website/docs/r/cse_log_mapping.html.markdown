---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_log_mapping"
description: |-
  Provides a CSE Log Mapping
---

# log_mapping
Provides a CSE Log Mapping.

## Example Usage
```hcl
resource "sumologic_cse_log_mapping" "log_mapping" {
  name = "New Log Mapping"
  product_guid = "003d35b3-3ba8-4e93-8776-e5810b4e243e"
  record_type = "Audit"
  enabled = "true"
  relates_entities = "true"
  skipped_values = ["skipped"]
  fields {
    name = "action"
    value = "action"
    value_type = "constant"
    skipped_values = ["-"]
    default_value = ""
    format = "JSON"
    case_insensitive = "false"
    alternate_values = ["altValue"]
    time_zone = "UTC"
    split_delimiter = ","
    split_index = "index"
    field_join = ["and"]
    join_delimiter = ""
    format_parameters = ["param"]
    lookup {
      key = "tunnel-up"
      value = "true"
    }
  }
  structured_inputs  {
    event_id_pattern = "vpn"
    log_format = "JSON"
    product = "fortinate"
    vendor = "fortinate"
  }
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) The name of the insights status.
- `description` - (Required) The description of the insights status.
- `expression` - (Required) Expression to match.
- `enabled` - (Required) Enabled flag.
- `exclude` - (Required) Set to true to exclude records that also match expression.
- `is_global` - (Required) Set to true if tuning expression intended to be global.
- `rule_ids` - (Required) List of rule ids, for the tuning expression to be applied. ( Empty if is_global set to true)


The following attributes are exported:

- `id` - The internal ID of the rule tuning expression.


