---
layout: "sumologic"
page_title: "SumoLogic: sumologic_user"
description: |-
  Provides a Sumologic User
---

# sumologic_user
Provides a [Sumologic User][1].

A user belongs to the organization and gets its permissions from the roles listed in `role_ids`.

## Example Usage

### Basic
```hcl
resource "sumologic_role" "analyst" {
  name         = "Analyst"
  description  = "Read-only access to production logs"
  capabilities = ["viewCollectors"]
}

resource "sumologic_user" "jon_doe" {
  first_name  = "Jon"
  last_name   = "Doe"
  email       = "jon.doe@example.com"
  is_active   = true
  role_ids    = [sumologic_role.analyst.id]
  transfer_to = ""
}
```

### Deactivated user
Set `is_active = false` to deactivate a user without deleting them.
```hcl
resource "sumologic_role" "analyst" {
  name         = "Analyst"
  description  = "Read-only access to production logs"
  capabilities = ["viewCollectors"]
}

resource "sumologic_user" "former_employee" {
  first_name  = "Jane"
  last_name   = "Smith"
  email       = "jane.smith@example.com"
  is_active   = false
  role_ids    = [sumologic_role.analyst.id]
  transfer_to = ""
}
```

### Transfer content to another user on deletion
When `sumologic_user.jane_smith` is destroyed, its content is transferred to `sumologic_user.jon_doe`.
```hcl
resource "sumologic_role" "analyst" {
  name         = "Analyst"
  description  = "Read-only access to production logs"
  capabilities = ["viewCollectors"]
}

resource "sumologic_user" "jon_doe" {
  first_name  = "Jon"
  last_name   = "Doe"
  email       = "jon.doe@example.com"
  is_active   = true
  role_ids    = [sumologic_role.analyst.id]
  transfer_to = ""
}

resource "sumologic_user" "jane_smith" {
  first_name  = "Jane"
  last_name   = "Smith"
  email       = "jane.smith@example.com"
  is_active   = true
  role_ids    = [sumologic_role.analyst.id]
  transfer_to = sumologic_user.jon_doe.id
}
```

## Argument reference

The following arguments are supported:

- `first_name` - (Required) First name of the user. 1 to 128 characters.
- `last_name` - (Required) Last name of the user. Up to 128 characters.
- `email` - (Required) Email address of the user. Up to 255 characters. It cannot be changed after the user is created: changing it makes `terraform apply` fail with an error. See [Transferred content and email updates](#transferred-content-and-email-updates) for a workaround.
- `is_active` - (Required) `true` if the user is active, `false` if the user has been deactivated. Changing it requires the `manageUsersAndRoles` capability. Deactivating a user does not delete them, and they can be reactivated later. While a user is deactivated, their dashboards, scheduled searches, scheduled views and access keys are deactivated too, and monitors keep running (see [Deactivate or Activate a User][3]).
- `role_ids` - (Required) Set of IDs of the roles assigned to the user. Changing it requires the `manageUsersAndRoles` capability.
- `transfer_to` - (Required) ID of the user who receives this user's content when this user is deleted. It can be an empty string (`""`), in which case the content is transferred to the user who runs the delete (the user whose access key the provider uses). The value is only used on deletion, and the value saved in state is the one used, so it must be applied before the user is destroyed to take effect. The provider never deletes a user's content.

The following attributes are exported:

- `id` - The ID of the user.

See the [User Management API][4] for the underlying operations.

## Import
Users can be imported using the user id, e.g.:

```hcl
terraform import sumologic_user.jon_doe 1234567890
```

`transfer_to` is not stored by Sumo Logic. After an import, set it in your configuration and apply before destroying the user.

## Transferred content and email updates

When a user is deleted, all of that user's content is transferred to another user. If `transfer_to` is
set to another user's ID, then the content will be assigned to that user. If `transfer_to` is empty,
then it will instead be assigned to the user executing the delete operation.

Deleting a user is permanent, and also removes the user's access keys. The provider never deletes the
user's content. Transferred scheduled views, dashboards, scheduled searches and monitors run as the new owner,
and the new owner needs the Manage Views capability for scheduled views to transfer successfully. See
[Delete a User][2] for everything that happens on deletion.

A user's email address may not be changed through this resource (a change fails the apply). The Sumo Logic API and UI can request an email change, which sends an activation link to the new address, but the provider does not manage it. As a workaround, you may:

1. create a new `sumologic_user` with the desired email address
2. set the `transfer_to` of the existing user to the new user's ID
3. delete the user with the old email address

[1]: https://www.sumologic.com/help/docs/manage/users-roles/users/
[2]: https://www.sumologic.com/help/docs/manage/users-roles/users/delete-user/
[3]: https://www.sumologic.com/help/docs/manage/users-roles/users/deactivate-activate-user/
[4]: https://api.sumologic.com/docs/#tag/userManagement
