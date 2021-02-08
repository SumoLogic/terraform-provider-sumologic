---
layout: "sumologic"
page_title: "SumoLogic: sumologic_field_extraction_rule"
description: |-
  Provides a Sumologic Field Extraction Rule
---

# sumologic_field_extraction_rule
Provides a [Sumologic Field Extraction Rule][1].

## Example Usage
```hcl
resource "sumologic_field_extraction_rule" "fieldExtractionRule" {
      name = "exampleFieldExtractionRule"
      scope = "_sourceHost=127.0.0.1"
      parse_expression = "csv _raw extract 1 as f1"
      enabled = true
}
```

## Argument reference

The following arguments are supported:

- `name` - (Required) Name of the field extraction rule. Use a name that makes it easy to identify the rule.
- `scope` - (Required) Scope of the field extraction rule. This could be a sourceCategory, sourceHost, or any other metadata that describes the data you want to extract from. Think of the Scope as the first portion of an ad hoc search, before the first pipe ( | ). You'll use the Scope to run a search against the rule.
- `parse_expression` - (Required) Describes the fields to be parsed.
- `enabled` - (Required) Is the field extraction rule enabled.

## Attributes reference

The following attributes are exported:

- `id` - Unique identifier for the field extraction rule.

## Import
Extraction Rules can be imported using the extraction rule id, e.g.:

```hcl
terraform import sumologic_field_extraction_rule.fieldExtractionRule id
```

[1]: https://help.sumologic.com/Manage/Field-Extractions
