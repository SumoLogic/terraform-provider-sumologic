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
	groups = ["admin"]
	inventory_type = "username"
	name = "Inventory entity group configuration"
	suppressed = false
 	tags = ["tag"]
}
```

## Argument reference

The following arguments are supported:

- `criticality` - (Optional) The entity group configuration criticality Examples: "HIGH", "CRITICALITY".
- `description` - (Optional) The entity group configuration description.
- `groups` - (Optional) The entity group configuration groups list.
- `inventory_type` - (Optional) The entity type Examples: "computer", "username".
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