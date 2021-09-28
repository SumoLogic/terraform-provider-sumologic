---
layout: "sumologic"
page_title: "SumoLogic: sumologic_cse_custom_entity_type"
description: |-
  Provides a Sumologic CSE Custom Entity Type
---

# sumologic_cse_custom_entity_type
Provides a Sumologic CSE Custom Entity Type.

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

- `name` - (Required) Human friend and unique name. Example: "File Hash".
- `identifier` - (Required) Machine friendly and unique identifier. Example: "filehash".
- `fields` - (Required) Record schema fields. Examples: "file_hash_md5", "file_hash_sha1".".


The following attributes are exported:

- `id` - The internal ID of the custom entity type.

## Import

Custom entity type can be imported using the field id, e.g.:
```hcl
terraform import sumologic_cse_custom_entity_type.custom_entity_type id
```
