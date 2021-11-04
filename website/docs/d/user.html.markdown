---
layout: "sumologic"
page_title: "SumoLogic: sumologic_user"
description: |-
Provides a way to retrieve Sumo Logic user details (id, email, etc) for a user managed outside of terraform.
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
- `first_name` - (Required) First name of the user.
- `last_name` - (Required) Last name of the user.
- `is_active` - (Required) This has the value true if the user is active and false if they have been deactivated.
- `role_ids` - (Required) List of roleIds associated with the user.
