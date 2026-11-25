---
layout: "sumologic"
page_title: "SumoLogic: sumologic_service_account"
description: |-
  Provides a Sumologic Service Account
---

# sumologic_service_account
Provides a [Sumologic Service Account][1].

A service account is a non-human identity that owns access keys for automation. It gets its permissions from the roles listed in `role_ids`.

## Example Usage

### Basic
```hcl
resource "sumologic_role" "automation" {
  name         = "Automation"
  description  = "Read-only access for automation"
  capabilities = ["viewCollectors"]
}

resource "sumologic_service_account" "ci_bot" {
  name     = "ci-bot"
  email    = "ci-bot@example.com"
  role_ids = [sumologic_role.automation.id]
}
```

### Deactivated service account
The API always creates a service account as active. With `is_active = false` the provider creates it and then deactivates it.
```hcl
resource "sumologic_role" "automation" {
  name         = "Automation"
  description  = "Read-only access for automation"
  capabilities = ["viewCollectors"]
}

resource "sumologic_service_account" "paused_bot" {
  name      = "paused-bot"
  email     = "paused-bot@example.com"
  role_ids  = [sumologic_role.automation.id]
  is_active = false
}
```

### Transfer content to a user on deletion
When `sumologic_service_account.ci_bot` is destroyed, its content is transferred to `sumologic_user.owner`.
```hcl
resource "sumologic_role" "automation" {
  name         = "Automation"
  description  = "Read-only access for automation"
  capabilities = ["viewCollectors"]
}

resource "sumologic_user" "owner" {
  first_name  = "Jon"
  last_name   = "Doe"
  email       = "jon.doe@example.com"
  is_active   = true
  role_ids    = [sumologic_role.automation.id]
  transfer_to = ""
}

resource "sumologic_service_account" "ci_bot" {
  name        = "ci-bot"
  email       = "ci-bot@example.com"
  role_ids    = [sumologic_role.automation.id]
  transfer_to = sumologic_user.owner.id
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) Name of the service account. 1 to 128 characters.
- `email` - (Required) Email address of the service account. Up to 255 characters. It can be changed in place.
- `role_ids` - (Required) Set of IDs of the roles assigned to the service account. At least one role is required.
- `is_active` - (Optional) `true` if the service account is active, `false` to deactivate it. Defaults to `true`. A service account is always created active, so `false` makes the provider update it right after creation.
- `transfer_to` - (Optional) ID of a user or service account that receives this service account's content when it is destroyed. If both `transfer_to` and `delete_content` are unset, the content is transferred to the user who runs the delete (the owner of the access key the provider uses). Conflicts with `delete_content`. Only used on deletion; it is not stored by Sumo Logic.
- `delete_content` - (Optional) If `true`, the service account's content is deleted when it is destroyed. **The content is permanently deleted and cannot be recovered.** Defaults to `false`. Conflicts with `transfer_to`. Only used on deletion; it is not stored by Sumo Logic.

Deleting a service account also deletes its access keys.

The following attributes are exported:

- `id` - The ID of the service account.

See the [Service Account Management API][1] for the underlying operations.

## Import
Service accounts can be imported using the service account id, e.g.:

```hcl
terraform import sumologic_service_account.ci_bot 1234567890
```

`transfer_to` and `delete_content` are not stored by Sumo Logic. After an import, set them in your configuration and apply before destroying the service account.

[1]: https://api.sumologic.com/docs/#tag/serviceAccountManagement
