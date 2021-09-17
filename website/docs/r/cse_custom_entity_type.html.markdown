---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_custom_entity_type"
description: |-
  Provides a CSE Custom Entity Type
---

# custom_entity_type
Provides a CSE Custom Entity Type.

## Example Usage
```hcl
resource "sumologic_cse_custom_entity_type" "custom_entity_type" {
  name = "New Custom Entity Type"
  identifier = "identifier"
  fields =["file_hash_md5", "file_hash_sha1"]
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) Human friend and unique name. Examples: "Ip Address", "Username", "Mac Address".
- `identifier` - (Required) Machine friendly and unique identifier. Examples: "ip", "username", "mac".
- `fields` - (Required) Record schema fields. Examples: "file_hash_md5", "file_hash_sha1".".


The following attributes are exported:

- `id` - The internal ID of the custom entity type.


