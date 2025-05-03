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
- `email` - (Required) Email of the user.
- `is_active` - (Required) This has the value true if the user is active and false if they have been deactivated.
- `role_ids` - (Required) List of roleIds associated with the user.
- `transfer_to` - (Required) UserId of user to transfer this user's content to on deletion, can be empty. Must be applied prior to deletion to take effect.

The following attributes are exported:

- `id` - The internal ID of the user.

## Import
Users can be imported using the user id, e.g.:

```hcl
terraform import sumologic_user.user 1234567890
```

## Transfered content and email updates

When a user is deleted, all of that user's content is transferred to another user. If `transfer_to` is
set to another user's ID, then the content will be assigned to that user. If `transfer_to` is empty,
then it will instead be assigned to the user executing the delete operation.

A user's email address may not be changed. As a workaround, you may:

1. create a new `sumologic_user` with the desired email address
2. set the `transfer_to` of the existing user to the new user's ID
3. delete the user with the old email address

[1]: https://help.sumologic.com/Manage/Users-and-Roles/Manage-Users
