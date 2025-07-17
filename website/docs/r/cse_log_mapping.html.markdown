---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_log_mapping"
description: |-
  Provides a Sumologic CSE Log Mapping
---

# log_mapping
Provides a Sumologic CSE Log Mapping.

## Example Usage
```hcl
resource "sumologic_cse_log_mapping" "log_mapping" {
  name = "New Log Mapping"
  product_guid = "003d35b3-3ba8-4e93-8776-e5810b4e243e"
  record_type = "Audit"
  enabled = true
  relates_entities = true
  skipped_values = ["skipped"]
  fields {
    name = "action"
    value = "action"
    value_type = "constant"
    skipped_values = ["-"]
    default_value = ""
    format = "JSON"
    case_insensitive = false
    alternate_values = ["altValue"]
    time_zone = "UTC"
    split_delimiter = ","
    split_index = "0"
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

- `name` - (Required) The name of the log mapping.
- `parent_id` - (Optional) The id of the parent log mapping.
- `product_guid` - (Required) Product GUID.
- `record_type` - (Required) The record type to be created. (possible values: Audit, AuditChange, AuditFile, AuditResourceAccess, Authentication, AuthenticationPrivilegeEscalation, Canary, Email, Endpoint, EndpointModuleLoad, EndpointProcess, Network, NetworkDHCP, NetworkDNS, NetworkFlow, NetworkHTTP, NetworkProxy, Notification, NotificationVulnerability)
- `enabled` - (Required) Enabled flag.
- `relates_entities` - (Optional) Set to true to relate entities.
- `skipped_values` - (Optional) List of skipped values.
- `fields` - (Required) List of fields for the new log mapping. See [field_schema](#schema-for-field) for details.
- `structured_inputs` - (Optional, omit if unstructured_fields is defined) List of structured inputs for the new log mapping. See [structured_input_schema](#schema-for-structured_input) for details.
- `unstructured_fields` - (Optional, omit if structured_inputs is defined) Unstructured fields for the new log mapping. See [unstructured_field_schema](#schema-for-unstructured_field) for details.

### Schema for `field`
- `name` - (Required) Name of the field.
- `value` - (Optional) Value of the field.
- `value_type` - (Optional) The value type. Possible values: null (for standard mapping), constant, sumofield (for extracted mapping), format, join (for joined mapping), lookup, split, time
- `skipped_values` - (Optional) List of skipped values.
- `default_value` - (Optional) Default value of the field.
- `format` - (Optional) Format of the field. (JSON, Windows, Syslog, CEF, LEEF )
- `case_insensitive` - (Optional) Case insensitive flag.
- `alternate_values` - (Optional) List of alternate values.
- `time_zone` - (Optional) Time zone.
- `split_delimiter` - (Optional) Split delimiter to be used. (some example: ",", "-", "|")
- `split_index` - (Optional) The index value to select (starting at zero)
- `field_join` - (Optional) List of field join values.
- `join_delimiter` - (Optional) Join delimiter.
- `format_parameters` - (Optional) List of format parameters.
- `lookup` - (Optional) List of lookup key value pair for field. See [lookup_schema](#schema-for-lookup) for details.

### Schema for `lookup`
- `key` - (Required) Lookup key.
- `value` - (Required) Lookup value.

### Schema for `structured_input`
- `event_id_pattern` - (Required) Event id pattern.
- `log_format` - (Required) Log format. (JSON, Windows, Syslog, CEF, LEEF )
- `product` - (Required) Product name.
- `vendor` - (Required) Vendor name.

### Schema for `unstructured_field`
- `pattern_names` - (Required) List of grok pattern names.


The following attributes are exported:

- `id` - The internal ID of the log mapping.

## Import

Log Mapping can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_log_mapping.log_mapping id
```

