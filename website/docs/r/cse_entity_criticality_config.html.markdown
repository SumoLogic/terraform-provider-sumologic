---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_entity_criticality_config"
description: |-
  Provides a CSE Entity Criticality Configuration
---

# sumologic_cse_entity_criticality_config
Provides a CSE Entity Criticality Configuration.

## Example Usage
```hcl
resource "sumologic_cse_entity_criticality_config" "entity_criticality_config" {
  name = "New Name"
  severity_expression = "severity + 2"
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) Human friendly and unique name. Examples: "Executive Laptop", "Bastion Host".
- `severity_expression` - (Required) Algebraic expression representing this entity\'s criticality. Examples: "severity * 2", "severity - 5", "severity / 3".


The following attributes are exported:

- `id` - The internal ID of the entity criticality configuration.

## Import

Entity criticality configuration can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_entity_criticality_config.entity_criticality_config id
```