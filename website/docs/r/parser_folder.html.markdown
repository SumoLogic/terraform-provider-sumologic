---
layout: "sumologic"
page_title: "SumoLogic: sumologic_parser_folder"
description: |-
  Provides a Sumologic Parser Folder
---

# sumologic_parser_folder
Provides a folder in the Sumo Logic Parsers library. Use its `id` as the `parent_id` of a
[`sumologic_parser`](parser.html) or of another `sumologic_parser_folder` to organize parsers.

## Example Usage

```hcl
resource "sumologic_parser_folder" "firewalls" {
  name        = "Firewalls"
  description = "Parsers for firewall logs"
}

# A nested folder
resource "sumologic_parser_folder" "palo_alto" {
  name      = "Palo Alto"
  parent_id = sumologic_parser_folder.firewalls.id
}

# A parser inside the nested folder
resource "sumologic_parser" "pan_firewall" {
  name      = "PAN-firewall"
  parent_id = sumologic_parser_folder.palo_alto.id
  stanzas   = "[sourcetype:PAN-firewall]\nFORMAT = CSV"
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) Name of the folder. Between 1 and 255 characters.
- `description` - (Optional) Description of the folder. Up to 4096 characters. Defaults to an empty string.
- `parent_id` - (Optional) Identifier of the parent folder in the Parsers library. Defaults to the root folder. Changing it moves the folder.

~> **NOTE:** Locked folders cannot be updated or moved. This resource does not lock or unlock folders.

~> **WARNING:** Deleting a folder in the Parsers library also deletes everything inside it. To prevent accidental data loss, this resource refuses to delete a folder that still contains items, and returns an error with the number of items it contains. Parsers and folders managed with Terraform that reference this folder through `parent_id` are destroyed first, so a normal `terraform destroy` works. If the folder contains items that Terraform does not manage, delete or move those items first.

## Attributes reference

In addition to all arguments above, the following attributes are exported:

- `id` - The identifier of the folder.
- `version` - The current version of the folder.
- `content_type` - The content type of the object (`Folder`).
- `created_at` - Creation timestamp (RFC3339, UTC).
- `created_by` - Identifier of the user who created the folder.
- `modified_at` - Last modification timestamp (RFC3339, UTC).
- `modified_by` - Identifier of the user who last modified the folder.
- `is_locked` - Whether the folder is locked.
- `is_system` - Whether the folder is provided by Sumo Logic.
- `is_mutable` - Whether the folder can be modified.

## Import

Parser folders can be imported using the folder id, e.g.:

```
terraform import sumologic_parser_folder.firewalls 0000000000000001
```
