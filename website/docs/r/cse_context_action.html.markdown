---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_context_action"
description: |-
  Provides a Sumologic CSE Context Action
---

# sumologic_cse_context_action
Provides a Sumologic CSE Context Action.

## Example Usage
```hcl
resource "sumologic_cse_context_action" "context_action" {
	name = "Context Action Name"
	type = "URL"
	template = "https://bar.com/?q={{value}}"
	ioc_types = ["IP_ADDRESS"]
	entity_types = ["_hostname"]
	record_fields = ["request_url"]
	all_record_fields = false	
	enabled = true	
}

```

## Argument reference

The following arguments are supported:

- `name` - (Required) Context Action name.
- `type` - (Optional; defaults to URL) Context Action type. Valid values: "URL", "QUERY".
- `template` - (Optional) The URL/QUERY template.
- `ioc_types` - (Required) IOC Data types. Valid values: "ASN", "DOMAIN", "HASH", "IP_ADDRESS", "MAC_ADDRESS", "PORT", "RECORD_PROPERTY", "URL".
- `entity_types` - (Optional) Applicable to given entity types.
- `record_fields` - (Optional) Specific record fields.
- `all_record_fields` - (Optional; defaults to true) Use all record fields.
- `enabled` - (Optional; defaults to true) Whether the context action is enabled.

The following attributes are exported:

- `id` - The internal ID of the Context Action.

## Import

Context Action can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_contex_action.context_action id
```
