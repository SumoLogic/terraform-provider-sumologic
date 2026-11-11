---
layout: "sumologic"
page_title: "SumoLogic: sumologic_parser"
description: |-
  Provides a Sumologic Parser
---

# sumologic_parser
Provides a Sumologic Parser in the Parsers library (**Data Management > Logs > Parsers** in the Sumo Logic UI). Parsers are written in the Sumo Logic parsing language; see the [Parser Editor][1] and the [Parsers Library Management API][2] pages, and the [API reference][3].

## Example Usage

### Basic parser

```hcl
resource "sumologic_parser" "pan_firewall" {
  name        = "PAN-firewall"
  description = "Parser for Palo Alto Networks firewall logs"
  stanzas     = <<-EOT
    [parser]
    FORMAT = CSV
    MAPPER:vendor = Palo Alto Networks
    MAPPER:product = PAN-OS
    MAPPER:event_id = traffic
  EOT
}
```

### Parser in a folder managed by Terraform

Use the `id` of a [sumologic_parser_folder](parser_folder.html) as `parent_id`.

```hcl
resource "sumologic_parser_folder" "firewalls" {
  name = "Firewalls"
}

resource "sumologic_parser" "custom" {
  name        = "custom-parser"
  description = "Custom parser"
  parent_id   = sumologic_parser_folder.firewalls.id
  stanzas     = "[parser]\nFORMAT = CSV"
}
```

### Parser in an existing folder, with local stanzas

```hcl
resource "sumologic_parser" "custom_local" {
  name          = "custom-parser-local"
  description   = "Custom parser with local stanzas"
  parent_id     = "0000000000000001"
  stanzas       = "[parser]\nFORMAT = CSV"
  local_stanzas = "[sourcetype:custom]\nWRAPPER = BSD_SYSLOG"
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) Name of the parser. Between 1 and 255 characters.
- `description` - (Optional) Description of the parser. Up to 4096 characters. Defaults to an empty string.
- `stanzas` - (Required) Parser code, as a single string, written in the Sumo Logic parsing language (see the [Parsing Language Reference Guide](https://www.sumologic.com/help/docs/cse/schema/parsing-language-reference-guide/)). Code is organized into stanzas: one `[parser]` stanza, optional `[transform:<name>]` stanzas and an optional `[dependencies]` stanza. Attribute names are uppercase, for example `FORMAT = CSV`; `FORMAT` is one of `REGEX` (default), `CSV`, `JSON`, `XML`, `WINDOWS_XML`, `CEF` or `LEEF`. To map parsed logs to Cloud SIEM records, include `MAPPER:vendor`, `MAPPER:product` and `MAPPER:event_id` statements. The API reference examples use `[sourcetype:<name>]` stanzas. Leading and trailing whitespace is ignored when Terraform compares the configured value with the stored one.
- `local_stanzas` - (Optional) Localized stanzas, as a single string. Locked parsers and system parsers (provided by Sumo Logic) can only be changed through their localized fields, so use this argument to customize those. Leading and trailing whitespace is ignored when comparing.
- `model_path` - (Optional) The path to the Model a Model-Connector is associated with, for example `Models/firewalls/PAN`. Omit for a regular parser.
- `sourcetype_path` - (Optional) The path to the sourcetype a Model-Connector is associated with, for example `Parsers/firewalls/PAN-firewall`. Omit for a regular parser.
- `families` - (Optional) Comma-separated list of model families this parser belongs to or applies to, for example `firewalls,IDS`.
- `is_partial` - (Optional) Whether this is only a config fragment rather than a complete parser or Model-Connector. Defaults to `false`.
- `parent_id` - (Optional) Identifier of the parent folder in the Parsers library, for example `0000000000000001`. Defaults to the root folder of the Parsers library. Changing it moves the existing parser to the new folder; it does not recreate it. Use the `id` of a [sumologic_parser_folder](parser_folder.html) to manage the folder in Terraform.

~> **NOTE:** Locking requires the `LockParsers` capability. A locked parser cannot be moved or deleted, and only its localized fields can be modified; updating or moving it requires unlocking it first. This resource does not lock or unlock parsers, so `terraform destroy` fails for a locked parser until it is unlocked outside Terraform (the API states locked objects cannot be deleted).

~> **NOTE:** `terraform destroy` deletes the parser from the Parsers library.

## Attributes reference

In addition to all arguments above, the following attributes are exported:

- `id` - The identifier of the parser.
- `version` - The current version of the parser. Terraform sends the latest version automatically on each update.
- `content_type` - The content type of the object (`Parser`).
- `created_at` - Creation timestamp (RFC3339, UTC).
- `created_by` - Identifier of the user who created the parser.
- `modified_at` - Last modification timestamp (RFC3339, UTC).
- `modified_by` - Identifier of the user who last modified the parser.
- `is_locked` - Whether the parser is locked. Locked parsers cannot be moved or deleted, and only their localized fields can be modified.
- `is_system` - Whether the parser is provided by Sumo Logic. System parsers can only be localized; their non-local fields cannot be updated.
- `is_mutable` - Whether the parser can be modified. Immutable objects are read-only.

## Import

Parsers can be imported using the parser id, e.g.:

```
terraform import sumologic_parser.pan_firewall 0000000000000001
```

[1]: https://www.sumologic.com/help/docs/cse/schema/parser-editor/
[2]: https://www.sumologic.com/help/docs/api/parsers-library-management/
[3]: https://api.sumologic.com/docs/#tag/parsersLibraryManagement
