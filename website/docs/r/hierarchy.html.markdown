---
layout: "sumologic"
page_title: "SumoLogic: sumologic_hierarchy"
description: |-
  Provides a Sumologic Hierarchy
---

# sumologic_hierarchy
Provides a [Sumologic Hierarchy][1].

## Example Usage
```hcl
resource "sumologic_hierarchy" "example_hierarchy" {
  name = "testK8sHierarchy"
  filter {
    key   = "_origin"
    value = "kubernetes" 
  }
  level {
    entity_type = "cluster"
    next_levels_with_conditions {
      condition = "testCondition"
      level {
        entity_type = "namespace"
      }
    }
    next_level {
      entity_type = "node"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) Name of the hierarchy.
- `filter` - (Optional) An optional clause that a hierarchy requires to be matched.
  + `key` - (Required) Filtering key.
  + `value` - (Required) Value required for the filtering key.
- `level` - (Required) A hierarchy of entities. The order is up-down, left to right levels with condition, then level without condition. Maximum supported total depth is 6.
  + `entity_type` - (Required) Indicates the name and type for all entities at this hierarchy level, e.g. service or pod in case of kubernetes entities.
  + `next_levels_with_conditions` - (Optional) Zero or more next levels with conditions.
    + `condition` - (Required) Condition to be checked against for level.entityType value, for now full string match.
    + `level` - (Required)
  + `next_level` - (Optional) Next level without a condition.
  
The following attributes are exported:

- `id` - The internal ID of the hierarchy.

## Import
Hierarchies can be imported using the id, e.g.:

```hcl
terraform import sumologic_hierarchy.test id
```

[1]: https://help.sumologic.com/Visualizations-and-Alerts/Explore
