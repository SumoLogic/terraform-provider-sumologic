---
layout: "sumologic"
page_title: "SumoLogic: sumologic_token"
description: |-
  Provides a Sumologic Token
---

# sumologic_token
Provides a [Sumologic Token][1].

## Example Usage
```hcl
resource "sumologic_token" "example_token" {
  name          = "testToken"
  description   = "Testing resource sumologic_token"
  status        = "Active"
  type          = "CollectorRegistration"
}
```

## Argument Reference

The following arguments are supported:

  * `type` - (Required) Type of the token. Valid value:
    - `CollectorRegistration`.
  * `name` - (Required) Display name of the token. This must be unique across all of the tokens.
  * `description` - (Optional) The description of the token.
  * `status` - (Required) Status of the token. Valid values:
    - `Active`
    - `Inactive`
  
The following attributes are exported:

  * `id` - The internal ID of the token.
  * `encodedTokenAndUrl` - The encoded token for collector registration.

## Import
Tokens can be imported using the name, e.g.:

```hcl
terraform import sumologic_token.test id
```

[1]: https://help.sumologic.com/Manage/Security/Installation_Tokens
