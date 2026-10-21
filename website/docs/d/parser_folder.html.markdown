---
layout: "sumologic"
page_title: "SumoLogic: sumologic_parser_folder"
description: |-
  Provides an easy way to retrieve a folder in the Parsers library.
---

# sumologic_parser_folder
Provides an easy way to retrieve a folder in the Parsers library by its path, name or ID. Use the `id` to
reference the folder as the `parent_id` of a [`sumologic_parser`](../r/parser.html).

Exactly one of `path`, `id` or `name` must be specified.

## Example Usage
```hcl
# Look up by path
data "sumologic_parser_folder" "firewalls_by_path" {
  path = "/Parsers/Firewalls"
}

# Look up by name
data "sumologic_parser_folder" "firewalls" {
  name = "Firewalls"
}

# Look up by name, restricted to a specific parent folder (when several folders share the name)
data "sumologic_parser_folder" "firewalls_nested" {
  name      = "Firewalls"
  parent_id = "0000000000000001"
}

# Look up by ID
data "sumologic_parser_folder" "by_id" {
  id = "0000000000000002"
}

resource "sumologic_parser" "pan_firewall" {
  name      = "PAN-firewall"
  parent_id = data.sumologic_parser_folder.firewalls.id
  stanzas   = "[sourcetype:PAN-firewall]\nFORMAT = CSV"
}
```

## Argument reference

The following arguments are supported:

- `path` - (Optional) The absolute path of the folder in the Parsers library, for example `/Parsers/Firewalls`. The path is passed to the API as given. Conflicts with `id` and `name`.
- `id` - (Optional) The ID of the folder. Conflicts with `path` and `name`.
- `name` - (Optional) The exact name of the folder. The lookup fails if no folder, or more than one folder, has this name. Conflicts with `path` and `id`.
- `parent_id` - (Optional) When looking up by `name`, only consider folders directly under this parent folder.

## Attributes reference

In addition to the arguments above, the following attributes are exported:

- `id` - The ID of the folder.
- `name` - The name of the folder.
- `path` - The path of the folder, when looked up by `path`.
- `parent_id` - The ID of the parent folder.
- `description` - The description of the folder.
- `is_locked` - Whether the folder is locked.
- `is_system` - Whether the folder is provided by Sumo Logic.
