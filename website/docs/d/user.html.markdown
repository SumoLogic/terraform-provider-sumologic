---
layout: "sumologic"
page_title: "SumoLogic: sumologic_user"
description: |-
Provides a way to retrieve Sumo Logic user details (id, email, etc) for a user managed by another terraform stack.
---

# sumologic_role

Provides a way to retrieve Sumo Logic user details (id, email, etc) for a user
managed by another terraform stack.


## Example Usage
```hcl
data "sumologic_role" "this" {
  id = "1234567890"
}
```

```hcl
data "sumologic_role" "that" {
  email = "user@example.com"
}
```

A user can be looked up by either `id` or `email`. One of those attributes needs to be specified.

If both `id` and `email` have been specified, `id` takes precedence.

## Attributes reference

The following attributes are exported:

- `id` - The internal ID of the role. This can be used to create users having that role.
- `email` - (Required) Email of the user.
