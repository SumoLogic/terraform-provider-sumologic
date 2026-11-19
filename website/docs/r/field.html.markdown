---
layout: "sumologic"
page_title: "SumoLogic: sumologic_field"
description: |-
  Provides a Sumologic Field
---

# sumologic_field
Provides a [Sumologic Field][1]. See also the [Fields help page][2] and the [Field Management API][3].

## Example Usage

### Basic

```hcl
resource "sumologic_field" "basic" {
  field_name = "mystring"
}
```

### Disabled field

A disabled field is defined but its incoming values are dropped at ingest.

```hcl
resource "sumologic_field" "disabled" {
  field_name = "my_other_string"
  state      = "Disabled"
}
```

## Argument reference

The following arguments are supported:

- `field_name` - (Required, ForceNew) Name of the field. It is the key in the key-value pair attached to your logs. Between 1 and 255 characters. Changing it destroys and recreates the field.
- `state` - (Optional) Whether the field is enabled. Allowed values: `Enabled`, `Disabled`. A field has to be enabled for its values to be kept: while it is `Disabled`, Sumo Logic drops its incoming values at ingest (historical values are not removed and stay searchable). Fields created through the API are enabled by default. A disabled field is no longer assigned to new log data, searches that still use it do not return it, and built-in fields cannot be disabled. It can take up to 10 minutes for a change to take effect on your data. Any other value fails with an error when the resource is updated.

## Attributes reference

In addition to the arguments above, the following attributes are exported:

- `id` - Identifier of the field.
- `field_id` - Identifier of the field (the same value as `id`).
- `data_type` - Data type of the field, as reported by the Sumo Logic API (for example `String`, `Long`, `Int`, `Double` or `Boolean`). It is determined by Sumo Logic and is not meant to be set in configuration.

## Import
Fields can be imported using the field id, e.g.:

```hcl
terraform import sumologic_field.field 000000000ABC1234
```

## Destroy

Destroying the resource deletes the field. Deleting a field does not delete historical data assigned with that field. Fields share an account-level quota with fields created from field extraction rules; see [Fields][2] for current limits. Deleting a field fails while it is still referenced by a field extraction rule, role, partition, collector or source, so remove those references first. Consider setting `state = "Disabled"` and checking that nothing breaks before you delete a field.

[1]: https://help.sumologic.com/Manage/Fields
[2]: https://www.sumologic.com/help/docs/manage/fields/
[3]: https://api.sumologic.com/docs/#tag/fieldManagementV1
