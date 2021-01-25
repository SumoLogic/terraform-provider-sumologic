---
layout: "sumologic"
page_title: "SumoLogic: sumologic_user"
description: |-
  Provides a Sumologic User
---

# sumologic_user
Provides a [Sumologic User][1].

## Example Usage
```hcl
resource "sumologic_role" "example_role" {
  name        = "TestRole123"
  description = "Testing resource sumologic_role"

  lifecycle {
    ignore_changes = ["users"]
  }
}

resource "sumologic_user" "example_user1" {
  first_name   = "Jon"
  last_name    = "Doe"
  email        = "jon.doe@gmail.com"
  is_active    = true
  role_ids     = [sumologic_role.example_role.id]
  transfer_to  = ""
}

resource "sumologic_user" "example_user2" {
  first_name   = "Jane"
  last_name    = "Smith"
  email        = "jane.smith@gmail.com"
  role_ids     = [sumologic_role.example_role.id]
  transfer_to  = sumologic_user.example_user1.id
}
```

## Argument reference

The following arguments are supported:

- `first_name` - (Required) First name of the user.
- `last_name` - (Required) Last name of the user.
- `email` - (Required) Last name of the user.
- `is_active` - (Required) This has the value true if the user is active and false if they have been deactivated..
- `role_ids` - (Required) List of roleIds associated with the user.
- `transfer_to` - (Required) UserId of user to transfer this user's content to on deletion, can be empty. Must be applied prior to deletion to take effect.

The following attributes are exported:

- `id` - The internal ID of the user.

## Import
Users can be imported using the user id, e.g.:

```hcl
terraform import sumologic_user.user 1234567890
```

[1]: https://help.sumologic.com/Manage/Users-and-Roles/Manage-Users
