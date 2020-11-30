---
layout: "sumologic"
page_title: "SumoLogic: sumologic_subdomain"
description: |-
  Provides a Sumologic Subdomain
---

# sumologic_lookup_table
Provides a [Sumologic Subdomain][1].

## Example Usage
```hcl
resource "sumologic_subdomain" "exampleSubdomain" {
    subdomain = "my-company"
}
```

## Argument reference

The following arguments are supported:

- `subdomain` - (Required) The subdomain.

## Attributes reference

The following attributes are exported:

- `id` - Unique identifier for the subdomain.

[1]: https://help.sumologic.com/Manage/01Account_Usage/05Manage_Organization#change-account-subdomain
