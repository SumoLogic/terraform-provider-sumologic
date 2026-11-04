---
layout: "sumologic"
page_title: "SumoLogic: sumologic_parser"
description: |-
  Provides a Sumologic Parser
---

# sumologic_parser
Provides a [Sumologic Parser][1] in the Parsers library.

## Example Usage

```hcl
resource "sumologic_parser" "pan_firewall" {
  name        = "PAN-firewall"
  description = "Parser for Palo Alto Networks firewall logs"
  stanzas     = <<-EOT
    [sourcetype:PAN-firewall]
    EVENT_MULTILINE = False
    FORMAT = CSV
  EOT
}
```

Parser in an existing folder, with local stanzas:

```hcl
resource "sumologic_parser" "custom" {
  name          = "custom-parser"
  description   = "Custom parser"
  parent_id     = "0000000000000001"
  stanzas       = "[sourcetype:custom]\nFORMAT = CSV"
  local_stanzas = "[sourcetype:custom]\nWRAPPER = BSD_SYSLOG"
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) Name of the parser. Between 1 and 255 characters.
- `description` - (Optional) Description of the parser. Up to 4096 characters. Defaults to an empty string.
- `stanzas` - (Required) Collection of stanzas describing the parser.
- `local_stanzas` - (Optional) Localized stanzas.
- `model_path` - (Optional) The path to the Model a Model-Connector is associated with.
- `sourcetype_path` - (Optional) The path to the sourcetype a Model-Connector is associated with.
- `families` - (Optional) Comma-separated list of model families this parser belongs to or applies to.
- `is_partial` - (Optional) Whether this is only a config fragment rather than a complete parser or Model-Connector. Defaults to `false`.
- `parent_id` - (Optional) Identifier of the parent folder in the Parsers library, for example `0000000000000001`. Defaults to the root folder. Changing it moves the parser.

~> **NOTE:** Locked parsers cannot be updated or moved. This resource does not lock or unlock parsers.

## Attributes reference

In addition to all arguments above, the following attributes are exported:

- `id` - The identifier of the parser.
- `version` - The current version of the parser.
- `content_type` - The content type of the object (`Parser`).
- `created_at` - Creation timestamp (RFC3339, UTC).
- `created_by` - Identifier of the user who created the parser.
- `modified_at` - Last modification timestamp (RFC3339, UTC).
- `modified_by` - Identifier of the user who last modified the parser.
- `is_locked` - Whether the parser is locked.
- `is_system` - Whether the parser is provided by Sumo Logic.
- `is_mutable` - Whether the parser can be modified.

## Import

Parsers can be imported using the parser id, e.g.:

```
terraform import sumologic_parser.pan_firewall 0000000000000001
```

[1]: https://help.sumologic.com/docs/send-data/
