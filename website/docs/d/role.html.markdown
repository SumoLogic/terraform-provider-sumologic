---
layout: "sumologic"
page_title: "SumoLogic: sumologic_role"
description: |-
  Provides a way to retrieve Sumo Logic role details (id, names, etc) for a role managed by another terraform stack.
---

# sumologic_role

Provides a way to retrieve Sumo Logic role details (id, names, etc) for a role
managed by another terraform stack.


## Example Usage
```hcl
data "sumologic_role" "this" {
  name = "MyRole"
}
```

```hcl
data "sumologic_role" "that" {
  id = "1234567890"
}
```

A role can be looked up by either `id` or `name`. One of those attributes needs to be specified.

If both `id` and `name` have been specified, `id` takes precedence.

## Attributes reference

The following attributes are exported:

- `id` - The internal ID of the role. This can be used to create users having that role.
- `name` - The name of the role.
- `description` - The description of the role.
- `filter_predicate` - The search filter to restrict access to specific logs.
- `capabilities` - The list of capabilities associated with the role.
