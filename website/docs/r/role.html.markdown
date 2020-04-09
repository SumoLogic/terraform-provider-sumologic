---
layout: "sumologic"
page_title: "SumoLogic: sumologic_role"
description: |-
  Provides a Sumologic Role
---

# sumologic_role
Provides a [Sumologic Role][1].

## Example Usage
```hcl
resource "sumologic_role" "example_role" {
  name        = "TestRole123"
  description = "Testing resource sumologic_role"

  filter_predicate = "_sourceCategory=Test"
  users = [
    "0000000000000001",
    "0000000000000002"
  ]
  capabilities = [
    "manageCollectors"
  ]
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) The name of the role.
- `description` - (Optional) The description of the role.
- `filter_predicate` - (Optional) A search filter to restrict access to specific logs.
- `users` - (Optional) List of user identifiers to assign the role to.
- `capabilities` - (Optional) List of capabilities associated with this role.
- `lookup_by_name` - (Optional) Configures an existent role using the same 'name' or creates a new one if non existent. Defaults to false.
- `destroy` - (Optional) Whether or not to delete the role in Sumo when it is removed from Terraform.  Defaults to true.

The following attributes are exported:

- `id` - The internal ID of the role.

[1]: https://help.sumologic.com/Manage/Users-and-Roles/Manage-Roles
