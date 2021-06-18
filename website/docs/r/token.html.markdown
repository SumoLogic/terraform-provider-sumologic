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
  version       = 0
}
```

## Argument Reference

The following arguments are supported:

  * `name` - (Required) Display name of the token. This must be unique across all of the tokens.
  * `description` - (Optional) The description of the token.
  * `status` - (Required) Status of the token. Supported values are `Active`, and `Inactive`.
  * `type` - (Required) Type of the token. Supported value is `CollectorRegistration`.
  * `version` - (Required for update) Version of the token. This is only required for update. It starts from 0 when created and gets incremented with each update.
  
The following attributes are exported:

  * `id` - The internal ID of the token. 

## Import
Tokens can be imported using the name, e.g.:

```hcl
terraform import sumologic_token tokenName
```

[1]: https://help.sumologic.com/Manage/Security/Installation_Tokens
