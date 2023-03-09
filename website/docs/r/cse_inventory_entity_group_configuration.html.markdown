---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_inventory_entity_group_configuration"
description: |-
  Provides a CSE Inventory Entity Group Configuration
---

## Example Usage
```hcl
resource "sumologic_cse_inventory_entity_group_configuration" "inventory_entity_group_configuration" {
	criticality = "HIGH"
    description = "Inventory entity group description"
	inventory_type = "username"
	inventory_source = "Active Directory"
	inventory_key = "groups"
	inventory_value = "http_servers"
	name = "Inventory entity group configuration"
	suppressed = false
 	tags = ["tag"]
}
```

## Argument reference

The following arguments are supported:

- `criticality` - (Optional) The entity group configuration criticality Examples: "HIGH", "CRITICALITY".
- `description` - (Optional) The entity group configuration description.
- `group` - (Optional)(Deprecated) The entity group configuration inventory group. The field `group` is deprecated and will be removed in a future release of the provider -- please make usage of `inventory_key`, `inventory_value`  instead.
- `inventory_type` - (Required) The inventory type Examples: "computer", "username".
- `inventory_source` - (Required) The inventory source Examples: "Active Directory", "Okta".
- `inventory_key` - (Required) The inventory key to apply configuration Examples: "groups", "normalizedHostname", "normalizedComputerName".
- `inventory_value` - (Optional) The inventory value to match.
- `dynamic_tags` - (Optional) If dynamic tags are enabled for configuration.
- `tag_schema` - (Optional) The tag schema to be used for dynamic tags.
- `name` - (Required) The entity group configuration name.
- `suppresed` - (Optional) The entity group configuration suppressed value 
- `tags` - (Optional) The entity group configuration tags list.

The following attributes are exported:

- `id` - The internal ID of the entity group configuration.

## Import

Inventory Entity Group Configuration can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_inventory_entity_group_configuration.inventory_entity_group_configuration id
```