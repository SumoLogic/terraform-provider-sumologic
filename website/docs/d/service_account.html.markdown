---
layout: "sumologic"
page_title: "SumoLogic: sumologic_service_account"
description: |-
  Provides a way to retrieve Sumo Logic service account details for a service account managed outside of terraform.
---

# sumologic_service_account

Provides a way to retrieve Sumo Logic service account details (id, name, roles, etc) for a service account managed outside of terraform.

## Example Usage

### By id
```hcl
data "sumologic_service_account" "by_id" {
  id = "1234567890"
}
```

### By email
```hcl
data "sumologic_service_account" "by_email" {
  email = "ci-bot@example.com"
}
```

A service account is looked up by exactly one of `id` or `email`. Looking up by `email` fails if no service account has that exact email.

## Attributes reference

The following attributes are exported:

- `id` - The ID of the service account.
- `email` - Email address of the service account.
- `name` - Name of the service account.
- `is_active` - `true` if the service account is active, `false` if it has been deactivated.
- `role_ids` - Set of IDs of the roles assigned to the service account.
