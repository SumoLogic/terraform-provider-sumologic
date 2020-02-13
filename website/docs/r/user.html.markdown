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
  first_name = "Jon"
  last_name  = "Doe"
  email      = "jon.doe@gmail.com"
  active     = false
  role_ids   = ["${sumologic_role.example_role.id}"]
}

resource "sumologic_user" "example_user2" {
  first_name = "Jane"
  last_name  = "Smith"
  email      = "jane.smith@gmail.com"
  role_ids   = ["${sumologic_role.example_role.id}"]
}
```

## Argument reference

The following arguments are supported:

- `first_name` - (Required) First name of the user..
- `last_name` - (Required) Last name of the user.
- `email` - (Required) Last name of the user.
- `active` - (Optional) This has the value true if the user is active and false if they have been deactivated..
- `role_ids` - (Required) List of roleIds associated with the user.

The following attributes are exported:

- `id` - The internal ID of the user.


[Back to Index][0]

[0]: ../README.md
[1]: https://help.sumologic.com/Manage/Users-and-Roles/Manage-Users
