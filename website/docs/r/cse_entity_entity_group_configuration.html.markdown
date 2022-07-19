---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_entity_entity_group_configuration"
description: |-
  Provides a CSE Entity Entity Group Configuration
---

## Example Usage
```hcl
resource "sumologic_cse_entity_entity_group_configuration" "entity_entity_group_configuration" {
	criticality = "HIGH"
    description = "Entity Group description"
	entity_namespace = "namespace"
	entity_type = "_hostname"
	name = "Hostaname entity group configuration"
	suffix = "red.co"
	suppressed = true
 	tags = ["tag"]
}
```

## Argument reference

The following arguments are supported:

- `criticality` - (Optional) The entity group configuration criticality Examples: "HIGH", "CRITICALITY".
- `description` - (Optional) The entity group configuration description.
- `entity_namespace` - (Optional) The entity namespace.
- `entity_type` - (Optional) The entity type Examples: "_ip", "_mac", "_username", "_hostname".
- `name` - (Required) The entity group configuration name.
- `network_block` - (Optional) The entity group configuration network block value Example: "192.168.0.0/16".
- `prefix` - (Optional) The entity group configuration prefix value.
- `suffix` - (Optional) The entity group configuration suffix value.
- `suppresed` - (Optional) The entity group configuration suppressed value 
- `tags` - (Optional) The entity group configuration tags list.

The following attributes are exported:

- `id` - The internal ID of the entity group configuration.

## Import

Entity Entity Group Configuration can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_entity_entity_group_configuration.entity_entity_group_configuration id
```